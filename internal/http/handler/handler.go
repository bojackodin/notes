package handler

import (
	"log/slog"
	"net/http"
	"runtime"
	"strings"
	"time"

	authcontroller "github.com/bojackodin/notes/internal/http/handler/auth"
	contexthelper "github.com/bojackodin/notes/internal/http/handler/context"
	notecontroller "github.com/bojackodin/notes/internal/http/handler/note"
	"github.com/bojackodin/notes/internal/http/httperror"
	"github.com/bojackodin/notes/internal/log"
	"github.com/bojackodin/notes/internal/service"

	"github.com/rs/xid"
)

func New(services *service.Services, optFns ...OptionFn) http.Handler {
	options := &options{
		logger: slog.Default(),
	}
	for _, fn := range optFns {
		fn(options)
	}

	mux := http.NewServeMux()

	{
		authctrl := authcontroller.New(services.Auth)

		mux.Handle("POST /sign-up", errorHandler(authctrl.SignUp))
		mux.Handle("POST /sign-in", errorHandler(authctrl.SignIn))
	}

	authMiddleware := &authMiddleware{services.Auth}

	{
		notectrl := notecontroller.New(services.Note)

		mux.Handle("GET /notes", errorHandler(authMiddleware.authenticate(notectrl.ListNotes)))
		mux.Handle("POST /notes", errorHandler(authMiddleware.authenticate(notectrl.CreateNote)))
	}

	handler := loggingMiddleware(options.logger)(mux)
	handler = recoveryMiddleware(options.logger)(handler)

	return handler
}

type options struct {
	logger *slog.Logger
}

type OptionFn func(*options)

func WithLogger(logger *slog.Logger) OptionFn {
	return func(o *options) {
		o.logger = logger
	}
}

func errorHandler(next func(w http.ResponseWriter, r *http.Request) error) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := next(w, r); err != nil {
			if code := httperror.HTTPStatus(err); code != 0 {
				httperror.RespondWithError(w, err.Error(), code)
			}
		}
	}
}

type loggingResponseWriter struct {
	http.ResponseWriter

	statusCode int
}

func (lrw *loggingResponseWriter) WriteHeader(statusCode int) {
	lrw.statusCode = statusCode
	lrw.ResponseWriter.WriteHeader(statusCode)
}

func loggingMiddleware(logger *slog.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var (
				requestID = xid.New().String()
				start     = time.Now()
				logger    = logger.With(slog.Group("request", "id", requestID))
			)

			w.Header().Add("X-Request-Id", requestID)
			r = r.WithContext(log.WithContext(r.Context(), logger))
			lw := loggingResponseWriter{ResponseWriter: w, statusCode: http.StatusOK}
			next.ServeHTTP(&lw, r)

			logger.LogAttrs(r.Context(), slog.LevelInfo, "handle request", slog.Group("request",
				slog.Duration("duration", time.Since(start)),
				slog.String("method", r.Method),
				slog.String("url", r.URL.String()),
				slog.String("user_agent", r.Header.Get("User-Agent")),
				slog.String("address", r.RemoteAddr),
				slog.Int("status", lw.statusCode),
			))
		})
	}
}

type authMiddleware struct {
	auth service.Auth
}

func (h *authMiddleware) authenticate(next func(w http.ResponseWriter, r *http.Request) error) func(w http.ResponseWriter, r *http.Request) error {
	return func(w http.ResponseWriter, r *http.Request) error {
		authorizationHeader := r.Header.Get("Authorization")

		if authorizationHeader == "" {
			return httperror.WithStatus(http.StatusUnauthorized)
		}

		headerParts := strings.Split(authorizationHeader, " ")
		if len(headerParts) != 2 || headerParts[0] != "Bearer" {
			return httperror.WithStatus(http.StatusUnauthorized)
		}

		token := headerParts[1]

		if len(token) == 0 {
			return httperror.WithStatus(http.StatusUnauthorized)
		}

		userID, err := h.auth.ParseToken(token)
		if err != nil {
			return httperror.WithStatus(http.StatusUnauthorized)
		}

		r = contexthelper.ContextSetUserID(r, userID)

		return next(w, r)
	}
}

func recoveryMiddleware(logger *slog.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if p := recover(); p != nil {
					stack := make([]byte, 2048)
					stack = stack[:runtime.Stack(stack, false)]
					logger.Error("recovered from panic:\n"+string(stack), "panic", p)
					w.WriteHeader(http.StatusInternalServerError)
				}
			}()

			next.ServeHTTP(w, r)
		})
	}
}
