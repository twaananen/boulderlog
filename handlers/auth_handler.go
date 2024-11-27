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
	homeHandler *HomeHandler
}

func NewAuthHandler(userService *services.UserService, logService *services.LogService, homeHandler *HomeHandler) *AuthHandler {
	return &AuthHandler{
		userService: userService,
		logService:  logService,
		homeHandler: homeHandler,
	}
}

func (h *AuthHandler) LoginPage(w http.ResponseWriter, r *http.Request) {
	isHtmxRequest := r.Header.Get("HX-Request") == "true"
	redirectURL := r.URL.Query().Get("redirect")

	loginComponent := components.Login("", "", redirectURL)
	if isHtmxRequest {
		loginComponent.Render(r.Context(), w)
	} else {
		components.Layout("Login", loginComponent).Render(r.Context(), w)
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
	redirectURL := r.FormValue("redirect")

	token, err := h.userService.AuthenticateUser(username, password)
	if err != nil {
		if err == models.ErrInvalidCredentials {
			components.Login("Invalid username or password", username, redirectURL).Render(r.Context(), w)
		} else {
			utils.LogError("Error during authentication", err)
			components.Login("Something went wrong, try again later", username, redirectURL).Render(r.Context(), w)
		}
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    token,
		Expires:  time.Now().Add(14 * 24 * time.Hour),
		HttpOnly: true,
		Path:     "/",
	})

	w.Header().Set("HX-Trigger", "authStatusChanged")

	if redirectURL != "" {
		http.Redirect(w, r, redirectURL, http.StatusSeeOther)
	} else {
		h.homeHandler.HomeWithUserName(w, r, username)
	}
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

	// Call the HomeHandler to render the home page
	h.homeHandler.HomeWithUserName(w, r, "")
}
