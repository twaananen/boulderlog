package services

import (
	"time"

	"github.com/stretchr/testify/mock"
	"github.com/twaananen/boulderlog/models"
)

// MockDatabase is a mock implementation of the db.Database interface
type MockDatabase struct {
	mock.Mock
}

func (m *MockDatabase) GetUserByUsername(username string) (*models.User, error) {
	args := m.Called(username)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockDatabase) CreateUser(user *models.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockDatabase) SaveBoulderLog(log *models.BoulderLog) error {
	args := m.Called(log)
	return args.Error(0)
}

func (m *MockDatabase) GetTodayGradeCounts(username string) (map[string]int, map[string]int, error) {
	args := m.Called(username)
	return args.Get(0).(map[string]int), args.Get(1).(map[string]int), args.Error(2)
}

func (m *MockDatabase) GetBoulderLogs(username string) ([]models.BoulderLog, error) {
	args := m.Called(username)
	return args.Get(0).([]models.BoulderLog), args.Error(1)
}

func (m *MockDatabase) GetBoulderLogByUsernameAndDate(username string, date time.Time) (*models.BoulderLog, error) {
	args := m.Called(username, date)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.BoulderLog), args.Error(1)
}

func (m *MockDatabase) GetGradeCounts(username string) ([]string, []int, error) {
	args := m.Called(username)
	return args.Get(0).([]string), args.Get(1).([]int), args.Error(2)
}

func (m *MockDatabase) GetProgressData(username string) ([]string, map[string][]int, error) {
	args := m.Called(username)
	return args.Get(0).([]string), args.Get(1).(map[string][]int), args.Error(2)
}

// Add any other methods from the db.Database interface that you need for testing
