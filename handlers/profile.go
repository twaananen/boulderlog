package handlers

import (
	"net/http"

	"github.com/twaananen/boulderlog/components"
	"github.com/twaananen/boulderlog/services"
	"github.com/twaananen/boulderlog/utils"
)

type ProfileHandler struct {
	userService *services.UserService
}

func NewProfileHandler(userService *services.UserService) *ProfileHandler {
	return &ProfileHandler{userService: userService}
}

func (h *ProfileHandler) ProfilePage(w http.ResponseWriter, r *http.Request) {
	isHtmxRequest := r.Header.Get("HX-Request") == "true"
	username, err := h.userService.GetUsernameFromToken(r)
	if err != nil {
		utils.LogError("Failed to get username from token", err)
		return
	}
	if isHtmxRequest {
		components.Profile(username).Render(r.Context(), w)
	} else {
		components.Layout("Profile", components.Profile(username)).Render(r.Context(), w)
	}
}

// Add more profile-related methods here as needed
