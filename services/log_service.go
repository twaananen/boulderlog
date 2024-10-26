package services

import (
	"fmt"
	"sort"
	"time"

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

func (s *LogService) GetTodayGradeCounts(userID string) (map[string]int, map[string]int, error) {
	counts, toppedCounts, err := s.db.GetTodayGradeCounts(userID)
	if err != nil {
		return nil, nil, err
	}
	if counts == nil {
		return map[string]int{}, map[string]int{}, nil
	}
	return counts, toppedCounts, nil
}

func (s *LogService) GetBoulderLogsBetweenDates(username string, startDate, endDate time.Time) ([]models.BoulderLog, error) {
	return s.db.GetBoulderLogsBetweenDates(username, startDate, endDate)
}

func (s *LogService) GetGradeCountsFromLogs(logs []models.BoulderLog) ([]string, map[string][]int, error) {
	gradeCounts := make(map[string]map[string]int)
	grades := make([]string, 0)
	
	for _, log := range logs {
		if _, exists := gradeCounts[log.Grade]; !exists {
			gradeCounts[log.Grade] = make(map[string]int)
				grades = append(grades, log.Grade)
		}
		if log.Difficulty >= 1 && log.Difficulty <= 4 {
			gradeCounts[log.Grade]["Topped"]++
		} else {
			gradeCounts[log.Grade]["Untopped"]++
		}
		if log.Flash {
			gradeCounts[log.Grade]["Flashed"]++
		}
		if log.NewRoute {
			gradeCounts[log.Grade]["New"]++
		}
	}

	sort.Strings(grades)

	datasets := map[string][]int{
		"Topped":   make([]int, len(grades)),
		"Untopped": make([]int, len(grades)),
		"Flashed":  make([]int, len(grades)),
		"New":      make([]int, len(grades)),
	}

	for i, grade := range grades {
		datasets["Topped"][i] = gradeCounts[grade]["Topped"]
		datasets["Untopped"][i] = gradeCounts[grade]["Untopped"]
		datasets["Flashed"][i] = gradeCounts[grade]["Flashed"]
		datasets["New"][i] = gradeCounts[grade]["New"]
	}

	return grades, datasets, nil
}

func (s *LogService) GetGradeCounts(username string, startDate, endDate *time.Time) ([]string, map[string][]int, error) {
	var logs []models.BoulderLog
	var err error

	if startDate == nil || endDate == nil {
		logs, err = s.db.GetBoulderLogs(username)
	} else {
		logs, err = s.db.GetBoulderLogsBetweenDates(username, *startDate, *endDate)
	}

	if err != nil {
		return nil, nil, err
	}

	return s.GetGradeCountsFromLogs(logs)
}

func (s *LogService) GetProgressData(username string) ([]string, map[string][]int, error) {
	logs, err := s.db.GetBoulderLogs(username)
	if err != nil {
		return nil, nil, err
	}

	weeklyData := make(map[string]map[string]int)
	for _, log := range logs {
		weekLabel := log.CreatedAt.Truncate(7 * 24 * time.Hour).Format("2006-01-02")
		if _, exists := weeklyData[weekLabel]; !exists {
			weeklyData[weekLabel] = make(map[string]int)
		}
		weeklyData[weekLabel][log.Grade]++
	}

	var labels []string
	for label := range weeklyData {
		labels = append(labels, label)
	}
	sort.Strings(labels)

	gradeDatasets := make(map[string][]int)
	for _, label := range labels {
		for grade, count := range weeklyData[label] {
			gradeDatasets[grade] = append(gradeDatasets[grade], count)
		}
	}

	return labels, gradeDatasets, nil
}

func (s *LogService) GetWeekBounds(date time.Time) (time.Time, time.Time) {
	year, week := date.ISOWeek()
	startDate := time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC).AddDate(0, 0, (week-1)*7)
	// Adjust startDate to the correct Monday of the week
	for startDate.Weekday() != time.Monday {
		startDate = startDate.AddDate(0, 0, -1)
	}
	endDate := startDate.AddDate(0, 0, 6)
	return startDate, endDate
}

func (s *LogService) SaveLog(log *models.BoulderLog) (*models.BoulderLog, error) {
	return s.db.SaveBoulderLog(log)
}

func (s *LogService) GetBoulderLogByID(username string, logID uint) (*models.BoulderLog, error) {
	return s.db.GetBoulderLogByID(username, logID)
}

func (s *LogService) UpdateBoulderLog(log *models.BoulderLog) (*models.BoulderLog, error) {
	return s.db.UpdateBoulderLog(log)
}

func (s *LogService) DeleteBoulderLog(username string, logID uint) error {
	return s.db.DeleteBoulderLog(username, logID)
}

func (s *LogService) GetBoulderLogs(username string) ([]models.BoulderLog, error) {
	return s.db.GetBoulderLogs(username)
}

