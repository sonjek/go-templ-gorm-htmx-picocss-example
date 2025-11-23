package middleware

import (
	"net/http"
	"time"
)

func SlowdownMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Slow down request by 10 milliseconds for lazy loading demonstration
		time.Sleep(10 * time.Millisecond)
		next.ServeHTTP(w, r)
	})
}
