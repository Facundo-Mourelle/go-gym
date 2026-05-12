package postgres

import (
	"database/sql"

	"github.com/Facundo-Mourelle/go-gym/internal/domain"
	"github.com/Facundo-Mourelle/go-gym/internal/repository"
	"github.com/lib/pq"
)

type RoutinePostgresRepository struct {
	db *sql.DB
}

func NewRoutinePostgresRepository(db *sql.DB) *RoutinePostgresRepository {
	return &RoutinePostgresRepository{db: db}
}

func (r *RoutinePostgresRepository) Create(routine domain.Routine) error {
	query := `
		INSERT INTO routines (id, user_id, name, description, movement_patterns, is_preset, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`

	_, err := r.db.Exec(
		query,
		routine.ID,
		routine.UserID,
		routine.Name,
		routine.Description,
		pq.Array(routine.MovementPatterns),
		routine.IsPreset,
		routine.CreatedAt,
		routine.UpdatedAt,
	)

	return err
}

func (r *RoutinePostgresRepository) scanRoutine(scanner interface {
	Scan(dest ...any) error
}) (domain.Routine, error) {
	var routine domain.Routine
	var patternStrings []string

	err := scanner.Scan(
		&routine.ID,
		&routine.UserID,
		&routine.Name,
		&routine.Description,
		pq.Array(&patternStrings),
		&routine.IsPreset,
		&routine.CreatedAt,
		&routine.UpdatedAt,
	)
	if err != nil {
		return domain.Routine{}, err
	}

	// Convert []string to []domain.MovementPattern
	routine.MovementPatterns = make([]domain.MovementPattern, len(patternStrings))
	for i, s := range patternStrings {
		routine.MovementPatterns[i] = domain.MovementPattern(s)
	}

	return routine, nil
}

func (r *RoutinePostgresRepository) FindByID(id domain.RoutineID) (domain.Routine, error) {
	query := `
		SELECT id, user_id, name, description, movement_patterns, is_preset, created_at, updated_at
		FROM routines
		WHERE id = $1
	`

	routine, err := r.scanRoutine(r.db.QueryRow(query, id))
	if err == sql.ErrNoRows {
		return domain.Routine{}, repository.ErrNotFound
	}
	if err != nil {
		return domain.Routine{}, err
	}

	return routine, nil
}

func (r *RoutinePostgresRepository) FindByUser(userID string) ([]domain.Routine, error) {
	query := `
		SELECT id, user_id, name, description, movement_patterns, is_preset, created_at, updated_at
		FROM routines
		WHERE user_id = $1
		ORDER BY created_at DESC
	`

	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	routines := make([]domain.Routine, 0)
	for rows.Next() {
		routine, err := r.scanRoutine(rows)
		if err != nil {
			return nil, err
		}
		routines = append(routines, routine)
	}

	return routines, rows.Err()
}

func (r *RoutinePostgresRepository) Delete(id domain.RoutineID) error {
	query := `DELETE FROM routines WHERE id = $1`

	result, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return repository.ErrNotFound
	}

	return nil
}