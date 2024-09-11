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

	// Static files
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	// Routes
	http.HandleFunc("/", handlers.Home)
	http.HandleFunc("/login", handlers.LoginPage)
	http.HandleFunc("/profile", handlers.ProfilePage)
	http.HandleFunc("/auth/status", handlers.AuthStatus)
	http.HandleFunc("/auth/login", handlers.Login)
	http.HandleFunc("/auth/logout", handlers.Logout)

	// New logging routes
	http.HandleFunc("/log/grade", handlers.NewLogHandler().GetGradeSelection)
	http.HandleFunc("/log/difficulty/", handlers.NewLogHandler().GetPerceivedDifficulty)
	http.HandleFunc("/log/submit/", handlers.NewLogHandler().SubmitLog)

	utils.LogInfo("Server starting on http://localhost:8080")

	// Start server
	log.Fatal(http.ListenAndServe(":8080", nil))
}
