package services

import (
	"encoding/json"
	"strings"
	"testing"
	"time"
	"unicode"
	"unicode/utf8"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/twaananen/boulderlog/models"
	"gorm.io/gorm"
)

func TestResponseTier(t *testing.T) {
	for _, tc := range []struct {
		grade string
		band  int
	}{
		{"7B", 0}, {"7A+", 0}, {"7A", 0},
		{"6C+", 1}, {"6C", 1}, {"6B+", 2}, {"6B", 2}, {"4", 2},
	} {
		for difficulty := 1; difficulty <= 8; difficulty++ {
			outcome := 0
			if difficulty > 4 {
				outcome = 1
			}
			assert.Equal(t, responseTiers[outcome][tc.band], responseTier(models.GradeIndex("7A+"), models.GradeIndex(tc.grade), difficulty), "grade %s difficulty %d", tc.grade, difficulty)
		}
	}
	assert.Equal(t, "victory", responseTier(-1, 0, 1))
	assert.Equal(t, "pep_talk", responseTier(-1, 0, 8))
	assert.Equal(t, "pep_talk", responseTier(20, -1, 8))
	assert.Equal(t, 1, models.GradeIndex("6A")-models.GradeIndex("5+"))
	assert.Equal(t, 1, models.GradeIndex("7A")-models.GradeIndex("6C+"))
	assert.Equal(t, models.GradeIndex("6B+"), models.GradeIndex(" 6b+ "))
	assert.Equal(t, -1, models.GradeIndex("V5"))
}

func TestGetLogResponseBaseline(t *testing.T) {
	for _, tc := range []struct {
		name       string
		history    []models.BoulderLog
		err        error
		grade      string
		difficulty int
		tier       string
	}{
		{name: "first top", grade: "4", difficulty: 1, tier: "victory"},
		{name: "no tops", grade: "4", difficulty: 8, tier: "pep_talk"},
		{name: "lookup failure", err: assert.AnError, grade: "4", difficulty: 8, tier: "pep_talk"},
		{name: "lookup failure on top", err: assert.AnError, grade: "4", difficulty: 1, tier: "victory"},
		{name: "only failures", history: []models.BoulderLog{{Username: "alice", Grade: "8C", Difficulty: 5}}, grade: "4", difficulty: 8, tier: "pep_talk"},
		{name: "invalid history", history: []models.BoulderLog{{Username: "alice", Grade: "V12", Difficulty: 1}, {Username: "alice", Grade: "8C", Difficulty: 0}}, grade: "4", difficulty: 8, tier: "pep_talk"},
		{name: "other user", history: []models.BoulderLog{{Username: "bob", Grade: "8C", Difficulty: 1}}, grade: "4", difficulty: 8, tier: "pep_talk"},
		{name: "deleted top", history: []models.BoulderLog{{Username: "alice", Grade: "8C", Difficulty: 1, Model: gorm.Model{DeletedAt: gorm.DeletedAt{Time: time.Now(), Valid: true}}}}, grade: "4", difficulty: 8, tier: "pep_talk"},
		{name: "normalised top", history: []models.BoulderLog{{Username: "alice", Grade: " 7a+ ", Difficulty: 4}}, grade: "6B+", difficulty: 5, tier: "full_roast"},
		{name: "new high", history: []models.BoulderLog{{Username: "alice", Grade: "4", Difficulty: 1}}, grade: "8C", difficulty: 4, tier: "victory"},
		{name: "above best failure", history: []models.BoulderLog{{Username: "alice", Grade: "4", Difficulty: 1}}, grade: "8C", difficulty: 8, tier: "pep_talk"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			database := new(MockDatabase)
			database.On("GetBoulderLogs", "alice").Return(tc.history, tc.err).Once()
			message := NewLogService(database).GetLogResponse(&models.BoulderLog{Username: "alice", Grade: tc.grade, Difficulty: tc.difficulty})
			assert.Contains(t, responseMessages[tc.tier], message)
			database.AssertExpectations(t)
		})
	}
}

func TestResponseBaselineUsesCurrentHistory(t *testing.T) {
	database := new(MockDatabase)
	service := NewLogService(database)
	log := &models.BoulderLog{Username: "alice", Grade: "6B+", Difficulty: 5}
	database.On("GetBoulderLogs", "alice").Return([]models.BoulderLog{{Username: "alice", Grade: "7A+", Difficulty: 1}}, nil).Once()
	assert.Contains(t, responseMessages["full_roast"], service.GetLogResponse(log))
	// The previous high was edited or deleted; no cached best may survive it.
	database.On("GetBoulderLogs", "alice").Return([]models.BoulderLog{{Username: "alice", Grade: "6B+", Difficulty: 1}}, nil).Once()
	assert.Contains(t, responseMessages["pep_talk"], service.GetLogResponse(log))
	database.AssertExpectations(t)
}

func TestResponseCatalogue(t *testing.T) {
	var catalogue map[string]map[string][]string
	require.NoError(t, json.Unmarshal(responseMessagesJSON, &catalogue))
	require.Len(t, catalogue, 6)
	rich := map[string]bool{"mythology": true, "middle_earth": true, "fantasy_worlds": true, "gaming": true, "anime": true}
	seen := make(map[string]bool)
	for _, tiers := range responseTiers {
		for _, tier := range tiers {
			require.Len(t, catalogue[tier], 12)
			require.Len(t, responseMessages[tier], 170)
			for theme, lines := range catalogue[tier] {
				count := 10
				if rich[theme] {
					count = 20
				}
				require.Len(t, lines, count, "%s/%s", tier, theme)
				for _, line := range lines {
					require.NotEmpty(t, line)
					assert.Equal(t, strings.TrimSpace(line), line)
					assert.LessOrEqual(t, utf8.RuneCountInString(line), 180)
					key := strings.Map(func(r rune) rune {
						if unicode.IsLetter(r) || unicode.IsDigit(r) {
							return unicode.ToLower(r)
						}
						return -1
					}, line)
					require.False(t, seen[key], "duplicate: %s", line)
					seen[key] = true
				}
			}
		}
	}
	assert.Len(t, seen, 1020)
}
