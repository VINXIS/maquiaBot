package osucommands

import (
	"regexp"
	"strconv"
	"strings"

	osutools "../../osu-functions"
	structs "../../structs"
	tools "../../tools"
	"github.com/bwmarrin/discordgo"
	"github.com/thehowl/go-osuapi"
)

// ProfileMessage gets the information for the specified profile linked
func ProfileMessage(s *discordgo.Session, m *discordgo.MessageCreate, profileRegex *regexp.Regexp, osuAPI *osuapi.Client, cache []structs.MapData) {
	modeRegex, _ := regexp.Compile(`-m (\S+)`)
	bestRegex, _ := regexp.Compile(` -b`)
	var mode osuapi.Mode

	if modeRegex.MatchString(m.Content) {
		mValue := strings.ToLower(modeRegex.FindStringSubmatch(m.Content)[1])
		switch mValue {
		case "osu", "osu!", "std", "osu!std", "0", "o":
			mode = osuapi.ModeOsu
		case "taiko", "tko", "1", "t":
			mode = osuapi.ModeTaiko
		case "catch", "ctb", "2", "c":
			mode = osuapi.ModeCatchTheBeat
		case "mania", "man", "3", "m":
			mode = osuapi.ModeOsuMania
		default:
			mode = osuapi.ModeOsu
		}
	} else {
		mode = osuapi.ModeOsu
	}

	// Obtain username/user ID and assign variables
	value := profileRegex.FindStringSubmatch(m.Content)[3]
	id, err := strconv.Atoi(value)
	user := &osuapi.User{}

	// Get user
	if err != nil {
		user, err = osuAPI.GetUser(osuapi.GetUserOpts{
			Username: value,
			Mode:     mode,
		})
		if err != nil {
			return
		}
	} else {
		user, err = osuAPI.GetUser(osuapi.GetUserOpts{
			UserID: id,
			Mode:   mode,
		})
		if err != nil {
			return
		}
	}

	// Get the user's best scores
	userBest, err := osuAPI.GetUserBest(osuapi.GetUserScoresOpts{
		UserID: user.UserID,
		Mode:   mode,
	})
	tools.ErrRead(err)

	// Assign embed values
	Color := osutools.ModeColour(osuapi.ModeOsu)
	PP := "**PP:** " + strconv.FormatFloat(user.PP, 'f', 2, 64) + " "
	rank := "**Rank:** #" + strconv.Itoa(user.Rank) + " (" + user.Country + "#" + strconv.Itoa(user.CountryRank) + ")"
	accuracy := "**Acc:** " + strconv.FormatFloat(user.Accuracy, 'f', 2, 64) + "% "
	pc := "**Playcount:** " + strconv.Itoa(user.Playcount)
	topPlayFooter := ""

	var mapList []*discordgo.MessageEmbedField
	if bestRegex.MatchString(m.Content) {
		topPlayFooter = "**Top plays:**" + "\n" + `\_\_\_\_\_\_\_\_\_\_`
		for i := 0; i < 5; i++ {
			score := userBest[i]
			beatmaps, err := osuAPI.GetBeatmaps(osuapi.GetBeatmapsOpts{
				BeatmapID: score.BeatmapID,
			})
			beatmap := beatmaps[0]

			mapField := &discordgo.MessageEmbedField{
				Name: beatmap.Artist + " - " + beatmap.Title,
				Value: "**Link:** https://osu.ppy.sh/beatmaps/" + strconv.Itoa(beatmap.BeatmapID) + "\n" +
					"**Score:** " + strconv.FormatInt(score.Score.Score, 10) + "\n" +
					"**Combo:** " + strconv.Itoa(score.MaxCombo) + "/" + strconv.Itoa(beatmap.MaxCombo) + "x" + "\n" +
					"**PP:** " + strconv.FormatFloat(score.PP, 'f', 2, 64) + "\n" +
					`\_\_\_\_\_\_\_\_\_\_\_\_\_\_\_\_\_\_\_\_\_\_\_\_\_\_\_\_\_\_\_\_\_\_\_\_\_\_\_\_\_\_\_\_\_\_\_\_\_\_\_\_\_\_\_\_\_\_\_\_`,
			}

			mapList = append(mapList, mapField)
			tools.ErrRead(err)

		}
	}

	// Create embed
	embed := &discordgo.MessageEmbed{
		Author: &discordgo.MessageEmbedAuthor{
			URL:     "https://osu.ppy.sh/users/" + strconv.Itoa(user.UserID),
			Name:    user.Username,
			IconURL: "https://osu.ppy.sh/images/flags/" + user.Country + ".png",
		},
		Color: Color,
		Description: PP + "\n" +
			rank + "\n" +
			accuracy + "\n" +
			pc + "\n" +
			topPlayFooter,
		Thumbnail: &discordgo.MessageEmbedThumbnail{
			URL: "https://a.ppy.sh/" + strconv.Itoa(user.UserID),
		},
		Fields: mapList,
	}
	s.ChannelMessageSendEmbed(m.ChannelID, embed)
	return
}
