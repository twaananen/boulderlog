package handlers

import (
	"net/http"

	"github.com/twaananen/boulderlog/components"
	"github.com/twaananen/boulderlog/services"
	"github.com/twaananen/boulderlog/utils"
)

type StatsHandler struct {
	userService *services.UserService
	logService  *services.LogService
}

func NewStatsHandler(userService *services.UserService, logService *services.LogService) *StatsHandler {
	return &StatsHandler{userService: userService, logService: logService}
}

func (h *StatsHandler) StatsPage(w http.ResponseWriter, r *http.Request) {
	username, err := h.userService.GetUsernameFromToken(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Fetch grade counts
	gradeLabels, gradeCounts, err := h.logService.GetGradeCounts(username)
	if err != nil {
		utils.LogError("Failed to get grade counts", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	isHtmxRequest := r.Header.Get("HX-Request") == "true"
	content := components.Stats(gradeLabels, gradeCounts)

	if isHtmxRequest {
		content.Render(r.Context(), w)
	} else {
		components.Layout("Stats", content).Render(r.Context(), w)
	}
}

// Add more stats-related methods here as needed
