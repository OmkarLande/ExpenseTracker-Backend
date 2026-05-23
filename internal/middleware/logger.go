package middleware

import (
	"log"
	"net/http"
	"time"
)

func LoggerMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			
			// Call the next handler
			next.ServeHTTP(w, r)
			
			// Log the request
			log.Printf("%s %s %v", r.Method, r.URL.Path, time.Since(start))
		})
	}
}
