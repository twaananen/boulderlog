package models

import (
	"slices"
	"strings"
)

var BoulderGrades = [...]string{
	"4", "4+", "5", "5+", "6A", "6A+", "6B", "6B+", "6C", "6C+", "7A",
	"7A+", "7B", "7B+", "7C", "7C+", "8A", "8A+", "8B", "8B+", "8C",
}

func GradeIndex(grade string) int {
	return slices.Index(BoulderGrades[:], strings.ToUpper(strings.TrimSpace(grade)))
}
