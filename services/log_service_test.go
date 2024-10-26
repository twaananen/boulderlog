package services

import (
	"testing"
	"time"
	"github.com/stretchr/testify/assert"
	"github.com/twaananen/boulderlog/models"
)

func TestSaveLog(t *testing.T) {
	mockDB := new(MockDatabase)
	logService := NewLogService(mockDB)

	log := &models.BoulderLog{
			Username:   "testuser",
			Grade:      "V5",
			Difficulty: 3,
			Flash:      true,
			NewRoute:   false,
	}

	// Test case 1: Successful log save
	mockDB.On("SaveBoulderLog", log).Return(nil).Once()

	err := logService.SaveLog(log)
	assert.NoError(t, err)

	// Test case 2: Database error
	mockDB.On("SaveBoulderLog", log).Return(assert.AnError).Once()

	err = logService.SaveLog(log)
	assert.Error(t, err)

	mockDB.AssertExpectations(t)
}

func TestGetTodayGradeCounts(t *testing.T) {
	mockDB := new(MockDatabase)
	logService := NewLogService(mockDB)

	username := "testuser"
	gradeCounts := map[string]int{"V1": 2, "V2": 1}
	toppedCounts := map[string]int{"V1": 1, "V2": 1}

	// Test case 1: Successful retrieval
	mockDB.On("GetTodayGradeCounts", username).Return(gradeCounts, toppedCounts, nil).Once()

	gc, tc, err := logService.GetTodayGradeCounts(username)
	assert.NoError(t, err)
	assert.Equal(t, gradeCounts, gc)
	assert.Equal(t, toppedCounts, tc)

	// Test case 2: Database error
	mockDB.On("GetTodayGradeCounts", username).Return(make(map[string]int), make(map[string]int), assert.AnError).Once()

	gc, tc, err = logService.GetTodayGradeCounts(username)
	assert.Error(t, err)
	assert.Empty(t, gc) // Use assert.Empty to check for empty map
	assert.Empty(t, tc) // Use assert.Empty to check for empty map

	mockDB.AssertExpectations(t)
}

func TestGetGradeCounts(t *testing.T) {
	mockDB := new(MockDatabase)
	logService := NewLogService(mockDB)

	username := "testuser"
	logs := []models.BoulderLog{
		{Username: username, Grade: "V1"},
		{Username: username, Grade: "V2"},
		{Username: username, Grade: "V1"},
		{Username: username, Grade: "V3"},
	}

	// Mock the GetBoulderLogs method to return the logs
	mockDB.On("GetBoulderLogs", username).Return(logs, nil).Once()

	// Test case 1: Successful retrieval
	grades := []string{"V1", "V2", "V3"}
	counts := []int{2, 1, 1} // V1: 2, V2: 1, V3: 1

	g, c, err := logService.GetGradeCounts(username, nil, nil)
	assert.NoError(t, err)
	assert.Equal(t, grades, g)
	assert.Equal(t, counts, c)

	// Test case 2: Database error
	mockDB.On("GetBoulderLogs", username).Return([]models.BoulderLog{}, assert.AnError).Once()

	g, c, err = logService.GetGradeCounts(username, nil, nil)
	assert.Error(t, err)
	assert.Nil(t, g)
	assert.Nil(t, c)

	mockDB.AssertExpectations(t)
}

// func TestGetProgressData(t *testing.T) {
// 	mockDB := new(MockDatabase)
// 	logService := NewLogService(mockDB)

// 	username := "testuser"
// 	labels := []string{"2023-01-01", "2023-01-08", "2023-01-15"}
// 	gradeDatasets := map[string][]int{
// 		"V1": {3, 2, 1},
// 		"V2": {1, 2, 3},
// 	}

// 	logs := []models.BoulderLog{
// 		{Username: username, Grade: "V1"},
// 		{Username: username, Grade: "V1"},
// 		{Username: username, Grade: "V1"},
// 		{Username: username, Grade: "V2"},
// 		{Username: username, Grade: "V2"},
// 		{Username: username, Grade: "V2"},
// 	}

