package handlers

import (
	"encoding/csv"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/twaananen/boulderlog/components"
	"github.com/twaananen/boulderlog/utils"
)

type LogHandler struct {
	// Add any dependencies here, like a database connection
}

func NewLogHandler() *LogHandler {
	return &LogHandler{}
}

func (h *LogHandler) GetGradeSelection(w http.ResponseWriter, r *http.Request) {
	components.BoulderGradeSelection().Render(r.Context(), w)
}

func (h *LogHandler) GetPerceivedDifficulty(w http.ResponseWriter, r *http.Request) {
	grade := r.URL.Path[len("/log/difficulty/"):]
	components.PerceivedDifficulty(grade).Render(r.Context(), w)
}

func (h *LogHandler) SubmitLog(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) != 5 {
		http.Error(w, "Invalid URL", http.StatusBadRequest)
		return
	}
	grade := parts[3]
	difficultyStr := parts[4]
	difficulty, _ := strconv.Atoi(difficultyStr)

	username, err := GetUsernameFromSession(r)
	if err != nil {
		http.Error(w, "User not authenticated", http.StatusUnauthorized)
		return
	}

	flash := r.FormValue("flash") == "on"
	newRoute := r.FormValue("new") == "on"

	// Save the log to CSV file
	if err := saveLogToCSV(username, grade, difficulty, flash, newRoute); err != nil {
		utils.LogError("Failed to save log", err)
		http.Error(w, "Failed to save log", http.StatusInternalServerError)
		return
	}
	utils.LogInfo(fmt.Sprintf("Log saved successfully for user: %s, grade: %s, difficulty: %d, flash: %v, new: %v", username, grade, difficulty, flash, newRoute))

	// Get today's grade counts and topped counts
	gradeCounts, toppedCounts, err := getTodayGradeCounts(username)
	if err != nil {
		utils.LogError("Failed to get grade counts", err)
		http.Error(w, "Failed to get grade counts", http.StatusInternalServerError)
		return
	}

	components.LogSummary(gradeCounts, toppedCounts, true, difficulty).Render(r.Context(), w)
}

func saveLogToCSV(username, grade string, difficulty int, flash, newRoute bool) error {
	filename := fmt.Sprintf("%s-log.csv", username)
	filepath := filepath.Join("data", filename)

	// Ensure the logs directory exists
	if err := os.MkdirAll("data", os.ModePerm); err != nil {
		utils.LogError("Failed to create data directory", err)
		return err
	}

	fileExists := false
	if _, err := os.Stat(filepath); err == nil {
		fileExists = true
	}

	file, err := os.OpenFile(filepath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		utils.LogError("Failed to open log file", err)
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	if !fileExists {
		// Write header row if the file is new
		header := []string{"Timestamp", "Grade", "Difficulty", "Flash", "NewRoute"}
		if err := writer.Write(header); err != nil {
			utils.LogError("Failed to write CSV header", err)
			return err
		}
	}

	record := []string{
		time.Now().Format(time.RFC3339),
		grade,
		strconv.Itoa(difficulty),
		strconv.FormatBool(flash),
		strconv.FormatBool(newRoute),
	}

	return writer.Write(record)
}

func getTodayGradeCounts(username string) (map[string]int, map[string]int, error) {
	filename := fmt.Sprintf("%s-log.csv", username)
	filepath := filepath.Join("data", filename)

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

	// Skip the header row
	for _, record := range records[1:] {
		if len(record) != 5 {
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
