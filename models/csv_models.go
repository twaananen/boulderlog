package models

import (
	"time"
)

type CsvUser struct {
	Username string
	Password string
}

type CsvBoulderLog struct {
	Username   string
	Date       time.Time
	Grade      string
	Difficulty int
	Flash      bool
	NewRoute   bool
}
