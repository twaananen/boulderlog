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

func (s *LogService) SaveLog(log *models.BoulderLog) error {
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

func (s *LogService) GetGradeCounts(username string, startDate, endDate *time.Time) ([]string, []int, error) {
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

	gradeCounts := make(map[string]int)
	for _, log := range logs {
		gradeCounts[log.Grade]++
	}

	grades := make([]string, 0, len(gradeCounts))
	for grade := range gradeCounts {
		grades = append(grades, grade)
	}
	sort.Strings(grades)

	counts := make([]int, len(grades))
	for i, grade := range grades {
		counts[i] = gradeCounts[grade]
	}

	return grades, counts, nil
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
