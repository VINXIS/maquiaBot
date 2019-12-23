package structs

import (
	"time"
)

// PenisRecordData stores information regarding the largest and smallest penis sizes
type PenisRecordData struct {
	Largest  PenisData
	Smallest PenisData
}

// PenisData stores information regarding the largest and smallest penis sizes
type PenisData struct {
	Size   float64
	UserID string
	Date   time.Time
}
