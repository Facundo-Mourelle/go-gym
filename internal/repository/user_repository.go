package repository

import (
	"github.com/Facundo-Mourelle/go-gym/internal/domain"
)

type UserRepository interface {
	Create(user domain.User) error
	FindByID(id domain.UserID) (*domain.User, error)
	FindByEmail(email string) (*domain.User, error)
	Update(user domain.User) error
}
