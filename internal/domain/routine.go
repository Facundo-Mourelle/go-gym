package domain

import "time"

type RoutineID string

type Routine struct {
	ID               RoutineID
	UserID           string
	Name             string
	Description      string
	MovementPatterns []MovementPattern
	IsPreset         bool
	CreatedAt        time.Time
	UpdatedAt        time.Time
}