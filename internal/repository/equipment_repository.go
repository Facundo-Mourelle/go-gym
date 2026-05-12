package repository

import (
	"github.com/Facundo-Mourelle/go-gym/internal/domain"
)

type EquipmentRepository interface {
	Create(equipment domain.Equipment) error

	FindByID(id domain.EquipmentID) (domain.Equipment, error)

	FindAll() ([]domain.Equipment, error)

	FindByType(equipmentType domain.EquipmentType) ([]domain.Equipment, error)

	Update(equipment domain.Equipment) error

	Delete(id domain.EquipmentID) error
}
