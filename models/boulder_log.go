package models

import "time"

type BoulderLog struct {
	ID         int
	Username   string
	Grade      string
	Difficulty int
	Flash      bool
	NewRoute   bool
	Date       time.Time
}
