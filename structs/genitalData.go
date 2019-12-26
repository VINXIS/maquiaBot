package structs

import (
	"time"
)

// GenitalRecordData stores information regarding the largest and smallest genital sizes
type GenitalRecordData struct {
	Penis struct {
		Largest  GenitalData
		Smallest GenitalData
	}
	Vagina struct {
		Largest  GenitalData
		Smallest GenitalData
	}
}

// GenitalData stores information regarding the largest and smallest genital sizes
type GenitalData struct {
	Size   float64
	UserID string
	Date   time.Time
}
