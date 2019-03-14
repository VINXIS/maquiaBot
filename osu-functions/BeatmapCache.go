package osutools

import (
	"encoding/json"
	"io/ioutil"
	"sort"
	"strconv"
	"time"

	structs "../structs"
	tools "../tools"

	"github.com/thehowl/go-osuapi"
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
		starRating = "**SR:** " + PPData.SR + " "
		ppSS = "**100%:** " + PPData.PPSS + "pp | "
		pp99 = "**99%:** " + PPData.PP99 + "pp | "
		pp98 = "**98%:** " + PPData.PP98 + "pp | "
		pp97 = "**97%:** " + PPData.PP97 + "pp | "
		pp95 = "**95%:** " + PPData.PP95 + "pp"
	} else {
		if beatmap.Mode != osuapi.ModeCatchTheBeat {
			ppValues := make(chan string, 5)
			var ppValueArray [5]string
			go PPCalc(beatmap, 100.0, "", "", mods, ppValues)
			go PPCalc(beatmap, 99.0, "", "", mods, ppValues)
			go PPCalc(beatmap, 98.0, "", "", mods, ppValues)
			go PPCalc(beatmap, 97.0, "", "", mods, ppValues)
			go PPCalc(beatmap, 95.0, "", "", mods, ppValues)
			for v := 0; v < 5; v++ {
				ppValueArray[v] = <-ppValues
			}
			sort.Slice(ppValueArray[:], func(i, j int) bool {
				pp1, _ := strconv.Atoi(ppValueArray[i])
				pp2, _ := strconv.Atoi(ppValueArray[j])
				return pp1 > pp2
			})
			starRating = "**SR:** " + strconv.FormatFloat(beatmap.DifficultyRating, 'f', 2, 64) + " "
			ppSS = "**100%:** " + ppValueArray[0] + "pp | "
			pp99 = "**99%:** " + ppValueArray[1] + "pp | "
			pp98 = "**98%:** " + ppValueArray[2] + "pp | "
			pp97 = "**97%:** " + ppValueArray[3] + "pp | "
			pp95 = "**95%:** " + ppValueArray[4] + "pp"
			if update {
				CacheData.Beatmap = beatmap
				CacheData.Time = time.Now()
				modExist := false
				for j := range CacheData.PP {
					if CacheData.PP[j].Mods == mods {
						modExist = true
						CacheData.PP[j].SR = strconv.FormatFloat(beatmap.DifficultyRating, 'f', 2, 64)
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
						SR:   strconv.FormatFloat(beatmap.DifficultyRating, 'f', 2, 64),
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
					SR:   strconv.FormatFloat(beatmap.DifficultyRating, 'f', 2, 64),
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
		} else {
			ppSS = "pp is not available for ctb yet"
			pp99 = ""
			pp98 = ""
			pp97 = ""
			pp95 = ""
		}
	}
	return starRating, ppSS, pp99, pp98, pp97, pp95
}
