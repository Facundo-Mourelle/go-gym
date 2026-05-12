package postgres

import (
	"database/sql"
	"errors"

	"github.com/Facundo-Mourelle/go-gym/internal/domain"
	"github.com/Facundo-Mourelle/go-gym/internal/repository"
)

type EquipmentPostgresRepository struct {
	db *sql.DB
}

func NewEquipmentPostgresRepository(db *sql.DB) *EquipmentPostgresRepository {
	return &EquipmentPostgresRepository{db: db}
}

func (r *EquipmentPostgresRepository) FindByID(id domain.EquipmentID) (domain.Equipment, error) {
	query := `SELECT id, name, type, resistance_profile_id FROM equipment WHERE id = $1`
	var e domain.Equipment
	var profileID sql.NullString
	err := r.db.QueryRow(query, id).Scan(&e.ID, &e.Name, &e.Type, &profileID)
	if errors.Is(err, sql.ErrNoRows) {
		return domain.Equipment{}, repository.ErrNotFound
	}
	if profileID.Valid {
		e.ResistanceProfileID = profileID.String
	}
	return e, err
}

func (r *EquipmentPostgresRepository) Create(equipment domain.Equipment) error {
	query := `INSERT INTO equipment (id, name, type, resistance_profile_id) VALUES ($1, $2, $3, $4)`
	var profileID sql.NullString
	if equipment.ResistanceProfileID != "" {
		profileID = sql.NullString{String: equipment.ResistanceProfileID, Valid: true}
	}
	_, err := r.db.Exec(query, equipment.ID, equipment.Name, equipment.Type, profileID)
	return err
}

func (r *EquipmentPostgresRepository) FindAll() ([]domain.Equipment, error) {
	query := `SELECT id, name, type, resistance_profile_id FROM equipment`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []domain.Equipment
	for rows.Next() {
		var e domain.Equipment
		var profileID sql.NullString
		if err := rows.Scan(&e.ID, &e.Name, &e.Type, &profileID); err != nil {
			return nil, err
		}
		if profileID.Valid {
			e.ResistanceProfileID = profileID.String
		}
		result = append(result, e)
	}
	return result, nil
}

func (r *EquipmentPostgresRepository) Delete(id domain.EquipmentID) error {
	// TODO: implement
	return nil
}

func (r *EquipmentPostgresRepository) FindByType(t domain.EquipmentType) ([]domain.Equipment, error) {
	return nil, nil
}

func (r *EquipmentPostgresRepository) Update(equipment domain.Equipment) error {
	return nil
}
