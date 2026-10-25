package handlers

import (
	"encoding/json"
	"errors"
	"html"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"slices"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/twaananen/boulderlog/db"
	"github.com/twaananen/boulderlog/models"
	"github.com/twaananen/boulderlog/services"
	"github.com/twaananen/boulderlog/utils"
)

type logTestDatabase struct {
	db.Database
	logs       []models.BoulderLog
	historyErr error
	saveErr    error
	events     []string
}

func (d *logTestDatabase) GetBoulderLogs(username string) ([]models.BoulderLog, error) {
	d.events = append(d.events, "history:"+username)
	return d.logs, d.historyErr
}

func (d *logTestDatabase) SaveBoulderLog(log *models.BoulderLog) (*models.BoulderLog, error) {
	d.events = append(d.events, "save")
	if d.saveErr != nil {
		return nil, d.saveErr
	}
	d.logs = append(d.logs, *log)
	return log, nil
}

func (d *logTestDatabase) GetBoulderLogsBetweenDates(_ string, _, _ time.Time) ([]models.BoulderLog, error) {
	d.events = append(d.events, "chart")
	return d.logs, nil
}

func logRequest(t *testing.T, target, body string) *http.Request {
	t.Helper()
	r := httptest.NewRequest(http.MethodPost, target, strings.NewReader(body))
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{"username": "alice", "exp": time.Now().Add(time.Hour).Unix()})
	signed, err := token.SignedString(utils.JWTSecret)
	require.NoError(t, err)
	r.AddCookie(&http.Cookie{Name: "token", Value: signed})
	return r
}

func TestSubmitLogResponses(t *testing.T) {
	data, err := os.ReadFile("../services/response-messages.json")
	require.NoError(t, err)
	var catalogue map[string]map[string][]string
	require.NoError(t, json.Unmarshal(data, &catalogue))
	for _, htmx := range []bool{false, true} {
		for _, tc := range []struct {
			grade      string
			difficulty int
			tier       string
		}{
			{"7A", 1, "victory"}, {"6C", 4, "mildly_impressed"}, {"6B+", 2, "unimpressed"},
			{"7A", 8, "pep_talk"}, {"6C+", 7, "light_roast"}, {"6B+", 5, "full_roast"},
			{"6B", 8, "full_roast"}, {"7B", 4, "victory"}, {"7B", 8, "pep_talk"},
		} {
			t.Run(tc.tier+"/"+tc.grade+"/htmx="+strconv.FormatBool(htmx), func(t *testing.T) {
				database := &logTestDatabase{logs: []models.BoulderLog{{Username: "alice", Grade: "7A+", Difficulty: 1}}}
				handler := NewLogHandler(services.NewLogService(database), services.NewUserService(database))
				query := url.Values{"grade": {strings.ToLower(tc.grade)}, "difficulty": {strconv.Itoa(tc.difficulty)}, "flash": {"true"}}
				r := logRequest(t, "/log/submit/?"+query.Encode(), "")
				r.Header.Set("HX-Request", strconv.FormatBool(htmx))
				w := httptest.NewRecorder()
				handler.SubmitLog(w, r)
				require.Equal(t, http.StatusOK, w.Code)
				assert.Equal(t, []string{"history:alice", "save", "chart"}, database.events)
				saved := database.logs[len(database.logs)-1]
				assert.Equal(t, tc.grade, saved.Grade)
				assert.Equal(t, tc.difficulty <= 4, saved.Flash)
				assert.Equal(t, tc.difficulty <= 4, saved.NewRoute)
				found := false
				for _, lines := range catalogue[tc.tier] {
					for _, line := range lines {
						found = found || strings.Contains(html.UnescapeString(w.Body.String()), line)
					}
				}
				assert.True(t, found, "response must belong to %s", tc.tier)
				assert.Equal(t, !htmx, strings.Contains(strings.ToLower(w.Body.String()), "<!doctype html>"))
			})
		}
	}
}

func TestSubmitLogFailures(t *testing.T) {
	for _, tc := range []struct {
		name, grade, difficulty string
		historyErr, saveErr     error
		status                  int
		events                  []string
	}{
		{name: "unsupported grade", grade: "V5", difficulty: "1", status: 400},
		{name: "missing grade", difficulty: "1", status: 400},
		{name: "low difficulty", grade: "6B", difficulty: "0", status: 400},
		{name: "high difficulty", grade: "6B", difficulty: "9", status: 400},
		{name: "invalid difficulty", grade: "6B", difficulty: "wat", status: 400},
		{name: "save failure", grade: "6B", difficulty: "1", saveErr: errors.New("save failed"), status: 500, events: []string{"history:alice", "save"}},
		{name: "baseline failure still saves", grade: "6B", difficulty: "8", historyErr: errors.New("history failed"), status: 200, events: []string{"history:alice", "save", "chart"}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			database := &logTestDatabase{historyErr: tc.historyErr, saveErr: tc.saveErr}
			handler := NewLogHandler(services.NewLogService(database), services.NewUserService(database))
			r := logRequest(t, "/log/submit/?"+url.Values{"grade": {tc.grade}, "difficulty": {tc.difficulty}}.Encode(), "")
			w := httptest.NewRecorder()
			handler.SubmitLog(w, r)
			assert.Equal(t, tc.status, w.Code)
			assert.Equal(t, tc.events, database.events)
			if tc.status != 200 {
				assert.NotContains(t, w.Body.String(), "text-3xl")
			}
		})
	}
}

func TestUpdateLogRejectsInvalidInputBeforeDatabase(t *testing.T) {
	for _, form := range []url.Values{
		{"grade": {"V5"}, "difficulty": {"1"}},
		{"grade": {"6B"}, "difficulty": {"0"}},
		{"grade": {"6B"}, "difficulty": {"9"}},
	} {
		database := &logTestDatabase{}
		handler := NewLogHandler(services.NewLogService(database), services.NewUserService(database))
		w := httptest.NewRecorder()
		handler.UpdateLog(w, logRequest(t, "/log/update/1", form.Encode()))
		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Empty(t, database.events)
	}
}

func TestHomeDoesNotChooseResponse(t *testing.T) {
	database := &logTestDatabase{}
	handler := NewHomeHandler(services.NewUserService(database), services.NewLogService(database))
	w := httptest.NewRecorder()
	r := logRequest(t, "/", "")
	handler.Home(w, r)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.False(t, slices.Contains(database.events, "history:alice"))
	assert.NotContains(t, w.Body.String(), "bg-gradient-to-br")
}
