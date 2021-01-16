package osucommands

import (
	"fmt"
	"math"
	"math/rand"
	"regexp"
	"strconv"
	"strings"

	osuapi "maquiaBot/osu-api"
	osutools "maquiaBot/osu-tools"

	"github.com/bwmarrin/discordgo"
)

// BeatmapMessage is a handler executed when a message contains a beatmap link
func BeatmapMessage(s *discordgo.Session, m *discordgo.MessageCreate, regex *regexp.Regexp) {
	modRegex, _ := regexp.Compile(`(?i)-m\s*(\S+)`)
	accRegex, _ := regexp.Compile(`(?i)-acc\s*(\S+)`)
	comboRegex, _ := regexp.Compile(`(?i)-c\s*(\S+)`)
	goodRegex, _ := regexp.Compile(`(?i)-100\s*(\d+)`)
	mehRegex, _ := regexp.Compile(`(?i)-50\s*(\d+)`)
	missRegex, _ := regexp.Compile(`(?i)-x\s*(\S+)`)
	mapRegex, _ := regexp.Compile(`(?i)[^-]m`)
	scoreRegex, _ := regexp.Compile(`(?i)-s\s*(\S+)`)

	// See if map was linked or if the map command was used
	var submatches []string
	if (strings.Contains(m.Content, "map") || mapRegex.MatchString(m.Content)) && !regex.MatchString(m.Content) {
		// Get prev messages
		messages, err := s.ChannelMessages(m.ChannelID, -1, "", "", "")
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, "No map found!")
			return
		}

		// Look for a valid beatmap ID
		for _, msg := range messages {
			if len(msg.Embeds) > 0 && msg.Embeds[0].Author != nil {
				if regex.MatchString(msg.Embeds[0].URL) {
					submatches = regex.FindStringSubmatch(msg.Embeds[0].URL)
					break
				} else if regex.MatchString(msg.Embeds[0].Author.URL) {
					submatches = regex.FindStringSubmatch(msg.Embeds[0].Author.URL)
					break
				}
			} else if regex.MatchString(msg.Content) {
				submatches = regex.FindStringSubmatch(msg.Content)
				break
			}
		}
	} else {
		submatches = regex.FindStringSubmatch(m.Content)
	}
	if len(submatches) == 0 {
		s.ChannelMessageSend(m.ChannelID, "No map found!")
		return
	}

	message, err := s.ChannelMessageSend(m.ChannelID, "Processing beatmap...")
	if err != nil {
		return
	}
	var beatmap osuapi.Beatmap

	// Get requested mods
	mods := "NM"
	if modRegex.MatchString(m.Content) {
		mods = strings.ToUpper(modRegex.FindStringSubmatch(m.Content)[1])
		if strings.Contains(mods, "NC") && !strings.Contains(mods, "DT") {
			mods += "DT"
		}
	}

	// Check if the format uses a /b/, /s/, /beatmaps/, or /beatmapsets/ link
	diffMods := 338 & osuapi.ParseMods(mods)
	if diffMods&256 != 0 && diffMods&64 != 0 { // Remove DTHT
		diffMods -= 320
	}
	if diffMods&2 != 0 && diffMods&16 != 0 { // Remove EZHR
		diffMods -= 18
	}
	switch submatches[2] {
	case "s":
		beatmap = osutools.BeatmapParse(submatches[3], "set", &diffMods)
	case "b":
		beatmap = osutools.BeatmapParse(submatches[3], "map", &diffMods)
	case "beatmaps":
		beatmap = osutools.BeatmapParse(submatches[3], "map", &diffMods)
	case "beatmapsets":
		if len(submatches[6]) > 0 {
			beatmap = osutools.BeatmapParse(submatches[6], "map", &diffMods)
		} else {
			beatmap = osutools.BeatmapParse(submatches[3], "set", &diffMods)
		}
	}

	// Check if a beatmap was even obtained
	if beatmap.BeatmapID == 0 {
		s.ChannelMessageDelete(message.ChannelID, message.ID)
		return
	}

	// Assign embed colour for different modes
	Color := osutools.ModeColour(beatmap.Mode)

	// Obtain whole set
	beatmaps, err := OsuAPI.GetBeatmaps(osuapi.GetBeatmapsOpts{
		BeatmapSetID: beatmap.BeatmapSetID,
	})
	if err != nil {
		s.ChannelMessageEdit(message.ChannelID, message.ID, "")
		s.ChannelMessageEdit(message.ChannelID, message.ID, "The osu! API just owned me. Please try again!")
		return
	}

	// Assign variables for map specs
	totalMinutes := math.Floor(float64(beatmap.TotalLength / 60))
	totalSeconds := fmt.Sprint(math.Mod(float64(beatmap.TotalLength), float64(60)))
	if len(totalSeconds) == 1 {
		totalSeconds = "0" + totalSeconds
	}
	hitMinutes := math.Floor(float64(beatmap.HitLength / 60))
	hitSeconds := fmt.Sprint(math.Mod(float64(beatmap.HitLength), float64(60)))
	if len(hitSeconds) == 1 {
		hitSeconds = "0" + hitSeconds
	}

	sr := "**SR:** " + strconv.FormatFloat(beatmap.DifficultyRating, 'f', 2, 64)
	if beatmap.Mode == osuapi.ModeOsu {
		sr += " **Aim:** " + strconv.FormatFloat(beatmap.DifficultyAim, 'f', 2, 64) + " **Speed:** " + strconv.FormatFloat(beatmap.DifficultySpeed, 'f', 2, 64)
	}
	length := "**Length:** " + fmt.Sprint(totalMinutes) + ":" + totalSeconds + " (" + fmt.Sprint(hitMinutes) + ":" + hitSeconds + ") "
	bpm := "**BPM:** " + fmt.Sprint(beatmap.BPM) + " "
	combo := "**FC:** " + strconv.Itoa(beatmap.MaxCombo) + "x"
	mapStats := "**CS:** " + strconv.FormatFloat(beatmap.CircleSize, 'f', 1, 64) + " **AR:** " + strconv.FormatFloat(beatmap.ApproachRate, 'f', 1, 64) + " **OD:** " + strconv.FormatFloat(beatmap.OverallDifficulty, 'f', 1, 64) + " **HP:** " + strconv.FormatFloat(beatmap.HPDrain, 'f', 1, 64)
	mapObjs := "**Circles:** " + strconv.Itoa(beatmap.Circles) + " **Sliders:** " + strconv.Itoa(beatmap.Sliders) + " **Spinners:** " + strconv.Itoa(beatmap.Spinners)

	status := "**Rank Status:** " + strings.Title(beatmap.Approved.String())

	download := "**Download:** [osz link](https://osu.ppy.sh/beatmapsets/" + strconv.Itoa(beatmap.BeatmapSetID) + "/download)" + " | <osu://dl/" + strconv.Itoa(beatmap.BeatmapSetID) + ">"
	var diffs string
	if len(beatmaps) == 1 {
		diffs = "**1** difficulty <:ahFuck:550808614202245131>"
	} else {
		diffs = "**" + strconv.Itoa(len(beatmaps)) + "** difficulties <:ahFuck:550808614202245131>"
	}

	totalHits := beatmap.Circles + beatmap.Sliders + beatmap.Spinners

	// Get miss value
	missVal := ""
	missNum := 0
	if missRegex.MatchString(m.Content) {
		missVal = missRegex.FindStringSubmatch(m.Content)[1]
		missNum, err = strconv.Atoi(missVal)
		if err != nil || missNum <= 0 || missNum > totalHits {
			missVal = ""
		}
	}

	// Get acc value
	accVal := ""
	if accRegex.MatchString(m.Content) {
		accVal = accRegex.FindStringSubmatch(m.Content)[1]
		_, err = strconv.ParseFloat(accVal, 64)
		if err != nil {
			accVal = ""
		}
	}

	// Manual acc calc instead of auto decision on 100s and 50s if given
	greats := totalHits
	goods := 0
	mehs := 0
	if goodRegex.MatchString(m.Content) {
		goods, _ = strconv.Atoi(goodRegex.FindStringSubmatch(m.Content)[1])
		if goods > greats || goods < 0 {
			goods = 0
		}
	}

	if mehRegex.MatchString(m.Content) {
		mehs, _ = strconv.Atoi(mehRegex.FindStringSubmatch(m.Content)[1])
		if mehs > greats || mehs < 0 {
			mehs = 0
		}
	}
	if mehs+goods+missNum > greats { // Reset if bad values are given
		greats = 0
		mehs = 0
		goods = 0
	} else if mehs+goods+missNum == 0 { // Reset greats if no manual calculation is wanted
		greats = 0
	} else {
		greats -= goods
		greats -= mehs
		greats -= missNum
		acc := 100.0 * float64(mehs+2*goods+6*greats) / float64(6*(missNum+mehs+goods+greats))
		accVal = strconv.FormatFloat(acc, 'f', 2, 64)
	}

	// Get combo value
	comboVal := ""
	if comboRegex.MatchString(m.Content) {
		comboVal = comboRegex.FindStringSubmatch(m.Content)[1]
		comboNum, err := strconv.Atoi(comboVal)
		if err != nil || comboNum <= 0 || comboNum > beatmap.MaxCombo {
			comboVal = ""
		}

		if comboNum != beatmap.MaxCombo && accVal == "" && missVal == "" && greats == 0 {
			if beatmap.Sliders > 0 {
				greats = totalHits - 1
				goods = 1
			} else {
				greats = totalHits - 1
				missNum = 1
				missVal = "1"
			}
		}
	}

	// Get score value (only for osu!mania)
	if beatmap.Mode == osuapi.ModeOsuMania && scoreRegex.MatchString(m.Content) && accVal == "" {
		accVal = scoreRegex.FindStringSubmatch(m.Content)[1]
		_, err = strconv.ParseInt(accVal, 10, 64)
		if err != nil {
			accVal = ""
		}
	}

	// Calculate SR and PP
	values := osutools.BeatmapCalc(mods, accVal, comboVal, missVal, greats, goods, mehs, beatmap)
	ppText := ""
	if len(values) == 1 {
		ppText = values[0]
	} else if len(values) != 0 {
		ppText = values[0] + values[1] + values[2] + values[3] + values[4]
	}

	ppTextHeader := "**[" + beatmap.DiffName + "]** with mods: **" + strings.ToUpper(mods) + "**"
	if comboVal != "" {
		ppTextHeader += ", combo: **" + comboVal + "**"
	}
	if missVal != "" {
		ppTextHeader += ", misses: **" + missVal + "**"
	}

	// Create embed
	embed := &discordgo.MessageEmbed{
		Author: &discordgo.MessageEmbedAuthor{
			URL:     "https://osu.ppy.sh/beatmaps/" + strconv.Itoa(beatmap.BeatmapID),
			Name:    beatmap.Artist + " - " + beatmap.Title + " by " + beatmap.Creator,
			IconURL: "https://a.ppy.sh/" + strconv.Itoa(beatmap.CreatorID) + "?" + strconv.Itoa(rand.Int()) + ".jpeg",
		},
		Color: Color,
		Description: sr + "\n" +
			length + bpm + combo + "\n" +
			mapStats + "\n" +
			mapObjs + "\n" +
			status + "\n" +
			download + "\n" +
			diffs + "\n" + "\n" +
			ppTextHeader + "\n" +
			ppText,
		Thumbnail: &discordgo.MessageEmbedThumbnail{
			URL: "https://b.ppy.sh/thumb/" + strconv.Itoa(beatmap.BeatmapSetID) + "l.jpg",
		},
	}
	if strings.ToLower(beatmap.Title) == "crab rave" {
		embed.Image = &discordgo.MessageEmbedImage{
			URL: "https://cdn.discordapp.com/emojis/510169818893385729.gif",
		}
	}
	content := ""
	s.ChannelMessageEditComplex(&discordgo.MessageEdit{
		Content: &content,
		Embed:   embed,
		ID:      message.ID,
		Channel: message.ChannelID,
	})
	return
}
