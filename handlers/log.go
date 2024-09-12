package handlers

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/twaananen/boulderlog/components"
	"github.com/twaananen/boulderlog/models"
	"github.com/twaananen/boulderlog/services"
	"github.com/twaananen/boulderlog/utils"
)

type LogHandler struct {
	logService *services.LogService
}

func NewLogHandler(logService *services.LogService) *LogHandler {
	return &LogHandler{logService: logService}
}

func (h *LogHandler) GetGradeSelection(w http.ResponseWriter, r *http.Request) {
	isHtmxRequest := r.Header.Get("HX-Request") == "true"
	content := components.BoulderGradeSelection()

	var err error
	if isHtmxRequest {
		err = content.Render(r.Context(), w)
	} else {
		err = components.Layout("Grade Selection", content).Render(r.Context(), w)
	}

	if err != nil {
		utils.LogError("Failed to render grade selection", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

func (h *LogHandler) GetPerceivedDifficulty(w http.ResponseWriter, r *http.Request) {
	grade := r.URL.Path[len("/log/difficulty/"):]
	isHtmxRequest := r.Header.Get("HX-Request") == "true"
	content := components.PerceivedDifficulty(grade)

	var err error
	if isHtmxRequest {
		err = content.Render(r.Context(), w)
	} else {
		err = components.Layout("Perceived Difficulty", content).Render(r.Context(), w)
	}

	if err != nil {
		utils.LogError("Failed to render perceived difficulty", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

func (h *LogHandler) SubmitLog(w http.ResponseWriter, r *http.Request) {

	parts := strings.Split(r.URL.Path, "/")
	if len(parts) != 5 {
		http.Error(w, "Invalid URL", http.StatusBadRequest)
		return
	}
	grade := parts[3]
	difficultyStr := parts[4]
	difficulty, err := strconv.Atoi(difficultyStr)
	if err != nil {
		http.Error(w, "Invalid difficulty", http.StatusBadRequest)
		return
	}

	username, err := h.logService.GetUsernameFromToken(r)
	if err != nil {
		http.Error(w, "User not authenticated", http.StatusUnauthorized)
		return
	}

	flash := r.FormValue("flash") == "on"
	newRoute := r.FormValue("new") == "on"

	log := &models.BoulderLog{
		Username:   username,
		Grade:      grade,
		Difficulty: difficulty,
		Flash:      flash,
		NewRoute:   newRoute,
	}

	if err := h.logService.SaveLog(log); err != nil {
		utils.LogError("Failed to save log", err)
		http.Error(w, "Failed to save log", http.StatusInternalServerError)
		return
	}

	gradeCounts, toppedCounts, err := h.logService.GetTodayGradeCounts(username)
	if err != nil {
		utils.LogError("Failed to get grade counts", err)
		http.Error(w, "Failed to get grade counts", http.StatusInternalServerError)
		return
	}

	isHtmxRequest := r.Header.Get("HX-Request") == "true"
	content := components.LogSummary(gradeCounts, toppedCounts, true, difficulty)

	if isHtmxRequest {
		err = content.Render(r.Context(), w)
	} else {
		err = components.Layout("Log Summary", content).Render(r.Context(), w)
	}

	if err != nil {
		utils.LogError("Failed to render log summary", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}
