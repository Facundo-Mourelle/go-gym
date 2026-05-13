package postgres

import (
	"database/sql"
	"time"

	"github.com/Facundo-Mourelle/go-gym/internal/domain"
	"github.com/Facundo-Mourelle/go-gym/internal/repository"
)

type ExercisePostgresRepository struct {
	db *sql.DB
}

func NewExercisePostgresRepository(db *sql.DB) *ExercisePostgresRepository {
	return &ExercisePostgresRepository{db: db}
}

func (r *ExercisePostgresRepository) Create(exercise domain.Exercise) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Insert exercise
	query := `
        INSERT INTO exercises (id, name, description, user_id, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5, $6)
    `

	_, err = tx.Exec(
		query,
		exercise.ID,
		exercise.Name,
		exercise.Description,
		sql.NullString{String: exercise.UserID, Valid: exercise.UserID != ""},
		exercise.CreatedAt,
		exercise.UpdatedAt,
	)

	if err != nil {
		return err
	}

	// Insert primary patterns
	for _, pattern := range exercise.PrimaryPatterns {
		if err := r.insertPattern(tx, exercise.ID, pattern, true); err != nil {
			return err
		}
	}

	// Insert secondary patterns
	for _, pattern := range exercise.SecondaryPatterns {
		if err := r.insertPattern(tx, exercise.ID, pattern, false); err != nil {
			return err
		}
	}

	// Insert equipment suggestions
	for _, equipmentType := range exercise.SuggestedEquipment {
		equipQuery := `
            INSERT INTO exercise_equipment_suggestions (exercise_id, equipment_type)
            VALUES ($1, $2)
        `
		if _, err := tx.Exec(equipQuery, exercise.ID, equipmentType); err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (r *ExercisePostgresRepository) insertPattern(
	tx *sql.Tx,
	exerciseID domain.ExerciseID,
	pattern domain.MovementPatternContribution,
	isPrimary bool,
) error {
	query := `
        INSERT INTO exercise_pattern_contributions 
        (id, exercise_id, pattern, is_primary, contribution, range_of_motion, notes, created_at)
        VALUES (gen_random_uuid(), $1, $2, $3, $4, $5, $6, $7)
    `

	_, err := tx.Exec(
		query,
		exerciseID,
		pattern.Pattern,
		isPrimary,
		pattern.Contribution,
		pattern.RangeOfMotion,
		pattern.Notes,
		time.Now(),
	)

	return err
}

func (r *ExercisePostgresRepository) FindByID(id domain.ExerciseID) (domain.Exercise, error) {
	query := `
        SELECT id, name, description, user_id, created_at, updated_at
        FROM exercises
        WHERE id = $1
    `

	var exercise domain.Exercise
	var userID sql.NullString

	err := r.db.QueryRow(query, id).Scan(
		&exercise.ID,
		&exercise.Name,
		&exercise.Description,
		&userID,
		&exercise.CreatedAt,
		&exercise.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return domain.Exercise{}, repository.ErrNotFound
	}
	if err != nil {
		return domain.Exercise{}, err
	}

	if userID.Valid {
		exercise.UserID = userID.String
	}

	// Load patterns
	if err := r.loadPatterns(&exercise); err != nil {
		return domain.Exercise{}, err
	}

	// Load equipment
	if err := r.loadEquipment(&exercise); err != nil {
		return domain.Exercise{}, err
	}

	return exercise, nil
}

func (r *ExercisePostgresRepository) loadPatterns(exercise *domain.Exercise) error {
	query := `
        SELECT pattern, is_primary, contribution, range_of_motion, notes
        FROM exercise_pattern_contributions
        WHERE exercise_id = $1
        ORDER BY is_primary DESC, contribution DESC
    `

	rows, err := r.db.Query(query, exercise.ID)
	if err != nil {
		return err
	}
	defer rows.Close()

	exercise.PrimaryPatterns = make([]domain.MovementPatternContribution, 0)
	exercise.SecondaryPatterns = make([]domain.MovementPatternContribution, 0)

	for rows.Next() {
		var pattern domain.MovementPatternContribution
		var isPrimary bool

		var notes sql.NullString
		err := rows.Scan(
			&pattern.Pattern,
			&isPrimary,
			&pattern.Contribution,
			&pattern.RangeOfMotion,
			&notes,
		)
		if notes.Valid {
			pattern.Notes = notes.String
		}

		if err != nil {
			return err
		}

		if isPrimary {
			exercise.PrimaryPatterns = append(exercise.PrimaryPatterns, pattern)
		} else {
			exercise.SecondaryPatterns = append(exercise.SecondaryPatterns, pattern)
		}
	}

	return rows.Err()
}

func (r *ExercisePostgresRepository) loadEquipment(exercise *domain.Exercise) error {
	query := `
        SELECT equipment_type
        FROM exercise_equipment_suggestions
        WHERE exercise_id = $1
    `

	rows, err := r.db.Query(query, exercise.ID)
	if err != nil {
		return err
	}
	defer rows.Close()

	exercise.SuggestedEquipment = make([]domain.EquipmentType, 0)

	for rows.Next() {
		var equipType domain.EquipmentType
		if err := rows.Scan(&equipType); err != nil {
			return err
		}
		exercise.SuggestedEquipment = append(exercise.SuggestedEquipment, equipType)
	}

	return rows.Err()
}

func (r *ExercisePostgresRepository) FindAll() ([]domain.Exercise, error) {
	query := `
        SELECT id, name, description, user_id, created_at, updated_at
        FROM exercises
        ORDER BY name
    `

	return r.queryExercises(query)
}

func (r *ExercisePostgresRepository) FindByUser(userID string) ([]domain.Exercise, error) {
	query := `
        SELECT id, name, description, user_id, created_at, updated_at
        FROM exercises
        WHERE user_id = $1
        ORDER BY name
    `

	return r.queryExercises(query, userID)
}

func (r *ExercisePostgresRepository) FindByPattern(pattern domain.MovementPattern) ([]domain.Exercise, error) {
	query := `
        SELECT DISTINCT e.id, e.name, e.description, e.user_id, e.created_at, e.updated_at
        FROM exercises e
        JOIN exercise_pattern_contributions epc ON e.id = epc.exercise_id
        WHERE epc.pattern = $1 AND epc.is_primary = true
        ORDER BY e.name
    `

	return r.queryExercises(query, pattern)
}

func (r *ExercisePostgresRepository) queryExercises(query string, args ...any) ([]domain.Exercise, error) {
	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	exercises := make([]domain.Exercise, 0)

	for rows.Next() {
		var exercise domain.Exercise
		var userID sql.NullString

		err := rows.Scan(
			&exercise.ID,
			&exercise.Name,
			&exercise.Description,
			&userID,
			&exercise.CreatedAt,
			&exercise.UpdatedAt,
		)

		if err != nil {
			return nil, err
		}

		if userID.Valid {
			exercise.UserID = userID.String
		}

		// Load patterns and equipment
		if err := r.loadPatterns(&exercise); err != nil {
			return nil, err
		}
		if err := r.loadEquipment(&exercise); err != nil {
			return nil, err
		}

		exercises = append(exercises, exercise)
	}

	return exercises, rows.Err()
}

func (r *ExercisePostgresRepository) Update(exercise domain.Exercise) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Update exercise
	query := `
        UPDATE exercises
        SET name = $1, description = $2, updated_at = $3
        WHERE id = $4
    `

	result, err := tx.Exec(
		query,
		exercise.Name,
		exercise.Description,
		time.Now(),
		exercise.ID,
	)

	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return repository.ErrNotFound
	}

	// Delete and re-insert patterns
	if _, err := tx.Exec("DELETE FROM exercise_pattern_contributions WHERE exercise_id = $1", exercise.ID); err != nil {
		return err
	}

	for _, pattern := range exercise.PrimaryPatterns {
		if err := r.insertPattern(tx, exercise.ID, pattern, true); err != nil {
			return err
		}
	}

	for _, pattern := range exercise.SecondaryPatterns {
		if err := r.insertPattern(tx, exercise.ID, pattern, false); err != nil {
			return err
		}
	}

	// Delete and re-insert equipment
	if _, err := tx.Exec("DELETE FROM exercise_equipment_suggestions WHERE exercise_id = $1", exercise.ID); err != nil {
		return err
	}

	for _, equipmentType := range exercise.SuggestedEquipment {
		equipQuery := `
            INSERT INTO exercise_equipment_suggestions (exercise_id, equipment_type)
            VALUES ($1, $2)
        `
		if _, err := tx.Exec(equipQuery, exercise.ID, equipmentType); err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (r *ExercisePostgresRepository) Delete(id domain.ExerciseID) error {
	query := `DELETE FROM exercises WHERE id = $1`

	result, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return repository.ErrNotFound
	}

	return nil
}
