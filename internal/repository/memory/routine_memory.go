package memory

import (
	"sync"

	"github.com/Facundo-Mourelle/go-gym/internal/domain"
	"github.com/Facundo-Mourelle/go-gym/internal/repository"
)

type RoutineMemoryRepository struct {
	mu       sync.RWMutex
	routines map[domain.RoutineID]domain.Routine
}

func NewRoutineMemoryRepository() *RoutineMemoryRepository {
	return &RoutineMemoryRepository{
		routines: make(map[domain.RoutineID]domain.Routine),
	}
}

func (r *RoutineMemoryRepository) Create(routine domain.Routine) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.routines[routine.ID]; exists {
		return repository.ErrAlreadyExists
	}

	r.routines[routine.ID] = routine
	return nil
}

func (r *RoutineMemoryRepository) FindByID(id domain.RoutineID) (domain.Routine, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	routine, exists := r.routines[id]
	if !exists {
		return domain.Routine{}, repository.ErrNotFound
	}

	return routine, nil
}

func (r *RoutineMemoryRepository) FindByUser(userID string) ([]domain.Routine, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	routines := make([]domain.Routine, 0)
	for _, routine := range r.routines {
		if routine.UserID == userID {
			routines = append(routines, routine)
		}
	}

	return routines, nil
}

func (r *RoutineMemoryRepository) Delete(id domain.RoutineID) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.routines[id]; !exists {
		return repository.ErrNotFound
	}

	delete(r.routines, id)
	return nil
}