package memory

import (
	"sync"

	"github.com/Facundo-Mourelle/go-gym/internal/domain"
	"github.com/Facundo-Mourelle/go-gym/internal/repository"
)

type ExerciseMemoryRepository struct {
	mu        sync.RWMutex
	exercises map[domain.ExerciseID]domain.Exercise
}

func NewExerciseMemoryRepository() *ExerciseMemoryRepository {
	return &ExerciseMemoryRepository{
		exercises: make(map[domain.ExerciseID]domain.Exercise),
	}
}

func (r *ExerciseMemoryRepository) Create(exercise domain.Exercise) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.exercises[exercise.ID]; exists {
		return repository.ErrAlreadyExists
	}

	r.exercises[exercise.ID] = exercise
	return nil
}

func (r *ExerciseMemoryRepository) FindByID(id domain.ExerciseID) (domain.Exercise, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	exercise, exists := r.exercises[id]
	if !exists {
		return domain.Exercise{}, repository.ErrNotFound
	}

	return exercise, nil
}

func (r *ExerciseMemoryRepository) FindAll() ([]domain.Exercise, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	exercises := make([]domain.Exercise, 0, len(r.exercises))
	for _, exercise := range r.exercises {
		exercises = append(exercises, exercise)
	}

	return exercises, nil
}

func (r *ExerciseMemoryRepository) FindByUser(userID string) ([]domain.Exercise, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	exercises := make([]domain.Exercise, 0)
	for _, exercise := range r.exercises {
		if exercise.UserID == userID {
			exercises = append(exercises, exercise)
		}
	}

	return exercises, nil
}

func (r *ExerciseMemoryRepository) Update(exercise domain.Exercise) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.exercises[exercise.ID]; !exists {
		return repository.ErrNotFound
	}

	r.exercises[exercise.ID] = exercise
	return nil
}

func (r *ExerciseMemoryRepository) Delete(id domain.ExerciseID) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.exercises[id]; !exists {
		return repository.ErrNotFound
	}

	delete(r.exercises, id)
	return nil
}

func (r *ExerciseMemoryRepository) FindByPattern(pattern domain.MovementPattern) ([]domain.Exercise, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	exercises := make([]domain.Exercise, 0)
	for _, exercise := range r.exercises {
		for _, p := range exercise.PrimaryPatterns {
			if p.Pattern == pattern {
				exercises = append(exercises, exercise)
				break
			}
		}
	}

	return exercises, nil
}
