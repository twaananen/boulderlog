package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username string `gorm:"uniqueIndex"`
	Password string
}

type BoulderLog struct {
	gorm.Model
	Username   string
	Date       time.Time
	Grade      string
	Difficulty int
	Flash      bool
	NewRoute   bool
}
