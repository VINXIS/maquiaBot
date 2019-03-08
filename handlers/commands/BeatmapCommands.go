package commands

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"os/exec"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	structs "../../structs"
	tools "../../tools"

	"github.com/bwmarrin/discordgo"
	"github.com/thehowl/go-osuapi"
)

// BeatmapMessage is a handler executed when a message contains a beatmap link
func BeatmapMessage(s *discordgo.Session, m *discordgo.MessageCreate, regex *regexp.Regexp, osu *osuapi.Client, cache []structs.MapData) {
	submatches := regex.FindStringSubmatch(m.Content)

	// Check if message wants the bot to send details or not before doing anything
	if submatches[9] != "-n" {
		message, err := s.ChannelMessageSend(m.ChannelID, "Processing beatmap...")
		var beatmap osuapi.Beatmap

		// These if statements check if the format uses a /b/, /s/, /beatmaps/, or /beatmapsets/ link
		if len(submatches[3]) > 0 {
			if len(submatches[7]) > 0 {
				beatmap = beatmapParse(submatches[7], "map", osu)
			} else {
				beatmap = beatmapParse(submatches[4], "set", osu)
			}
		} else {
			if submatches[2] == "s" {
				beatmap = beatmapParse(submatches[4], "set", osu)
			} else {
				beatmap = beatmapParse(submatches[4], "map", osu)
			}
		}

		// Check if a beatmap was even obtained
		if beatmap == (osuapi.Beatmap{}) {
			s.ChannelMessageDelete(message.ChannelID, message.ID)
			return
		}

		log.Println("Someone linked a beatmap! The beatmap is " + strconv.Itoa(beatmap.BeatmapID) + " " + beatmap.Artist + " - " + beatmap.Title + " by " + beatmap.Creator)

		// Assign embed colour for different modes
		Color := tools.ModeColour(beatmap.Mode)

		// Temporary method to obtain mapper user id, once creator id is available, actual user avatars will be used for banned users
		mapper, err := osu.GetUser(osuapi.GetUserOpts{
			Username: beatmap.Creator,
		})
		if err != nil {
			mapper, err = osu.GetUser(osuapi.GetUserOpts{
				UserID: 3,
			})
			mapper.Username = beatmap.Creator
		}

		// Obtain whole set
		beatmaps, err := osu.GetBeatmaps(osuapi.GetBeatmapsOpts{
			BeatmapSetID: beatmap.BeatmapSetID,
		})
		tools.ErrRead(err)

		// Assign variables for map specs
		totalMinutes := math.Floor(float64(beatmap.TotalLength / 60))
		totalSeconds := math.Mod(float64(beatmap.TotalLength), float64(60))
		hitMinutes := math.Floor(float64(beatmap.HitLength / 60))
		hitSeconds := math.Mod(float64(beatmap.HitLength), float64(60))

		length := "**Length:** " + fmt.Sprint(totalMinutes) + ":" + fmt.Sprint(totalSeconds) + " (" + fmt.Sprint(hitMinutes) + ":" + fmt.Sprint(hitSeconds) + ") "
		bpm := "**BPM:** " + fmt.Sprint(beatmap.BPM) + " "
		combo := "**FC:** " + strconv.Itoa(beatmap.MaxCombo) + "x"

		status := "**Rank Status:** " + beatmap.Approved.String()

		download := "**Download:** [osz link](https://osu.ppy.sh/d/" + strconv.Itoa(beatmap.BeatmapSetID) + ")" + " | <osu://dl/" + strconv.Itoa(beatmap.BeatmapSetID) + ">"
		diffs := "**" + strconv.Itoa(len(beatmaps)) + `** difficulties <:ahFuck:550808614202245131>`

		// Get requested mods
		mods := "NM"
		if len(submatches[12]) > 0 {
			mods = submatches[12]
			if len(mods)%2 == 0 && len(osuapi.ParseMods(mods).String()) > 0 {
				mods = osuapi.ParseMods(mods).String()
			}
		}

		// Calculate SR
		//SRCalc(beatmap, mods)

		// Calculate pp
		var (
			starRating string
			ppSS       string
			pp99       string
			pp98       string
			pp97       string
			pp95       string
		)
		latest, update := false, false
		index := 0
		var PPData structs.PPData
		var CacheData structs.MapData
		twoDay, err := time.ParseDuration("48h")
		tools.ErrRead(err)
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
				go PPCalc(beatmap, 100.0, mods, ppValues)
				go PPCalc(beatmap, 99.0, mods, ppValues)
				go PPCalc(beatmap, 98.0, mods, ppValues)
				go PPCalc(beatmap, 97.0, mods, ppValues)
				go PPCalc(beatmap, 95.0, mods, ppValues)
				for v := 0; v < 5; v++ {
					ppValueArray[v] = <-ppValues
				}
				sort.Slice(ppValueArray[:], func(i, j int) bool {
					return ppValueArray[i] > ppValueArray[j]
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
				err = ioutil.WriteFile("./data/osuCache.json", jsonCache, 0644)
			} else {
				ppSS = "pp is not available for ctb yet"
				pp99 = ""
				pp98 = ""
				pp97 = ""
				pp95 = ""
			}
		}

		// Create embed
		embed := &discordgo.MessageEmbed{
			Author: &discordgo.MessageEmbedAuthor{
				URL:     "https://osu.ppy.sh/beatmaps/" + strconv.Itoa(beatmap.BeatmapID),
				Name:    beatmap.Artist + " - " + beatmap.Title + " by " + mapper.Username,
				IconURL: "https://a.ppy.sh/" + strconv.Itoa(mapper.UserID),
			},
			Color: Color,
			Description: starRating + length + bpm + combo + "\n" +
				status + "\n" +
				download + "\n" +
				diffs + "\n" + "\n" +
				"**[" + beatmap.DiffName + "]** with mods: " + mods + "\n" +
				//aimRating + speedRating + totalRating + "\n" +
				ppSS + pp99 + pp98 + pp97 + pp95,
			Thumbnail: &discordgo.MessageEmbedThumbnail{
				URL: "https://b.ppy.sh/thumb/" + strconv.Itoa(beatmap.BeatmapSetID) + "l.jpg",
			},
		}
		s.ChannelMessageEdit(message.ChannelID, message.ID, "")
		s.ChannelMessageEditEmbed(message.ChannelID, message.ID, embed)
		return
	}
}

