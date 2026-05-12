package handler

import (
	"encoding/json"
	"net/http"

	"github.com/Facundo-Mourelle/go-gym/internal/domain"
	"github.com/Facundo-Mourelle/go-gym/internal/service"
)

type WorkoutHandler struct {
	workoutService *service.WorkoutService
}

func NewWorkoutHandler(svc *service.WorkoutService) *WorkoutHandler {
	return &WorkoutHandler{
		workoutService: svc,
	}
}

// POST /api/v1/workouts
func (h *WorkoutHandler) CreateWorkout(w http.ResponseWriter, r *http.Request) {
	userID := GetUserIDFromContext(r.Context())

	var req struct {
		Name        string `json:"name"`
		Description string `json:"description"`
		Exercises   []struct {
			Order         int    `json:"order"`
			ExerciseID    string `json:"exercise_id"`
			Sets          int    `json:"sets"`
			Reps          int    `json:"reps"`
			RepsInReserve int    `json:"reps_in_reserve"`
			Notes         string `json:"notes"`
		} `json:"exercises"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate required fields
	if req.Name == "" {
		http.Error(w, "Workout name is required", http.StatusBadRequest)
		return
	}

	if len(req.Exercises) == 0 {
		http.Error(w, "At least one exercise is required", http.StatusBadRequest)
		return
	}

	// Convert to service request
	exercises := make([]service.WorkoutExerciseRequest, len(req.Exercises))
	for i, ex := range req.Exercises {
		exercises[i] = service.WorkoutExerciseRequest{
			Order:         ex.Order,
			ExerciseID:    domain.ExerciseID(ex.ExerciseID),
			Sets:          ex.Sets,
			Reps:          ex.Reps,
			RepsInReserve: ex.RepsInReserve,
			Notes:         ex.Notes,
		}
	}

	workout, err := h.workoutService.CreateWorkout(userID, service.CreateWorkoutRequest{
		Name:        req.Name,
		Description: req.Description,
		Exercises:   exercises,
	})

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(workout)
}

// GET /api/v1/workouts/{id}
func (h *WorkoutHandler) GetWorkout(w http.ResponseWriter, r *http.Request) {
	userID := GetUserIDFromContext(r.Context())
	workoutID := domain.WorkoutPlanID(r.PathValue("id"))

	if workoutID == "" {
		http.Error(w, "Workout ID is required", http.StatusBadRequest)
		return
	}

	workout, err := h.workoutService.GetWorkout(workoutID, userID)
	if err != nil {
		if err.Error() == "workout not found: record not found" {
			http.Error(w, "Workout not found", http.StatusNotFound)
			return
		}
		if err.Error() == "unauthorized: workout belongs to different user" {
			http.Error(w, "Unauthorized", http.StatusForbidden)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(workout)
}

// GET /api/v1/workouts
func (h *WorkoutHandler) ListWorkouts(w http.ResponseWriter, r *http.Request) {
	userID := GetUserIDFromContext(r.Context())

	workouts, err := h.workoutService.ListWorkouts(userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{
		"workouts": workouts,
		"count":    len(workouts),
	})
}

// PUT /api/v1/workouts/{id}
func (h *WorkoutHandler) UpdateWorkout(w http.ResponseWriter, r *http.Request) {
	userID := GetUserIDFromContext(r.Context())
	workoutID := domain.WorkoutPlanID(r.PathValue("id"))

	if workoutID == "" {
		http.Error(w, "Workout ID is required", http.StatusBadRequest)
		return
	}

	var req struct {
		Name        *string `json:"name"`
		Description *string `json:"description"`
		Exercises   []struct {
			Order         int    `json:"order"`
			ExerciseID    string `json:"exercise_id"`
			Sets          int    `json:"sets"`
			Reps          int    `json:"reps"`
			RepsInReserve int    `json:"reps_in_reserve"`
			Notes         string `json:"notes"`
		} `json:"exercises,omitempty"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Build update request
	updateReq := service.UpdateWorkoutRequest{
		Name:        req.Name,
		Description: req.Description,
	}

	// Convert exercises if provided
	if req.Exercises != nil {
		exercises := make([]service.WorkoutExerciseRequest, len(req.Exercises))
		for i, ex := range req.Exercises {
			exercises[i] = service.WorkoutExerciseRequest{
				Order:         ex.Order,
				ExerciseID:    domain.ExerciseID(ex.ExerciseID),
				Sets:          ex.Sets,
				Reps:          ex.Reps,
				RepsInReserve: ex.RepsInReserve,
				Notes:         ex.Notes,
			}
		}
		updateReq.Exercises = exercises
	}

	workout, err := h.workoutService.UpdateWorkout(workoutID, userID, updateReq)
	if err != nil {
		if err.Error() == "workout not found: record not found" {
			http.Error(w, "Workout not found", http.StatusNotFound)
			return
		}
		if err.Error() == "unauthorized: workout belongs to different user" {
			http.Error(w, "Unauthorized", http.StatusForbidden)
			return
		}
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(workout)
}

// DELETE /api/v1/workouts/{id}
func (h *WorkoutHandler) DeleteWorkout(w http.ResponseWriter, r *http.Request) {
	userID := GetUserIDFromContext(r.Context())
	workoutID := domain.WorkoutPlanID(r.PathValue("id"))

	if workoutID == "" {
		http.Error(w, "Workout ID is required", http.StatusBadRequest)
		return
	}

	err := h.workoutService.DeleteWorkout(workoutID, userID)
	if err != nil {
		if err.Error() == "workout not found: record not found" {
			http.Error(w, "Workout not found", http.StatusNotFound)
			return
		}
		if err.Error() == "unauthorized: workout belongs to different user" {
			http.Error(w, "Unauthorized", http.StatusForbidden)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
