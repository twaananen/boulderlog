package handlers

import (
	"net/http"

	"github.com/twaananen/boulderlog/internal/templates"
)

func Home(w http.ResponseWriter, r *http.Request) {
	templates.Home().Render(r.Context(), w)
}
