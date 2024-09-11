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

	components.LogSummary(gradeCounts, toppedCounts, true, difficulty).Render(r.Context(), w)
}
