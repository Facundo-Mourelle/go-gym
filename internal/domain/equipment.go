package domain

type EquipmentType string

const (
	EquipmentTypeFreeWeight EquipmentType = "freeweight"
	EquipmentTypeMachine    EquipmentType = "machine"
	EquipmentTypeCable      EquipmentType = "cable"
)

type EquipmentID string

type Equipment struct {
	ID                  EquipmentID
	Type                EquipmentType
	Name                string
	Manufacturer        string
	UserID              string
	ActualWeight        float64

	// Cable-specific fields
	StackWeights      []float64
	PulleyType        string
	WeightIncrement   float64 // weight increment for cable stack (1.25, 2.5, 4.5)

	// Machine-specific
	ResistanceProfileID   string
	ResistanceProfileName string
	MovementPattern       MovementPattern // which movement pattern this machine is for
}
