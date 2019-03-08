package structs

import (
	"time"

	"github.com/thehowl/go-osuapi"
)

// PPData stores the pp for the specific mod
type PPData struct {
	Mods string
	SR   string
	PPSS string
	PP99 string
	PP98 string
	PP97 string
	PP95 string
}

// MapData stores the beatmap information as well as pp for different mods
type MapData struct {
	Time    time.Time
	Beatmap osuapi.Beatmap
	PP      []PPData
}
