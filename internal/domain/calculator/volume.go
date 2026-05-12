package calculator

import (
	"github.com/Facundo-Mourelle/go-gym/internal/domain"
	"github.com/Facundo-Mourelle/go-gym/internal/domain/resistance"
)

type VolumeCalculator struct {
	registry *resistance.Registry
}

func NewVolumeCalculator(registry *resistance.Registry) *VolumeCalculator {
	return &VolumeCalculator{registry: registry}
}

func (vc *VolumeCalculator) CalculateEffectiveLoad(rawLoad float64, profile string) float64 {
	return rawLoad
}

func (vc *VolumeCalculator) CalculateMuscleGroupVolume(session domain.Session, exercises map[domain.ExerciseID]domain.Exercise) map[domain.MuscleGroup]float64 {
	volumes := make(map[domain.MuscleGroup]float64)
	return volumes
}
