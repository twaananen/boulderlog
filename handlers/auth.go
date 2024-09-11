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
}

func NewAuthHandler(userService *services.UserService) *AuthHandler {
	return &AuthHandler{userService: userService}
}

func (h *AuthHandler) LoginPage(w http.ResponseWriter, r *http.Request) {
	components.Login("", "").Render(r.Context(), w)
}

func (h *AuthHandler) ProfilePage(w http.ResponseWriter, r *http.Request) {
	components.Profile().Render(r.Context(), w)
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
			components.LoginForm("Invalid username or password", username).Render(r.Context(), w)
		} else {
			utils.LogError("Error during authentication", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    token,
		Expires:  time.Now().Add(24 * time.Hour),
		HttpOnly: true,
		Path:     "/",
	})

	w.Header().Set("HX-Redirect", "/")
}

func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    "",
		Path:     "/",
		Expires:  time.Unix(0, 0),
		MaxAge:   -1,
		HttpOnly: true,
	})
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
