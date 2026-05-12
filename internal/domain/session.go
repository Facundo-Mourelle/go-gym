package domain

import "time"

type SessionID string
type PerformedSetID string

type Session struct {
	ID            SessionID
	WorkoutPlanID WorkoutPlanID
	UserID        string
	StartedAt     time.Time
	CompletedAt   time.Time
	PerformedSets []PerformedSet
	Notes         string
}

type PerformedSet struct {
	ID                PerformedSetID
	WorkoutExerciseID WorkoutExerciseID
	ExerciseID        ExerciseID
	SetNumber         int
	Reps              int
	RepsInReserve     int
	RawLoad           float64 // Input value (e.g., stack weight)
	EquipmentID       EquipmentID
	EffectiveLoad     float64 // Calculated by resistance profile
	PerformedAt       time.Time
}

func (s *Session) SessionVolume() int {
	return len(s.PerformedSets)
}
