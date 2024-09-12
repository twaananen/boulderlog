package main

import (
	"log"
	"net/http"

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

	database, err := db.NewCSVDatabase("data")
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	userService := services.NewUserService(database)
	logService := services.NewLogService(database)

	router := http.NewServeMux()

	router.Handle("GET /static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	handlers.SetupRoutes(router, userService, logService)

	utils.LogInfo("Server starting on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}
