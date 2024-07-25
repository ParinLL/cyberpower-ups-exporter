package logger

import (
	"net/http"
	"time"

	"go.uber.org/zap"
)

// Initialize creates and returns a new zap logger
func Initialize() (*zap.Logger, error) {
	logger, err := zap.NewProduction()
	if err != nil {
		return nil, err
	}
	return logger, nil
}

// LoggingMiddleware wraps an http.Handler with logging functionality
func LoggingMiddleware(logger *zap.Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Create a custom ResponseWriter to capture the status code
		crw := &customResponseWriter{ResponseWriter: w, status: http.StatusOK}

		next.ServeHTTP(crw, r)

		duration := time.Since(start)

		logger.Info("Metrics request",
			zap.String("method", r.Method),
			zap.String("path", r.URL.Path),
			zap.Int("status", crw.status),
			zap.Duration("duration", duration),
			zap.String("remote_addr", r.RemoteAddr),
		)
	})
}

type customResponseWriter struct {
	http.ResponseWriter
	status int
}

func (crw *customResponseWriter) WriteHeader(code int) {
	crw.status = code
	crw.ResponseWriter.WriteHeader(code)
}
