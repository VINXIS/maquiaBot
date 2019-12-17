package osutools

import (
	"regexp"
	"sort"
	"strconv"

	osuapi "../osu-api"
	tools "../tools"
)

// BeatmapParse parses beatmap and obtains the .osu file
func BeatmapParse(id, format string, mods osuapi.Mods) (beatmap osuapi.Beatmap) {
	replacer, _ := regexp.Compile(`[^a-zA-Z0-9\s\(\)]`)

	mapID, err := strconv.Atoi(id)
	tools.ErrRead(err)

	if format == "map" {
		// Fetch the beatmap
		beatmaps, err := OsuAPI.GetBeatmaps(osuapi.GetBeatmapsOpts{
			BeatmapID: mapID,
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
	return beatmap
}
