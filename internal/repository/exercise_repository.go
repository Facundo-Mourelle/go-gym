package repository

import (
	"github.com/Facundo-Mourelle/go-gym/internal/domain"
)

type ExerciseRepository interface {
	Create(exercise domain.Exercise) error

	FindByID(id domain.ExerciseID) (domain.Exercise, error)

	FindByUser(userID string) ([]domain.Exercise, error)

	FindByPattern(pattern domain.MovementPattern) ([]domain.Exercise, error)

	FindAll() ([]domain.Exercise, error)

	Update(exercise domain.Exercise) error

	Delete(id domain.ExerciseID) error
}
