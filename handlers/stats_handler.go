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

	gradeLabels, datasets, err := h.logService.GetGradeCounts(username, nil, nil)
	if err != nil {
		utils.LogError("Failed to get grade counts", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	isHtmxRequest := r.Header.Get("HX-Request") == "true"
	content := components.Stats(gradeLabels, datasets, "all", time.Now().Format("2006-01-02"))

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
