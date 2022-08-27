package middleware

import (
	"net/http"
	"os"
	"time"

	"go.uber.org/zap"
)

func getLogger() *zap.Logger {
	option := zap.AddCallerSkip(1)
	logger, err := zap.NewProduction(option)
	if err != nil {
		panic("unable to initialize logger")
	}
	return logger
}

// responseWriter is a minimal wrapper for http.ResponseWriter that allows the
// written HTTP status code to be captured for logging.
type responseWriter struct {
	http.ResponseWriter
	status      int
	wroteHeader bool
}

func wrapResponseWriter(w http.ResponseWriter) *responseWriter {
	return &responseWriter{ResponseWriter: w}
}

func (rw *responseWriter) Status() int {
	return rw.status
}

func (rw *responseWriter) WriteHeader(code int) {
	if rw.wroteHeader {
		return
	}

	rw.status = code
	rw.ResponseWriter.WriteHeader(code)
	rw.wroteHeader = true
}

// RequestLoggerHandler logs the incoming HTTP request & its duration.
func RequestLoggerHandler(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				getLogger().Error(
					"error",
					zap.Any("err", err),
					zap.Stack("trace"),
				)
			}
		}()
		healthURL := "/health"
		if r.URL.Path == healthURL {
			return
		}
		start := time.Now()
		wrapped := wrapResponseWriter(w)
		handler.ServeHTTP(wrapped, r)
		rid := r.Header.Get("X-Request-ID")
		remoteAddr := r.Header.Get("X-Forwarded-For")
		if remoteAddr == "" {
			remoteAddr = r.RemoteAddr
		}
		hostname, err := os.Hostname()
		if err != nil {
			hostname = "localhost"
		}
		getLogger().Info(
			"http",
			zap.String("method", r.Method),
			zap.Int("status", wrapped.status),
			zap.String("path", r.URL.EscapedPath()),
			zap.Float64("duration", time.Since(start).Seconds()),
			zap.String("server", hostname),
			zap.String("remote", remoteAddr),
			zap.String("user-agent", r.Header.Get("User-Agent")),
			zap.String("request_id", rid),
		)
	})
}
