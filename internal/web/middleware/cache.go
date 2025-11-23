package middleware

import (
	"net/http"
	"strings"
)

// CacheStaticFiles wraps a file server handler to add cache headers for static assets (24 hours)
func CacheStaticFiles(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.URL.Path, "/static/") {
			w.Header().Set("Cache-Control", "public, max-age=86400")
		}

		next.ServeHTTP(w, r)
	})
}
