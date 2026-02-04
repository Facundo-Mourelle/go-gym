package types

type MovementPattern string

const (
	// Horizontal Patterns
	HorizontalPush MovementPattern = "horizontal_push" // Bench press
	HorizontalPull MovementPattern = "horizontal_pull" // Rows

	// Vertical Patterns
	VerticalPush MovementPattern = "vertical_push" // Overhead press
	VerticalPull MovementPattern = "vertical_pull" // lat pulldowns

	// Shoulder Specific
	ShoulderFlexion   MovementPattern = "shoulder_flexion"   // Front raises, pullover
	ShoulderExtension MovementPattern = "shoulder_extension" // Lat pulldowns, pullovers
	ShoulderAbduction MovementPattern = "shoulder_abduction" // Lateral raises
	ShoulderAdduction MovementPattern = "shoulder_adduction" // Cable crossovers, lat work

	// Hip and knee Patterns
	HipHinge           MovementPattern = "hip_hinge"      // Deadlifts, RDLs
	HipAdduction       MovementPattern = "hip_adduction"  // Adductor machine
	SquatPattern       MovementPattern = "squat_pattern"  // Squats, leg press
	KneeFixedExtension MovementPattern = "knee_extension" // Knee extension

	// Elbow Patterns
	ElbowFlexion   MovementPattern = "elbow_flexion"   // Curls
	ElbowExtension MovementPattern = "elbow_extension" // Tricep extensions

	// Torso Patterns
	SpinalFlexion   MovementPattern = "spinal_flexion"   // Crunches
	SpinalExtension MovementPattern = "spinal_extension" // Back extensions
)

// MovementPatternInfo provides metadata about a pattern
type MovementPatternInfo struct {
	Pattern          MovementPattern
	Name             string
	Plane            MovementPlane
	PrimaryMuscles   []MuscleGroup
	SecondaryMuscles []MuscleGroup
	JointActions     []string
	Description      string
}