func beatmapParse(id, format string, osu *osuapi.Client) (beatmap osuapi.Beatmap) {
	replacer, _ := regexp.Compile(`[^a-zA-Z0-9\s\(\)]`)

	mapID, err := strconv.Atoi(id)
	tools.ErrRead(err)
	if format == "map" {
		// Fetch the beatmap
		beatmaps, err := osu.GetBeatmaps(osuapi.GetBeatmapsOpts{
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
		// Fetch the set
		beatmaps, err := osu.GetBeatmaps(osuapi.GetBeatmapsOpts{
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

// SRCalc calcualtes the aim, speed, and total SR for a beatmap
func SRCalc(beatmap osuapi.Beatmap, mods string) (aim, speed, total string) {
	/*replacer, _ := regexp.Compile(`[^a-zA-Z0-9\s\(\)]`)
	//fileName := strconv.Itoa(rand.Intn(10000000))

	var commands []string
	commands = append(commands, "/C", "start", "dotnet", "run", "-p", "./osu-tools/PerformanceCalculator", "difficulty", "./data/osuFiles/"+strconv.Itoa(beatmap.BeatmapID)+" "+replacer.ReplaceAllString(beatmap.Artist, "")+" - "+replacer.ReplaceAllString(beatmap.Title, "")+".osu")
	cmd := exec.Command("cmd", commands[:]...)
	err := cmd.Start()
	tools.ErrRead(err)
	fmt.Println(string(res))
	for {
		if cmd.ProcessState.ExitCode() == 0 {
			text, err := ioutil.ReadFile(fileName + ".txt")
			tools.ErrRead(err)
			fmt.Println(string(text))
			tools.DeleteFile("./" + fileName + ".txt")
			break
		}
	}*/

	aim = "**Aim SR**: 0 "
	speed = "**Speed SR**: 0 "
	total = "**Total SR**: 0"
	return aim, speed, total
}

// PPCalc calculates the pp given by the beatmap with specified acc and mods TODO: More args
func PPCalc(beatmap osuapi.Beatmap, pp float64, mods string, store chan<- string) {
	replacer, _ := regexp.Compile(`[^a-zA-Z0-9\s\(\)]`)

	regex, err := regexp.Compile(`pp             : (\d+)(\.\d+)?`)
	tools.ErrRead(err)

	var data []string
	var commands []string
	var mode string
	switch beatmap.Mode {
	case osuapi.ModeOsu:
		mode = "osu"
	case osuapi.ModeOsuMania:
		mode = "mania"
	case osuapi.ModeTaiko:
		mode = "taiko"
	}
	commands = append(commands, "run", "-p", "./osu-tools/PerformanceCalculator", "simulate", mode, "./data/osuFiles/"+strconv.Itoa(beatmap.BeatmapID)+" "+replacer.ReplaceAllString(beatmap.Artist, "")+" - "+replacer.ReplaceAllString(beatmap.Title, "")+".osu", "-a", fmt.Sprint(pp))

	// Check mods
	if len(mods) > 0 && mods != "NM" {
		var modResult strings.Builder
		modList := tools.StringSplit(mods, 2)
		for i := range modList {
			modResult.WriteString("-m " + strings.ToLower(modList[i]) + " ")
		}
		commands = append(commands, strings.Split(modResult.String(), " ")[:]...)
	}

	out, err := exec.Command("dotnet", commands[:]...).Output()
	tools.ErrRead(err)
	data = strings.Split(string(out), "\n")

	var res []string
	for _, line := range data {
		if regex.MatchString(line) {
			res = regex.FindStringSubmatch(line)
		}
	}
	ppValue, err := strconv.ParseFloat(res[1]+res[2], 64)
	tools.ErrRead(err)

	value := strconv.FormatFloat(math.Round(ppValue), 'f', 0, 64)
	store <- value
}
