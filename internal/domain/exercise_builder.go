package domain

import (
	"fmt"
	"time"

)

type ExerciseBuilder struct {
	exercise Exercise
}

func NewExerciseBuilder() *ExerciseBuilder {
	return &ExerciseBuilder{
		exercise: Exercise{
			PrimaryPatterns:    make([]MovementPatternContribution, 0),
			SecondaryPatterns:  make([]MovementPatternContribution, 0),
			SuggestedEquipment: make([]EquipmentType, 0),
			CreatedAt:          time.Now(),
			UpdatedAt:          time.Now(),
		},
	}
}

func (b *ExerciseBuilder) WithName(name string) *ExerciseBuilder {
	b.exercise.Name = name
	return b
}

func (b *ExerciseBuilder) WithDescription(desc string) *ExerciseBuilder {
	b.exercise.Description = desc
	return b
}

func (b *ExerciseBuilder) WithUserID(userID string) *ExerciseBuilder {
	b.exercise.UserID = userID
	return b
}

func (b *ExerciseBuilder) WithPrimaryPattern(
	pattern MovementPattern,
	contribution float64,
	rom string,
	notes string,
) *ExerciseBuilder {
	b.exercise.PrimaryPatterns = append(b.exercise.PrimaryPatterns, MovementPatternContribution{
		Pattern:       pattern,
		Contribution:  contribution,
		RangeOfMotion: rom,
		Notes:         notes,
	})
	return b
}

func (b *ExerciseBuilder) WithSecondaryPattern(
	pattern MovementPattern,
	contribution float64,
	rom string,
	notes string,
) *ExerciseBuilder {
	b.exercise.SecondaryPatterns = append(b.exercise.SecondaryPatterns, MovementPatternContribution{
		Pattern:       pattern,
		Contribution:  contribution,
		RangeOfMotion: rom,
		Notes:         notes,
	})
	return b
}

func (b *ExerciseBuilder) WithEquipment(equipmentTypes ...EquipmentType) *ExerciseBuilder {
	b.exercise.SuggestedEquipment = append(b.exercise.SuggestedEquipment, equipmentTypes...)
	return b
}

func (b *ExerciseBuilder) Build() (Exercise, error) {
	b.exercise.ID = ExerciseID(generateID())

	if err := b.exercise.Validate(); err != nil {
		return Exercise{}, err
	}

	return b.exercise, nil
}

// Helper to generate exercise ID
func generateID() string {
	return fmt.Sprintf("ex_%d", time.Now().UnixNano())
}
