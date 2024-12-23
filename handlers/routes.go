package handlers

import (
	"net/http"

	"github.com/twaananen/boulderlog/middleware"
	"github.com/twaananen/boulderlog/services"
)

func SetupRoutes(router *http.ServeMux, userService *services.UserService, logService *services.LogService, migrationService *services.MigrationService) {
	homeHandler := NewHomeHandler(userService, logService)
	authHandler := NewAuthHandler(userService, logService, homeHandler)
	profileHandler := NewProfileHandler(userService, migrationService)
	logHandler := NewLogHandler(logService, userService)
	statsHandler := NewStatsHandler(userService, logService)
	authMiddleware := middleware.AuthMiddleware(userService)

	router.HandleFunc("GET /", homeHandler.Home)
	router.HandleFunc("GET /login", authHandler.LoginPage)
	router.HandleFunc("GET /auth/status", authHandler.AuthStatus)
	router.HandleFunc("POST /auth/login", authHandler.Login)
	router.HandleFunc("POST /auth/logout", authHandler.Logout)

	router.HandleFunc("GET /profile", authMiddleware(profileHandler.ProfilePage))
	router.HandleFunc("POST /profile/migrate", authMiddleware(profileHandler.MigrateData))
	router.HandleFunc("GET /profile/download-log", authMiddleware(profileHandler.DownloadLog)) // Ensure this line is present

	router.HandleFunc("GET /stats", authMiddleware(statsHandler.StatsPage))
	router.HandleFunc("GET /charts/grade-counts", authMiddleware(statsHandler.GradeCountsChart))
	router.HandleFunc("GET /charts/difficulty-progression", statsHandler.DifficultyProgressionChart)

	router.HandleFunc("GET /log/grade", authMiddleware(logHandler.GetGradeSelection))
	router.HandleFunc("GET /log/difficulty/", authMiddleware(logHandler.GetPerceivedDifficulty))
	router.HandleFunc("GET /log/confirm/", authMiddleware(logHandler.GetConfirmation))
	router.HandleFunc("POST /log/submit/", authMiddleware(logHandler.SubmitLog))

	router.HandleFunc("GET /log/history", authMiddleware(logHandler.GetLogHistory))
	router.HandleFunc("GET /log/edit/", authMiddleware(logHandler.GetEditLogRow))
	router.HandleFunc("POST /log/update/", authMiddleware(logHandler.UpdateLog))
	router.HandleFunc("GET /log/cancel-edit/", authMiddleware(logHandler.CancelEdit))
	router.HandleFunc("DELETE /log/delete/", authMiddleware(logHandler.DeleteLog))

}
