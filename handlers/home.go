package handlers

import (
	"net/http"

	"github.com/twaananen/boulderlog/components"
)

func Home(w http.ResponseWriter, r *http.Request) {
	components.Home().Render(r.Context(), w)
}
