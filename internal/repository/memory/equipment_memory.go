package memory

import (
	"sync"

	"github.com/Facundo-Mourelle/go-gym/internal/domain"
	"github.com/Facundo-Mourelle/go-gym/internal/repository"
)

type EquipmentMemoryRepository struct {
	mu        sync.RWMutex
	equipment map[domain.EquipmentID]domain.Equipment
}

func NewEquipmentMemoryRepository() *EquipmentMemoryRepository {
	return &EquipmentMemoryRepository{
		equipment: make(map[domain.EquipmentID]domain.Equipment),
	}
}

func (r *EquipmentMemoryRepository) Create(equipment domain.Equipment) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.equipment[equipment.ID]; exists {
		return repository.ErrAlreadyExists
	}

	r.equipment[equipment.ID] = equipment
	return nil
}

func (r *EquipmentMemoryRepository) FindByID(id domain.EquipmentID) (domain.Equipment, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	equipment, exists := r.equipment[id]
	if !exists {
		return domain.Equipment{}, repository.ErrNotFound
	}

	return equipment, nil
}

func (r *EquipmentMemoryRepository) FindAll() ([]domain.Equipment, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	equipment := make([]domain.Equipment, 0, len(r.equipment))
	for _, eq := range r.equipment {
		equipment = append(equipment, eq)
	}

	return equipment, nil
}

func (r *EquipmentMemoryRepository) FindByType(equipmentType domain.EquipmentType) ([]domain.Equipment, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	equipment := make([]domain.Equipment, 0)
	for _, eq := range r.equipment {
		if eq.Type == equipmentType {
			equipment = append(equipment, eq)
		}
	}

	return equipment, nil
}

func (r *EquipmentMemoryRepository) Update(equipment domain.Equipment) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.equipment[equipment.ID]; !exists {
		return repository.ErrNotFound
	}

	r.equipment[equipment.ID] = equipment
	return nil
}

func (r *EquipmentMemoryRepository) Delete(id domain.EquipmentID) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.equipment[id]; !exists {
		return repository.ErrNotFound
	}

	delete(r.equipment, id)
	return nil
}
