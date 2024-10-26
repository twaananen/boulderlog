package handlers

import (
	"net/http"
	"time"

	"github.com/twaananen/boulderlog/components"
	"github.com/twaananen/boulderlog/services"
	"github.com/twaananen/boulderlog/utils"
)

type StatsHandler struct {
	userService *services.UserService
	logService  *services.LogService
}

func NewStatsHandler(userService *services.UserService, logService *services.LogService) *StatsHandler {
	return &StatsHandler{userService: userService, logService: logService}
}

func (h *StatsHandler) StatsPage(w http.ResponseWriter, r *http.Request) {
	username, err := h.userService.GetUsernameFromToken(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Get all logs for the user
	logs, err := h.logService.GetBoulderLogs(username)
	if err != nil {
		utils.LogError("Failed to get boulder logs", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Generate grade counts data
	gradeLabels, datasets, err := h.logService.GetGradeCountsFromLogs(logs)
	if err != nil {
		utils.LogError("Failed to get grade counts", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Generate difficulty progression data
	difficultyLabels, difficultyData, err := h.logService.GetDifficultyProgressionData(logs, "week") // or whatever period you want
	if err != nil {
		utils.LogError("Failed to get difficulty progression data", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Get climbing statistics
	stats := h.logService.GetClimbingStats(logs)

	isHtmxRequest := r.Header.Get("HX-Request") == "true"
	content := components.Stats(gradeLabels, datasets, difficultyLabels, difficultyData, stats, "all", time.Now().Format("2006-01-02"))

	if isHtmxRequest {
		content.Render(r.Context(), w)
	} else {
		components.Layout("Stats", content).Render(r.Context(), w)
	}
}

func (h *StatsHandler) GradeCountsChart(w http.ResponseWriter, r *http.Request) {
	username, err := h.userService.GetUsernameFromToken(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	viewType := r.URL.Query().Get("view")
	dateStr := r.URL.Query().Get("date")

	var startDate, endDate *time.Time
	if viewType == "weekly" {
		if dateStr == "" {
			now := time.Now()
			dateStr = now.Format("2006-01-02")
		}
		date, _ := time.Parse("2006-01-02", dateStr)
		start, end := h.logService.GetWeekBounds(date)
		startDate, endDate = &start, &end
	}

	gradeLabels, datasets, err := h.logService.GetGradeCounts(username, startDate, endDate)
	if err != nil {
		utils.LogError("Failed to get grade counts", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	components.GradeCountsChart(gradeLabels, datasets, viewType, dateStr, true).Render(r.Context(), w)
}

// Add more stats-related methods here as needed

