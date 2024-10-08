package handlers

import (
	"net/http"
	"slices"
	"strconv"
	"time"

	"github.com/twaananen/boulderlog/components"
	"github.com/twaananen/boulderlog/models"
	"github.com/twaananen/boulderlog/services"
	"github.com/twaananen/boulderlog/utils"
)

type LogHandler struct {
	logService  *services.LogService
	userService *services.UserService
}

func NewLogHandler(logService *services.LogService, userService *services.UserService) *LogHandler {
	return &LogHandler{logService: logService, userService: userService}
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
	grade := r.URL.Query().Get("grade")
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

func (h *LogHandler) GetConfirmation(w http.ResponseWriter, r *http.Request) {
	grade := r.URL.Query().Get("grade")
	difficulty, _ := strconv.Atoi(r.URL.Query().Get("difficulty"))
	isHtmxRequest := r.Header.Get("HX-Request") == "true"

	content := components.BoulderConfirmation(grade, difficulty)

	var err error
	if isHtmxRequest {
		err = content.Render(r.Context(), w)
	} else {
		err = components.Layout("Boulder Confirmation", content).Render(r.Context(), w)
	}

	if err != nil {
		utils.LogError("Failed to render boulder confirmation", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

func (h *LogHandler) SubmitLog(w http.ResponseWriter, r *http.Request) {
	username, err := h.userService.GetUsernameFromToken(r)
	if err != nil {
		http.Error(w, "User not authenticated", http.StatusUnauthorized)
		return
	}

	grade := r.URL.Query().Get("grade")
	difficultyStr := r.URL.Query().Get("difficulty")
	difficulty, err := strconv.Atoi(difficultyStr)
	if err != nil {
		http.Error(w, "Invalid difficulty", http.StatusBadRequest)
		return
	}

	flash := r.URL.Query().Get("flash") == "true"
	newRoute := r.URL.Query().Get("new_route") == "true"

	if difficulty > 4 {
		flash = false
		newRoute = false
	}

	if flash {
		newRoute = true
	}

	log := &models.BoulderLog{
		Username:   username,
		Grade:      grade,
		Difficulty: difficulty,
		Flash:      flash,
		NewRoute:   newRoute,
	}

	_, err = h.logService.SaveLog(log)
	if err != nil {
		utils.LogError("Failed to save log", err)
		http.Error(w, "Failed to save log", http.StatusInternalServerError)
		return
	}

	startOfDay := time.Now().Truncate(24 * time.Hour)
	endOfDay := startOfDay.Add(24 * time.Hour)
	grades, datasets, err := h.logService.GetGradeCounts(username, &startOfDay, &endOfDay)
	if err != nil {
		utils.LogError("Failed to get grade counts for chart", err)
		http.Error(w, "Failed to get grade counts for chart", http.StatusInternalServerError)
		return
	}

	isHtmxRequest := r.Header.Get("HX-Request") == "true"
	content := components.LogSummary(true, difficulty, grades, datasets)

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

func (h *LogHandler) GetLogHistory(w http.ResponseWriter, r *http.Request) {
	isHtmxRequest := r.Header.Get("HX-Request") == "true"
	username, err := h.userService.GetUsernameFromToken(r)
	if err != nil {
		http.Error(w, "User not authenticated", http.StatusUnauthorized)
		return
	}

	logs, err := h.logService.GetBoulderLogs(username)
	if err != nil {
		utils.LogError("Failed to get boulder logs", err)
		http.Error(w, "Failed to get boulder logs", http.StatusInternalServerError)
		return
	}

	slices.Reverse(logs)

	if isHtmxRequest {
		err = components.LogHistory(logs).Render(r.Context(), w)
	} else {
		err = components.Layout("Log History", components.LogHistory(logs)).Render(r.Context(), w)
	}
	if err != nil {
		utils.LogError("Failed to render log history", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

func (h *LogHandler) GetEditLogRow(w http.ResponseWriter, r *http.Request) {
	username, err := h.userService.GetUsernameFromToken(r)
	if err != nil {
		http.Error(w, "User not authenticated", http.StatusUnauthorized)
		return
	}

	logID, err := strconv.ParseUint(r.URL.Path[len("/log/edit/"):], 10, 64)
	if err != nil {
		http.Error(w, "Invalid log ID", http.StatusBadRequest)
		return
	}

	log, err := h.logService.GetBoulderLogByID(username, uint(logID))
	if err != nil {
		utils.LogError("Failed to get boulder log", err)
		http.Error(w, "Failed to get boulder log", http.StatusInternalServerError)
		return
	}

	err = components.EditLogRow(*log).Render(r.Context(), w)
	if err != nil {
		utils.LogError("Failed to render edit log row", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

func (h *LogHandler) UpdateLog(w http.ResponseWriter, r *http.Request) {
	username, err := h.userService.GetUsernameFromToken(r)
	if err != nil {
		http.Error(w, "User not authenticated", http.StatusUnauthorized)
		return
	}

	logID, err := strconv.ParseUint(r.URL.Path[len("/log/update/"):], 10, 64)
	if err != nil {
		http.Error(w, "Invalid log ID", http.StatusBadRequest)
		return
	}

	err = r.ParseForm()
	if err != nil {
		http.Error(w, "Failed to parse form data", http.StatusBadRequest)
		return
	}

	grade := r.PostFormValue("grade")
	difficulty, err := strconv.Atoi(r.PostFormValue("difficulty"))
	if err != nil {
		http.Error(w, "Invalid difficulty", http.StatusBadRequest)
		return
	}
	flash := r.PostFormValue("flash") == "on"
	newRoute := r.PostFormValue("new_route") == "on"

	if difficulty > 4 {
		flash = false
		newRoute = false
	}

	if flash {
		newRoute = true
	}

	log, err := h.logService.GetBoulderLogByID(username, uint(logID))
	if err != nil {
		utils.LogError("Failed to get boulder log", err)
		http.Error(w, "Failed to get boulder log", http.StatusInternalServerError)
		return
	}

	log.Grade = grade
	log.Difficulty = difficulty
	log.Flash = flash
	log.NewRoute = newRoute

	updatedLog, err := h.logService.UpdateBoulderLog(log)
	if err != nil {
		utils.LogError("Failed to update boulder log", err)
		http.Error(w, "Failed to update boulder log", http.StatusInternalServerError)
		return
	}

	err = components.LogRow(*updatedLog).Render(r.Context(), w)
	if err != nil {
		utils.LogError("Failed to render updated log row", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

func (h *LogHandler) CancelEdit(w http.ResponseWriter, r *http.Request) {
	username, err := h.userService.GetUsernameFromToken(r)
	if err != nil {
		http.Error(w, "User not authenticated", http.StatusUnauthorized)
		return
	}

	logID, err := strconv.ParseUint(r.URL.Path[len("/log/cancel-edit/"):], 10, 64)
	if err != nil {
		http.Error(w, "Invalid log ID", http.StatusBadRequest)
		return
	}

	log, err := h.logService.GetBoulderLogByID(username, uint(logID))
	if err != nil {
		utils.LogError("Failed to get boulder log", err)
		http.Error(w, "Failed to get boulder log", http.StatusInternalServerError)
		return
	}

	err = components.LogRow(*log).Render(r.Context(), w)
	if err != nil {
		utils.LogError("Failed to render log row", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

func (h *LogHandler) DeleteLog(w http.ResponseWriter, r *http.Request) {
	username, err := h.userService.GetUsernameFromToken(r)
	if err != nil {
		http.Error(w, "User not authenticated", http.StatusUnauthorized)
		return
	}

	logID, err := strconv.ParseUint(r.URL.Path[len("/log/delete/"):], 10, 64)
	if err != nil {
		http.Error(w, "Invalid log ID", http.StatusBadRequest)
		return
	}

	err = h.logService.DeleteBoulderLog(username, uint(logID))
	if err != nil {
		utils.LogError("Failed to delete boulder log", err)
		http.Error(w, "Failed to delete boulder log", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
