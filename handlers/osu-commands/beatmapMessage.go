package osucommands

import (
	"fmt"
	"math"
	"math/rand"
	"regexp"
	"strconv"
	"strings"

	osuapi "../../osu-api"
	osutools "../../osu-functions"
	structs "../../structs"

	"github.com/bwmarrin/discordgo"
)

// BeatmapMessage is a handler executed when a message contains a beatmap link
func BeatmapMessage(s *discordgo.Session, m *discordgo.MessageCreate, regex *regexp.Regexp, osuAPI *osuapi.Client, cache []structs.MapData) {
	modRegex, _ := regexp.Compile(`-m\s*(\S+)`)
	submatches := regex.FindStringSubmatch(m.Content)

	message, err := s.ChannelMessageSend(m.ChannelID, "Processing beatmap...")
	var beatmap osuapi.Beatmap

	// Get requested mods
	mods := "NM"
	if modRegex.MatchString(m.Content) {
		modCheck := modRegex.FindStringSubmatch(m.Content)
		modsRaw := strings.ToUpper(modCheck[1])
		if len(modsRaw)%2 == 0 && len(osuapi.ParseMods(modsRaw).String()) > 0 {
			mods = osuapi.ParseMods(modsRaw).String()
		}
	}

	// Check if the format uses a /b/, /s/, /beatmaps/, or /beatmapsets/ link
	switch submatches[2] {
	case "s":
		beatmap = osutools.BeatmapParse(submatches[3], "set", osuapi.ParseMods(mods), osuAPI)
	case "b":
		beatmap = osutools.BeatmapParse(submatches[3], "map", osuapi.ParseMods(mods), osuAPI)
	case "beatmaps":
		beatmap = osutools.BeatmapParse(submatches[3], "map", osuapi.ParseMods(mods), osuAPI)
	case "beatmapsets":
		if len(submatches[6]) > 0 {
			beatmap = osutools.BeatmapParse(submatches[6], "map", osuapi.ParseMods(mods), osuAPI)
		} else {
			beatmap = osutools.BeatmapParse(submatches[3], "set", osuapi.ParseMods(mods), osuAPI)
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
	beatmaps, err := osuAPI.GetBeatmaps(osuapi.GetBeatmapsOpts{
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

	// Calculate SR and PP
	starRating, ppSS, pp99, pp98, pp97, pp95 := osutools.BeatmapCache(mods, beatmap, cache)

	// Create embed
	embed := &discordgo.MessageEmbed{
		Author: &discordgo.MessageEmbedAuthor{
			URL:     "https://osu.ppy.sh/beatmaps/" + strconv.Itoa(beatmap.BeatmapID),
			Name:    beatmap.Artist + " - " + beatmap.Title + " by " + beatmap.Creator,
			IconURL: "https://a.ppy.sh/" + strconv.Itoa(beatmap.CreatorID) + "?" + strconv.Itoa(rand.Int()) + ".jpeg",
		},
		Color: Color,
		Description: starRating + length + bpm + combo + "\n" +
			mapStats + "\n" +
			mapObjs + "\n" +
			status + "\n" +
			download + "\n" +
			diffs + "\n" + "\n" +
			"**[" + beatmap.DiffName + "]** with mods: **" + strings.ToUpper(mods) + "**\n" +
			//aimRating + speedRating + totalRating + "\n" + TODO: Make SR calc work
			ppSS + pp99 + pp98 + pp97 + pp95,
		Thumbnail: &discordgo.MessageEmbedThumbnail{
			URL: "https://b.ppy.sh/thumb/" + strconv.Itoa(beatmap.BeatmapSetID) + "l.jpg",
		},
	}
	if strings.ToLower(beatmap.Title) == "crab rave" {
		embed.Image = &discordgo.MessageEmbedImage{
			URL: "https://cdn.discordapp.com/emojis/510169818893385729.gif",
		}
	}
	s.ChannelMessageEdit(message.ChannelID, message.ID, "")
	s.ChannelMessageEditEmbed(message.ChannelID, message.ID, embed)
	return
}