// GetDifficultyProgressionData processes boulder logs to generate difficulty progression data
type DifficultyDataPoint struct {
	Value float64
	Label string
}

func (s *LogService) GetDifficultyProgressionData(logs []models.BoulderLog, period string) (map[string][]DifficultyDataPoint, []string, error) {
	if len(logs) == 0 {
		return nil, nil, nil
	}

	if period == "" {
		period = "week"
	}

	// Sort logs by date to ensure consistent processing
	sort.Slice(logs, func(i, j int) bool {
		return logs[i].CreatedAt.Before(logs[j].CreatedAt)
	})

	// Track periods and their labels in order of appearance
	periodMap := make(map[string]string)        // periodKey -> label
	periodOrder := make([]string, 0)            // ordered periodKeys
	gradeMap := make(map[string]bool)           // track unique grades

	// First pass: collect all periods and grades
	for _, log := range logs {
		gradeMap[log.Grade] = true
		
		var periodKey, label string
		switch period {
		case "day":
			periodKey = log.CreatedAt.Format("2006-01-02")
			label = log.CreatedAt.Format("02.01")
		case "week":
			year, week := log.CreatedAt.ISOWeek()
			periodKey = fmt.Sprintf("%d-W%02d", year, week)
			label = fmt.Sprintf("W%d", week)
		case "month":
			periodKey = log.CreatedAt.Format("2006-01")
			label = log.CreatedAt.Format("Jan 06")
		case "year":
			periodKey = log.CreatedAt.Format("2006")
			label = periodKey
		}

		if _, exists := periodMap[periodKey]; !exists {
			periodMap[periodKey] = label
			periodOrder = append(periodOrder, periodKey)
		}
	}

	// Create ordered slices
	grades := make([]string, 0, len(gradeMap))
	for grade := range gradeMap {
		grades = append(grades, grade)
	}
	sort.Strings(grades)

	// Create data accumulator
	periodData := make(map[string]map[string]struct {
		Sum   float64
		Count int
	})

	// Initialize all period-grade combinations
	for _, periodKey := range periodOrder {
		periodData[periodKey] = make(map[string]struct {
			Sum   float64
			Count int
		})
	}

	// Second pass: accumulate data
	for _, log := range logs {
		var periodKey string
		switch period {
		case "day":
			periodKey = log.CreatedAt.Format("2006-01-02")
		case "week":
			year, week := log.CreatedAt.ISOWeek()
			periodKey = fmt.Sprintf("%d-W%02d", year, week)
		case "month":
			periodKey = log.CreatedAt.Format("2006-01")
		case "year":
			periodKey = log.CreatedAt.Format("2006")
		}

		data := periodData[periodKey][log.Grade]
		data.Sum += float64(log.Difficulty)
		data.Count++
		periodData[periodKey][log.Grade] = data
	}

	// Create final dataset with proper alignment
	progressionData := make(map[string][]DifficultyDataPoint)
	for _, grade := range grades {
		// Initialize array with same length as labels for proper alignment
		progressionData[grade] = make([]DifficultyDataPoint, len(periodOrder))
		
		// Fill in data points where they exist
		for i, periodKey := range periodOrder {
			label := periodMap[periodKey]
			if data, exists := periodData[periodKey][grade]; exists && data.Count > 0 {
				avgDifficulty := data.Sum / float64(data.Count)
				progressionData[grade][i] = DifficultyDataPoint{
					Value: avgDifficulty,
					Label: label,
				}
			} else {
				// Create a null data point to maintain alignment
				progressionData[grade][i] = DifficultyDataPoint{
					Value: 0,  // This will be filtered out in the chart
					Label: label,
				}
			}
		}
	}

	// Extract labels in order
	labels := make([]string, len(periodOrder))
	for i, periodKey := range periodOrder {
		labels[i] = periodMap[periodKey]
	}

	return progressionData, labels, nil
}

// ClimbingStats holds the summary statistics for climbing activities
type ClimbingStats struct {
	ClimbingDays int
	Flashed      int
	New          int
	Topped       int
	Untopped     int
}

// GetClimbingStats calculates summary statistics from boulder logs
func (s *LogService) GetClimbingStats(logs []models.BoulderLog) ClimbingStats {
	stats := ClimbingStats{}
	
	// Use a map to count unique climbing days
	climbingDays := make(map[string]bool)
	
	for _, log := range logs {
		// Count unique days
		dateStr := log.CreatedAt.Truncate(24 * time.Hour).Format("2006-01-02")
		climbingDays[dateStr] = true
		
		// Count different types of climbs
		if log.Flash {
			stats.Flashed++
		}
		if log.NewRoute {
			stats.New++
		}
		if log.Difficulty >= 1 && log.Difficulty <= 4 {
			stats.Topped++
		} else {
			stats.Untopped++
		}
	}
	
	stats.ClimbingDays = len(climbingDays)
	return stats
}
