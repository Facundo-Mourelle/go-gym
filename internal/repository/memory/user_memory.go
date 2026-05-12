package memory

import (
	"github.com/Facundo-Mourelle/go-gym/internal/domain"
	"github.com/Facundo-Mourelle/go-gym/internal/repository"
	"sync"
)

type UserMemoryRepository struct {
	mu     sync.RWMutex
	users  map[domain.UserID]domain.User
	emails map[string]domain.UserID
}

func NewUserMemoryRepository() *UserMemoryRepository {
	return &UserMemoryRepository{
		users:  make(map[domain.UserID]domain.User),
		emails: make(map[string]domain.UserID),
	}
}

func (r *UserMemoryRepository) Create(user domain.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.users[user.ID]; exists {
		return repository.ErrAlreadyExists
	}

	if _, exists := r.emails[user.Email]; exists {
		return repository.ErrAlreadyExists
	}

	r.users[user.ID] = user
	r.emails[user.Email] = user.ID

	return nil
}

func (r *UserMemoryRepository) FindByID(id domain.UserID) (*domain.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	user, exists := r.users[id]
	if !exists {
		return nil, repository.ErrNotFound
	}

	return &user, nil
}

func (r *UserMemoryRepository) FindByEmail(email string) (*domain.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	userID, exists := r.emails[email]
	if !exists {
		return nil, repository.ErrNotFound
	}

	user := r.users[userID]
	return &user, nil
}

func (r *UserMemoryRepository) Update(user domain.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.users[user.ID]; !exists {
		return repository.ErrNotFound
	}

	r.users[user.ID] = user
	return nil
}
