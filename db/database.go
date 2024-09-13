package db

import "github.com/twaananen/boulderlog/models"

type Database interface {
	GetUserByUsername(username string) (*models.User, error)
	CreateUser(user *models.User) error
	SaveBoulderLog(log *models.BoulderLog) error
	GetTodayGradeCounts(username string) (map[string]int, map[string]int, error)
	GetBoulderLogs(username string) ([]models.BoulderLog, error) // Add this line
}
