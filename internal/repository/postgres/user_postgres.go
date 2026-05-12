package postgres

import (
	"database/sql"
	"errors"

	"github.com/Facundo-Mourelle/go-gym/internal/domain"
	"github.com/Facundo-Mourelle/go-gym/internal/repository"
)

type UserPostgresRepository struct {
	db *sql.DB
}

func NewUserPostgresRepository(db *sql.DB) *UserPostgresRepository {
	return &UserPostgresRepository{db: db}
}

func (r *UserPostgresRepository) Create(user domain.User) error {
	query := `
		INSERT INTO users (id, email, password_hash, name, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`
	_, err := r.db.Exec(query, user.ID, user.Email, user.PasswordHash, user.Name, user.CreatedAt, user.UpdatedAt)
	// TODO: Handle unique constraint violation -> repository.ErrAlreadyExists
	return err
}

func (r *UserPostgresRepository) FindByID(id domain.UserID) (*domain.User, error) {
	query := `SELECT id, email, password_hash, name, created_at, updated_at FROM users WHERE id = $1`
	var u domain.User
	err := r.db.QueryRow(query, id).Scan(&u.ID, &u.Email, &u.PasswordHash, &u.Name, &u.CreatedAt, &u.UpdatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, repository.ErrNotFound
	}
	return &u, err
}

func (r *UserPostgresRepository) FindByEmail(email string) (*domain.User, error) {
	query := `SELECT id, email, password_hash, name, created_at, updated_at FROM users WHERE email = $1`
	var u domain.User
	err := r.db.QueryRow(query, email).Scan(&u.ID, &u.Email, &u.PasswordHash, &u.Name, &u.CreatedAt, &u.UpdatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, repository.ErrNotFound
	}
	return &u, err
}

func (r *UserPostgresRepository) Update(user domain.User) error {
	query := `
		UPDATE users 
		SET email = $1, password_hash = $2, name = $3, updated_at = $4
		WHERE id = $5
	`
	res, err := r.db.Exec(query, user.Email, user.PasswordHash, user.Name, user.UpdatedAt, user.ID)
	if err != nil {
		return err
	}
	rows, _ := res.RowsAffected()
	if rows == 0 {
		return repository.ErrNotFound
	}
	return nil
}
