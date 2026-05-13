package domain

import "time"

type WorkoutPlanID string
type WorkoutExerciseID string

type WorkoutPlan struct {
	ID          WorkoutPlanID
	UserID      string
	Name        string
	Description string
	CreatedAt   time.Time
	UpdatedAt   time.Time
	Exercises   []WorkoutExercise
}

type WorkoutExercise struct {
	ID         WorkoutExerciseID
	ExerciseID ExerciseID
	Order      int
	Sets       int
	Reps       int
	RepsInReserve int
	Notes      string
}
