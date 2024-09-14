package middleware

import (
	"net/http"

	"github.com/twaananen/boulderlog/services"
)

func AuthMiddleware(userService *services.UserService) func(http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			_, err := userService.GetUsernameFromToken(r)
			if err != nil {
				http.Error(w, "User not authenticated", http.StatusUnauthorized)
				return
			}
			next.ServeHTTP(w, r)
		}
	}
}
