package log

import (
	"context"
	"log/slog"
)

type loggerCtx struct{}

// FromContext returns the logger associated with this context,
// or default logger if no value is associated.
func FromContext(ctx context.Context) *slog.Logger {
	if logger, ok := ctx.Value(loggerCtx{}).(*slog.Logger); ok {
		return logger
	}
	return slog.Default()
}

// WithContext returns a copy of parent in which logger is stored.
func WithContext(parent context.Context, logger *slog.Logger) context.Context {
	return context.WithValue(parent, loggerCtx{}, logger)
}

// Err returns an Attr with the key "error" and the error message as the value.
func Err(err error) slog.Attr {
	return slog.String("error", err.Error())
}
