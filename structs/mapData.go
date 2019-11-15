package structs

import (
	"time"

	osuapi "../osu-api"
)

// MapData stores the beatmap information as well as pp for different mods
type MapData struct {
	Time    time.Time
	Beatmap osuapi.Beatmap
	PP      []PPData
}

// PPData stores the pp for the specific mod
type PPData struct {
	Mods string
	SR   float64
	PPSS float64
	PP99 float64
	PP98 float64
	PP97 float64
	PP95 float64
}
