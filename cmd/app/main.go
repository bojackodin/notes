package main

import (
	"context"
	"database/sql"
	"errors"
	"flag"
	"fmt"
	"io"
	"log/slog"
	"net"
	"os"
	"os/signal"
	"time"

	httphandler "github.com/bojackodin/notes/internal/http/handler"
	httpserver "github.com/bojackodin/notes/internal/http/server"
	"github.com/bojackodin/notes/internal/log"
	"github.com/bojackodin/notes/internal/repository"
	"github.com/bojackodin/notes/internal/service"
	"github.com/bojackodin/notes/internal/yandex/speller"

	"github.com/kelseyhightower/envconfig"
	_ "github.com/lib/pq"
	"gopkg.in/yaml.v2"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	if err := run(ctx, os.Stdout, os.Args); err != nil {
		slog.Error("run", log.Err(err))
		os.Exit(1)
	}
}

type config struct {
	Server struct {
		Host string `yaml:"host"`
		Port string `yaml:"port"`
		HTTP struct {
			ShutdownTimeout time.Duration `yaml:"shutdown_timeout" split_words:"true"`
			ReadTimeout     time.Duration `yaml:"read_timeout" split_words:"true"`
			WriteTimeout    time.Duration `yaml:"write_timeout" split_words:"true"`
			IdleTimeout     time.Duration `yaml:"idle_timeout" split_words:"true"`
		} `yaml:"http"`
	} `yaml:"server"`
	Logger struct {
		Level     string `yaml:"level"`
		AddSource bool   `yaml:"add_source" split_words:"true"`
	} `yaml:"logger"`
	Postgres struct {
		DSN          string        `yaml:"dsn"`
		MaxOpenConns int           `yaml:"max_open_connections" split_words:"true"`
		MaxIdleConns int           `yaml:"max_idle_connections" split_words:"true"`
		MaxIdleTime  time.Duration `yaml:"max_idle_time" split_words:"true"`
	} `yaml:"postgres"`
	JWT struct {
		Secret   string        `yaml:"secret"`
		TokenTTL time.Duration `yaml:"token_ttl" split_words:"true"`
	} `yaml:"jwt"`
}

func run(ctx context.Context, w io.Writer, args []string) (err error) {
	var configPath string

	flag := flag.NewFlagSet(args[0], flag.ExitOnError)
	flag.StringVar(&configPath, "config", "./etc/config.yml", "path")
	err = flag.Parse(args[1:])

	f, err := os.Open(configPath)
	if err != nil {
		return fmt.Errorf("open config file: %w", err)
	}
	defer f.Close()

	var cfg config

	if err = yaml.NewDecoder(f).Decode(&cfg); err != nil {
		return fmt.Errorf("parse config: %w", err)
	}
	if err = envconfig.Process("", &cfg); err != nil {
		var parseErr *envconfig.ParseError
		if errors.As(err, &parseErr) {
			err = fmt.Errorf("%v: expected value of type %v: failed to parse '%v'",
				parseErr.KeyName, parseErr.TypeName, parseErr.Value)
		}
		return fmt.Errorf("populate config with environment variables: %w", err)
	}

	logger, err := initLogger(w, &cfg)
	if err != nil {
		return err
	}

	db, err := connectDB(ctx, cfg)
	if err != nil {
		return fmt.Errorf("db: %w", err)
	}
	defer db.Close()

	repositories := repository.NewRepositories(db)

	deps := service.ServicesDependencies{
		Repositories: repositories,
		Speller:      speller.NewYandexSpeller(),
		Secret:       cfg.JWT.Secret,
		TokenTTL:     cfg.JWT.TokenTTL,
	}

	services := service.NewServices(deps)

	address := net.JoinHostPort(cfg.Server.Host, cfg.Server.Port)

	err = httpserver.New(
		address,
		httphandler.New(services, httphandler.WithLogger(logger)),
		httpserver.WithLogger(logger),
		httpserver.WithShutdownTimeout(cfg.Server.HTTP.ShutdownTimeout),
		httpserver.WithReadTimeout(cfg.Server.HTTP.ReadTimeout),
		httpserver.WithWriteTimeout(cfg.Server.HTTP.WriteTimeout),
		httpserver.WithIdleTimeout(cfg.Server.HTTP.IdleTimeout),
	).Run(ctx)

	return err
}

func connectDB(ctx context.Context, cfg config) (*sql.DB, error) {
	db, err := sql.Open("postgres", cfg.Postgres.DSN)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(cfg.Postgres.MaxOpenConns)
	db.SetMaxIdleConns(cfg.Postgres.MaxIdleConns)
	db.SetConnMaxIdleTime(cfg.Postgres.MaxIdleTime)

	err = db.PingContext(ctx)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func initLogger(w io.Writer, cfg *config) (*slog.Logger, error) {
	logOpts := &slog.HandlerOptions{
		AddSource: cfg.Logger.AddSource,
	}

	switch cfg.Logger.Level {
	case "debug":
		logOpts.Level = slog.LevelDebug
	case "info":
		logOpts.Level = slog.LevelInfo
	case "warn":
		logOpts.Level = slog.LevelWarn
	case "error":
		logOpts.Level = slog.LevelError
	default:
		return nil, fmt.Errorf("logger.level value must be one of [debug, info, warn, error]: '%v'", cfg.Logger.Level)

	}

	h := slog.NewJSONHandler(w, logOpts)

	logger := slog.New(h)
	slog.SetDefault(logger)

	return logger, nil
}
