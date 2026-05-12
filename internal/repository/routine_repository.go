package repository

import "github.com/Facundo-Mourelle/go-gym/internal/domain"

type RoutineRepository interface {
	Create(routine domain.Routine) error
	FindByID(id domain.RoutineID) (domain.Routine, error)
	FindByUser(userID string) ([]domain.Routine, error)
	Delete(id domain.RoutineID) error
}