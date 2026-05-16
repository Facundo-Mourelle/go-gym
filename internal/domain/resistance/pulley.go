package resistance

import (
	"errors"
	"math"
)

type PulleyProfile struct {
	id           string
	pulleyConfig PulleyConfiguration
	// Effective load multiplier (1-n)
	mechanicalRatio float64
	// Percentage loss (0-1)
	frictionLoss float64
}

func NewPulleyProfile(id string, pulleyConfig PulleyConfiguration, mechanicalRatio, frictionLoss float64) (*PulleyProfile, error) {

	if mechanicalRatio < 1 {
		return nil, errors.New("invalid mechanicalRatio value: must be greater than 0")
	}
	if frictionLoss < 0 || frictionLoss > 1 {
		return nil, errors.New("invalid frictionLoss value: must be in range (0,1)")
	}

	return &PulleyProfile{
		id:              id,
		pulleyConfig:    pulleyConfig,
		mechanicalRatio: mechanicalRatio,
		frictionLoss:    frictionLoss,
	}, nil
}

func (p *PulleyProfile) PulleyCalculateEffectiveLoad(rawLoad float64) (float64, error) {
	if rawLoad < 0 {
		return 0, errors.New("invalid rawLoad: must be greater than 0")
	}

	// Apply mechanical advantage
	effectiveLoad := rawLoad * p.mechanicalRatio

	// Apply friction loss
	effectiveLoad = effectiveLoad * (1 - p.frictionLoss)

	return math.Round(effectiveLoad*100) / 100, nil
}
