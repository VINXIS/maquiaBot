package osucommands

import (
	"bytes"
	"math/rand"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	config "../../config"
	osuapi "../../osu-api"
	osutools "../../osu-tools"
	structs "../../structs"
	tools "../../tools"
	"github.com/bwmarrin/discordgo"
	"github.com/wcharczuk/go-chart"
)

// ProfileMessage gets the information for the specified profile linked
func ProfileMessage(s *discordgo.Session, m *discordgo.MessageCreate, profileRegex *regexp.Regexp, cache []structs.PlayerData) {
	profileCmd1Regex, _ := regexp.Compile(`osu(top|detail)?\s+(.+)`)
	profileCmd2Regex, _ := regexp.Compile(`profile\s+(.+)`)
	profileCmd3Regex, _ := regexp.Compile(`osutop`)
	profileCmd4Regex, _ := regexp.Compile(`osudetail`)
	modeRegex, _ := regexp.Compile(`-m\s+(.+)`)
	recentRegex, _ := regexp.Compile(`-r(ecent)?`)
	mode := osuapi.ModeOsu

	mValue := ""
	if modeRegex.MatchString(m.Content) {
		mValue = strings.ToLower(modeRegex.FindStringSubmatch(m.Content)[1])
		switch mValue {
		case "taiko", "tko", "1", "t":
			mode = osuapi.ModeTaiko
		case "catch", "ctb", "2", "c":
			mode = osuapi.ModeCatchTheBeat
		case "mania", "man", "3", "m":
			mode = osuapi.ModeOsuMania
		}
		mValue = strings.ToLower(modeRegex.FindStringSubmatch(m.Content)[0])
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

	if recentRegex.MatchString(m.Content) {
		value = strings.TrimSpace(strings.Replace(value, recentRegex.FindStringSubmatch(m.Content)[0], "", -1))
	}
	value = strings.TrimSpace(strings.Replace(value, mValue, "", -1))

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
	buffer := bytes.NewBuffer([]byte{})
	if profileCmd3Regex.MatchString(m.Content) {
		// Get the user's best scores
		userBest, err := OsuAPI.GetUserBest(osuapi.GetUserScoresOpts{
			UserID: user.UserID,
			Mode:   mode,
			Limit:  100,
		})
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, "The osu! API just owned me. Please try again!")
			return
		} else if len(userBest) == 0 {
			s.ChannelMessageSend(m.ChannelID, "This user has no top scores!")
			return
		}

		userRecent := userBest
		if recentRegex.MatchString(m.Content) {
			// Sort scores by date and get score
			sort.Slice(userRecent, func(i, j int) bool {
				time1, err := time.Parse("2006-01-02 15:04:05", userRecent[i].Date.String())
				tools.ErrRead(s, err)
				time2, err := time.Parse("2006-01-02 15:04:05", userRecent[j].Date.String())
				tools.ErrRead(s, err)

				return time1.Unix() > time2.Unix()
			})
		}

		amount := 5
		if amount > len(userBest) {
			amount = len(userBest)
		}
		var mapList []*discordgo.MessageEmbedField
		for i := 0; i < amount; i++ {
			score := userBest[i]
			if recentRegex.MatchString(m.Content) {
				score = userRecent[i]
			}

			beatmaps, err := OsuAPI.GetBeatmaps(osuapi.GetBeatmapsOpts{
				BeatmapID: score.BeatmapID,
			})
			if err != nil {
				s.ChannelMessageSend(m.ChannelID, "The osu! API just owned me. Please try again!")
				return
			}
			beatmap := beatmaps[0]

			score.Rank = strings.Replace(score.Rank, "X", "SS", -1)
			scoreRank := ""
			for _, emoji := range g.Emojis {
				if emoji.Name == score.Rank+"_" {
					scoreRank = emoji.MessageFormat()
				}
			}
			scorePrint := " **" + tools.Comma(score.Score.Score) + "** "
			mods := score.Mods.String()
			if strings.Contains(mods, "DTNC") {
				mods = strings.Replace(mods, "DTNC", "NC", -1)
			}
			accCalc := 100.0 * float64(score.Count50+2*score.Count100+6*score.Count300) / float64(6*(score.CountMiss+score.Count50+score.Count100+score.Count300))
			var combo string
			if score.MaxCombo == beatmap.MaxCombo {
				if accCalc == 100.0 {
					combo = " **SS** "
				} else {
					combo = " **FC** "
				}
			} else {
				combo = " **" + strconv.Itoa(score.MaxCombo) + "**/" + strconv.Itoa(beatmap.MaxCombo) + "x "
			}
			acc := "** " + strconv.FormatFloat(accCalc, 'f', 2, 64) + "%** "
			replay := ""
			if score.Replay {
				replay = "| [**Replay**](https://osu.ppy.sh/scores/osu/" + strconv.FormatInt(score.ScoreID, 10) + "/download)"
				reader, _ := OsuAPI.GetReplay(osuapi.GetReplayOpts{
					Username:  user.Username,
					Mode:      beatmap.Mode,
					BeatmapID: beatmap.BeatmapID,
					Mods:      &score.Mods,
				})
				buf := new(bytes.Buffer)
				buf.ReadFrom(reader)
				replayData := structs.ReplayData{
					Mode:    beatmap.Mode,
					Beatmap: beatmap,
					Score:   score.Score,
					Player:  *user,
					Data:    buf.Bytes(),
				}
				replayData.PlayData = replayData.GetPlayData(true)
				UR := replayData.GetUnstableRate()
				replay += " | " + strconv.FormatFloat(UR, 'f', 2, 64)
				if strings.Contains(mods, "DT") || strings.Contains(mods, "NC") || strings.Contains(mods, "HT") {
					replay += " cv. UR"
				} else {
					replay += " UR"
				}
			}
			var pp string
			totalObjs := beatmap.Circles + beatmap.Sliders + beatmap.Spinners
			if score.Score.FullCombo { // If play was a perfect combo
				pp = "**" + strconv.FormatFloat(score.PP, 'f', 2, 64) + "pp**/" + strconv.FormatFloat(score.PP, 'f', 2, 64) + "pp "
			} else { // If map was finished, but play was not a perfect combo
				ppValues := make(chan string, 1)
				go osutools.PPCalc(beatmap, osuapi.Score{
					MaxCombo: beatmap.MaxCombo,
					Count50:  score.Count50,
					Count100: score.Count100,
					Count300: totalObjs - score.Count50 - score.Count100,
					Mods:     score.Mods,
				}, ppValues)
				pp = "**" + strconv.FormatFloat(score.PP, 'f', 2, 64) + "pp**/" + <-ppValues + "pp "
			}
			hits := "[" + strconv.Itoa(score.Count300) + "/" + strconv.Itoa(score.Count100) + "/" + strconv.Itoa(score.Count50) + "/" + strconv.Itoa(score.CountMiss) + "]"
			timeParse, _ := time.Parse("2006-01-02 15:04:05", score.Date.String())
			scoreTime := tools.TimeSince(timeParse)
			mods = " **+" + mods + "** "

			mapField := &discordgo.MessageEmbedField{
				Name: "#" + strconv.Itoa(i+1) + " " + beatmap.Artist + " - " + beatmap.Title + " [" + beatmap.DiffName + "]",
				Value: "[**Link**](https://osu.ppy.sh/beatmaps/" + strconv.Itoa(beatmap.BeatmapID) + ") | <osu://dl/" + strconv.Itoa(beatmap.BeatmapSetID) + ">\n" +
					scoreRank + scorePrint + mods + combo + acc + replay + "\n" +
					pp + hits + "\n" +
					scoreTime,
			}
			if recentRegex.MatchString(m.Content) {
				// Sort userBest back to original
				sort.Slice(userBest, func(i, j int) bool { return userBest[i].PP > userBest[j].PP })
				for j, bestScore := range userBest {
					if score.BeatmapID == bestScore.BeatmapID {
						mapField.Name = "#" + strconv.Itoa(j+1) + " " + beatmap.Artist + " - " + beatmap.Title + " [" + beatmap.DiffName + "]"
						break
					}
				}
				// Sort scores back to chronological
				sort.Slice(userRecent, func(i, j int) bool {
					time1, err := time.Parse("2006-01-02 15:04:05", userRecent[i].Date.String())
					tools.ErrRead(s, err)
					time2, err := time.Parse("2006-01-02 15:04:05", userRecent[j].Date.String())
					tools.ErrRead(s, err)

					return time1.Unix() > time2.Unix()
				})
			}

			mapList = append(mapList, mapField)
		}
		embed.Description += "**Top plays:**" + "\n" + `\_\_\_\_\_\_\_\_\_\_`
		embed.Fields = mapList

		var x, y []float64
		for i, play := range userBest {
			x = append(x, float64(i+1))
			y = append(y, play.PP)
		}
		graph := chart.Chart{
			Series: []chart.Series{
				chart.ContinuousSeries{
					XValueFormatter: chart.IntValueFormatter,
					Style: chart.Style{
						StrokeColor: chart.GetDefaultColor(0).WithAlpha(64),
						FillColor:   chart.GetDefaultColor(0).WithAlpha(64),
					},
					XValues: x,
					YValues: y,
				},
			},
		}
		err = graph.Render(chart.PNG, buffer)
	} else if profileCmd4Regex.MatchString(m.Content) && totalHits != 0 {
		// Get the user's recent scores
		userRecent, err := OsuAPI.GetUserRecent(osuapi.GetUserScoresOpts{
			UserID: user.UserID,
			Mode:   mode,
			Limit:  100,
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
		embed.Image = &discordgo.MessageEmbedImage{
			URL: "https://osu.ppy.sh/pages/include/profile-graphactivity.php?_jpg&u=" + strconv.Itoa(user.UserID),
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
	if buffer.Len() != 0 {
		s.ChannelMessageSendComplex(m.ChannelID, &discordgo.MessageSend{
			Files: []*discordgo.File{
				&discordgo.File{
					Name:   "tops.png",
					Reader: buffer,
				},
			},
		})
	}
	return
}
