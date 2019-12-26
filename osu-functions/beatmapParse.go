package osutools

import (
	"math"
	"regexp"
	"sort"
	"strconv"

	osuapi "../osu-api"
	tools "../tools"
)

// BeatmapParse parses beatmap and obtains the .osu file
func BeatmapParse(id, format string, mods *osuapi.Mods) (beatmap osuapi.Beatmap) {
	replacer, _ := regexp.Compile(`[^a-zA-Z0-9\s\(\)]`)

	mapID, err := strconv.Atoi(id)
	tools.ErrRead(err)

	if format == "map" {
		// Fetch the beatmap
		beatmaps, err := OsuAPI.GetBeatmaps(osuapi.GetBeatmapsOpts{
			BeatmapID: mapID,
			Mods:      mods,
		})
		tools.ErrRead(err)
		if len(beatmaps) > 0 {
			beatmap = beatmaps[0]
		}

		// Download the .osu file for the map
		tools.DownloadFile(
			"./data/osuFiles/"+
				strconv.Itoa(beatmap.BeatmapID)+
				" "+
				replacer.ReplaceAllString(beatmap.Artist, "")+
				" - "+
				replacer.ReplaceAllString(beatmap.Title, "")+
				".osu",
			"https://osu.ppy.sh/osu/"+
				strconv.Itoa(beatmap.BeatmapID))
	} else if format == "set" {
		// Fetch the beatmap
		beatmaps, err := OsuAPI.GetBeatmaps(osuapi.GetBeatmapsOpts{
			BeatmapSetID: mapID,
			Mods:         mods,
		})
		tools.ErrRead(err)

		// Reorder the maps so that it returns the highest difficulty in the set
		sort.Slice(beatmaps, func(i, j int) bool {
			return beatmaps[i].DifficultyRating > beatmaps[j].DifficultyRating
		})

		// Download the .osu files for the set
		for _, diff := range beatmaps {
			tools.DownloadFile(
				"./data/osuFiles/"+
					strconv.Itoa(diff.BeatmapID)+
					" "+
					replacer.ReplaceAllString(diff.Artist, "")+
					" - "+
					replacer.ReplaceAllString(diff.Title, "")+
					".osu",
				"https://osu.ppy.sh/osu/"+
					strconv.Itoa(diff.BeatmapID))
		}
		if len(beatmaps) > 0 {
			beatmap = beatmaps[0]
		}
	}

	// Mod scaling
	scaleMods := *mods

	// HR / EZ scaling
	if scaleMods&osuapi.ModHardRock != 0 {
		beatmap.CircleSize = math.Min(10.0, beatmap.CircleSize*1.3)
		beatmap.ApproachRate = math.Min(10.0, beatmap.ApproachRate*1.4)
		beatmap.OverallDifficulty = math.Min(10.0, beatmap.OverallDifficulty*1.4)
		beatmap.HPDrain = math.Min(10.0, beatmap.HPDrain*1.4)
	} else if scaleMods&osuapi.ModEasy != 0 {
		beatmap.CircleSize /= 2.0
		beatmap.ApproachRate /= 2.0
		beatmap.OverallDifficulty /= 2.0
		beatmap.HPDrain /= 2.0
	}

	// DT / HT scaling
	clock := float64(1)
	if scaleMods&osuapi.ModDoubleTime != 0 {
		clock = 1.5
	} else if scaleMods&osuapi.ModHalfTime != 0 {
		clock = 0.75
	}

	beatmap.BPM *= clock
	beatmap.TotalLength = int(float64(beatmap.TotalLength) / clock)
	beatmap.HitLength = int(float64(beatmap.HitLength) / clock)
	ARMS := diffRange(beatmap.ApproachRate) / clock
	ODScale := (80.0 - 6.0*beatmap.OverallDifficulty) / clock
	HPMS := diffRange(beatmap.HPDrain) / clock
	beatmap.OverallDifficulty = (80.0 - ODScale) / 6.0
	beatmap.ApproachRate = diffValue(ARMS)
	beatmap.HPDrain = diffValue(HPMS)

	return beatmap
}

func diffRange(value float64) float64 {
	if value > 5.0 {
		return 1200 + (450-1200)*(value-5)/5
	} else if value < 5.0 {
		return 1200 - (1200-1800)*(5-value)/5
	}
	return 1200
}

func diffValue(value float64) float64 {
	if value > 1200 {
		return (1800 - value) / 120
	}
	return (1200-value)/150 + 5
}
