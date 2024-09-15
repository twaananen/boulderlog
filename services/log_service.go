package services

import (
	"fmt"
	"sort"
	"time"

	"github.com/twaananen/boulderlog/db"
	"github.com/twaananen/boulderlog/models"
	"github.com/twaananen/boulderlog/utils"
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
	utils.LogInfo(fmt.Sprintf("Saving log: %+v", log))
	return s.db.SaveBoulderLog(log)
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
