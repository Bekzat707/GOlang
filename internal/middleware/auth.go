package middleware

import (
	"encoding/json"
	"net/http"

	"github.com/bekzatsaparbekov/task-api/internal/models"
)

const (
	APIKeyHeader = "X-API-KEY"
	ValidAPIKey  = "secret12345"
)

func APIKeyAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		apiKey := r.Header.Get(APIKeyHeader)

		if apiKey == "" || apiKey != ValidAPIKey {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(models.ErrorResponse{
				Error: "unauthorized",
			})
			return
		}

		next.ServeHTTP(w, r)
	})
}
