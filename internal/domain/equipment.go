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
	StackWeights []float64
	PulleyType   string

	// Machine-specific
	ResistanceProfileID   string
	ResistanceProfileName string
}
