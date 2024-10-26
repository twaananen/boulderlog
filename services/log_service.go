package services

import (
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
func (s *LogService) GetDifficultyProgressionData(logs []models.BoulderLog, period string) ([]string, map[string][]struct {
    Time  time.Time
    Value float64
}, error) {
	if len(logs) == 0 {
		return nil, nil, nil
	}

	// Default to weekly if period is not specified
	if period == "" {
		period = "week"
	}

	// Track unique grades
	gradeMap := make(map[string]bool)
	
	// Group logs by time period and grade
	periodData := make(map[string]map[string]struct {
		Sum   int
		Count int
	})

	for _, log := range logs {
		gradeMap[log.Grade] = true
		
		// Get the period key based on the specified period
		var periodKey time.Time
		switch period {
		case "day":
			periodKey = log.CreatedAt.Truncate(24 * time.Hour)
		case "week":
			year, week := log.CreatedAt.ISOWeek()
			// Get the Monday of the week
			periodKey = time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC).
				AddDate(0, 0, (week-1)*7)
			// Adjust to Monday if needed
			for periodKey.Weekday() != time.Monday {
				periodKey = periodKey.AddDate(0, 0, -1)
			}
		case "month":
			periodKey = time.Date(log.CreatedAt.Year(), log.CreatedAt.Month(), 1, 0, 0, 0, 0, time.UTC)
		case "year":
			periodKey = time.Date(log.CreatedAt.Year(), 1, 1, 0, 0, 0, 0, time.UTC)
		default:
			// Default to weekly if invalid period specified
			year, week := log.CreatedAt.ISOWeek()
			periodKey = time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC).
				AddDate(0, 0, (week-1)*7)
			for periodKey.Weekday() != time.Monday {
				periodKey = periodKey.AddDate(0, 0, -1)
			}
		}

		periodStr := periodKey.Format("2006-01-02")
		if _, exists := periodData[periodStr]; !exists {
			periodData[periodStr] = make(map[string]struct {
				Sum   int
				Count int
			})
		}

		data := periodData[periodStr][log.Grade]
		data.Sum += log.Difficulty
		data.Count++
		periodData[periodStr][log.Grade] = data
	}

	// Convert grades to sorted slice
	grades := make([]string, 0, len(gradeMap))
	for grade := range gradeMap {
		grades = append(grades, grade)
	}
	sort.Strings(grades)

	// Create progression data with period averages
	progressionData := make(map[string][]struct {
		Time  time.Time
		Value float64
	})

	for periodStr, gradeData := range periodData {
		date, _ := time.Parse("2006-01-02", periodStr)
		
		for grade, data := range gradeData {
			avgDifficulty := float64(data.Sum) / float64(data.Count)
			progressionData[grade] = append(progressionData[grade], struct {
				Time  time.Time
				Value float64
			}{
				Time:  date,
				Value: avgDifficulty,
			})
		}
	}

	// Sort data points by time for each grade
	for grade := range progressionData {
		sort.Slice(progressionData[grade], func(i, j int) bool {
			return progressionData[grade][i].Time.Before(progressionData[grade][j].Time)
		})
	}

	return grades, progressionData, nil
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
