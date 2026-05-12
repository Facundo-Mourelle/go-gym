package repository

import (
	"github.com/Facundo-Mourelle/go-gym/internal/domain"
)

type WorkoutRepository interface {
	Create(workout domain.WorkoutPlan) error

	FindByID(id domain.WorkoutPlanID) (domain.WorkoutPlan, error)

	FindByUser(userID string) ([]domain.WorkoutPlan, error)

	Update(workout domain.WorkoutPlan) error

	Delete(id domain.WorkoutPlanID) error
}
