package db

import (
	"time"

	"github.com/twaananen/boulderlog/models"
)

type Database interface {
	GetUserByUsername(username string) (*models.User, error)
	CreateUser(user *models.User) error
	SaveBoulderLog(log *models.BoulderLog) error
	GetTodayGradeCounts(username string) (map[string]int, map[string]int, error)
	GetBoulderLogs(username string) ([]models.BoulderLog, error)
	GetGradeCounts(username string) ([]string, []int, error)
	GetProgressData(username string) ([]string, map[string][]int, error)
	GetBoulderLogByUsernameAndDate(username string, date time.Time) (*models.BoulderLog, error)
	GetBoulderLogsBetweenDates(username string, startDate, endDate time.Time) ([]models.BoulderLog, error)
}
