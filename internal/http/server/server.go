package server

import (
	"context"
	"log/slog"
	"net/http"
	"time"
)

type Server struct {
	server  *http.Server
	options *options
}

func New(address string, handler http.Handler, optFns ...OptionFn) *Server {
	options := &options{
		logger:          slog.Default(),
		shutdownTimeout: 5 * time.Second,
		readTimeout:     0,
		writeTimeout:    0,
		idleTimeout:     0,
	}

	for _, fn := range optFns {
		fn(options)
	}

	return &Server{
		server: &http.Server{
			Addr:         address,
			Handler:      handler,
			ReadTimeout:  options.readTimeout,
			WriteTimeout: options.writeTimeout,
			IdleTimeout:  options.idleTimeout,
		},
		options: options,
	}
}

// Run runs the server until ctx is done, or an error occurs.
func (s *Server) Run(ctx context.Context) error {
	logger := s.options.logger.
		With(slog.String("address", s.server.Addr), slog.String("server", "http"))

	serveErr := make(chan error, 1)
	go func() {
		logger.Info("listening")
		serveErr <- s.server.ListenAndServe()
	}()

	select {
	case err := <-serveErr:
		return err
	case <-ctx.Done():
	}

	logger.Info("shutting down", "timeout", s.options.shutdownTimeout)
	shutdownCtx, cancel := context.WithTimeout(context.Background(), s.options.shutdownTimeout)
	defer cancel()
	if err := s.server.Shutdown(shutdownCtx); err != nil {
		return err
	}
	logger.Info("shutdown completed")

	return nil
}

type options struct {
	logger          *slog.Logger
	shutdownTimeout time.Duration
	readTimeout     time.Duration
	writeTimeout    time.Duration
	idleTimeout     time.Duration
}

type OptionFn func(*options)

func WithLogger(logger *slog.Logger) OptionFn {
	return func(o *options) {
		o.logger = logger
	}
}

func WithShutdownTimeout(d time.Duration) OptionFn {
	return func(o *options) {
		o.shutdownTimeout = d
	}
}

func WithReadTimeout(d time.Duration) OptionFn {
	return func(o *options) {
		o.readTimeout = d
	}
}

func WithWriteTimeout(d time.Duration) OptionFn {
	return func(o *options) {
		o.writeTimeout = d
	}
}

func WithIdleTimeout(d time.Duration) OptionFn {
	return func(o *options) {
		o.idleTimeout = d
	}
}
