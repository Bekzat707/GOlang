package middleware

import (
	"log"
	"net/http"
	"time"
)

func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		timestamp := time.Now().Format("2006-01-02T15:04:05")
		log.Printf("%s %s %s", timestamp, r.Method, r.URL.Path)

		next.ServeHTTP(w, r)
	})
}
