package handlers

import (
	"net/http"
	"time"

	"github.com/twaananen/boulderlog/components"
	"github.com/twaananen/boulderlog/models"
	"github.com/twaananen/boulderlog/services"
	"github.com/twaananen/boulderlog/utils"
)

type AuthHandler struct {
	userService *services.UserService
	logService  *services.LogService
}

func NewAuthHandler(userService *services.UserService, logService *services.LogService) *AuthHandler {
	return &AuthHandler{userService: userService, logService: logService}
}

func (h *AuthHandler) LoginPage(w http.ResponseWriter, r *http.Request) {
	isHtmxRequest := r.Header.Get("HX-Request") == "true"
	if isHtmxRequest {
		components.Login("", "").Render(r.Context(), w)
	} else {
		components.Layout("Login", components.Login("", "")).Render(r.Context(), w)
	}
}

func (h *AuthHandler) AuthStatus(w http.ResponseWriter, r *http.Request) {
	isLoggedIn := h.userService.IsUserLoggedIn(r)

	if isLoggedIn {
		components.AuthStatusLoggedIn().Render(r.Context(), w)
	} else {
		components.AuthStatusLoggedOut().Render(r.Context(), w)
	}
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	password := r.FormValue("password")

	token, err := h.userService.AuthenticateUser(username, password)
	if err != nil {
		if err == models.ErrInvalidCredentials {
			components.Login("Invalid username or password", username).Render(r.Context(), w)
		} else {
			utils.LogError("Error during authentication", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			components.Login("Something went wrong, try again later", username).Render(r.Context(), w)
		}
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    token,
		Expires:  time.Now().Add(7 * 24 * time.Hour),
		HttpOnly: true,
		Path:     "/",
	})

	// Render the home content for logged-in users
	gradeCounts, toppedCounts, err := h.logService.GetTodayGradeCounts(username)
	if err != nil {
		utils.LogError("Failed to get grade counts", err)
	}

	w.Header().Set("HX-Trigger", "authStatusChanged")
	components.Home(true, gradeCounts, toppedCounts, false, -1).Render(r.Context(), w)
}

func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	// Clear the token cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    "",
		Path:     "/",
		Expires:  time.Now().Add(-1 * time.Hour),
		HttpOnly: true,
	})

	w.Header().Set("HX-Trigger", "authStatusChanged")
	// Render the home content for logged out users
	components.Home(false, nil, nil, false, -1).Render(r.Context(), w)
}
