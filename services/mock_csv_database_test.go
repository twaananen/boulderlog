package services

import (
	"github.com/stretchr/testify/mock"
	"github.com/twaananen/boulderlog/models"
)

// MockCSVDatabase is a mock implementation of the db.CSV_Database interface
type MockCSVDatabase struct {
	mock.Mock
}

func (m *MockCSVDatabase) GetUserByUsername(username string) (*models.CsvUser, error) {
	args := m.Called(username)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.CsvUser), args.Error(1)
}

func (m *MockCSVDatabase) CreateUser(user *models.CsvUser) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockCSVDatabase) SaveBoulderLog(log *models.CsvBoulderLog) error {
	args := m.Called(log)
	return args.Error(0)
}

func (m *MockCSVDatabase) GetTodayGradeCounts(username string) (map[string]int, map[string]int, error) {
	args := m.Called(username)
	return args.Get(0).(map[string]int), args.Get(1).(map[string]int), args.Error(2)
}

func (m *MockCSVDatabase) GetBoulderLogs(username string) ([]models.CsvBoulderLog, error) {
	args := m.Called(username)
	return args.Get(0).([]models.CsvBoulderLog), args.Error(1)
}
