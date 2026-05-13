package postgres

import (
	"database/sql"
	"time"

	"github.com/Facundo-Mourelle/go-gym/internal/domain"
	"github.com/Facundo-Mourelle/go-gym/internal/repository"
)

type WorkoutPostgresRepository struct {
	db *sql.DB
}

func NewWorkoutPostgresRepository(db *sql.DB) *WorkoutPostgresRepository {
	return &WorkoutPostgresRepository{db: db}
}

func (r *WorkoutPostgresRepository) Create(workout domain.WorkoutPlan) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Insert workout plan
	query := `
        INSERT INTO workout_plans (id, user_id, name, description) VALUES ($1, $2, $3, $4)`
	_, err = tx.Exec(
		query,
		workout.ID,
		workout.UserID,
		workout.Name,
		workout.Description,
	)
	if err != nil {
		return err
	}

	// Insert exercises
	for _, exercise := range workout.Exercises {
		exQuery := `
        INSERT INTO workout_exercises 
        (id, workout_plan_id, exercise_id, "order", sets, reps, reps_in_reserve, notes)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
    `

		_, err = tx.Exec(
			exQuery,
			exercise.ID,
			workout.ID,
			exercise.ExerciseID,
			exercise.Order,
			exercise.Sets,
			exercise.Reps,
			exercise.RepsInReserve,
			exercise.Notes,
		)

		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (r *WorkoutPostgresRepository) FindByID(id domain.WorkoutPlanID) (domain.WorkoutPlan, error) {
	query :=
		`SELECT id, user_id, name, description, created_at, updated_at
     FROM workout_plans 
     WHERE id = $1`
	var workout domain.WorkoutPlan

	err := r.db.QueryRow(query, id).Scan(
		&workout.ID,
		&workout.UserID,
		&workout.Name,
		&workout.Description,
		&workout.CreatedAt,
		&workout.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return domain.WorkoutPlan{}, repository.ErrNotFound
	}
	if err != nil {
		return domain.WorkoutPlan{}, err
	}

	// Load exercises
	if err := r.loadExercises(&workout); err != nil {
		return domain.WorkoutPlan{}, err
	}

	return workout, nil
}

func (r *WorkoutPostgresRepository) loadExercises(workout *domain.WorkoutPlan) error {
	query :=
		`SELECT id, exercise_id, "order", sets, reps, reps_in_reserve, notes
         FROM workout_exercises 
         WHERE workout_plan_id = $1 
         ORDER BY "order"`
	rows, err := r.db.Query(query, workout.ID)
	if err != nil {
		return err
	}
	defer rows.Close()

	workout.Exercises = make([]domain.WorkoutExercise, 0)

	for rows.Next() {
		var exercise domain.WorkoutExercise

		err := rows.Scan(
			&exercise.ID,
			&exercise.ExerciseID,
			&exercise.Order,
			&exercise.Sets,
			&exercise.Reps,
			&exercise.RepsInReserve,
			&exercise.Notes,
		)

		if err != nil {
			return err
		}

		workout.Exercises = append(workout.Exercises, exercise)
	}

	return rows.Err()
}

func (r *WorkoutPostgresRepository) FindByUser(userID string) ([]domain.WorkoutPlan, error) {
	query :=
		`SELECT id, user_id, name, description, created_at, updated_at
         FROM workout_plans 
         WHERE user_id = $1 
         ORDER BY created_at DESC`
	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	workouts := make([]domain.WorkoutPlan, 0)

	for rows.Next() {
		var workout domain.WorkoutPlan

		err := rows.Scan(
			&workout.ID,
			&workout.UserID,
			&workout.Name,
			&workout.Description,
			&workout.CreatedAt,
			&workout.UpdatedAt,
		)

		if err != nil {
			return nil, err
		}

		if err := r.loadExercises(&workout); err != nil {
			return nil, err
		}

		workouts = append(workouts, workout)
	}

	return workouts, rows.Err()
}

func (r *WorkoutPostgresRepository) Update(workout domain.WorkoutPlan) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	query :=
		`UPDATE workout_plans 
    SET name = $1, description = $2, updated_at = $3 
    WHERE id = $4`

	result, err := tx.Exec(
		query,
		workout.Name,
		workout.Description,
		time.Now(),
		workout.ID,
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

	// Delete and re-insert exercises
	if _, err := tx.Exec("DELETE FROM workout_exercises WHERE workout_plan_id = $1", workout.ID); err != nil {
		return err
	}

	for _, exercise := range workout.Exercises {
		exQuery := `
        INSERT INTO workout_exercises 
        (id, workout_plan_id, exercise_id, "order", sets, reps, reps_in_reserve, notes) 
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`

		_, err = tx.Exec(
			exQuery,
			exercise.ID,
			workout.ID,
			exercise.ExerciseID,
			exercise.Order,
			exercise.Sets,
			exercise.Reps,
			exercise.RepsInReserve,
			exercise.Notes,
		)

		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (r *WorkoutPostgresRepository) Delete(id domain.WorkoutPlanID) error {
	_, err := r.db.Exec("DELETE FROM workout_exercises WHERE workout_plan_id = $1", id)
	if err != nil {
		return err
	}
	
	query := `DELETE FROM workout_plans WHERE id = $1`
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
