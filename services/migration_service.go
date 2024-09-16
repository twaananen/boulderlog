package services

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"strconv"

	"github.com/twaananen/boulderlog/db"
	"github.com/twaananen/boulderlog/models"
	"github.com/twaananen/boulderlog/utils"
)

type MigrationService struct {
	postgresDB db.Database
	csvDB      db.CSV_Database
}

func NewMigrationService(postgresDB db.Database, csvDB db.CSV_Database) *MigrationService {
	return &MigrationService{
		postgresDB: postgresDB,
		csvDB:      csvDB,
	}
}

func (s *MigrationService) MigrateUserData(username string) (int, error) {
	csvLogs, err := s.csvDB.GetBoulderLogs(username)
	if err != nil {
		return 0, err
	}

	migratedCount := 0
	for _, csvLog := range csvLogs {
		existingLog, err := s.postgresDB.GetBoulderLogByUsernameAndDate(username, csvLog.Date)
		if err != nil {
			utils.LogError(fmt.Sprintf("Error checking for existing log: %v", err), err)
			continue
		}
		if existingLog != nil {
			// Log already exists, skip
			utils.LogInfo(fmt.Sprintf("Skipping log for user %s at date %s : grade %s, difficulty %v, flash %t, new route %t", username, csvLog.Date, csvLog.Grade, csvLog.Difficulty, csvLog.Flash, csvLog.NewRoute))
			continue
		}

		newLog := &models.BoulderLog{
			Username:   username,
			Grade:      csvLog.Grade,
			Difficulty: csvLog.Difficulty,
			Flash:      csvLog.Flash,
			NewRoute:   csvLog.NewRoute,
		}
		newLog.CreatedAt = csvLog.Date

		if _, err := s.postgresDB.SaveBoulderLog(newLog); err != nil {
			utils.LogError(fmt.Sprintf("Failed to migrate log for user %s at date %s", username, csvLog.Date), err)
			continue
		}

		utils.LogInfo(fmt.Sprintf("Migrated log for user %s at date %s : grade %s, difficulty %v, flash %t, new route %t", username, newLog.CreatedAt, newLog.Grade, newLog.Difficulty, newLog.Flash, newLog.NewRoute))

		migratedCount++
	}

	return migratedCount, nil
}

func (s *MigrationService) GetBoulderLogCSV(username string) (string, error) {
	logs, err := s.postgresDB.GetBoulderLogs(username)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	writer := csv.NewWriter(&buf)

	// Write header
	if err := writer.Write([]string{"Date", "Grade", "Difficulty", "Flash", "New Route"}); err != nil {
		return "", err
	}

	// Write log entries
	for _, log := range logs {
		record := []string{
			log.CreatedAt.Format("2006-01-02 15:04:05"),
			log.Grade,
			strconv.Itoa(log.Difficulty),
			strconv.FormatBool(log.Flash),
			strconv.FormatBool(log.NewRoute),
		}
		if err := writer.Write(record); err != nil {
			return "", err
		}
	}

	writer.Flush()
	if err := writer.Error(); err != nil {
		return "", err
	}

	return buf.String(), nil
}
