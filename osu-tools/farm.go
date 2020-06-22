package osutools

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"strconv"
	"time"

	structs "maquiaBot/structs"
	tools "maquiaBot/tools"
	"github.com/bwmarrin/discordgo"
)

// FarmUpdate gets the new data from grumd's site: https://grumd.github.io/osu-pps
func FarmUpdate(s *discordgo.Session) {
	// Obtain farm data
	farmData := structs.FarmData{}
	f, err := ioutil.ReadFile("./data/osuData/mapFarm.json")
	tools.ErrRead(s, err)
	_ = json.Unmarshal(f, &farmData)
	if time.Since(farmData.Time) > 24*time.Hour {
		UpdateFarmSystem(s)
	}

	// Loop everyday
	ticker := time.NewTicker(24 * time.Hour)
	for {
		select {
		case <-ticker.C:
			UpdateFarmSystem(s)
		}
	}
}

// UpdateFarmSystem updates the whole farm data
func UpdateFarmSystem(s *discordgo.Session) {
	log.Println("Fetching data as more than 24 hours have passed...")

	// Obtain data
	var byteArray []byte
	for {
		res, err := http.Get("https://raw.githubusercontent.com/grumd/osu-pps/master/data-osu.json")
		byteArray, err = ioutil.ReadAll(res.Body)
		if err == nil {
			break
		}
	}

	// Convert to readable data
	info := []structs.RawData{}
	err := json.Unmarshal(byteArray, &info)
	tools.ErrRead(s, err)

	log.Println("Obtained data! Now parsing...")

	// grumd's Overweightness formula implementation
	data := []structs.MapFarm{}
	max := 0.00
	for _, raw := range info {
		ow := raw.X / math.Pow(raw.Adj, 0.65) / math.Pow(float64(raw.Passcount), 0.2) / math.Pow(float64(raw.Age), 0.5)
		if max < ow && !math.IsNaN(ow) && !math.IsInf(ow, 1) && !math.IsInf(ow, -1) {
			max = ow
		}
	}

	for _, raw := range info {
		ow := raw.X / math.Pow(raw.Adj, 0.65) / math.Pow(float64(raw.Passcount), 0.2) / math.Pow(float64(raw.Age), 0.5)
		if !math.IsNaN(ow) && !math.IsInf(ow, 1) && !math.IsInf(ow, -1) {
			data = append(data, structs.MapFarm{
				BeatmapID:      raw.BeatmapID,
				Artist:         string(raw.Artist),
				Title:          string(raw.Title),
				DiffName:       string(raw.DiffName),
				Overweightness: ow / max * 800.0,
				Mods:           raw.Mods,
			})
		}
	}

	farmData := structs.FarmData{
		Time: time.Now(),
		Maps: data,
	}
	jsonCache, err := json.Marshal(farmData)
	tools.ErrRead(s, err)

	// Save map farm data
	err = ioutil.WriteFile("./data/osuData/mapFarm.json", jsonCache, 0644)
	tools.ErrRead(s, err)

	// Obtain profile cache data
	profileCache := []structs.PlayerData{}
	f, err := ioutil.ReadFile("./data/osuData/profileCache.json")
	tools.ErrRead(s, err)
	_ = json.Unmarshal(f, &profileCache)

	log.Println("Saved data! Updating all " + strconv.Itoa(len(profileCache)) + " players...")

	for i, player := range profileCache {
		if player.Osu.Username != "" {
			player.FarmCalc(OsuAPI, farmData)
			profileCache[i] = player
			log.Println("Updated player #" + strconv.Itoa(i+1) + ": " + player.Osu.Username + " Farm Rating " + fmt.Sprint(player.Farm.Rating))
		}
	}

	jsonCache, err = json.Marshal(profileCache)
	tools.ErrRead(s, err)

	err = ioutil.WriteFile("./data/osuData/profileCache.json", jsonCache, 0644)
	tools.ErrRead(s, err)
	log.Println("Updated all players!")
}
