package middleware

import (
	"encoding/json"
	"log"
	"net/http"
	"runtime/debug"
)

type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}

func Recovery() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if err := recover(); err != nil {
					// Log the panic with stack trace
					log.Printf("Panic recovered: %v\n%s", err, debug.Stack())

					// Return error response
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusInternalServerError)

					response := ErrorResponse{
						Error:   "internal_server_error",
						Message: "An internal server error occurred",
					}

					json.NewEncoder(w).Encode(response)
				}
			}()

			next.ServeHTTP(w, r)
		})
	}
}
