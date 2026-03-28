package middleware

import (
	"fmt"
	"net/http"
	"server/internal/log"
	"time"
)

type loggingResponseWriter struct {
	http.ResponseWriter
	statusCode   int
	bytesWritten int
}

func newLoggingResponseWriter(w http.ResponseWriter) *loggingResponseWriter {
	return &loggingResponseWriter{
		ResponseWriter: w,
		statusCode:     http.StatusOK,
	}
}

func (w *loggingResponseWriter) WriteHeader(statusCode int) {
	w.statusCode = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}

func (w *loggingResponseWriter) Write(b []byte) (int, error) {
	n, err := w.ResponseWriter.Write(b)
	w.bytesWritten += n
	return n, err
}

func statusTextCategory(status int) string {
	switch {
	case status >= 500:
		return "server_error"
	case status >= 400:
		return "client_error"
	case status >= 300:
		return "redirect"
	default:
		return "success"
	}
}

func LoggingMiddleware(logger log.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			loggingResponseWriter := newLoggingResponseWriter(w)

			next.ServeHTTP(loggingResponseWriter, r)

			duration := time.Since(start)

			msg := fmt.Sprintf(
				"%s %s %d %s %dB category: %s",
				r.Method,
				r.URL.RequestURI(),
				loggingResponseWriter.statusCode,
				duration,
				loggingResponseWriter.bytesWritten,
				statusTextCategory(loggingResponseWriter.statusCode),
			)

			if loggingResponseWriter.statusCode >= 500 {
				logger.Error(msg)
				return
			}

			logger.Info(msg)
		})
	}
}
