package service

import (
	"fmt"
	"time"

	"github.com/Facundo-Mourelle/go-gym/internal/domain"
	"github.com/Facundo-Mourelle/go-gym/internal/repository"
)

type EquipmentService struct {
	equipmentRepo repository.EquipmentRepository
}

func NewEquipmentService(equipmentRepo repository.EquipmentRepository) *EquipmentService {
	return &EquipmentService{equipmentRepo: equipmentRepo}
}

func (s *EquipmentService) ListAll() ([]domain.Equipment, error) {
	return s.equipmentRepo.FindAll()
}

func (s *EquipmentService) Create(equipment domain.Equipment) (domain.Equipment, error) {
	equipment.ID = domain.EquipmentID(fmt.Sprintf("eq_%d", time.Now().UnixNano()))
	if err := s.equipmentRepo.Create(equipment); err != nil {
		return domain.Equipment{}, err
	}
	return equipment, nil
}

func (s *EquipmentService) GetEquipment(id domain.EquipmentID) (domain.Equipment, error) {
	return s.equipmentRepo.FindByID(id)
}

func (s *EquipmentService) Update(equipment domain.Equipment) error {
	return s.equipmentRepo.Update(equipment)
}

func (s *EquipmentService) Delete(id domain.EquipmentID) error {
	return s.equipmentRepo.Delete(id)
}
