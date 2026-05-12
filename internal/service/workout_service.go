package service

import (
	"fmt"
	"time"

	"github.com/Facundo-Mourelle/go-gym/internal/domain"
	"github.com/Facundo-Mourelle/go-gym/internal/repository"
)

type WorkoutService struct {
	workoutRepo  repository.WorkoutRepository
	exerciseRepo repository.ExerciseRepository
}

func NewWorkoutService(
	workoutRepo repository.WorkoutRepository,
	exerciseRepo repository.ExerciseRepository,
) *WorkoutService {
	return &WorkoutService{
		workoutRepo:  workoutRepo,
		exerciseRepo: exerciseRepo,
	}
}

func (s *WorkoutService) CreateWorkout(
	userID string,
	req CreateWorkoutRequest,
) (WorkoutResponse, error) {

	for _, ex := range req.Exercises {
		_, err := s.exerciseRepo.FindByID(ex.ExerciseID)
		if err != nil {
			return WorkoutResponse{}, fmt.Errorf("exercise %s not found: %w", ex.ExerciseID, err)
		}
	}

	workoutExercises := make([]domain.WorkoutExercise, len(req.Exercises))
	for i, exReq := range req.Exercises {
		workoutExercises[i] = domain.WorkoutExercise{
			ID:         domain.WorkoutExerciseID(generateID()),
			Order:      exReq.Order,
			ExerciseID: exReq.ExerciseID,
			Sets:       exReq.Sets,

			Notes: exReq.Notes,
		}
	}

	workout := domain.WorkoutPlan{
		ID:          domain.WorkoutPlanID(generateID()),
		UserID:      userID,
		Name:        req.Name,
		Description: req.Description,
		Exercises:   workoutExercises,
	}

	if err := s.workoutRepo.Create(workout); err != nil {
		return WorkoutResponse{}, fmt.Errorf("failed to create workout: %w", err)
	}

	return s.toWorkoutResponse(workout)
}

func (s *WorkoutService) GetWorkout(
	workoutID domain.WorkoutPlanID,
	userID string,
) (WorkoutResponse, error) {

	workout, err := s.workoutRepo.FindByID(workoutID)
	if err != nil {
		return WorkoutResponse{}, fmt.Errorf("workout not found: %w", err)
	}

	// Verify ownership
	if workout.UserID != userID {
		return WorkoutResponse{}, fmt.Errorf("unauthorized: workout belongs to different user")
	}

	return s.toWorkoutResponse(workout)
}

func (s *WorkoutService) ListWorkouts(userID string) ([]WorkoutSummaryResponse, error) {
	workouts, err := s.workoutRepo.FindByUser(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch workouts: %w", err)
	}

	responses := make([]WorkoutSummaryResponse, len(workouts))
	for i, workout := range workouts {
		responses[i] = WorkoutSummaryResponse{
			ID:            workout.ID,
			Name:          workout.Name,
			Description:   workout.Description,
			ExerciseCount: len(workout.Exercises),
		}
	}

	return responses, nil
}

// UpdateWorkout updates an existing workout plan
func (s *WorkoutService) UpdateWorkout(
	workoutID domain.WorkoutPlanID,
	userID string,
	req UpdateWorkoutRequest,
) (WorkoutResponse, error) {

	// Fetch existing workout
	workout, err := s.workoutRepo.FindByID(workoutID)
	if err != nil {
		return WorkoutResponse{}, fmt.Errorf("workout not found: %w", err)
	}

	// Verify ownership
	if workout.UserID != userID {
		return WorkoutResponse{}, fmt.Errorf("unauthorized: workout belongs to different user")
	}

	// Update fields if provided
	if req.Name != nil {
		workout.Name = *req.Name
	}

	if req.Description != nil {
		workout.Description = *req.Description
	}

	// Update exercises if provided
	if req.Exercises != nil {
		// Validate all exercises exist
		for _, ex := range req.Exercises {
			_, err := s.exerciseRepo.FindByID(ex.ExerciseID)
			if err != nil {
				return WorkoutResponse{}, fmt.Errorf("exercise %s not found: %w", ex.ExerciseID, err)
			}
		}

		// Create new workout exercises
		workoutExercises := make([]domain.WorkoutExercise, len(req.Exercises))
		for i, exReq := range req.Exercises {
			workoutExercises[i] = domain.WorkoutExercise{
				ID:            domain.WorkoutExerciseID(generateID()),
				Order:         exReq.Order,
				ExerciseID:    exReq.ExerciseID,
				Sets:          exReq.Sets,
				Reps:          exReq.Reps,
				RepsInReserve: exReq.RepsInReserve,
				Notes:         exReq.Notes,
			}
		}

		workout.Exercises = workoutExercises
	}

	if err := s.workoutRepo.Update(workout); err != nil {
		return WorkoutResponse{}, fmt.Errorf("failed to update workout: %w", err)
	}

	return s.toWorkoutResponse(workout)
}

