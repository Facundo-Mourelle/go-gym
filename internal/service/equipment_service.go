package service

import (
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
