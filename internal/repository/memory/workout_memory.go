package memory

import (
	"sync"

	"github.com/Facundo-Mourelle/go-gym/internal/domain"
	"github.com/Facundo-Mourelle/go-gym/internal/repository"
)

type WorkoutMemoryRepository struct {
	mu       sync.RWMutex
	workouts map[domain.WorkoutPlanID]domain.WorkoutPlan
}

func NewWorkoutMemoryRepository() *WorkoutMemoryRepository {
	return &WorkoutMemoryRepository{
		workouts: make(map[domain.WorkoutPlanID]domain.WorkoutPlan),
	}
}

func (r *WorkoutMemoryRepository) Create(workout domain.WorkoutPlan) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.workouts[workout.ID]; exists {
		return repository.ErrAlreadyExists
	}

	r.workouts[workout.ID] = workout
	return nil
}

func (r *WorkoutMemoryRepository) FindByID(id domain.WorkoutPlanID) (domain.WorkoutPlan, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	workout, exists := r.workouts[id]
	if !exists {
		return domain.WorkoutPlan{}, repository.ErrNotFound
	}

	return workout, nil
}

func (r *WorkoutMemoryRepository) FindByUser(userID string) ([]domain.WorkoutPlan, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	workouts := make([]domain.WorkoutPlan, 0)
	for _, workout := range r.workouts {
		if workout.UserID == userID {
			workouts = append(workouts, workout)
		}
	}

	return workouts, nil
}

func (r *WorkoutMemoryRepository) Update(workout domain.WorkoutPlan) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.workouts[workout.ID]; !exists {
		return repository.ErrNotFound
	}

	r.workouts[workout.ID] = workout
	return nil
}

func (r *WorkoutMemoryRepository) Delete(id domain.WorkoutPlanID) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.workouts[id]; !exists {
		return repository.ErrNotFound
	}

	delete(r.workouts, id)
	return nil
}
