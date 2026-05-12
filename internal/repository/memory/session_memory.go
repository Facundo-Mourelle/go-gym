package memory

import (
	"sync"
	"time"

	"github.com/Facundo-Mourelle/go-gym/internal/domain"
	"github.com/Facundo-Mourelle/go-gym/internal/repository"
)

type SessionMemoryRepository struct {
	mu       sync.RWMutex
	sessions map[domain.SessionID]domain.Session
}

func NewSessionMemoryRepository() *SessionMemoryRepository {
	return &SessionMemoryRepository{
		sessions: make(map[domain.SessionID]domain.Session),
	}
}

func (r *SessionMemoryRepository) Create(session domain.Session) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.sessions[session.ID]; exists {
		return repository.ErrAlreadyExists
	}

	r.sessions[session.ID] = session
	return nil
}

func (r *SessionMemoryRepository) FindByID(id domain.SessionID) (domain.Session, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	session, exists := r.sessions[id]
	if !exists {
		return domain.Session{}, repository.ErrNotFound
	}

	return session, nil
}

func (r *SessionMemoryRepository) FindByUser(userID string) ([]domain.Session, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	sessions := make([]domain.Session, 0)
	for _, session := range r.sessions {
		if session.UserID == userID {
			sessions = append(sessions, session)
		}
	}

	return sessions, nil
}

func (r *SessionMemoryRepository) FindByUserAndDateRange(
	userID string,
	startDate, endDate time.Time,
) ([]domain.Session, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	sessions := make([]domain.Session, 0)
	for _, session := range r.sessions {
		if session.UserID != userID {
			continue
		}

		// Check if session is within date range
		if session.StartedAt.Before(startDate) {
			continue
		}

		// Use completed time if available, otherwise use start time
		sessionEndTime := session.StartedAt
		if session.CompletedAt.IsZero() {
			sessionEndTime = session.CompletedAt
		}

		if sessionEndTime.After(endDate) {
			continue
		}

		sessions = append(sessions, session)
	}

	return sessions, nil
}

func (r *SessionMemoryRepository) AddPerformedSet(
	sessionID domain.SessionID,
	set domain.PerformedSet,
) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	session, exists := r.sessions[sessionID]
	if !exists {
		return repository.ErrNotFound
	}

	session.PerformedSets = append(session.PerformedSets, set)
	r.sessions[sessionID] = session

	return nil
}

func (r *SessionMemoryRepository) Complete(
	sessionID domain.SessionID,
	completedAt time.Time,
) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	session, exists := r.sessions[sessionID]
	if !exists {
		return repository.ErrNotFound
	}

	session.CompletedAt = completedAt
	r.sessions[sessionID] = session

	return nil
}

func (r *SessionMemoryRepository) Update(session domain.Session) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.sessions[session.ID]; !exists {
		return repository.ErrNotFound
	}

	r.sessions[session.ID] = session
	return nil
}

func (r *SessionMemoryRepository) Delete(id domain.SessionID) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.sessions[id]; !exists {
		return repository.ErrNotFound
	}

	delete(r.sessions, id)
	return nil
}
