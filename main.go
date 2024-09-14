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
	requiredEnvVars := []string{"DB_HOST", "DB_USER", "DB_PASSWORD", "DB_NAME", "DB_PORT"}
	envVars := make(map[string]string)

	for _, envVar := range requiredEnvVars {
		value := os.Getenv(envVar)
		if value == "" {
			log.Fatalf("Missing required environment variable: %s", envVar)
		}
		envVars[envVar] = value
	}

	dbPort, err := strconv.Atoi(envVars["DB_PORT"])
	if err != nil {
		log.Fatalf("Invalid DB_PORT: %v", err)
	}

	// Create a new PostgreSQL database instance
	database, err := db.NewPostgresDatabase(envVars["DB_HOST"], envVars["DB_USER"], envVars["DB_PASSWORD"], envVars["DB_NAME"], dbPort)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	userService := services.NewUserService(database)
	logService := services.NewLogService(database)
	migrationService, err := services.NewMigrationService(database, "data")
	if err != nil {
		log.Fatal(err)
	}

	router := http.NewServeMux()

	router.Handle("GET /static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	handlers.SetupRoutes(router, userService, logService, migrationService)

	utils.LogInfo("Server starting on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}
