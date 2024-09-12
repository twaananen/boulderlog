package handlers

import (
	"net/http"

	"github.com/twaananen/boulderlog/components"
	"github.com/twaananen/boulderlog/services"
	"github.com/twaananen/boulderlog/utils"
)

type HomeHandler struct {
	userService *services.UserService
	logService  *services.LogService
}

func NewHomeHandler(userService *services.UserService, logService *services.LogService) *HomeHandler {
	return &HomeHandler{
		userService: userService,
		logService:  logService,
	}
}

func (h *HomeHandler) Home(w http.ResponseWriter, r *http.Request) {
	isLoggedIn := h.userService.IsUserLoggedIn(r)
	isHtmxRequest := r.Header.Get("HX-Request") == "true"

	var gradeCounts, toppedCounts map[string]int
	var err error

	if isLoggedIn {
		username, err := h.userService.GetUsernameFromToken(r)
		if err != nil {
			utils.LogError("Failed to get username from token", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		gradeCounts, toppedCounts, err = h.logService.GetTodayGradeCounts(username)
		if err != nil {
			utils.LogError("Failed to get grade counts", err)
			http.Error(w, "Failed to get grade counts", http.StatusInternalServerError)
			return
		}
	}

	content := components.Home(isLoggedIn, gradeCounts, toppedCounts, false, -1)
	if isHtmxRequest {
		err = content.Render(r.Context(), w)
	} else {
		err = components.Layout("Home", content).Render(r.Context(), w)
	}

	if err != nil {
		utils.LogError("Failed to render home page", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}
