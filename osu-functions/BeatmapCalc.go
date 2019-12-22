package osutools

import (
	"sort"
	"strconv"

	osuapi "../osu-api"
)

// BeatmapCalc calculates PP for the map
func BeatmapCalc(mods, acc, combo, misses string, beatmap osuapi.Beatmap) (values []string) {
	if beatmap.Mode != osuapi.ModeCatchTheBeat && acc != "N/A" {
		accVal, err := strconv.ParseFloat(acc, 64)
		if err != nil || accVal <= 0 || accVal >= 100 {
			ppValues := make(chan string, 5)
			var ppValueArray [5]float64
			go PPCalc(beatmap, 100.0, combo, misses, mods, ppValues)
			go PPCalc(beatmap, 99.0, combo, misses, mods, ppValues)
			go PPCalc(beatmap, 98.0, combo, misses, mods, ppValues)
			go PPCalc(beatmap, 97.0, combo, misses, mods, ppValues)
			go PPCalc(beatmap, 95.0, combo, misses, mods, ppValues)
			for v := 0; v < 5; v++ {
				ppValueArray[v], _ = strconv.ParseFloat(<-ppValues, 64)
			}
			sort.Slice(ppValueArray[:], func(i, j int) bool {
				return ppValueArray[i] > ppValueArray[j]
			})
			ppSS := "**100%:** " + strconv.FormatFloat(ppValueArray[0], 'f', 0, 64) + "pp | "
			pp99 := "**99%:** " + strconv.FormatFloat(ppValueArray[1], 'f', 0, 64) + "pp | "
			pp98 := "**98%:** " + strconv.FormatFloat(ppValueArray[2], 'f', 0, 64) + "pp | "
			pp97 := "**97%:** " + strconv.FormatFloat(ppValueArray[3], 'f', 0, 64) + "pp | "
			pp95 := "**95%:** " + strconv.FormatFloat(ppValueArray[4], 'f', 0, 64) + "pp"
			values = []string{ppSS, pp99, pp98, pp97, pp95}
		} else {
			ppValues := make(chan string, 1)
			go PPCalc(beatmap, accVal, combo, misses, mods, ppValues)
			ppNum, _ := strconv.ParseFloat(<-ppValues, 64)
			ppVal := "**" + acc + "%:** " + strconv.FormatFloat(ppNum, 'f', 0, 64) + "pp"
			values = []string{ppVal}
		}
	}
	return values
}
