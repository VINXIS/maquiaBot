package osutools

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	osuapi "../osu-api"
	structs "../structs"
	tools "../tools"
)

// FarmUpdate gets the new data from grumd's site: https://grumd.github.io/osu-pps
func FarmUpdate() {
	// Obtain farm data
	farmData := structs.FarmData{}
	f, err := ioutil.ReadFile("./data/osuData/mapFarm.json")
	tools.ErrRead(err)
	_ = json.Unmarshal(f, &farmData)
	if time.Since(farmData.Time) > 24*time.Hour {
		updateSystem()
	}

	// Loop everyday
	ticker := time.NewTicker(24 * time.Hour)
	for {
		select {
		case <-ticker.C:
			updateSystem()
		}
	}
}

func updateSystem() {
	fmt.Println("Fetching data as more than 24 hours have passed...")

	// Obtain data
	res, err := http.Get("https://raw.githubusercontent.com/grumd/osu-pps/master/data-osu.json")
	tools.ErrRead(err)

	byteArray, err := ioutil.ReadAll(res.Body)
	tools.ErrRead(err)

	// Convert to readable data
	info := []structs.RawData{}
	err = json.Unmarshal(byteArray, &info)
	tools.ErrRead(err)

	fmt.Println("Obtained data! Now parsing...")

	// grumd's Overweightness formula implementation
	data := []structs.MapFarm{}
	max := 0.00
	for _, raw := range info {
		ow := raw.X / math.Pow(raw.Adj, 0.65) / math.Pow(float64(raw.Passcount), 0.2) / math.Pow(float64(raw.Age), 0.5)
		if max < ow {
			max = ow
		}
	}

	for _, raw := range info {
		ow := raw.X / math.Pow(raw.Adj, 0.65) / math.Pow(float64(raw.Passcount), 0.2) / math.Pow(float64(raw.Age), 0.5)
		data = append(data, structs.MapFarm{
			BeatmapID:      raw.BeatmapID,
			Artist:         string(raw.Artist),
			Title:          string(raw.Title),
			DiffName:       string(raw.DiffName),
			Overweightness: ow / max * 800.0,
			Mods:           raw.Mods,
		})
	}

	farmData := structs.FarmData{
		Time: time.Now(),
		Maps: data,
	}
	jsonCache, err := json.Marshal(farmData)
	tools.ErrRead(err)

	// Save map farm data
	err = ioutil.WriteFile("./data/osuData/mapFarm.json", jsonCache, 0644)
	tools.ErrRead(err)

	osuAPI := osuapi.NewClient(os.Getenv("OSU_API"))

	// Obtain profile cache data
	profileCache := []structs.PlayerData{}
	f, err := ioutil.ReadFile("./data/osuData/profileCache.json")
	tools.ErrRead(err)
	_ = json.Unmarshal(f, &profileCache)

	fmt.Println("Saved data! Updating all " + strconv.Itoa(len(profileCache)) + " players...")

	for i, player := range profileCache {
		if player.Osu.Username != "" {
			player.Farm = structs.FarmerdogData{}

			scoreList, _ := osuAPI.GetUserBest(osuapi.GetUserScoresOpts{
				Username: player.Osu.Username,
				Limit:    100,
			})

			for j, score := range scoreList {
				var HDVer osuapi.Mods
				var playerFarmScore = structs.PlayerScore{}

				if strings.Contains(score.Mods.String(), "NC") {
					stringMods := strings.Replace(score.Mods.String(), "NC", "", 1)
					score.Mods = osuapi.ParseMods(stringMods)
				}
				if strings.Contains(score.Mods.String(), "HD") {
					HDVer = score.Mods
					stringMods := strings.Replace(score.Mods.String(), "HD", "", 1)
					score.Mods = osuapi.ParseMods(stringMods)
				} else {
					stringMods := score.Mods.String() + "HD"
					HDVer = osuapi.ParseMods(stringMods)
				}
				for _, farmMap := range farmData.Maps {
					if score.BeatmapID == farmMap.BeatmapID && (score.Mods == farmMap.Mods || HDVer == farmMap.Mods) {
						playerFarmScore.BeatmapSet = score.BeatmapID
						playerFarmScore.PP = score.PP
						playerFarmScore.FarmScore = math.Max(playerFarmScore.FarmScore, math.Pow(0.95, float64(j))*farmMap.Overweightness)
						playerFarmScore.Name = farmMap.Artist + " - " + farmMap.Title + " [" + farmMap.DiffName + "]"
					}
				}
				player.Farm.List = append(player.Farm.List, playerFarmScore)
				player.Farm.Rating += playerFarmScore.FarmScore
			}

			profileCache[i] = player
			fmt.Println("Updated player #" + strconv.Itoa(i+1) + ": " + player.Osu.Username + " Farm Rating " + fmt.Sprint(player.Farm.Rating))
		}
	}

	jsonCache, err = json.Marshal(profileCache)
	tools.ErrRead(err)

	err = ioutil.WriteFile("./data/osuData/profileCache.json", jsonCache, 0644)
	tools.ErrRead(err)
	fmt.Println("Updated all players!")
}
