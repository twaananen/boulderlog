package middleware

import (
	"net/http"
	"net/url"

	"github.com/twaananen/boulderlog/services"
)

func AuthMiddleware(userService *services.UserService) func(http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			_, err := userService.GetUsernameFromToken(r)
			if err != nil {
				currentURL := url.QueryEscape(r.URL.RequestURI())
				loginURL := "/login?redirect=" + currentURL
				http.Redirect(w, r, loginURL, http.StatusSeeOther)
				return
			}
			next.ServeHTTP(w, r)
		}
	}
}
