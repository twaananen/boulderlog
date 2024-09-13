package handlers

import (
	"net/http"

	"github.com/twaananen/boulderlog/middleware"
	"github.com/twaananen/boulderlog/services"
)

func SetupRoutes(router *http.ServeMux, userService *services.UserService, logService *services.LogService) {
	homeHandler := NewHomeHandler(userService, logService)
	authHandler := NewAuthHandler(userService, logService)
	profileHandler := NewProfileHandler(userService)
	logHandler := NewLogHandler(logService, userService)
	statsHandler := NewStatsHandler(userService, logService)
	authMiddleware := middleware.AuthMiddleware(userService)

	router.HandleFunc("GET /", homeHandler.Home)
	router.HandleFunc("GET /login", authHandler.LoginPage)
	router.HandleFunc("GET /auth/status", authHandler.AuthStatus)
	router.HandleFunc("POST /auth/login", authHandler.Login)
	router.HandleFunc("POST /auth/logout", authHandler.Logout)

	router.HandleFunc("GET /profile", authMiddleware(profileHandler.ProfilePage))

	router.HandleFunc("GET /stats", authMiddleware(statsHandler.StatsPage))
	router.HandleFunc("GET /log/grade", authMiddleware(logHandler.GetGradeSelection))
	router.HandleFunc("GET /log/difficulty/", authMiddleware(logHandler.GetPerceivedDifficulty))
	router.HandleFunc("POST /log/submit/", authMiddleware(logHandler.SubmitLog))

}
