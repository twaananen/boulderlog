package db

import (
	"fmt"
	"time"

	"github.com/twaananen/boulderlog/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type PostgresDatabase struct {
	db *gorm.DB
}

func NewPostgresDatabase(host, user, password, dbname string, port int) (*PostgresDatabase, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable TimeZone=UTC", host, user, password, dbname, port)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// Auto Migrate the schema
	err = db.AutoMigrate(&models.User{}, &models.BoulderLog{})
	if err != nil {
		return nil, err
	}

	return &PostgresDatabase{db: db}, nil
}

func (pdb *PostgresDatabase) GetUserByUsername(username string) (*models.User, error) {
	var user models.User
	result := pdb.db.Where("username = ?", username).First(&user)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, result.Error
	}
	return &user, nil
}

func (pdb *PostgresDatabase) CreateUser(user *models.User) error {
	return pdb.db.Create(user).Error
}

func (pdb *PostgresDatabase) SaveBoulderLog(log *models.BoulderLog) error {
	return pdb.db.Create(log).Error
}

func (pdb *PostgresDatabase) GetTodayGradeCounts(username string) (map[string]int, map[string]int, error) {
	var logs []models.BoulderLog
	today := time.Now().UTC().Truncate(24 * time.Hour)
	result := pdb.db.Where("username = ? AND date >= ?", username, today).Find(&logs)
	if result.Error != nil {
		return nil, nil, result.Error
	}

	gradeCounts := make(map[string]int)
	toppedCounts := make(map[string]int)

	for _, log := range logs {
		gradeCounts[log.Grade]++
		if log.Difficulty <= 4 {
			toppedCounts[log.Grade]++
		}
	}

	return gradeCounts, toppedCounts, nil
}

func (pdb *PostgresDatabase) GetBoulderLogs(username string) ([]models.BoulderLog, error) {
	var logs []models.BoulderLog
	result := pdb.db.Where("username = ?", username).Order("date").Find(&logs)
	return logs, result.Error
}

func (pdb *PostgresDatabase) GetGradeCounts(username string) ([]string, []int, error) {
	var results []struct {
		Grade string
		Count int
	}

	err := pdb.db.Model(&models.BoulderLog{}).
		Select("grade, count(*) as count").
		Where("username = ?", username).
		Group("grade").
		Order("grade").
		Find(&results).Error

	if err != nil {
		return nil, nil, err
	}

	grades := make([]string, len(results))
	counts := make([]int, len(results))

	for i, result := range results {
		grades[i] = result.Grade
		counts[i] = result.Count
	}

	return grades, counts, nil
}

func (pdb *PostgresDatabase) GetProgressData(username string) ([]string, map[string][]int, error) {
	var logs []models.BoulderLog
	err := pdb.db.Where("username = ?", username).Order("date").Find(&logs).Error
	if err != nil {
		return nil, nil, err
	}

	weeklyData := make(map[time.Time]map[string]int)
	var weeks []time.Time

	for _, log := range logs {
		weekStart := log.Date.Truncate(7 * 24 * time.Hour)
		if _, exists := weeklyData[weekStart]; !exists {
			weeklyData[weekStart] = make(map[string]int)
			weeks = append(weeks, weekStart)
		}
		weeklyData[weekStart][log.Grade]++
	}

	labels := make([]string, len(weeks))
	for i, week := range weeks {
		labels[i] = week.Format("2006-01-02")
	}

	gradeDatasets := make(map[string][]int)
	for _, week := range weeks {
		for grade := range weeklyData[week] {
			if _, exists := gradeDatasets[grade]; !exists {
				gradeDatasets[grade] = make([]int, len(weeks))
			}
		}
	}

	for i, week := range weeks {
		for grade := range gradeDatasets {
			gradeDatasets[grade][i] = weeklyData[week][grade]
		}
	}

	return labels, gradeDatasets, nil
}
