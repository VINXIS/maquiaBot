package osutools

import (
	"sort"
	"strconv"

	osuapi "../osu-api"
)

// BeatmapCalc calculates PP for the map
// 300 100 and 50 calculations are based off of https://github.com/ppy/osu-tools/blob/master/PerformanceCalculator/Simulate/OsuSimulateCommand.cs#L76
func BeatmapCalc(mods, acc, combo, misses string, beatmap osuapi.Beatmap) (values []string) {
	if beatmap.Mode != osuapi.ModeCatchTheBeat && acc != "N/A" {
		totalHits := beatmap.Circles + beatmap.Sliders + beatmap.Spinners
		missCount, _ := strconv.Atoi(misses)
		accVal, err := strconv.ParseFloat(acc, 64)

		if err != nil || accVal <= 0 || accVal >= 100 {
			ppValues := make(chan string, 5)
			var ppValueArray [5]float64

			go PPCalc(beatmap, osuapi.Score{
				MaxCombo: beatmap.MaxCombo,
				Count300: totalHits,
				Mods:     osuapi.ParseMods(mods),
			}, ppValues)
			go PPCalc(beatmap, osuapi.Score{
				MaxCombo: beatmap.MaxCombo,
				Count300: (int(0.99*float64(totalHits)*6) - totalHits + missCount) / 5,
				Count100: (int(0.99*float64(totalHits)*6) - totalHits + missCount) % 5,
				Count50:  totalHits - (int(0.99*float64(totalHits)*6)-totalHits+missCount)/5 - (int(0.99*float64(totalHits)*6)-totalHits+missCount)%5 - missCount,
				Mods:     osuapi.ParseMods(mods),
			}, ppValues)
			go PPCalc(beatmap, osuapi.Score{
				MaxCombo: beatmap.MaxCombo,
				Count300: (int(0.98*float64(totalHits)*6) - totalHits + missCount) / 5,
				Count100: (int(0.98*float64(totalHits)*6) - totalHits + missCount) % 5,
				Count50:  totalHits - (int(0.98*float64(totalHits)*6)-totalHits+missCount)/5 - (int(0.98*float64(totalHits)*6)-totalHits+missCount)%5 - missCount,
				Mods:     osuapi.ParseMods(mods),
			}, ppValues)
			go PPCalc(beatmap, osuapi.Score{
				MaxCombo: beatmap.MaxCombo,
				Count300: (int(0.97*float64(totalHits)*6) - totalHits + missCount) / 5,
				Count100: (int(0.97*float64(totalHits)*6) - totalHits + missCount) % 5,
				Count50:  totalHits - (int(0.97*float64(totalHits)*6)-totalHits+missCount)/5 - (int(0.97*float64(totalHits)*6)-totalHits+missCount)%5 - missCount,
				Mods:     osuapi.ParseMods(mods),
			}, ppValues)
			go PPCalc(beatmap, osuapi.Score{
				MaxCombo: beatmap.MaxCombo,
				Count300: (int(0.95*float64(totalHits)*6) - totalHits + missCount) / 5,
				Count100: (int(0.95*float64(totalHits)*6) - totalHits + missCount) % 5,
				Count50:  totalHits - (int(0.95*float64(totalHits)*6)-totalHits+missCount)/5 - (int(0.95*float64(totalHits)*6)-totalHits+missCount)%5 - missCount,
				Mods:     osuapi.ParseMods(mods),
			}, ppValues)

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

			go PPCalc(beatmap, osuapi.Score{
				MaxCombo: beatmap.MaxCombo,
				Count300: (int(accVal*float64(totalHits)*6) - totalHits + missCount) / 5,
				Count100: (int(accVal*float64(totalHits)*6) - totalHits + missCount) % 5,
				Count50:  totalHits - (int(accVal*float64(totalHits)*6)-totalHits+missCount)/5 - (int(accVal*float64(totalHits)*6)-totalHits+missCount)%5 - missCount,
				Mods:     osuapi.ParseMods(mods),
			}, ppValues)

			ppNum, _ := strconv.ParseFloat(<-ppValues, 64)
			ppVal := "**" + acc + "%:** " + strconv.FormatFloat(ppNum, 'f', 0, 64) + "pp"
			values = []string{ppVal}
		}
	}
	return values
}
