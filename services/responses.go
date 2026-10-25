package services

import (
	_ "embed"
	"encoding/json"
	"math/rand"

	"github.com/twaananen/boulderlog/models"
	"github.com/twaananen/boulderlog/utils"
)

//go:embed response-messages.json
var responseMessagesJSON []byte

var responseTiers = [2][3]string{
	{"victory", "mildly_impressed", "unimpressed"},
	{"pep_talk", "light_roast", "full_roast"},
}

var responseMessages = loadResponseMessages()

func loadResponseMessages() map[string][]string {
	var catalogue map[string]map[string][]string
	if err := json.Unmarshal(responseMessagesJSON, &catalogue); err != nil {
		panic(err)
	}
	messages := make(map[string][]string)
	for _, tiers := range responseTiers {
		for _, tier := range tiers {
			for _, lines := range catalogue[tier] {
				messages[tier] = append(messages[tier], lines...)
			}
			if len(messages[tier]) == 0 {
				panic("empty response tier: " + tier)
			}
		}
	}
	return messages
}

func responseTier(best, current, difficulty int) string {
	outcome, band := 0, 0
	if difficulty > 4 {
		outcome = 1
	}
	if best >= 0 && current >= 0 {
		gap := best - current
		if gap >= 4 {
			band = 2
		} else if gap >= 2 {
			band = 1
		}
	}
	return responseTiers[outcome][band]
}

// GetLogResponse uses history before the new log is saved so a new high gets a victory.
func (s *LogService) GetLogResponse(log *models.BoulderLog) string {
	best := -1
	logs, err := s.db.GetBoulderLogs(log.Username)
	if err != nil {
		utils.LogError("Failed to get response baseline", err)
	} else {
		// shortcut: scan user history; use a ranked database maximum if histories make this slow.
		for _, previous := range logs {
			if previous.Username != log.Username || previous.DeletedAt.Valid || previous.Difficulty < 1 || previous.Difficulty > 4 {
				continue
			}
			if rank := models.GradeIndex(previous.Grade); rank > best {
				best = rank
			}
		}
	}
	lines := responseMessages[responseTier(best, models.GradeIndex(log.Grade), log.Difficulty)]
	return lines[rand.Intn(len(lines))]
}
