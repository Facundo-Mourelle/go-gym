package types

type WorkoutPlanID string
type WorkoutExerciseID string

type WorkoutPlan struct {
	ID        WorkoutPlanID
	UserID    string
	Name      string
	Exercises []WorkoutExercise
}

type WorkoutExercise struct {
	ID         WorkoutExerciseID
	ExerciseID ExerciseID
	Sets       int
	Notes      string
}
