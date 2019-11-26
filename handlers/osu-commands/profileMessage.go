package osucommands

import (
	"math/rand"
	"regexp"
	"strconv"
	"strings"
	"time"

	osuapi "../../osu-api"
	osutools "../../osu-functions"
	structs "../../structs"
	"github.com/bwmarrin/discordgo"
)

// ProfileMessage gets the information for the specified profile linked
func ProfileMessage(s *discordgo.Session, m *discordgo.MessageCreate, profileRegex *regexp.Regexp, osuAPI *osuapi.Client, cache []structs.PlayerData) {
	profileCmd1Regex, _ := regexp.Compile(`osu\s+(.+)`)
	profileCmd2Regex, _ := regexp.Compile(`profile\s+(.+)`)
	modeRegex, _ := regexp.Compile(`-m\s+(.+)`)
	mode := osuapi.ModeOsu

	if modeRegex.MatchString(m.Content) {
		mValue := strings.ToLower(modeRegex.FindStringSubmatch(m.Content)[1])
		switch mValue {
		case "taiko", "tko", "1", "t":
			mode = osuapi.ModeTaiko
		case "catch", "ctb", "2", "c":
			mode = osuapi.ModeCatchTheBeat
		case "mania", "man", "3", "m":
			mode = osuapi.ModeOsuMania
		}
	}

	// Obtain username/user ID and assign variables
	value := ""
	cmdMode := "link"
	if profileRegex.MatchString(m.Content) {
		value = profileRegex.FindStringSubmatch(m.Content)[3]
	} else if profileCmd2Regex.MatchString(m.Content) {
		value = profileCmd2Regex.FindStringSubmatch(m.Content)[1]
		cmdMode = "command"
	} else if profileCmd1Regex.MatchString(m.Content) {
		value = profileCmd1Regex.FindStringSubmatch(m.Content)[1]
		cmdMode = "command"
	} else {
		for _, player := range cache {
			if m.Author.ID == player.Discord.ID && player.Osu.Username != "" {
				value = player.Osu.Username
				cmdMode = "command"
				break
			}
		}
		if value == "" {
			s.ChannelMessageSend(m.ChannelID, "Could not find any osu! account linked for "+m.Author.Mention()+" ! Please use `set` or `link` to link an osu! account to you!")
			return
		}
	}
	id, err := strconv.Atoi(value)
	user := &osuapi.User{}

	// Get user
	if err != nil {
		user, err = osuAPI.GetUser(osuapi.GetUserOpts{
			Username: value,
			Mode:     mode,
		})
		if err != nil {
			if cmdMode == "command" {
				s.ChannelMessageSend(m.ChannelID, "User not found!")
			}
			return
		}
	} else {
		user, err = osuAPI.GetUser(osuapi.GetUserOpts{
			UserID: id,
			Mode:   mode,
		})
		if err != nil {
			user, err = osuAPI.GetUser(osuapi.GetUserOpts{
				Username: value,
				Mode:     mode,
			})
			if err != nil {
				if cmdMode == "command" {
					s.ChannelMessageSend(m.ChannelID, "User not found!")
				}
				return
			}
		}
	}

	if user.UserID == 0 {
		return
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
	timePlayed := "**Time played:** " + (time.Duration(user.TimePlayed) * time.Second).String()
	topPlayFooter := ""

	var mapList []*discordgo.MessageEmbedField
	if strings.Contains(m.Content, "-t") {
		g, _ := s.Guild("556243477084635170")
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

			mods := score.Mods.String()
			if strings.Contains(mods, "DTNC") {
				mods = strings.Replace(mods, "DTNC", "NC", -1)
			}
			scoreRank := ""
			for _, emoji := range g.Emojis {
				if emoji.Name == score.Rank+"_" {
					scoreRank = emoji.MessageFormat()
				}
			}
			accCalc := (50.0*float64(score.Count50) + 100.0*float64(score.Count100) + 300.0*float64(score.Count300)) / (300.0 * float64(score.CountMiss+score.Count50+score.Count100+score.Count300)) * 100.0

			mapField := &discordgo.MessageEmbedField{
				Name: beatmap.Artist + " - " + beatmap.Title + " **+" + mods + "**",
				Value: "**Link:** https://osu.ppy.sh/beatmaps/" + strconv.Itoa(beatmap.BeatmapID) + "\n" +
					"**Score:** " + strconv.FormatInt(score.Score.Score, 10) + " " + scoreRank + "\n" +
					"**Acc:** " + strconv.FormatFloat(accCalc, 'f', 2, 64) + "%\n" +
					"**Combo:** " + strconv.Itoa(score.MaxCombo) + "/" + strconv.Itoa(beatmap.MaxCombo) + "x\n" +
					"**PP:** " + strconv.FormatFloat(score.PP, 'f', 2, 64) + "\n",
			}

			mapList = append(mapList, mapField)
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
			timePlayed + "\n" +
			topPlayFooter,
		Fields: mapList,
	}
	if len(mapList) > 0 {
		embed.Thumbnail = &discordgo.MessageEmbedThumbnail{
			URL: "https://a.ppy.sh/" + strconv.Itoa(user.UserID) + "?" + strconv.Itoa(rand.Int()),
		}
	} else {
		embed.Image = &discordgo.MessageEmbedImage{
			URL: "https://a.ppy.sh/" + strconv.Itoa(user.UserID) + "?" + strconv.Itoa(rand.Int()),
		}
	}

	s.ChannelMessageSendEmbed(m.ChannelID, embed)
	return
}
