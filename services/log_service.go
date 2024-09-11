package services

import (
	"net/http"

	"github.com/twaananen/boulderlog/db"
	"github.com/twaananen/boulderlog/models"
)

type LogService struct {
	db          db.Database
	userService *UserService
}

func NewLogService(db db.Database) *LogService {
	return &LogService{
		db:          db,
		userService: NewUserService(db),
	}
}

func (s *LogService) SaveLog(log *models.BoulderLog) error {
	return s.db.SaveBoulderLog(log)
}

func (s *LogService) GetTodayGradeCounts(username string) (map[string]int, map[string]int, error) {
	return s.db.GetTodayGradeCounts(username)
}

func (s *LogService) GetUsernameFromToken(r *http.Request) (string, error) {
	return s.userService.GetUsernameFromToken(r)
}
