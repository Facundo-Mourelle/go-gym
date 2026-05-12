package postgres

import (
	"encoding/json"
	"time"

	"database/sql"

	"github.com/Facundo-Mourelle/go-gym/internal/domain"
	"github.com/Facundo-Mourelle/go-gym/internal/repository"
)

type SessionPostgresRepository struct {
	db *sql.DB
}

func NewSessionPostgresRepository(db *sql.DB) *SessionPostgresRepository {
	return &SessionPostgresRepository{db: db}
}

func (r *SessionPostgresRepository) Create(session domain.Session) error {
	query := `
        INSERT INTO sessions (id, user_id, workout_plan_id, started_at, completed_at, performed_sets, notes)
        VALUES ($1, $2, $3, $4, $5, $6, $7)
    `

	performedSetsJSON, err := json.Marshal(session.PerformedSets)
	if err != nil {
		return err
	}

	_, err = r.db.Exec(
		query,
		session.ID,
		session.UserID,
		session.WorkoutPlanID,
		session.StartedAt,
		session.CompletedAt,
		performedSetsJSON,
		session.Notes,
	)

	return err
}

func (r *SessionPostgresRepository) FindByID(id domain.SessionID) (domain.Session, error) {
	query := `
        SELECT id, user_id, workout_plan_id, started_at, completed_at, performed_sets, notes
        FROM sessions
        WHERE id = $1
    `

	var session domain.Session
	var workoutPlanID sql.NullString
	var completedAt sql.NullTime
	var performedSetsJSON []byte

	err := r.db.QueryRow(query, id).Scan(
		&session.ID,
		&session.UserID,
		&workoutPlanID,
		&session.StartedAt,
		&completedAt,
		&performedSetsJSON,
		&session.Notes,
	)

	if err == sql.ErrNoRows {
		return domain.Session{}, repository.ErrNotFound
	}
	if err != nil {
		return domain.Session{}, err
	}

	if workoutPlanID.Valid {
		wpID := domain.WorkoutPlanID(workoutPlanID.String)
		session.WorkoutPlanID = wpID
	}

	if completedAt.Valid {
		session.CompletedAt = completedAt.Time
	}

	if performedSetsJSON != nil {
		if err := json.Unmarshal(performedSetsJSON, &session.PerformedSets); err != nil {
			return domain.Session{}, err
		}
	}

	return session, nil
}

func (r *SessionPostgresRepository) AddPerformedSet(sessionID domain.SessionID, set domain.PerformedSet) error {
	session, err := r.FindByID(sessionID)
	if err != nil {
		return err
	}

	session.PerformedSets = append(session.PerformedSets, set)

	return r.Update(session)
}

func (r *SessionPostgresRepository) Complete(id domain.SessionID, completedAt time.Time) error {
	query := `
        UPDATE sessions SET completed_at = $1 WHERE id = $2
    `

	_, err := r.db.Exec(query, completedAt, id)
	return err
}

func (r *SessionPostgresRepository) Delete(id domain.SessionID) error {
	return nil
}

func (r *SessionPostgresRepository) FindByUser(id string) ([]domain.Session, error) {
	query := `
		SELECT id, user_id, workout_plan_id, started_at, completed_at, performed_sets, notes
		FROM sessions
		WHERE user_id = $1
		ORDER BY started_at DESC
	`

	rows, err := r.db.Query(query, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var sessions []domain.Session
	for rows.Next() {
		var session domain.Session
		var workoutPlanID sql.NullString
		var completedAt sql.NullTime
		var performedSetsJSON []byte

		err := rows.Scan(
			&session.ID,
			&session.UserID,
			&workoutPlanID,
			&session.StartedAt,
			&completedAt,
			&performedSetsJSON,
			&session.Notes,
		)
		if err != nil {
			return nil, err
		}

		if workoutPlanID.Valid {
			session.WorkoutPlanID = domain.WorkoutPlanID(workoutPlanID.String)
		}

		if completedAt.Valid {
			session.CompletedAt = completedAt.Time
		}

		if performedSetsJSON != nil {
			// Handle legacy data where performed_sets is stored as {} instead of []
			if len(performedSetsJSON) > 0 && performedSetsJSON[0] == '{' {
				session.PerformedSets = make([]domain.PerformedSet, 0)
			} else if err := json.Unmarshal(performedSetsJSON, &session.PerformedSets); err != nil {
				return nil, err
			}
		}

		sessions = append(sessions, session)
	}

	return sessions, nil
}

func (r *SessionPostgresRepository) FindByUserAndDateRange(userID string, start, end time.Time) ([]domain.Session, error) {
	return nil, nil
}

func (r *SessionPostgresRepository) Update(session domain.Session) error {
	query := `
        UPDATE sessions SET workout_plan_id = $1, completed_at = $2, performed_sets = $3, notes = $4
        WHERE id = $5
    `

	performedSetsJSON, err := json.Marshal(session.PerformedSets)
	if err != nil {
		return err
	}

	_, err = r.db.Exec(
		query,
		session.WorkoutPlanID,
		session.CompletedAt,
		performedSetsJSON,
		session.Notes,
		session.ID,
	)

	return err
}
