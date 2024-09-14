package handlers

import (
	"fmt"
	"net/http"

	"github.com/twaananen/boulderlog/components"
	"github.com/twaananen/boulderlog/services"
	"github.com/twaananen/boulderlog/utils"
)

type ProfileHandler struct {
	userService      *services.UserService
	migrationService *services.MigrationService
}

func NewProfileHandler(userService *services.UserService, migrationService *services.MigrationService) *ProfileHandler {
	return &ProfileHandler{
		userService:      userService,
		migrationService: migrationService,
	}
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

func (h *ProfileHandler) MigrateData(w http.ResponseWriter, r *http.Request) {
	username, err := h.userService.GetUsernameFromToken(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	migratedCount, err := h.migrationService.MigrateUserData(username)
	if err != nil {
		utils.LogError("Failed to migrate user data", err)
		http.Error(w, "Failed to migrate data", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(fmt.Sprintf("<p>Successfully migrated %d records.</p>", migratedCount)))
}

// Add more profile-related methods here as needed
