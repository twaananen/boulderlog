package services

import (
	"fmt"

	"github.com/twaananen/boulderlog/db"
	"github.com/twaananen/boulderlog/models"
	"github.com/twaananen/boulderlog/utils"
)

type MigrationService struct {
	postgresDB *db.PostgresDatabase
	csvDB      *db.CSVDatabase
}

func NewMigrationService(postgresDB *db.PostgresDatabase, dataDir string) (*MigrationService, error) {
	csvDB, err := db.NewCSVDatabase(dataDir)
	if err != nil {
		return nil, err
	}
	return &MigrationService{
		postgresDB: postgresDB,
		csvDB:      csvDB,
	}, nil
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
			utils.LogInfo(fmt.Sprintf("Skipping log for user %s at date %s : grade %s, difficulty %i, flash %t, new route %t", username, csvLog.Date, csvLog.Grade, csvLog.Difficulty, csvLog.Flash, csvLog.NewRoute))
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

		if err := s.postgresDB.SaveBoulderLog(newLog); err != nil {
			utils.LogError(fmt.Sprintf("Failed to migrate log for user %s at date %s", username, csvLog.Date), err)
			continue
		}

		utils.LogInfo(fmt.Sprintf("Migrated log for user %s at date %s : grade %s, difficulty %i, flash %t, new route %t", username, newLog.CreatedAt, newLog.Grade, newLog.Difficulty, newLog.Flash, newLog.NewRoute))

		migratedCount++
	}

	return migratedCount, nil
}
