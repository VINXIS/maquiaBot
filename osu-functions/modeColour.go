package osutools

import (
	osuapi "../osu-api"
)

// ModeColour assigns a colour based on the beatmap's mode
func ModeColour(mode osuapi.Mode) (Colour int) {
	switch mode {
	case osuapi.ModeOsu:
		Colour = 0xD65288
	case osuapi.ModeTaiko:
		Colour = 0xFF0000
	case osuapi.ModeCatchTheBeat:
		Colour = 0x007419
	case osuapi.ModeOsuMania:
		Colour = 0xff6200
	}
	return Colour
}
