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
	return s.db.SaveBoulderLog(log)
}

func (s *LogService) GetTodayGradeCounts(username string) (map[string]int, map[string]int, error) {
	return s.db.GetTodayGradeCounts(username)
}

func (s *LogService) GetGradeCounts(username string) ([]string, []int, error) {
	logs, err := s.db.GetBoulderLogs(username)
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

	utils.LogInfo(fmt.Sprintf("Grade counts: %v", gradeCounts))

	return grades, counts, nil
}

func (s *LogService) GetProgressData(username string) ([]string, map[string][]int, error) {
	logs, err := s.db.GetBoulderLogs(username)
	if err != nil {
		return nil, nil, err
	}

	weeklyData := make(map[string]map[string]int)
	for _, log := range logs {
		weekLabel := log.Date.Truncate(7 * 24 * time.Hour).Format("2006-01-02")
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
