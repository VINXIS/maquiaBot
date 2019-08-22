package osutools

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"net/http"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/thehowl/go-osuapi"

	structs "../structs"
	tools "../tools"
)

// FarmUpdate gets the new data from grumd's site: https://grumd.github.io/osu-pps
func FarmUpdate(s *discordgo.Session) {
	// Obtain farm data
	farmData := structs.FarmData{}
	f, err := ioutil.ReadFile("./data/osuData/mapFarm.json")
	tools.ErrRead(err)
	_ = json.Unmarshal(f, &farmData)
	if time.Since(farmData.Time) < 24*time.Hour {
		return
	}

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

	jsonCache, err := json.Marshal(structs.FarmData{
		Time: time.Now(),
		Maps: data,
	})
	tools.ErrRead(err)

	err = ioutil.WriteFile("./data/osuData/mapFarm.json", jsonCache, 0644)
	tools.ErrRead(err)

	osuAPI := osuapi.NewClient(os.Getenv("OSU_API"))
	// Obtain profile cache data
	profileCache := []structs.PlayerData{}
	f, err = ioutil.ReadFile("./data/osuData/profileCache.json")
	tools.ErrRead(err)
	_ = json.Unmarshal(f, &profileCache)

	fmt.Println("Saved data! Updating all " + strconv.Itoa(len(profileCache)) + " players...")

	for i, player := range profileCache {
		if player.Osu.Username != "" {
			FarmRating := 0.00
			FarmValues := make(map[float64]string)
			list := ""
			scoreList, _ := osuAPI.GetUserBest(osuapi.GetUserScoresOpts{
				Username: player.Osu.Username,
				Limit:    100,
			})

			for j, score := range scoreList {
				var HDVer osuapi.Mods
				var weight float64
				var artist string
				var title string
				var diffName string

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
						weight = math.Max(weight, math.Pow(0.95, float64(j))*farmMap.Overweightness)
						artist = farmMap.Artist
						title = farmMap.Title
						diffName = farmMap.DiffName
					}
				}

				if weight != 0 {
					FarmRating = FarmRating + weight
					FarmValues[weight] = "`" + artist + " - " + title + " [" + diffName + "]`: " + strconv.FormatFloat(weight, 'f', 2, 64) + " Farmerdog rating (" + strconv.FormatFloat(score.PP, 'f', 2, 64) + "pp) \n"
				}
			}

			keys := []float64{}
			for FarmVal := range FarmValues {
				keys = append(keys, FarmVal)
			}
			sort.Slice(keys, func(i, j int) bool {
				return keys[i] > keys[j]
			})

			for _, FarmVal := range keys {
				lines := strings.Split(list, "\n")
				if len(lines) > 5 {
					break
				}

				list = list + FarmValues[FarmVal]
			}

			player.Farm.Rating = FarmRating
			player.Farm.Time = time.Now()
			player.Farm.List = list
			profileCache[i] = player
			fmt.Println("Updated player #" + strconv.Itoa(i+1) + ": " + player.Osu.Username + " Farm Rating " + fmt.Sprint(FarmRating))
		}
	}

	jsonCache, err = json.Marshal(profileCache)
	tools.ErrRead(err)

	err = ioutil.WriteFile("./data/osuData/profileCache.json", jsonCache, 0644)
	tools.ErrRead(err)
	fmt.Println("Updated all players!")

	// Loop everyday
	ticker := time.NewTicker(24 * time.Hour)
	for {
		select {
		case <-ticker.C:
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

			jsonCache, err := json.Marshal(structs.FarmData{
				Time: time.Now(),
				Maps: data,
			})
			tools.ErrRead(err)

			err = ioutil.WriteFile("./data/osuData/mapFarm.json", jsonCache, 0644)
			tools.ErrRead(err)

			osuAPI := osuapi.NewClient(os.Getenv("OSU_API"))
			// Obtain profile cache data
			profileCache := []structs.PlayerData{}
			f, err = ioutil.ReadFile("./data/osuData/profileCache.json")
			tools.ErrRead(err)
			_ = json.Unmarshal(f, &profileCache)

			fmt.Println("Saved data! Updating all " + strconv.Itoa(len(profileCache)) + " players...")

			for i, player := range profileCache {
				if player.Osu.Username != "" {
					FarmRating := 0.00
					FarmValues := make(map[float64]string)
					list := ""
					scoreList, _ := osuAPI.GetUserBest(osuapi.GetUserScoresOpts{
						Username: player.Osu.Username,
						Limit:    100,
					})

					for j, score := range scoreList {
						if strings.Contains(score.Mods.String(), "NC") {
							stringMods := strings.Replace(score.Mods.String(), "NC", "", 1)
							score.Mods = osuapi.ParseMods(stringMods)
						}
						if strings.Contains(score.Mods.String(), "HD") {
							stringMods := strings.Replace(score.Mods.String(), "HD", "", 1)
							score.Mods = osuapi.ParseMods(stringMods)
						}
						for _, farmMap := range data {
							if score.BeatmapID == farmMap.BeatmapID && score.Mods == farmMap.Mods {
								FarmRating = FarmRating + math.Pow(0.95, float64(j))*farmMap.Overweightness
								FarmValues[math.Pow(0.95, float64(j))*farmMap.Overweightness] = "`" + farmMap.Artist + " - " + farmMap.Title + " [" + farmMap.DiffName + "]`: " + strconv.FormatFloat(math.Pow(0.95, float64(j))*farmMap.Overweightness, 'f', 2, 64) + " Farmerdog rating (" + strconv.FormatFloat(score.PP, 'f', 2, 64) + "pp) \n"
								break
							}
						}
					}
					keys := []float64{}
					for FarmVal := range FarmValues {
						keys = append(keys, FarmVal)
					}
					sort.Slice(keys, func(i, j int) bool {
						return keys[i] > keys[j]
					})

					for _, FarmVal := range keys {
						lines := strings.Split(list, "\n")
						if len(lines) > 5 {
							break
						}

						list = list + FarmValues[FarmVal]
					}

					player.Farm.Rating = FarmRating
					player.Farm.Time = time.Now()
					player.Farm.List = list
					profileCache[i] = player
					fmt.Println("Updated player #" + strconv.Itoa(i+1) + ": " + player.Osu.Username + " Farm Rating " + fmt.Sprint(FarmRating))
				}
			}

			jsonCache, err = json.Marshal(profileCache)
			tools.ErrRead(err)

			err = ioutil.WriteFile("./data/osuData/profileCache.json", jsonCache, 0644)
			tools.ErrRead(err)
			fmt.Println("Updated all players!")
		}
	}
}
