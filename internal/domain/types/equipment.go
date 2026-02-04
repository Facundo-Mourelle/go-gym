package types

type EquipmentType string

const (
	EquipmentTypeFreeWeight EquipmentType = "freeweight"
	EquipmentTypeMachine    EquipmentType = "machine"
)

type EquipmentID string

type Equipment struct {
	ID                  EquipmentID
	Type                EquipmentType
	Name                string
	Manufacturer        string
	ResistanceProfileID string // References a resistance profile
}
