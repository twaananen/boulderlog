package services

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/twaananen/boulderlog/models"
)

func TestMigrateUserData(t *testing.T) {
	mockPostgresDB := new(MockDatabase)
	mockCSVDB := new(MockCSVDatabase)
	migrationService := NewMigrationService(mockPostgresDB, mockCSVDB)

	username := "testuser"
	now := time.Now()

	csvLogs := []models.CsvBoulderLog{
		{Username: username, Date: now, Grade: "V1", Difficulty: 3, Flash: true, NewRoute: false},
		{Username: username, Date: now.Add(time.Hour), Grade: "V2", Difficulty: 4, Flash: false, NewRoute: true},
	}

	// Test case 1: Successful migration
	mockCSVDB.On("GetBoulderLogs", username).Return(csvLogs, nil).Once()
	mockPostgresDB.On("GetBoulderLogByUsernameAndDate", username, csvLogs[0].Date).Return(nil, nil).Once()
	mockPostgresDB.On("GetBoulderLogByUsernameAndDate", username, csvLogs[1].Date).Return(nil, nil).Once()
	mockPostgresDB.On("SaveBoulderLog", mock.AnythingOfType("*models.BoulderLog")).Return(nil).Twice()

	count, err := migrationService.MigrateUserData(username)
	assert.NoError(t, err)
	assert.Equal(t, 2, count)

	// Test case 2: Some logs already exist
	mockCSVDB.On("GetBoulderLogs", username).Return(csvLogs, nil).Once()
	mockPostgresDB.On("GetBoulderLogByUsernameAndDate", username, csvLogs[0].Date).Return(&models.BoulderLog{}, nil).Once()
	mockPostgresDB.On("GetBoulderLogByUsernameAndDate", username, csvLogs[1].Date).Return(nil, nil).Once()
	mockPostgresDB.On("SaveBoulderLog", mock.AnythingOfType("*models.BoulderLog")).Return(nil).Once()

	count, err = migrationService.MigrateUserData(username)
	assert.NoError(t, err)
	assert.Equal(t, 1, count)

	// Test case 3: Error getting CSV logs
	mockCSVDB.On("GetBoulderLogs", username).Return([]models.CsvBoulderLog{}, assert.AnError).Once()

	count, err = migrationService.MigrateUserData(username)
	assert.Error(t, err)
	assert.Equal(t, 0, count)

	// Test case 4: Error saving to Postgres
	mockCSVDB.On("GetBoulderLogs", username).Return(csvLogs, nil).Once()
	mockPostgresDB.On("GetBoulderLogByUsernameAndDate", username, csvLogs[0].Date).Return(nil, nil).Once()
	mockPostgresDB.On("GetBoulderLogByUsernameAndDate", username, csvLogs[1].Date).Return(nil, nil).Once()
	mockPostgresDB.On("SaveBoulderLog", mock.AnythingOfType("*models.BoulderLog")).Return(assert.AnError).Once()
	mockPostgresDB.On("SaveBoulderLog", mock.AnythingOfType("*models.BoulderLog")).Return(assert.AnError).Once()

	count, err = migrationService.MigrateUserData(username)
	assert.NoError(t, err)
	assert.Equal(t, 0, count)

	mockCSVDB.AssertExpectations(t)
	mockPostgresDB.AssertExpectations(t)
}
