package resistance

import (
	"errors"
	"math"
)

// For machines that do not use pulleys, should always use this configuration
var DirectDrive = &PulleyProfile{
	id:              "direct-drive",
	pulleyConfig:    Pulley1to1,
	mechanicalRatio: 1,
	frictionLoss:    0,
}

func (m *ResistanceProfile) MachineCalculateEffectiveLoad(rawLoad float64) (float64, error) {
	if rawLoad < 0 {
		return 0, errors.New("invalid rawLoad")
	}
	if m.position < 0 || m.position > 1 {
		return 0, errors.New("invalid position: must be in range [0,1]")
	}
	if m.pulley == nil {
		return 0, errors.New("pulley profile not set")
	}

	// First apply pulley mechanics
	load, err := m.pulley.PulleyCalculateEffectiveLoad(rawLoad)
	if err != nil {
		return 0, err
	}

	var multiplier float64

	switch m.profileType {
	case Ascending:
		multiplier = m.position
	case Descending:
		multiplier = 1 - m.position
	case Bell:
		// Peak in the middle
		multiplier = 4 * m.position * (1 - m.position)
	case Uniform:
		multiplier = 1
	default:
		return 0, errors.New("unknown resistance profile")
	}

	effectiveLoad := load * multiplier
	return math.Round(effectiveLoad*100) / 100, nil
}