// DeleteWorkout deletes a workout plan
func (s *WorkoutService) DeleteWorkout(
	workoutID domain.WorkoutPlanID,
	userID string,
) error {

	// Fetch existing workout to verify ownership
	workout, err := s.workoutRepo.FindByID(workoutID)
	if err != nil {
		return fmt.Errorf("workout not found: %w", err)
	}

	if workout.UserID != userID {
		return fmt.Errorf("unauthorized: workout belongs to different user")
	}

	if err := s.workoutRepo.Delete(workoutID); err != nil {
		return fmt.Errorf("failed to delete workout: %w", err)
	}

	return nil
}

// Helper to convert domain workout to response
func (s *WorkoutService) toWorkoutResponse(workout domain.WorkoutPlan) (WorkoutResponse, error) {
	// Fetch exercise details
	exerciseMap := make(map[domain.ExerciseID]domain.Exercise)
	for _, we := range workout.Exercises {
		if _, exists := exerciseMap[we.ExerciseID]; !exists {
			exercise, err := s.exerciseRepo.FindByID(we.ExerciseID)
			if err != nil {
				// Exercise might have been deleted, use ID as name
				exerciseMap[we.ExerciseID] = domain.Exercise{
					ID:   we.ExerciseID,
					Name: string(we.ExerciseID),
				}
			} else {
				exerciseMap[we.ExerciseID] = exercise
			}
		}
	}

	// Convert exercises to response format
	exerciseResponses := make([]WorkoutExerciseResponse, len(workout.Exercises))
	for i, we := range workout.Exercises {
		exercise := exerciseMap[we.ExerciseID]

		exerciseResponses[i] = WorkoutExerciseResponse{
			ID:            we.ID,
			Order:         we.Order,
			ExerciseID:    exercise.ID,
			Sets:          we.Sets,
			Reps:          we.Reps,
			RepsInReserve: we.RepsInReserve,
			Notes:         we.Notes,
		}
	}

	return WorkoutResponse{
		ID:          workout.ID,
		Name:        workout.Name,
		Description: workout.Description,
		Exercises:   exerciseResponses,
	}, nil
}

// Request/Response Types

type CreateWorkoutRequest struct {
	Name        string
	Description string
	Exercises   []WorkoutExerciseRequest
}

type UpdateWorkoutRequest struct {
	Name        *string
	Description *string
	Exercises   []WorkoutExerciseRequest
}

type WorkoutExerciseRequest struct {
	Order         int
	ExerciseID    domain.ExerciseID
	Sets          int
	Reps          int
	RepsInReserve int
	Notes         string
}

type WorkoutResponse struct {
	ID          domain.WorkoutPlanID
	Name        string
	Description string
	Exercises   []WorkoutExerciseResponse
}

type WorkoutExerciseResponse struct {
	ID            domain.WorkoutExerciseID
	Order         int
	ExerciseID    domain.ExerciseID
	Sets          int
	Reps          int
	RepsInReserve int
	Notes         string
}

type WorkoutSummaryResponse struct {
	ID            domain.WorkoutPlanID
	Name          string
	Description   string
	ExerciseCount int
	CreatedAt     time.Time
	UpdatedAt     time.Time
}
