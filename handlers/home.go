package handlers

import (
	"net/http"
	"time"

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
	username, err := h.userService.GetUsernameFromToken(r)
	if err != nil {
		utils.LogError("Failed to get username from token", err)
		username = ""
	}

	h.HomeWithUserName(w, r, username)
}

func (h *HomeHandler) HomeWithUserName(w http.ResponseWriter, r *http.Request, username string) {
	isHtmxRequest := r.Header.Get("HX-Request") == "true"

	var err error
	if username == "" {
		if isHtmxRequest {
			err = components.Home(false, false, -1, nil, nil).Render(r.Context(), w)
		} else {
			err = components.Layout("Home", components.Home(false, false, -1, nil, nil)).Render(r.Context(), w)
		}
		if err != nil {
			utils.LogError("Failed to render home page", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	startOfDay := time.Now().Truncate(24 * time.Hour)
	endOfDay := startOfDay.Add(24 * time.Hour)
	grades, datasets, err := h.logService.GetGradeCounts(username, &startOfDay, &endOfDay)
	if err != nil {
		utils.LogError("Failed to get grade counts for chart", err)
		return
	}

	content := components.Home(true, false, -1, grades, datasets)
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
