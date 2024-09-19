package db

import "github.com/twaananen/boulderlog/models"

type CSV_Database interface {
	GetUserByUsername(username string) (*models.CsvUser, error)
	CreateUser(user *models.CsvUser) error
	SaveBoulderLog(log *models.CsvBoulderLog) error
	GetTodayGradeCounts(username string) (map[string]int, map[string]int, error)
	GetBoulderLogs(username string) ([]models.CsvBoulderLog, error)
}
