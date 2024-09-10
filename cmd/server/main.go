package main

import (
	"log"
	"net/http"

	"github.com/joho/godotenv"
	"github.com/twaananen/boulderlog/internal/handlers"
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

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	http.HandleFunc("/", handlers.Home)
	http.HandleFunc("/login", handlers.LoginPage)
	http.HandleFunc("/profile", handlers.ProfilePage) // Add this line
	http.HandleFunc("/auth/status", handlers.AuthStatus)
	http.HandleFunc("/auth/login", handlers.Login)
	http.HandleFunc("/auth/logout", handlers.Logout) // Change this to GET request

	log.Println("Server starting on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
