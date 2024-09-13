package db

import (
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/twaananen/boulderlog/models"
	"github.com/twaananen/boulderlog/utils" // Assuming this is where utils.LogError is defined
)

type CSVDatabase struct {
	dataDir string
}

func NewCSVDatabase(dataDir string) (*CSVDatabase, error) {
	if err := os.MkdirAll(dataDir, os.ModePerm); err != nil {
		return nil, err
	}
	return &CSVDatabase{dataDir: dataDir}, nil
}

func (db *CSVDatabase) GetUserByUsername(username string) (*models.User, error) {
	filename := filepath.Join(db.dataDir, "users.csv")
	file, err := os.Open(filename)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil // User not found
		}
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	for _, record := range records {
		if len(record) != 2 {
			continue
		}
		if record[0] == username {
			return &models.User{
				Username: record[0],
				Password: record[1],
			}, nil
		}
	}

	return nil, nil // User not found
}

func (db *CSVDatabase) CreateUser(user *models.User) error {
	filename := filepath.Join(db.dataDir, "users.csv")
	file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	record := []string{
		user.Username,
		user.Password, // Note: This should already be hashed in the UserService
	}

	return writer.Write(record)
}

func (db *CSVDatabase) SaveBoulderLog(log *models.BoulderLog) error {
	filename := fmt.Sprintf("%s-log.csv", log.Username)
	filepath := filepath.Join(db.dataDir, filename)

	file, err := os.OpenFile(filepath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	record := []string{
		time.Now().Format(time.RFC3339),
		log.Grade,
		strconv.Itoa(log.Difficulty),
		strconv.FormatBool(log.Flash),
		strconv.FormatBool(log.NewRoute),
	}

	return writer.Write(record)
}

func (db *CSVDatabase) GetTodayGradeCounts(username string) (map[string]int, map[string]int, error) {
	filename := fmt.Sprintf("%s-log.csv", username)
	filepath := filepath.Join(db.dataDir, filename)

	file, err := os.Open(filepath)
	if err != nil {
		if os.IsNotExist(err) {
			return make(map[string]int), make(map[string]int), nil
		}
		return nil, nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, nil, err
	}

	gradeCounts := make(map[string]int)
	toppedCounts := make(map[string]int)
	today := time.Now().Format("2006-01-02")

	for _, record := range records {
		if len(record) != 5 {
			utils.LogError("Invalid record format", fmt.Errorf("expected 5 fields, got %d", len(record)))
			continue
		}

		logDate, err := time.Parse(time.RFC3339, record[0])
		if err != nil {
			continue
		}

		if logDate.Format("2006-01-02") == today {
			grade := record[1]
			difficulty, _ := strconv.Atoi(record[2])
			gradeCounts[grade]++
			if difficulty <= 4 {
				toppedCounts[grade]++
			}
		}
	}

	return gradeCounts, toppedCounts, nil
}

func (db *CSVDatabase) GetBoulderLogs(username string) ([]models.BoulderLog, error) {
	filename := fmt.Sprintf("%s-log.csv", username)
	filepath := filepath.Join(db.dataDir, filename)

	file, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	var logs []models.BoulderLog
	for i, record := range records {
		if i == 0 && record[0] == "Timestamp" {
			// Skip the header row
			continue
		}
		if len(record) != 5 {
			utils.LogError("Invalid record format", fmt.Errorf("expected 5 fields, got %d", len(record)))
			continue
		}
		date, err := time.Parse(time.RFC3339, record[0])
		if err != nil {
			utils.LogError("Invalid date format", err)
			continue
		}
		difficulty, err := strconv.Atoi(record[2])
		if err != nil {
			utils.LogError("Invalid difficulty format", err)
			continue
		}
		flash, err := strconv.ParseBool(record[3])
		if err != nil {
			utils.LogError("Invalid flash format", err)
			continue
		}
		newRoute, err := strconv.ParseBool(record[4])
		if err != nil {
			utils.LogError("Invalid newRoute format", err)
			continue
		}
		log := models.BoulderLog{
			Username:   username,
			Date:       date,
			Grade:      record[1],
			Difficulty: difficulty,
			Flash:      flash,
			NewRoute:   newRoute,
		}
		logs = append(logs, log)
	}

	return logs, nil
}
