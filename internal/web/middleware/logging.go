package middleware

import (
	"log"
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
		log.Println(statusCode, r.Method, r.URL.Path, time.Since(start))
	})
}
