package postgres

import (
	"database/sql"
	"errors"

	"github.com/Facundo-Mourelle/go-gym/internal/domain"
	"github.com/Facundo-Mourelle/go-gym/internal/repository"
	"github.com/lib/pq"
)

type EquipmentPostgresRepository struct {
	db *sql.DB
}

func NewEquipmentPostgresRepository(db *sql.DB) *EquipmentPostgresRepository {
	return &EquipmentPostgresRepository{db: db}
}

func (r *EquipmentPostgresRepository) FindByID(id domain.EquipmentID) (domain.Equipment, error) {
	query := `SELECT id, name, type, manufacturer, user_id, actual_weight, 
		cable_pulley_type, cable_stack_weights, resistance_profile_id, resistance_profile_name 
		FROM equipment WHERE id = $1`
	var e domain.Equipment
	var profileID, profileName, manufacturer, userID, pulleyType sql.NullString
	var actualWeight sql.NullFloat64
	var stackWeights pq.Float64Array

	err := r.db.QueryRow(query, id).Scan(
		&e.ID, &e.Name, &e.Type, &manufacturer, &userID, &actualWeight,
		&pulleyType, &stackWeights, &profileID, &profileName,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return domain.Equipment{}, repository.ErrNotFound
	}
	if err != nil {
		return domain.Equipment{}, err
	}
	if manufacturer.Valid {
		e.Manufacturer = manufacturer.String
	}
	if userID.Valid {
		e.UserID = userID.String
	}
	if actualWeight.Valid {
		e.ActualWeight = actualWeight.Float64
	}
	if pulleyType.Valid {
		e.PulleyType = pulleyType.String
	}
	e.StackWeights = stackWeights
	if profileID.Valid {
		e.ResistanceProfileID = profileID.String
	}
	if profileName.Valid {
		e.ResistanceProfileName = profileName.String
	}
	return e, nil
}

func (r *EquipmentPostgresRepository) Create(equipment domain.Equipment) error {
	query := `INSERT INTO equipment (id, name, type, manufacturer, user_id, actual_weight, 
		cable_pulley_type, cable_stack_weights, resistance_profile_id, resistance_profile_name) 
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`
	var profileID, profileName sql.NullString
	if equipment.ResistanceProfileID != "" {
		profileID = sql.NullString{String: equipment.ResistanceProfileID, Valid: true}
	}
	if equipment.ResistanceProfileName != "" {
		profileName = sql.NullString{String: equipment.ResistanceProfileName, Valid: true}
	}
	_, err := r.db.Exec(query,
		equipment.ID, equipment.Name, equipment.Type, equipment.Manufacturer,
		equipment.UserID, equipment.ActualWeight, equipment.PulleyType,
		pq.Array(equipment.StackWeights), profileID, profileName,
	)
	return err
}

func (r *EquipmentPostgresRepository) FindAll() ([]domain.Equipment, error) {
	query := `SELECT id, name, type, manufacturer, user_id, actual_weight, 
		cable_pulley_type, cable_stack_weights, resistance_profile_id, resistance_profile_name 
		FROM equipment`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []domain.Equipment
	for rows.Next() {
		var e domain.Equipment
		var profileID, profileName, manufacturer, userID, pulleyType sql.NullString
		var actualWeight sql.NullFloat64
		var stackWeights pq.Float64Array

		if err := rows.Scan(
			&e.ID, &e.Name, &e.Type, &manufacturer, &userID, &actualWeight,
			&pulleyType, &stackWeights, &profileID, &profileName,
		); err != nil {
			return nil, err
		}
		if manufacturer.Valid {
			e.Manufacturer = manufacturer.String
		}
		if userID.Valid {
			e.UserID = userID.String
		}
		if actualWeight.Valid {
			e.ActualWeight = actualWeight.Float64
		}
		if pulleyType.Valid {
			e.PulleyType = pulleyType.String
		}
		e.StackWeights = stackWeights
		if profileID.Valid {
			e.ResistanceProfileID = profileID.String
		}
		if profileName.Valid {
			e.ResistanceProfileName = profileName.String
		}
		result = append(result, e)
	}
	return result, nil
}

func (r *EquipmentPostgresRepository) FindByType(t domain.EquipmentType) ([]domain.Equipment, error) {
	query := `SELECT id, name, type, manufacturer, user_id, actual_weight, 
		cable_pulley_type, cable_stack_weights, resistance_profile_id, resistance_profile_name 
		FROM equipment WHERE type = $1`
	rows, err := r.db.Query(query, t)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []domain.Equipment
	for rows.Next() {
		var e domain.Equipment
		var profileID, profileName, manufacturer, userID, pulleyType sql.NullString
		var actualWeight sql.NullFloat64
		var stackWeights pq.Float64Array

		if err := rows.Scan(
			&e.ID, &e.Name, &e.Type, &manufacturer, &userID, &actualWeight,
			&pulleyType, &stackWeights, &profileID, &profileName,
		); err != nil {
			return nil, err
		}
		if manufacturer.Valid {
			e.Manufacturer = manufacturer.String
		}
		if userID.Valid {
			e.UserID = userID.String
		}
		if actualWeight.Valid {
			e.ActualWeight = actualWeight.Float64
		}
		if pulleyType.Valid {
			e.PulleyType = pulleyType.String
		}
		e.StackWeights = stackWeights
		if profileID.Valid {
			e.ResistanceProfileID = profileID.String
		}
		if profileName.Valid {
			e.ResistanceProfileName = profileName.String
		}
		result = append(result, e)
	}
	return result, nil
}

func (r *EquipmentPostgresRepository) Update(equipment domain.Equipment) error {
	query := `UPDATE equipment SET name = $1, type = $2, manufacturer = $3, user_id = $4, 
		actual_weight = $5, cable_pulley_type = $6, cable_stack_weights = $7, 
		resistance_profile_id = $8, resistance_profile_name = $9 WHERE id = $10`
	var profileID, profileName sql.NullString
	if equipment.ResistanceProfileID != "" {
		profileID = sql.NullString{String: equipment.ResistanceProfileID, Valid: true}
	}
	if equipment.ResistanceProfileName != "" {
		profileName = sql.NullString{String: equipment.ResistanceProfileName, Valid: true}
	}
	result, err := r.db.Exec(query,
		equipment.Name, equipment.Type, equipment.Manufacturer, equipment.UserID,
		equipment.ActualWeight, equipment.PulleyType, pq.Array(equipment.StackWeights),
		profileID, profileName, equipment.ID,
	)
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

func (r *EquipmentPostgresRepository) Delete(id domain.EquipmentID) error {
	result, err := r.db.Exec("DELETE FROM equipment WHERE id = $1", id)
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