// 	// set these at midday UTC
// 	logs[0].CreatedAt, _ = time.Parse("2006-01-02 15:04:05", "2023-01-01 12:00:00")
// 	logs[1].CreatedAt, _ = time.Parse("2006-01-02 15:04:05", "2023-01-08 12:00:00")
// 	logs[2].CreatedAt, _ = time.Parse("2006-01-02 15:04:05", "2023-01-15 12:00:00")
// 	logs[3].CreatedAt, _ = time.Parse("2006-01-02 15:04:05", "2023-01-01 12:00:00")
// 	logs[4].CreatedAt, _ = time.Parse("2006-01-02 15:04:05", "2023-01-08 12:00:00")
// 	logs[5].CreatedAt, _ = time.Parse("2006-01-02 15:04:05", "2023-01-15 12:00:00")

// 	// Test case 1: Successful retrieval
// 	mockDB.On("GetBoulderLogs", username).Return(logs, nil).Once()

// 	l, gd, err := logService.GetProgressData(username)
// 	assert.NoError(t, err)
// 	assert.Equal(t, labels, l)
// 	assert.Equal(t, gradeDatasets, gd)

// 	// Test case 2: Database error
// 	mockDB.On("GetBoulderLogs", username).Return(nil, assert.AnError).Once()

// 	l, gd, err = logService.GetProgressData(username)
// 	assert.Error(t, err)
// 	assert.Nil(t, l)
// 	assert.Nil(t, gd)

// 	mockDB.AssertExpectations(t)
// }

func TestGetDifficultyProgressionData(t *testing.T) {
	// Create test logs spanning multiple periods
	now := time.Now()
	logs := []models.BoulderLog{
		{
			Grade:      "6A",
			Difficulty: 3,
			CreatedAt:  now,
		},
		{
			Grade:      "6A",
			Difficulty: 4,
			CreatedAt:  now,
		},
		{
			Grade:      "6B",
			Difficulty: 2,
			CreatedAt:  now,
		},
		{
			Grade:      "6A",
			Difficulty: 2,
			CreatedAt:  now.AddDate(0, 0, -7), // Previous week
		},
		{
			Grade:      "6B",
			Difficulty: 3,
			CreatedAt:  now.AddDate(0, 0, -7),
		},
	}

	service := &LogService{}

	t.Run("Weekly averages", func(t *testing.T) {
		data, labels, err := service.GetDifficultyProgressionData(logs, "week")
		assert.NoError(t, err)
		assert.Len(t, labels, 2) // Should have 2 weeks

		// Check 6A grade averages
		assert.Len(t, data["6A"], 2)
		assert.InDelta(t, 2.0, data["6A"][0].Value, 0.01) // First week: 2
		assert.InDelta(t, 3.5, data["6A"][1].Value, 0.01) // Second week: (3+4)/2

		// Check 6B grade averages
		assert.Len(t, data["6B"], 2)
		assert.InDelta(t, 3.0, data["6B"][0].Value, 0.01) // First week: 3
		assert.InDelta(t, 2.0, data["6B"][1].Value, 0.01) // Second week: 2
	})

	t.Run("Daily averages", func(t *testing.T) {
		data, labels, err := service.GetDifficultyProgressionData(logs, "day")
		assert.NoError(t, err)
		assert.Len(t, labels, 2) // Should have 2 days

		// Verify averages for each day
		assert.Len(t, data["6A"], 2)
		assert.Len(t, data["6B"], 2)
	})

	t.Run("Empty logs", func(t *testing.T) {
		data, labels, err := service.GetDifficultyProgressionData([]models.BoulderLog{}, "week")
		assert.NoError(t, err)
		assert.Nil(t, data)
		assert.Nil(t, labels)
	})
}
