package models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username string `gorm:"uniqueIndex"`
	Password string
}

type BoulderLog struct {
	gorm.Model
	Username   string `gorm:"index"`
	Grade      string
	Difficulty int
	Flash      bool
	NewRoute   bool
}
