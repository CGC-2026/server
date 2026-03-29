package middleware

import (
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

func (w *loggingResponseWriter) Unwrap() http.ResponseWriter {
	return w.ResponseWriter
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

// wrapLoggingResponseWriter preserves optional ResponseWriter capabilities
// only when the underlying writer implements them.
func wrapLoggingResponseWriter(w *loggingResponseWriter) http.ResponseWriter {
	flusher, hasFlusher := w.ResponseWriter.(http.Flusher)
	hijacker, hasHijacker := w.ResponseWriter.(http.Hijacker)
	pusher, hasPusher := w.ResponseWriter.(http.Pusher)

	switch {
	case hasFlusher && hasHijacker && hasPusher:
		return struct {
			*loggingResponseWriter
			http.Flusher
			http.Hijacker
			http.Pusher
		}{w, flusher, hijacker, pusher}
	case hasFlusher && hasHijacker:
		return struct {
			*loggingResponseWriter
			http.Flusher
			http.Hijacker
		}{w, flusher, hijacker}
	case hasFlusher && hasPusher:
		return struct {
			*loggingResponseWriter
			http.Flusher
			http.Pusher
		}{w, flusher, pusher}
	case hasHijacker && hasPusher:
		return struct {
			*loggingResponseWriter
			http.Hijacker
			http.Pusher
		}{w, hijacker, pusher}
	case hasFlusher:
		return struct {
			*loggingResponseWriter
			http.Flusher
		}{w, flusher}
	case hasHijacker:
		return struct {
			*loggingResponseWriter
			http.Hijacker
		}{w, hijacker}
	case hasPusher:
		return struct {
			*loggingResponseWriter
			http.Pusher
		}{w, pusher}
	default:
		return w
	}
}

func statusTextCategory(status int) string {
	switch {
	case status >= 500:
		return "server_error"
	case status >= 400:
		return "client_error"
	case status >= 300:
		return "redirect"
	case status >= 200:
		return "success"
	case status >= 100:
		return "informational"
	default:
		return "unknown"
	}
}

func LoggingMiddleware(logger log.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			loggingResponseWriter := newLoggingResponseWriter(w)
			wrappedWriter := wrapLoggingResponseWriter(loggingResponseWriter)

			next.ServeHTTP(wrappedWriter, r)

			duration := time.Since(start)

			if loggingResponseWriter.statusCode >= 400 {
				logger.Error(
					"%s %s %d %s %dB category: %s",
					r.Method,
					r.URL.RequestURI(),
					loggingResponseWriter.statusCode,
					duration,
					loggingResponseWriter.bytesWritten,
					statusTextCategory(loggingResponseWriter.statusCode),
				)
				return
			}

			logger.Info(
				"%s %s %d %s %dB category: %s",
				r.Method,
				r.URL.RequestURI(),
				loggingResponseWriter.statusCode,
				duration,
				loggingResponseWriter.bytesWritten,
				statusTextCategory(loggingResponseWriter.statusCode),
			)
		})
	}
}
