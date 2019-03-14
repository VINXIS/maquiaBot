package osucommands

import (
	"fmt"
	"log"
	"math"
	"regexp"
	"strconv"

	osutools "../../osu-functions"
	structs "../../structs"
	tools "../../tools"

	"github.com/bwmarrin/discordgo"
	"github.com/thehowl/go-osuapi"
)

// BeatmapMessage is a handler executed when a message contains a beatmap link
func BeatmapMessage(s *discordgo.Session, m *discordgo.MessageCreate, regex *regexp.Regexp, osuAPI *osuapi.Client, cache []structs.MapData) {
	submatches := regex.FindStringSubmatch(m.Content)

	// Check if message wants the bot to send details or not before doing anything
	if submatches[9] != "-n" {
		message, err := s.ChannelMessageSend(m.ChannelID, "Processing beatmap...")
		var beatmap osuapi.Beatmap

		// These if statements check if the format uses a /b/, /s/, /beatmaps/, or /beatmapsets/ link
		if len(submatches[3]) > 0 {
			if len(submatches[7]) > 0 {
				beatmap = osutools.BeatmapParse(submatches[7], "map", osuAPI)
			} else {
				beatmap = osutools.BeatmapParse(submatches[4], "set", osuAPI)
			}
		} else {
			if submatches[2] == "s" {
				beatmap = osutools.BeatmapParse(submatches[4], "set", osuAPI)
			} else {
				beatmap = osutools.BeatmapParse(submatches[4], "map", osuAPI)
			}
		}

		// Check if a beatmap was even obtained
		if beatmap == (osuapi.Beatmap{}) {
			s.ChannelMessageDelete(message.ChannelID, message.ID)
			return
		}

		log.Println("Someone linked a beatmap! The beatmap is " + strconv.Itoa(beatmap.BeatmapID) + " " + beatmap.Artist + " - " + beatmap.Title + " by " + beatmap.Creator)

		// Assign embed colour for different modes
		Color := osutools.ModeColour(beatmap.Mode)

		// Temporary method to obtain mapper user id, once creator id is available, actual user avatars will be used for banned users
		mapper, err := osuAPI.GetUser(osuapi.GetUserOpts{
			Username: beatmap.Creator,
		})
		if err != nil {
			mapper, err = osuAPI.GetUser(osuapi.GetUserOpts{
				UserID: 3,
			})
			mapper.Username = beatmap.Creator
		}

		// Obtain whole set
		beatmaps, err := osuAPI.GetBeatmaps(osuapi.GetBeatmapsOpts{
			BeatmapSetID: beatmap.BeatmapSetID,
		})
		tools.ErrRead(err)

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

		status := "**Rank Status:** " + beatmap.Approved.String()

		download := "**Download:** [osz link](https://osu.ppy.sh/d/" + strconv.Itoa(beatmap.BeatmapSetID) + ")" + " | <osu://dl/" + strconv.Itoa(beatmap.BeatmapSetID) + ">"
		var diffs string
		if len(beatmaps) == 1 {
			diffs = "**" + strconv.Itoa(len(beatmaps)) + `** difficulty <:ahFuck:550808614202245131>`
		} else {
			diffs = "**" + strconv.Itoa(len(beatmaps)) + `** difficulties <:ahFuck:550808614202245131>`
		}

		// Get requested mods
		mods := "NM"
		if len(submatches[12]) > 0 {
			mods = submatches[12]
			if len(mods)%2 == 0 && len(osuapi.ParseMods(mods).String()) > 0 {
				mods = osuapi.ParseMods(mods).String()
			}
		}

		// Calculate SR and PP
		starRating, ppSS, pp99, pp98, pp97, pp95 := osutools.BeatmapCache(mods, beatmap, cache)

		// Create embed
		embed := &discordgo.MessageEmbed{
			Author: &discordgo.MessageEmbedAuthor{
				URL:     "https://osu.ppy.sh/beatmaps/" + strconv.Itoa(beatmap.BeatmapID),
				Name:    beatmap.Artist + " - " + beatmap.Title + " by " + mapper.Username,
				IconURL: "https://a.ppy.sh/" + strconv.Itoa(mapper.UserID),
			},
			Color: Color,
			Description: starRating + length + bpm + combo + "\n" +
				mapStats + "\n" +
				status + "\n" +
				download + "\n" +
				diffs + "\n" + "\n" +
				"**[" + beatmap.DiffName + "]** with mods: " + mods + "\n" +
				//aimRating + speedRating + totalRating + "\n" + TODO: Make SR calc work
				ppSS + pp99 + pp98 + pp97 + pp95,
			Thumbnail: &discordgo.MessageEmbedThumbnail{
				URL: "https://b.ppy.sh/thumb/" + strconv.Itoa(beatmap.BeatmapSetID) + "l.jpg",
			},
		}
		if beatmap.Title == "Crab Rave" {
			embed.Image = &discordgo.MessageEmbedImage{
				URL: "https://cdn.discordapp.com/emojis/510169818893385729.gif",
			}
		}
		s.ChannelMessageEdit(message.ChannelID, message.ID, "")
		s.ChannelMessageEditEmbed(message.ChannelID, message.ID, embed)
		return
	}
}
