package handlers

import (
	"net/http"

	"github.com/twaananen/boulderlog/components"
	"github.com/twaananen/boulderlog/services"
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
	if isLoggedIn {
		username, _ := h.userService.GetUsernameFromToken(r)
		gradeCounts, toppedCounts, err := h.logService.GetTodayGradeCounts(username)
		if err != nil {
			http.Error(w, "Failed to get grade counts", http.StatusInternalServerError)
			return
		}
		components.Home(isLoggedIn, gradeCounts, toppedCounts, false).Render(r.Context(), w)
	} else {
		components.Home(isLoggedIn, nil, nil, false).Render(r.Context(), w)
	}
}
