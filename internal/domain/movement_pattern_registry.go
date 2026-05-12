package domain

import (
)

var MovementPatternRegistry = map[MovementPattern]MovementPatternInfo{
	// Horizontal Patterns
	HorizontalPush: {
		Pattern:          HorizontalPush,
		Name:             "Horizontal Push",
		Plane:            Transverse,
		PrimaryMuscles:   []MuscleGroup{Chest},
		SecondaryMuscles: []MuscleGroup{Triceps},
		JointActions:     []string{"shoulder_horizontal_adduction", "elbow_extension", "scapular_protraction"},
		Description:      "Pushing movement away from body in horizontal plane - bench press, push-ups",
	},
	HorizontalPull: {
		Pattern:          HorizontalPull,
		Name:             "Horizontal Pull",
		Plane:            Transverse,
		PrimaryMuscles:   []MuscleGroup{Traps},
		SecondaryMuscles: []MuscleGroup{Lats},
		JointActions:     []string{"shoulder_horizontal_abduction", "scapular_retraction"},
		Description:      "Pulling movement toward body in horizontal plane - rows, face pulls",
	},

	// Vertical Patterns
	VerticalPush: {
		Pattern:          VerticalPush,
		Name:             "Vertical Push",
		Plane:            Frontal,
		PrimaryMuscles:   []MuscleGroup{Shoulders},
		SecondaryMuscles: []MuscleGroup{Triceps},
		JointActions:     []string{"shoulder_abduction", "elbow_extension", "scapular_upward_rotation"},
		Description:      "Pressing movement overhead - overhead press, military press",
	},
	VerticalPull: {
		Pattern:          VerticalPull,
		Name:             "Vertical Pull",
		Plane:            Frontal,
		PrimaryMuscles:   []MuscleGroup{Lats},
		SecondaryMuscles: []MuscleGroup{Shoulders},
		JointActions:     []string{"shoulder_adduction", "scapular_depression"},
		Description:      "Pull should be done with spine stacked, avoiding leaning backwards - lat pulldown",
	},

	// Shoulder Specific
	ShoulderFlexion: {
		Pattern:          ShoulderFlexion,
		Name:             "Shoulder Flexion",
		Plane:            Sagittal,
		PrimaryMuscles:   []MuscleGroup{Shoulders, Chest},
		SecondaryMuscles: []MuscleGroup{},
		JointActions:     []string{"shoulder_flexion"},
		Description:      "Raising arm forward and upward - front raises, converging upper press",
	},
	ShoulderExtension: {
		Pattern:          ShoulderExtension,
		Name:             "Shoulder Extension",
		Plane:            Sagittal,
		PrimaryMuscles:   []MuscleGroup{Lats, Shoulders},
		SecondaryMuscles: []MuscleGroup{Chest, Triceps},
		JointActions:     []string{"shoulder_extension"},
		Description:      "Moving arm backward from flexed position - pullovers, high to low rows",
	},
	ShoulderAbduction: {
		Pattern:          ShoulderAbduction,
		Name:             "Shoulder Abduction",
		Plane:            Frontal,
		PrimaryMuscles:   []MuscleGroup{Shoulders},
		SecondaryMuscles: []MuscleGroup{Traps},
		JointActions:     []string{"shoulder_abduction"},
		Description:      "Raising arm away from body to the side - lateral raises",
	},
	ShoulderAdduction: {
		Pattern:          ShoulderAdduction,
		Name:             "Shoulder Adduction",
		Plane:            Frontal,
		PrimaryMuscles:   []MuscleGroup{Lats},
		SecondaryMuscles: []MuscleGroup{Chest},
		JointActions:     []string{"shoulder_adduction"},
		Description:      "Bringing arm toward body from abducted position - cable crossovers, lat work",
	},

	// Hip and Knee Patterns
	HipHinge: {
		Pattern:          HipHinge,
		Name:             "Hip Hinge",
		Plane:            Sagittal,
		PrimaryMuscles:   []MuscleGroup{Hamstrings, Glutes},
		SecondaryMuscles: []MuscleGroup{Adductors, Erectors},
		JointActions:     []string{"hip_extension", "spinal_extension"},
		Description:      "Hip-dominant extension pattern with minimal knee bend - glute bridges, RDLs",
	},
	HipAdduction: {
		Pattern:        HipAdduction,
		Name:           "Hip Adduction",
		Plane:          Frontal,
		PrimaryMuscles: []MuscleGroup{Adductors},
		JointActions:   []string{"hip_adduction"},
		Description:    "Bringing legs together against resistance - adductor machine, cable adductions",
	},
	SquatPattern: {
		Pattern:        SquatPattern,
		Name:           "Squat Pattern",
		Plane:          Sagittal,
		PrimaryMuscles: []MuscleGroup{Quads, Glutes, Adductors},
		JointActions:   []string{"hip_extension", "knee_extension"},
		Description:    "Combined hip and knee extension - squats, leg press, lunges",
	},
	KneeFixedExtension: {
		Pattern:        KneeFixedExtension,
		Name:           "Knee Extension",
		Plane:          Sagittal,
		PrimaryMuscles: []MuscleGroup{Quads},
		JointActions:   []string{"knee_extension"},
		Description:    "Isolated knee extension with fixed hip - leg extensions",
	},

	// Elbow Patterns
	ElbowFlexion: {
		Pattern:        ElbowFlexion,
		Name:           "Elbow Flexion",
		Plane:          Sagittal,
		PrimaryMuscles: []MuscleGroup{Biceps},
		JointActions:   []string{"elbow_flexion"},
		Description:    "Bending elbow to bring hand toward shoulder - curls",
	},
	ElbowExtension: {
		Pattern:        ElbowExtension,
		Name:           "Elbow Extension",
		Plane:          Sagittal,
		PrimaryMuscles: []MuscleGroup{Triceps},
		JointActions:   []string{"elbow_extension"},
		Description:    "Straightening elbow against resistance - tricep extensions, pushdowns",
	},

	// Torso Patterns
	SpinalFlexion: {
		Pattern:        SpinalFlexion,
		Name:           "Spinal Flexion",
		Plane:          Sagittal,
		PrimaryMuscles: []MuscleGroup{Abs},
		JointActions:   []string{"spinal_flexion"},
		Description:    "Flexing spine to bring ribcage toward pelvis - crunches, sit-ups",
	},
	SpinalExtension: {
		Pattern:        SpinalExtension,
		Name:           "Spinal Extension",
		Plane:          Sagittal,
		PrimaryMuscles: []MuscleGroup{Erectors},
		JointActions:   []string{"spinal_extension"},
		Description:    "Extending spine against resistance - back extensions, reverse hypers",
	},
}
