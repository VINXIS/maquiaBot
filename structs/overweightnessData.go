package structs

import (
	"encoding/json"
	"time"

	osuapi "../osu-api"
)

// RawData is the raw data obtained from https://raw.githubusercontent.com/grumd/osu-pps/master/data.json
type RawData struct {
	BeatmapID    int             `json:"b"`
	BeatmapSetID int             `json:"s"`
	X            float64         `json:"x"`
	PP99         float64         `json:"pp99"`
	Adj          float64         `json:"adj"`
	Artist       json.Number     `json:"art"`
	Title        json.Number     `json:"t"`
	DiffName     json.Number     `json:"v"`
	HitLength    int             `json:"l"`
	BPM          float64         `json:"bpm"`
	SR           float64         `json:"d"`
	Passcount    int             `json:"p"`
	Age          int             `json:"h"`
	Genre        osuapi.Genre    `json:"g"`
	Language     osuapi.Language `json:"ln"`
	Mods         osuapi.Mods     `json:"m"`
}

// FarmData is the farm data of all maps
type FarmData struct {
	Time time.Time
	Maps []MapFarm
}

// MapFarm holds the farm data for each map
type MapFarm struct {
	BeatmapID      int
	Artist         string
	Title          string
	DiffName       string
	Overweightness float64
	Mods           osuapi.Mods
}
