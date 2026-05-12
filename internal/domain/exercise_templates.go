package domain

import (
	"fmt"

)

// ExerciseTemplate provides a starting point for creating exercises
type ExerciseTemplate struct {
	Name               string
	Description        string
	PrimaryPatterns    []MovementPatternContribution
	SecondaryPatterns  []MovementPatternContribution
	SuggestedEquipment []EquipmentType
	Category           string
}

var ExerciseTemplates = map[string]ExerciseTemplate{
	"barbell_row": {
		Name:        "Barbell Row",
		Description: "Hip-hinged horizontal pulling movement",
		PrimaryPatterns: []MovementPatternContribution{
			{Pattern: HorizontalPull, Contribution: 1.0, RangeOfMotion: "full"},
			{Pattern: ShoulderAdduction, Contribution: 0.7, RangeOfMotion: "full"},
		},
		SecondaryPatterns: []MovementPatternContribution{
			{Pattern: ElbowFlexion, Contribution: 0.4, RangeOfMotion: "full"},
			{Pattern: HipHinge, Contribution: 0.3, RangeOfMotion: "partial"},
		},
		SuggestedEquipment: []EquipmentType{EquipmentTypeFreeWeight},
		Category:           "back",
	},

	"cable_row": {
		Name:        "Cable Row",
		Description: "Seated horizontal pulling movement",
		PrimaryPatterns: []MovementPatternContribution{
			{Pattern: HorizontalPull, Contribution: 1.0, RangeOfMotion: "full"},
			{Pattern: ShoulderAdduction, Contribution: 0.6, RangeOfMotion: "full"},
		},
		SecondaryPatterns: []MovementPatternContribution{
			{Pattern: ElbowFlexion, Contribution: 0.3, RangeOfMotion: "full"},
		},
		SuggestedEquipment: []EquipmentType{EquipmentTypeMachine},
		Category:           "back",
	},

	"lat_pulldown": {
		Name:        "Lat Pulldown",
		Description: "Vertical pulling from overhead",
		PrimaryPatterns: []MovementPatternContribution{
			{Pattern: VerticalPull, Contribution: 1.0, RangeOfMotion: "full"},
		},
		SecondaryPatterns: []MovementPatternContribution{
			{Pattern: ElbowFlexion, Contribution: 0.3, RangeOfMotion: "full"},
		},
		SuggestedEquipment: []EquipmentType{EquipmentTypeMachine},
		Category:           "back",
	},

	"bench_press": {
		Name:        "Bench Press",
		Description: "Horizontal pressing movement",
		PrimaryPatterns: []MovementPatternContribution{
			{Pattern: HorizontalPush, Contribution: 1.0, RangeOfMotion: "full"},
		},
		SecondaryPatterns:  []MovementPatternContribution{},
		SuggestedEquipment: []EquipmentType{EquipmentTypeFreeWeight},
		Category:           "chest",
	},

	"dumbbell_press": {
		Name:        "Dumbbell Press",
		Description: "Horizontal pressing with dumbbells",
		PrimaryPatterns: []MovementPatternContribution{
			{Pattern: HorizontalPush, Contribution: 1.0, RangeOfMotion: "full"},
		},
		SecondaryPatterns:  []MovementPatternContribution{},
		SuggestedEquipment: []EquipmentType{EquipmentTypeFreeWeight},
		Category:           "chest",
	},
}

func CreateFromTemplate(templateID string, userID string) (*ExerciseBuilder, error) {
	template, exists := ExerciseTemplates[templateID]
	if !exists {
		return nil, fmt.Errorf("template not found: %s", templateID)
	}

	builder := NewExerciseBuilder().
		WithName(template.Name).
		WithDescription(template.Description).
		WithUserID(userID).
		WithEquipment(template.SuggestedEquipment...)

	for _, pattern := range template.PrimaryPatterns {
		builder.WithPrimaryPattern(
			pattern.Pattern,
			pattern.Contribution,
			pattern.RangeOfMotion,
			pattern.Notes,
		)
	}

	for _, pattern := range template.SecondaryPatterns {
		builder.WithSecondaryPattern(
			pattern.Pattern,
			pattern.Contribution,
			pattern.RangeOfMotion,
			pattern.Notes,
		)
	}

	return builder, nil
}
