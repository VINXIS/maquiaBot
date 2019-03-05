package commands

import (
	"bytes"
	"fmt"
	"log"
	"math"
	"os/exec"
	"regexp"
	"sort"
	"strconv"
	"strings"

	tools "../../tools"

	"github.com/bwmarrin/discordgo"
	"github.com/thehowl/go-osuapi"
)

// BeatmapMessage is a handler executed when a message contains a beatmap link
func BeatmapMessage(s *discordgo.Session, m *discordgo.MessageCreate, regex *regexp.Regexp, osu *osuapi.Client) {
	submatches := regex.FindStringSubmatch(m.Content)

	// Check if message wants the bot to send details or not before doing anything
	if submatches[9] != "-n" {
		message, err := s.ChannelMessageSend(m.ChannelID, "Processing beatmap...")
		var beatmap osuapi.Beatmap
		var Color int

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

		// Assign embed colour for different modes
		switch beatmap.Mode {
		case osuapi.ModeOsu:
			Color = 0xD65288
		case osuapi.ModeTaiko:
			Color = 0xFF0000
		case osuapi.ModeCatchTheBeat:
			Color = 0x007419
		case osuapi.ModeOsuMania:
			Color = 0xff6200
		}

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

		starRating := "**SR:** " + fmt.Sprintf("%.2f", beatmap.DifficultyRating) + " "
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
		s.ChannelMessageEdit(message.ChannelID, message.ID, "Calculating SR...")
		//aimRating, speedRating, totalRating := SRCalc(beatmap, mods)

		// Calculate pp
		s.ChannelMessageEdit(message.ChannelID, message.ID, "Calculating pp...")
		ppSS := "**100%:** " + PPCalc(beatmap, 100.0, mods) + " | "
		pp99 := "**99%:** " + PPCalc(beatmap, 99.0, mods) + " | "
		pp98 := "**98%:** " + PPCalc(beatmap, 98.0, mods) + " | "
		pp97 := "**97%:** " + PPCalc(beatmap, 97.0, mods) + " | "
		pp95 := "**95%:** " + PPCalc(beatmap, 95.0, mods)

		// Create embed
		s.ChannelMessageEdit(message.ChannelID, message.ID, "Creating embed...")
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
				"**" + beatmap.DiffName + "** diff with mods: " + mods + "\n" +
				//aimRating + speedRating + totalRating + "\n" +
				ppSS + pp99 + pp98 + pp97 + pp95,
			Thumbnail: &discordgo.MessageEmbedThumbnail{
				URL: "https://b.ppy.sh/thumb/" + strconv.Itoa(beatmap.BeatmapSetID) + "l.jpg",
			},
		}
		s.ChannelMessageEdit(message.ChannelID, message.ID, "")
		s.ChannelMessageEditEmbed(message.ChannelID, message.ID, embed)
	}
}

func beatmapParse(id string, format string, osu *osuapi.Client) (beatmap osuapi.Beatmap) {
	mapID, err := strconv.Atoi(id)
	tools.ErrRead(err)
	if format == "map" {
		// Fetch the beatmap
		beatmaps, err := osu.GetBeatmaps(osuapi.GetBeatmapsOpts{
			BeatmapID: mapID,
		})
		tools.ErrRead(err)
		beatmap = beatmaps[0]

		// Download the .osu file for the map
		err = tools.DownloadFile(
			"./data/osuFiles/"+
				strconv.Itoa(beatmap.BeatmapID)+
				" "+
				beatmap.Artist+
				" - "+
				beatmap.Title+
				".osu",
			"https://osu.ppy.sh/osu/"+
				strconv.Itoa(beatmap.BeatmapID))
		tools.ErrRead(err)
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
			err = tools.DownloadFile(
				"./data/osuFiles/"+
					strconv.Itoa(diff.BeatmapID)+
					" "+
					diff.Artist+
					" - "+
					diff.Title+
					".osu",
				"https://osu.ppy.sh/osu/"+
					strconv.Itoa(diff.BeatmapID))
			tools.ErrRead(err)
		}
		beatmap = beatmaps[0]
	}
	return beatmap
}

// SRCalc calcualtes the aim, speed, and total SR for a beatmap
func SRCalc(beatmap osuapi.Beatmap, mods string) (aim string, speed string, total string) {
	_, err := regexp.Compile(`pp             : (\d+)(\.\d+)?`)
	tools.ErrRead(err)

	var commands []string
	fmt.Println(beatmap.Mode.String())
	commands = append(commands, "run", "-p", "./osu-tools/PerformanceCalculator", "difficulty", "./data/osuFiles/"+strconv.Itoa(beatmap.BeatmapID)+" "+beatmap.Artist+" - "+beatmap.Title+".osu")

	// Check mods
	if len(mods) > 0 && mods != "NM" {
		var modResult strings.Builder
		modList := tools.StringSplit(mods, 2)
		for i := 0; i < len(modList); i++ {
			modResult.WriteString("-m " + strings.ToLower(modList[i]) + " ")
		}
		commands = append(commands, strings.Split(modResult.String(), " ")[:]...)
	}

	cmd := exec.Command("dotnet", commands[:]...)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err = cmd.Run()
	if err != nil {
		log.Fatal(fmt.Sprint(err) + ": " + stderr.String())
		return
	}
	data := strings.Split(out.String(), "\n")
	fmt.Println(data)

	aim = "**Aim SR**: 0 "
	speed = "**Speed SR**: 0 "
	total = "**Total SR**: 0"
	return aim, speed, total
}

// PPCalc calculates the pp given by the beatmap with specified acc and mods TODO: More args
func PPCalc(beatmap osuapi.Beatmap, pp float64, mods string) (value string) {
	regex, err := regexp.Compile(`pp             : (\d+)(\.\d+)?`)
	tools.ErrRead(err)

	var data []string
	var commands []string
	commands = append(commands, "run", "-p", "./osu-tools/PerformanceCalculator", "simulate", "osu", "./data/osuFiles/"+strconv.Itoa(beatmap.BeatmapID)+" "+beatmap.Artist+" - "+beatmap.Title+".osu", "-a", fmt.Sprint(pp))

	// Check mods
	if len(mods) > 0 && mods != "NM" {
		var modResult strings.Builder
		modList := tools.StringSplit(mods, 2)
		for i := 0; i < len(modList); i++ {
			modResult.WriteString("-m " + strings.ToLower(modList[i]) + " ")
		}
		commands = append(commands, strings.Split(modResult.String(), " ")[:]...)
	}

	out, err := exec.Command("dotnet", commands[:]...).Output()
	tools.ErrRead(err)
	data = strings.Split(string(out), "\n")

	res := regex.FindStringSubmatch(data[14])
	ppValue, err := strconv.ParseFloat(res[1]+res[2], 64)
	tools.ErrRead(err)

	value = fmt.Sprint(math.Round(ppValue)) + "pp"
	return value
}
