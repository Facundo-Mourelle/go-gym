package domain

import (
	"errors"
	"time"
)

type ExerciseID string

type Exercise struct {
	ID          ExerciseID
	Name        string
	Description string
	UserID      string

	// Muscles are derived from patterns automatically
	PrimaryPatterns    []MovementPatternContribution
	SecondaryPatterns  []MovementPatternContribution
	EquipmentUsed      Equipment
	SuggestedEquipment []EquipmentType

	CreatedAt time.Time
	UpdatedAt time.Time
}

func (e *Exercise) Validate() error {
	if e.Name == "" {
		return errors.New("exercise name is required")
	}
	return nil
}

type MovementPatternContribution struct {
	Pattern MovementPattern

	// 0.0 - 1.0, how much this pattern is emphasized
	Contribution float64

	// "full", "partial", "lengthened", "shortened"
	RangeOfMotion string

	Notes string
}
