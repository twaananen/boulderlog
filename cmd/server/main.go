package main

import (
	"log"
	"net/http"

	"github.com/joho/godotenv"
	"github.com/twaananen/boulderlog/handlers"
	"github.com/twaananen/boulderlog/utils"
)

func init() {
	// Load .env file
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file")
	}
}

func main() {
	// Initialize JWT secret
	if err := handlers.InitJWTSecret(); err != nil {
		log.Fatalf("Failed to initialize JWT secret: %v", err)
	}

	utils.InitLogger()

	router := http.NewServeMux()

	// Static files
	router.Handle("GET /static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	// Routes
	router.HandleFunc("GET /", handlers.Home)
	router.HandleFunc("GET /login", handlers.LoginPage)
	router.HandleFunc("GET /profile", handlers.ProfilePage)
	router.HandleFunc("GET /auth/status", handlers.AuthStatus)
	router.HandleFunc("POST /auth/login", handlers.Login)
	router.HandleFunc("GET /auth/logout", handlers.Logout)

	// New logging routes
	router.HandleFunc("GET /log/grade", handlers.NewLogHandler().GetGradeSelection)
	router.HandleFunc("POST /log/difficulty/", handlers.NewLogHandler().GetPerceivedDifficulty)
	router.HandleFunc("POST /log/submit/", handlers.NewLogHandler().SubmitLog)

	utils.LogInfo("Server starting on http://localhost:8080")

	// Start server
	log.Fatal(http.ListenAndServe(":8080", router))
}
