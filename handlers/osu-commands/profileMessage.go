package osucommands

import (
	"math/rand"
	"regexp"
	"strconv"
	"strings"
	"time"

	config "../../config"
	osuapi "../../osu-api"
	osutools "../../osu-functions"
	structs "../../structs"
	tools "../../tools"
	"github.com/bwmarrin/discordgo"
)

// ProfileMessage gets the information for the specified profile linked
func ProfileMessage(s *discordgo.Session, m *discordgo.MessageCreate, profileRegex *regexp.Regexp, cache []structs.PlayerData) {
	profileCmd1Regex, _ := regexp.Compile(`osu(top|detail)?\s+(.+)`)
	profileCmd2Regex, _ := regexp.Compile(`profile\s+(.+)`)
	profileCmd3Regex, _ := regexp.Compile(`osutop`)
	profileCmd4Regex, _ := regexp.Compile(`osudetail`)
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
	cmdMode := "command"
	if profileRegex.MatchString(m.Content) {
		value = profileRegex.FindStringSubmatch(m.Content)[3]
		cmdMode = "link"
	} else if profileCmd2Regex.MatchString(m.Content) {
		value = profileCmd2Regex.FindStringSubmatch(m.Content)[1]
	} else if profileCmd1Regex.MatchString(m.Content) {
		value = profileCmd1Regex.FindStringSubmatch(m.Content)[2]
	}

	if value == "" {
		for _, player := range cache {
			if m.Author.ID == player.Discord.ID && player.Osu.Username != "" {
				value = player.Osu.Username
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
		user, err = OsuAPI.GetUser(osuapi.GetUserOpts{
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
		user, err = OsuAPI.GetUser(osuapi.GetUserOpts{
			UserID: id,
			Mode:   mode,
		})
		if err != nil {
			user, err = OsuAPI.GetUser(osuapi.GetUserOpts{
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

	// Create embed
	embed := &discordgo.MessageEmbed{
		Author: &discordgo.MessageEmbedAuthor{
			URL:     "https://osu.ppy.sh/users/" + strconv.Itoa(user.UserID),
			Name:    user.Username + " (" + strconv.Itoa(user.UserID) + ")",
			IconURL: "https://osu.ppy.sh/images/flags/" + user.Country + ".png",
		},
		Thumbnail: &discordgo.MessageEmbedThumbnail{
			URL: "https://a.ppy.sh/" + strconv.Itoa(user.UserID) + "?" + strconv.Itoa(rand.Int()),
		},
		Color: osutools.ModeColour(mode),
	}

	// Get tops / details if asked
	totalHits := user.Count50 + user.Count100 + user.Count300
	g, _ := s.Guild(config.Conf.Server)
	if profileCmd3Regex.MatchString(m.Content) {
		// Get the user's best scores
		userBest, err := OsuAPI.GetUserBest(osuapi.GetUserScoresOpts{
			UserID: user.UserID,
			Mode:   mode,
		})
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, "The osu! API just owned me. Please try again!")
			return
		}

		var mapList []*discordgo.MessageEmbedField
		for i := 0; i < 5; i++ {
			score := userBest[i]

			beatmaps, err := OsuAPI.GetBeatmaps(osuapi.GetBeatmapsOpts{
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
				Name: beatmap.Artist + " - " + beatmap.Title + " [" + beatmap.DiffName + "] **+" + mods + "**",
				Value: "[**Link**](https://osu.ppy.sh/beatmaps/" + strconv.Itoa(beatmap.BeatmapID) + ") | <osu://dl/" + strconv.Itoa(beatmap.BeatmapSetID) + ">\n" +
					"**PP:** " + strconv.FormatFloat(score.PP, 'f', 2, 64) + " " + scoreRank + "\n" +
					"**Acc:** " + strconv.FormatFloat(accCalc, 'f', 2, 64) + "%\n" +
					"**Score:** " + strconv.FormatInt(score.Score.Score, 10) + "\n" +
					"**Combo:** " + strconv.Itoa(score.MaxCombo) + "/" + strconv.Itoa(beatmap.MaxCombo) + "x\n",
			}

			mapList = append(mapList, mapField)
		}
		embed.Description += "**Top plays:**" + "\n" + `\_\_\_\_\_\_\_\_\_\_`
		embed.Fields = mapList
	} else if profileCmd4Regex.MatchString(m.Content) && totalHits != 0 {
		// Get the user's recent scores
		userRecent, err := OsuAPI.GetUserRecent(osuapi.GetUserScoresOpts{
			UserID: user.UserID,
			Mode:   mode,
		})
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, "The osu! API just owned me. Please try again!")
			return
		}
		scoreCount := 0
		for _, score := range userRecent {
			if score.Rank != "F" {
				scoreCount++
			}
		}

		pp := "**PP:** " + strconv.FormatFloat(user.PP, 'f', 2, 64)
		rank := "**Rank:** #" + strconv.Itoa(user.Rank) + " (" + user.Country + "#" + strconv.Itoa(user.CountryRank) + ")"

		percent50 := float64(user.Count50) / float64(user.Count50+user.Count100+user.Count300)
		percent100 := float64(user.Count100) / float64(user.Count50+user.Count100+user.Count300)
		percent300 := float64(user.Count300) / float64(user.Count50+user.Count100+user.Count300)
		scoreRank := osutools.ScoreRank(percent50, percent300, 0, false)
		accuracy := "**Acc:** " + strconv.FormatFloat(user.Accuracy, 'f', 2, 64) + "%"
		count50 := "**50:** " + tools.Comma(int64(user.Count50)) + " (" + strconv.FormatFloat(percent50*100, 'f', 2, 64) + "%)"
		count100 := "**100:** " + tools.Comma(int64(user.Count100)) + " (" + strconv.FormatFloat(percent100*100, 'f', 2, 64) + "%)"
		count300 := "**300:** " + tools.Comma(int64(user.Count300)) + " (" + strconv.FormatFloat(percent300*100, 'f', 2, 64) + "%)"
		for _, emoji := range g.Emojis {
			if emoji.Name == scoreRank+"_" {
				accuracy += emoji.MessageFormat()
			}
		}

		pc := "**Playcount:** " + tools.Comma(int64(user.Playcount))
		timePlayed := "**Time played:** " + tools.TimeSince(time.Now().Add(time.Duration(-user.TimePlayed)*time.Second))
		hitsperPlay := "**Hits/play:** " + strconv.FormatFloat(float64(totalHits)/float64(user.Playcount), 'f', 2, 64)
		timePlayed = strings.Replace(timePlayed, "ago.", "", -1)

		level := "**Level:** " + strconv.FormatFloat(user.Level, 'f', 2, 64)
		rankedScore := "**Ranked Score:** " + tools.Comma(user.RankedScore)
		totalScore := "**Total Score:** " + tools.Comma(user.TotalScore)
		recentPlays := "**Recent Plays:** " + strconv.Itoa(len(userRecent))
		fullPlays := "**Recent Full Plays:** " + strconv.Itoa(scoreCount)

		ssh := "**SSH:** " + strconv.Itoa(user.CountSSH)
		ss := "**SS:** " + strconv.Itoa(user.CountSS)
		sh := "**SH:** " + strconv.Itoa(user.CountSH)
		s := "**S:** " + strconv.Itoa(user.CountS)
		a := "**A:** " + strconv.Itoa(user.CountA)
		embed.Fields = []*discordgo.MessageEmbedField{
			&discordgo.MessageEmbedField{
				Name:  "Placing",
				Value: pp + "\n" + rank,
			},
			&discordgo.MessageEmbedField{
				Name:  "Accuracy",
				Value: accuracy + "\n" + count50 + "\n" + count100 + "\n" + count300,
			},
			&discordgo.MessageEmbedField{
				Name:  "Playtime",
				Value: pc + "\n" + timePlayed + "\n" + hitsperPlay,
			},
			&discordgo.MessageEmbedField{
				Name:  "Score",
				Value: level + "\n" + rankedScore + "\n" + totalScore + "\n" + recentPlays + "\n" + fullPlays,
			},
			&discordgo.MessageEmbedField{
				Name:  "Ranks",
				Value: ssh + " " + ss + "\n" + sh + " " + s + "\n" + a,
			},
		}
	} else {
		pp := "**PP:** " + strconv.FormatFloat(user.PP, 'f', 2, 64)
		rank := "**Rank:** #" + strconv.Itoa(user.Rank) + " (" + user.Country + "#" + strconv.Itoa(user.CountryRank) + ")"
		accuracy := "**Acc:** " + strconv.FormatFloat(user.Accuracy, 'f', 2, 64) + "%"
		pc := "**Playcount:** " + tools.Comma(int64(user.Playcount))
		timePlayed := "**Time played:** " + tools.TimeSince(time.Now().Add(time.Duration(-user.TimePlayed)*time.Second))
		timePlayed = strings.Replace(timePlayed, "ago.", "", -1)
		embed.Description = pp + "\n" +
			rank + "\n" +
			accuracy + "\n" +
			pc + "\n" +
			timePlayed
	}

	s.ChannelMessageSendEmbed(m.ChannelID, embed)
	return
}
