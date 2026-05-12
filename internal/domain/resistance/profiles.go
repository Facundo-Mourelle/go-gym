package resistance

import "errors"


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

type Profile interface {
	Name() string
	CalculateEffectiveLoad(rawLoad float64) (float64, error)
}

type Registry struct {
	profiles map[string]Profile
}

func NewRegistry() *Registry {
	return &Registry{
		profiles: make(map[string]Profile),
	}
}

func (r *Registry) Register(p Profile) {
	r.profiles[p.Name()] = p
}

type FreeWeightProfile struct{}

func NewFreeWeightProfile() Profile {
	return &FreeWeightProfile{}
}

func (f *FreeWeightProfile) Name() string { return "free_weight" }

func NewMachineProfile(name string, pulley PulleyConfiguration, advantage, friction float64) Profile {
	return &MachineProfile{name: name}
}

type MachineProfile struct {
	name string
}

func (m *MachineProfile) Name() string { return m.name }



func (r *Registry) Get(name string) (Profile, error) {
	p, ok := r.profiles[name]
	if !ok { return nil, errors.New("not found") }
	return p, nil
}

func (f *FreeWeightProfile) CalculateEffectiveLoad(rawLoad float64) (float64, error) { return rawLoad, nil }
func (m *MachineProfile) CalculateEffectiveLoad(rawLoad float64) (float64, error) { return rawLoad, nil }
