package domain

type MovementPlane string

const (
	Sagittal   MovementPlane = "sagittal"   // Forward/backward
	Frontal    MovementPlane = "frontal"    // Side to side
	Transverse MovementPlane = "transverse" // Rotational
)
