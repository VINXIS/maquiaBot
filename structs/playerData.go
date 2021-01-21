package structs

import (
	"math"
	"strings"
	"time"

	osuapi "maquiaBot/osu-api"
)

// PlayerData stores information regarding the discord user, and the osu user
type PlayerData struct {
	Time     time.Time
	Discord  string
	Osu      osuapi.User
	Farm     FarmerdogData
	Currency CurrencyData
}

// FarmerdogData is how much of a farmerdog the player is
type FarmerdogData struct {
	Rating float64
	List   []PlayerScore
}

// PlayerScore is the score by the player, it tells you how farmy the score is as well
type PlayerScore struct {
	BeatmapSet int
	PP         float64
	FarmScore  float64
	Name       string
}

// CurrencyData is the amount of in-bot currency the user has
type CurrencyData struct {
	Amount    float64
	LastDaily time.Time
}

// FarmCalc does the actual calculations of the farm values and everything for the player
func (player *PlayerData) FarmCalc(osuAPI *osuapi.Client, farmData FarmData) {
	player.Farm = FarmerdogData{}

	scoreList, err := osuAPI.GetUserBest(osuapi.GetUserScoresOpts{
		Username: player.Osu.Username,
		Limit:    100,
	})
	if err != nil {
		return
	}

	for j, score := range scoreList {
		var HDVer osuapi.Mods
		var playerFarmScore = PlayerScore{}

		// Remove NC
		if strings.Contains(score.Mods.String(), "NC") {
			stringMods := strings.Replace(score.Mods.String(), "NC", "", 1)
			score.Mods = osuapi.ParseMods(stringMods)
		}

		// Treat HD and no HD the same
		if strings.Contains(score.Mods.String(), "HD") {
			HDVer = score.Mods
			stringMods := strings.Replace(score.Mods.String(), "HD", "", 1)
			score.Mods = osuapi.ParseMods(stringMods)
		} else {
			stringMods := score.Mods.String() + "HD"
			HDVer = osuapi.ParseMods(stringMods)
		}

		// Actual farm calc for the map
		for _, farmMap := range farmData.Maps {
			if score.BeatmapID == farmMap.BeatmapID && (score.Mods == farmMap.Mods || HDVer == farmMap.Mods) {
				playerFarmScore.BeatmapSet = score.BeatmapID
				playerFarmScore.PP = score.PP
				playerFarmScore.FarmScore = math.Max(playerFarmScore.FarmScore, math.Pow(0.95, float64(j))*farmMap.Overweightness)
				playerFarmScore.Name = farmMap.Artist + " - " + farmMap.Title + " [" + farmMap.DiffName + "]"
			}
		}

		if playerFarmScore.BeatmapSet != 0 {
			player.Farm.List = append(player.Farm.List, playerFarmScore)
			player.Farm.Rating += playerFarmScore.FarmScore
		}
	}
}
