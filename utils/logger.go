package utils

import (
	"log"
	"os"
	"path/filepath"
)

var Logger *log.Logger

func InitLogger() {
	dataDir := filepath.Join(".", "data")
	if err := os.MkdirAll(dataDir, 0755); err != nil {
		log.Fatal("Failed to create data directory:", err)
	}

	logFile, err := os.OpenFile(filepath.Join(dataDir, "log.txt"), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal("Failed to open log file:", err)
	}

	Logger = log.New(logFile, "", log.LstdFlags)
	LogInfo("Logger initialized")
}

func LogInfo(message string) {
	if Logger != nil {
		Logger.Println("INFO:", message)
	}
	log.Println("INFO:", message)
}

func LogError(message string, err error) {
	if Logger != nil {
		Logger.Printf("ERROR: %s: %v", message, err)
	}
	log.Printf("ERROR: %s: %v", message, err)
}
