package main

import (
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"github.com/twaananen/boulderlog/db"
	"github.com/twaananen/boulderlog/handlers"
	"github.com/twaananen/boulderlog/services"
	"github.com/twaananen/boulderlog/utils"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Println("Error loading .env file")
	}
}

func main() {
	if err := utils.InitJWTSecret(); err != nil {
		log.Fatalf("Failed to initialize JWT secret: %v", err)
	}

	utils.InitLogger()

	// Get database connection details from environment variables
	dbHost := os.Getenv("DB_HOST")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	dbPortStr := os.Getenv("DB_PORT")
	dbPort, err := strconv.Atoi(dbPortStr)
	if err != nil {
		log.Fatalf("Invalid DB_PORT: %v", err)
	}

	// Create a new PostgreSQL database instance
	database, err := db.NewPostgresDatabase(dbHost, dbUser, dbPassword, dbName, dbPort)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	userService := services.NewUserService(database)
	logService := services.NewLogService(database)

	router := http.NewServeMux()

	router.Handle("GET /static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	handlers.SetupRoutes(router, userService, logService)

	utils.LogInfo("Server starting on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}
