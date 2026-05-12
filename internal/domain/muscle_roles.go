package domain


type MuscleRole string

const (
	MuscleRolePrimary   MuscleRole = "primary"
	MuscleRoleSecondary MuscleRole = "secondary"
)

func (e *Exercise) GetPrimaryMuscles() []MuscleGroup {
	muscleMap := make(map[MuscleGroup]bool)

	for _, patternContrib := range e.PrimaryPatterns {
		patternInfo, exists := MovementPatternRegistry[patternContrib.Pattern]
		if !exists {
			continue
		}

		// Add all primary muscles from this pattern
		for _, muscle := range patternInfo.PrimaryMuscles {
			muscleMap[muscle] = true
		}
	}

	muscles := make([]MuscleGroup, 0, len(muscleMap))
	for muscle := range muscleMap {
		muscles = append(muscles, muscle)
	}

	return muscles
}

// GetSecondaryMuscles derives secondary muscles from secondary patterns
func (e *Exercise) GetSecondaryMuscles() []MuscleGroup {
	muscleMap := make(map[MuscleGroup]bool)
	primaryMuscles := make(map[MuscleGroup]bool)

	// Track primary muscles to avoid duplication
	for _, muscle := range e.GetPrimaryMuscles() {
		primaryMuscles[muscle] = true
	}

	for _, patternContrib := range e.SecondaryPatterns {
		patternInfo, exists := MovementPatternRegistry[patternContrib.Pattern]
		if !exists {
			continue
		}

		for _, muscle := range patternInfo.PrimaryMuscles {
			// Only add if not already a primary muscle
			if !primaryMuscles[muscle] {
				muscleMap[muscle] = true
			}
		}
	}

	muscles := make([]MuscleGroup, 0, len(muscleMap))
	for muscle := range muscleMap {
		muscles = append(muscles, muscle)
	}

	return muscles
}

func (e *Exercise) GetAllMuscles() map[MuscleGroup]MuscleRole {
	result := make(map[MuscleGroup]MuscleRole)

	// Add primary muscles
	for _, muscle := range e.GetPrimaryMuscles() {
		result[muscle] = MuscleRolePrimary
	}

	// Add secondary muscles
	for _, muscle := range e.GetSecondaryMuscles() {
		result[muscle] = MuscleRoleSecondary
	}

	return result
}
