package handlers

import (
	"net/http"

	"github.com/twaananen/boulderlog/components"
)

func Home(w http.ResponseWriter, r *http.Request) {
	isLoggedIn := IsUserLoggedIn(r)
	components.Home(isLoggedIn).Render(r.Context(), w)
}
