package osutools

import (
	"encoding/json"
	"io/ioutil"
	"sort"
	"strconv"
	"time"

	osuapi "../osu-api"
	structs "../structs"
	tools "../tools"
)

// BeatmapCache checks to see if the latest beatmap information is already saved, otherwise it will calculate the SR and PP of the map
func BeatmapCache(mods string, beatmap osuapi.Beatmap, cache []structs.MapData) (starRating, ppSS, pp99, pp98, pp97, pp95 string) {

	latest, update := false, false
	index := 0
	var PPData structs.PPData
	var CacheData structs.MapData
	twoDay := 48 * time.Hour

	for i := range cache {
		if cache[i].Beatmap.BeatmapID == beatmap.BeatmapID {
			if cache[i].Beatmap == beatmap && time.Now().Sub(cache[i].Time) < twoDay {
				for j := range cache[i].PP {
					if cache[i].PP[j].Mods == mods {
						PPData = cache[i].PP[j]
						latest = true
						update = false
						break
					}
				}
				if !latest {
					index = i
					CacheData = cache[i]
					update = true
					break
				}
			} else {
				index = i
				CacheData = cache[i]
				update = true
				break
			}
		}
	}
	if latest {
		starRating = "**SR:** " + strconv.FormatFloat(PPData.SR, 'f', 2, 64) + " "
		ppSS = "**100%:** " + strconv.FormatFloat(PPData.PPSS, 'f', 0, 64) + "pp | "
		pp99 = "**99%:** " + strconv.FormatFloat(PPData.PP99, 'f', 0, 64) + "pp | "
		pp98 = "**98%:** " + strconv.FormatFloat(PPData.PP98, 'f', 0, 64) + "pp | "
		pp97 = "**97%:** " + strconv.FormatFloat(PPData.PP97, 'f', 0, 64) + "pp | "
		pp95 = "**95%:** " + strconv.FormatFloat(PPData.PP95, 'f', 0, 64) + "pp"
	} else {
		if beatmap.Mode != osuapi.ModeCatchTheBeat {
			ppValues := make(chan string, 5)
			var ppValueArray [5]float64
			go PPCalc(beatmap, 100.0, "", "", mods, ppValues)
			go PPCalc(beatmap, 99.0, "", "", mods, ppValues)
			go PPCalc(beatmap, 98.0, "", "", mods, ppValues)
			go PPCalc(beatmap, 97.0, "", "", mods, ppValues)
			go PPCalc(beatmap, 95.0, "", "", mods, ppValues)
			for v := 0; v < 5; v++ {
				ppValueArray[v], _ = strconv.ParseFloat(<-ppValues, 64)
			}
			sort.Slice(ppValueArray[:], func(i, j int) bool {
				return ppValueArray[i] > ppValueArray[j]
			})
			starRating = "**SR:** " + strconv.FormatFloat(beatmap.DifficultyRating, 'f', 2, 64) + " "
			ppSS = "**100%:** " + strconv.FormatFloat(ppValueArray[0], 'f', 0, 64) + "pp | "
			pp99 = "**99%:** " + strconv.FormatFloat(ppValueArray[1], 'f', 0, 64) + "pp | "
			pp98 = "**98%:** " + strconv.FormatFloat(ppValueArray[2], 'f', 0, 64) + "pp | "
			pp97 = "**97%:** " + strconv.FormatFloat(ppValueArray[3], 'f', 0, 64) + "pp | "
			pp95 = "**95%:** " + strconv.FormatFloat(ppValueArray[4], 'f', 0, 64) + "pp"
			if update {
				CacheData.Beatmap = beatmap
				CacheData.Time = time.Now()
				modExist := false
				for j := range CacheData.PP {
					if CacheData.PP[j].Mods == mods {
						modExist = true
						CacheData.PP[j].SR = beatmap.DifficultyRating
						CacheData.PP[j].PPSS = ppValueArray[0]
						CacheData.PP[j].PP99 = ppValueArray[1]
						CacheData.PP[j].PP98 = ppValueArray[2]
						CacheData.PP[j].PP97 = ppValueArray[3]
						CacheData.PP[j].PP95 = ppValueArray[4]
					}
				}
				if !modExist {
					CacheData.PP = append(CacheData.PP, structs.PPData{
						Mods: mods,
						SR:   beatmap.DifficultyRating,
						PPSS: ppValueArray[0],
						PP99: ppValueArray[1],
						PP98: ppValueArray[2],
						PP97: ppValueArray[3],
						PP95: ppValueArray[4],
					})
				}
				cache[index] = CacheData
			} else {
				var cachePPData []structs.PPData
				cachePPData = append(cachePPData, structs.PPData{
					Mods: mods,
					SR:   beatmap.DifficultyRating,
					PPSS: ppValueArray[0],
					PP99: ppValueArray[1],
					PP98: ppValueArray[2],
					PP97: ppValueArray[3],
					PP95: ppValueArray[4],
				})
				cache = append(cache, structs.MapData{
					Time:    time.Now(),
					Beatmap: beatmap,
					PP:      cachePPData,
				})
			}
			jsonCache, err := json.Marshal(cache)
			tools.ErrRead(err)
			err = ioutil.WriteFile("./data/osuData/mapCache.json", jsonCache, 0644)
			tools.ErrRead(err)
		}
	}
	return starRating, ppSS, pp99, pp98, pp97, pp95
}
