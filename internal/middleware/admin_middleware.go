package middleware

import (
	"monoex_backend/internal/services"
	"net/http"
)

// AdminMiddleware проверяет BasicAuth и админские права
// Принимает обычную функцию хэндлера и возвращает http.HandlerFunc
func AdminMiddleware(service *services.AdminService) func(func(http.ResponseWriter, *http.Request)) http.HandlerFunc {
	return func(next func(http.ResponseWriter, *http.Request)) http.HandlerFunc {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			username, password, ok := r.BasicAuth()
			if !ok {
				w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			if !service.ValidateAdmin(username, password) {
				http.Error(w, "Forbidden", http.StatusForbidden)
				return
			}

			next(w, r)
		})
	}
}
