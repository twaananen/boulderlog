package handlers

import (
	"net/http"

	"github.com/twaananen/boulderlog/services"
)

func SetupRoutes(router *http.ServeMux, userService *services.UserService, logService *services.LogService) {
	homeHandler := NewHomeHandler(userService, logService)
	authHandler := NewAuthHandler(userService)
	logHandler := NewLogHandler(logService)

	router.HandleFunc("/", homeHandler.Home)
	router.HandleFunc("/login", authHandler.LoginPage)
	router.HandleFunc("/profile", authHandler.ProfilePage)
	router.HandleFunc("/auth/status", authHandler.AuthStatus)
	router.HandleFunc("/auth/login", authHandler.Login)
	router.HandleFunc("/auth/logout", authHandler.Logout)

	router.HandleFunc("/log/grade", logHandler.GetGradeSelection)
	router.HandleFunc("/log/difficulty/", logHandler.GetPerceivedDifficulty)
	router.HandleFunc("/log/submit/", logHandler.SubmitLog)
}
