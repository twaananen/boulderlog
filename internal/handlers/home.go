package handlers

import (
	"net/http"

	"github.com/twaananen/boulderlog/internal/templates"
)

func Home(w http.ResponseWriter, r *http.Request) {
	component := templates.Home()
	component.Render(r.Context(), w)
}
