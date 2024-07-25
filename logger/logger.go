package logger

import (
	"net/http"
	"time"

	"go.uber.org/zap"
)

// Logger is a wrapper around zap.Logger
type Logger struct {
	zapLogger *zap.Logger
}

// NewLogger creates and returns a new Logger
func NewLogger() (*Logger, error) {
	zapLogger, err := zap.NewProduction()
	if err != nil {
		return nil, err
	}
	return &Logger{zapLogger: zapLogger}, nil
}

// Info logs an info message
func (l *Logger) Info(msg string, fields ...Field) {
	l.zapLogger.Info(msg, fields...)
}

// Error logs an error message
func (l *Logger) Error(msg string, fields ...Field) {
	l.zapLogger.Error(msg, fields...)
}

// Fatal logs a fatal message and then calls os.Exit(1)
func (l *Logger) Fatal(msg string, fields ...Field) {
	l.zapLogger.Fatal(msg, fields...)
}

// Sync flushes any buffered log entries
func (l *Logger) Sync() error {
	return l.zapLogger.Sync()
}

// Field is a wrapper around zap.Field
type Field = zap.Field

// Common field constructors
var (
	String   = zap.String
	Int      = zap.Int
	Duration = zap.Duration
	Error    = zap.Error
)

// LoggingMiddleware wraps an http.Handler and logs request details
func LoggingMiddleware(logger *Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		crw := &customResponseWriter{ResponseWriter: w, status: http.StatusOK}

		next.ServeHTTP(crw, r)

		duration := time.Since(start)

		logger.Info("Metrics request",
			String("method", r.Method),
			String("path", r.URL.Path),
			Int("status", crw.status),
			Duration("duration", duration),
			String("remote_addr", r.RemoteAddr),
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

// LogCollectionStart logs the start of a metrics collection
func LogCollectionStart(logger *Logger) {
	logger.Info("Starting metrics collection")
}

// LogCollectionEnd logs the end of a metrics collection
func LogCollectionEnd(logger *Logger, start time.Time) {
	duration := time.Since(start)
	logger.Info("Finished metrics collection", Duration("duration", duration))
}
