package handlers

import (
	"net/http"

	"github.com/twaananen/boulderlog/services"
)

func SetupRoutes(router *http.ServeMux, userService *services.UserService, logService *services.LogService) {
	homeHandler := NewHomeHandler(userService, logService)
	authHandler := NewAuthHandler(userService, logService)
	logHandler := NewLogHandler(logService)

	router.HandleFunc("GET /", homeHandler.Home)
	router.HandleFunc("GET /login", authHandler.LoginPage)
	router.HandleFunc("GET /profile", authHandler.ProfilePage)
	router.HandleFunc("GET /auth/status", authHandler.AuthStatus)
	router.HandleFunc("POST /auth/login", authHandler.Login)
	router.HandleFunc("POST /auth/logout", authHandler.Logout)

	router.HandleFunc("GET /log/grade", logHandler.GetGradeSelection)
	router.HandleFunc("GET /log/difficulty/", logHandler.GetPerceivedDifficulty)
	router.HandleFunc("POST /log/submit/", logHandler.SubmitLog)
}
