package osucommands

import (
	"math/rand"
	"regexp"
	"strconv"
	"strings"

	osuapi "../../osu-api"
	osutools "../../osu-functions"
	structs "../../structs"
	tools "../../tools"
	"github.com/bwmarrin/discordgo"
)

// ProfileMessage gets the information for the specified profile linked
func ProfileMessage(s *discordgo.Session, m *discordgo.MessageCreate, profileRegex *regexp.Regexp, osuAPI *osuapi.Client, cache []structs.PlayerData) {
	modeRegex, _ := regexp.Compile(`-m (\S+)`)
	var mode osuapi.Mode
	imgMode := "0"

	if modeRegex.MatchString(m.Content) {
		mValue := strings.ToLower(modeRegex.FindStringSubmatch(m.Content)[1])
		switch mValue {
		case "osu", "osu!", "std", "osu!std", "0", "o":
			imgMode = "0"
			mode = osuapi.ModeOsu
		case "taiko", "tko", "1", "t":
			imgMode = "1"
			mode = osuapi.ModeTaiko
		case "catch", "ctb", "2", "c":
			imgMode = "2"
			mode = osuapi.ModeCatchTheBeat
		case "mania", "man", "3", "m":
			imgMode = "3"
			mode = osuapi.ModeOsuMania
		default:
			imgMode = "0"
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
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "The osu! API just owned me. Please try again!")
		return
	}

	// Assign embed values
	Color := osutools.ModeColour(osuapi.ModeOsu)
	PP := "**PP:** " + strconv.FormatFloat(user.PP, 'f', 2, 64) + " "
	rank := "**Rank:** #" + strconv.Itoa(user.Rank) + " (" + user.Country + "#" + strconv.Itoa(user.CountryRank) + ")"
	accuracy := "**Acc:** " + strconv.FormatFloat(user.Accuracy, 'f', 2, 64) + "% "
	pc := "**Playcount:** " + strconv.Itoa(user.Playcount)
	imgURL := "https://lemmmy.pw/osusig/sig.php?uname=" + strconv.Itoa(user.UserID) + "&pp=2&mode=" + imgMode
	topPlayFooter := ""

	var mapList []*discordgo.MessageEmbedField
	if strings.Contains(m.Content, " -b") {
		imgURL = ""
		topPlayFooter = "**Top plays:**" + "\n" + `\_\_\_\_\_\_\_\_\_\_`
		for i := 0; i < 5; i++ {
			score := userBest[i]
			beatmaps, err := osuAPI.GetBeatmaps(osuapi.GetBeatmapsOpts{
				BeatmapID: score.BeatmapID,
			})
			if err != nil {
				s.ChannelMessageSend(m.ChannelID, "The osu! API just owned me. Please try again!")
				return
			}
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
			URL: "https://a.ppy.sh/" + strconv.Itoa(user.UserID) + "?" + strconv.Itoa(rand.Int()) + ".jpeg",
		},
		Fields: mapList,
		Image: &discordgo.MessageEmbedImage{
			URL: imgURL,
		},
	}
	s.ChannelMessageSendEmbed(m.ChannelID, embed)
	return
}
