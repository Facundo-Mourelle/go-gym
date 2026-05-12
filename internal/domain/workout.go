package domain

type WorkoutPlanID string
type WorkoutExerciseID string

type WorkoutPlan struct {
	ID          WorkoutPlanID
	UserID      string
	Name        string
	Description string
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
