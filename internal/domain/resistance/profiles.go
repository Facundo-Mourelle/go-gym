package resistance

type ResistanceProfileType string
type PulleyConfiguration string

// We distinguish pulleys from uniform resistance profiles to
// simplify effective load calculations and easily
// create new pulley profiles
const (
	Pulley1to1 PulleyConfiguration = "1:1"
	Pulley2to1 PulleyConfiguration = "2:1"

	Ascending  ResistanceProfileType = "ascending"
	Descending ResistanceProfileType = "descending"
	Bell       ResistanceProfileType = "bell"
	Uniform    ResistanceProfileType = "uniform"
)

type ResistanceProfile struct {
	profileType ResistanceProfileType
	pulley      *PulleyProfile
	// Normalized position in movement:
	// 0.0 = start, 1.0 = end
	position float64
}
