package middleware

import (
	"net/http"
	"net/url"

	"github.com/twaananen/boulderlog/services"
	"github.com/twaananen/boulderlog/utils"
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

			// Refresh the token
			err = userService.RefreshToken(r, w)
			if err != nil {
				// Log the error but don't fail the request
				// The user is still authenticated
				utils.LogError("Error refreshing token", err)
			}

			next.ServeHTTP(w, r)
		}
	}
}
