package repository

import (
	"time"

	"github.com/Facundo-Mourelle/go-gym/internal/domain"
)

type SessionRepository interface {
	Create(session domain.Session) error

	FindByID(id domain.SessionID) (domain.Session, error)

	FindByUser(userID string) ([]domain.Session, error)

	FindByUserAndDateRange(userID string, startDate, endDate time.Time) ([]domain.Session, error)

	AddPerformedSet(sessionID domain.SessionID, set domain.PerformedSet) error

	Complete(sessionID domain.SessionID, completedAt time.Time) error

	Update(session domain.Session) error

	Delete(id domain.SessionID) error
}
