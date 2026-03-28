package middleware

import (
	"log/slog"
	"net/http"
	"time"
)

type wrappedWriter struct {
	http.ResponseWriter
	statusCode int
}

func (w *wrappedWriter) WriteHeader(sc int) {
	if w.statusCode == 0 {
		w.ResponseWriter.WriteHeader(sc)
	}
	w.statusCode = sc
}

func (w *wrappedWriter) Write(data []byte) (int, error) {
	if w.statusCode == 0 {
		w.statusCode = http.StatusOK
	}
	return w.ResponseWriter.Write(data)
}

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		wrapped := &wrappedWriter{
			ResponseWriter: w,
			statusCode:     0,
		}

		next.ServeHTTP(wrapped, r)

		statusCode := wrapped.statusCode
		if statusCode == 0 {
			statusCode = http.StatusOK
		}

		// #nosec G706 - input is sanitized via slog's internal escaping
		slog.Info("Request",
			"status", statusCode,
			"method", r.Method,
			"path", r.URL.Path,
			"duration", time.Since(start),
		)
	})
}
