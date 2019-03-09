package osucommands

import (
	"fmt"
	"log"
	"math"
	"regexp"
	"sort"
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
				beatmap = beatmapParse(submatches[7], "map", osuAPI)
			} else {
				beatmap = beatmapParse(submatches[4], "set", osuAPI)
			}
		} else {
			if submatches[2] == "s" {
				beatmap = beatmapParse(submatches[4], "set", osuAPI)
			} else {
				beatmap = beatmapParse(submatches[4], "map", osuAPI)
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
