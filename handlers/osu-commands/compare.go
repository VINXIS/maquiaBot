package osucommands

import (
	"fmt"
	"math"
	"math/rand"
	"regexp"
	"sort"
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

// Compare compares finds a score from the current user on the previous map linked by the bot
func Compare(s *discordgo.Session, m *discordgo.MessageCreate, cache []structs.PlayerData) {
	mapRegex, _ := regexp.Compile(`(https:\/\/)?(osu|old)\.ppy\.sh\/(s|b|beatmaps|beatmapsets)\/(\d+)(#(osu|taiko|fruits|mania)\/(\d+))?`)
	modRegex, _ := regexp.Compile(`-m\s*(\S+)`)
	compareRegex, _ := regexp.Compile(`(c|compare)\s*(.+)?`)
	strictRegex, _ := regexp.Compile(`-nostrict`)
	allRegex, _ := regexp.Compile(`-all`)
	scorePostRegex, _ := regexp.Compile(`-sp`)

	// Obtain username and mods
	username := ""
	mods := "NM"
	parsedMods := osuapi.Mods(0)
	strict := true
	if compareRegex.MatchString(m.Content) {
		username = compareRegex.FindStringSubmatch(m.Content)[2]
		if modRegex.MatchString(username) {
			mods = strings.ToUpper(modRegex.FindStringSubmatch(username)[1])
			if strings.Contains(mods, "NC") && !strings.Contains(mods, "DT") {
				mods += "DT"
			}
			parsedMods = osuapi.ParseMods(mods)

			username = strings.TrimSpace(strings.Replace(username, modRegex.FindStringSubmatch(username)[0], "", 1))
		}
		if strictRegex.MatchString(m.Content) {
			strict = false
			username = strings.TrimSpace(strings.Replace(username, strictRegex.FindStringSubmatch(m.Content)[0], "", 1))
		}
		if allRegex.MatchString(m.Content) {
			username = strings.TrimSpace(strings.Replace(username, allRegex.FindStringSubmatch(m.Content)[0], "", 1))
		}
		if scorePostRegex.MatchString(m.Content) {
			username = strings.TrimSpace(strings.Replace(username, scorePostRegex.FindStringSubmatch(m.Content)[0], "", 1))
		}
	}

	// Get the map
	var beatmap osuapi.Beatmap
	var submatches []string
	if mapRegex.MatchString(m.Content) {
		submatches = mapRegex.FindStringSubmatch(m.Content)
	} else {
		// Get prev messages
		messages, err := s.ChannelMessages(m.ChannelID, -1, "", "", "")
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, "No map to compare to!")
			return
		}

		// Look for a valid beatmap ID
		for _, msg := range messages {
			if len(msg.Embeds) > 0 && msg.Embeds[0].Author != nil {
				if mapRegex.MatchString(msg.Embeds[0].URL) {
					submatches = mapRegex.FindStringSubmatch(msg.Embeds[0].URL)
					break
				} else if mapRegex.MatchString(msg.Embeds[0].Author.URL) {
					submatches = mapRegex.FindStringSubmatch(msg.Embeds[0].Author.URL)
					break
				}
			} else if mapRegex.MatchString(msg.Content) {
				submatches = mapRegex.FindStringSubmatch(msg.Content)
				break
			}
		}
	}

	// Check if found
	if len(submatches) == 0 {
		s.ChannelMessageSend(m.ChannelID, "No map to compare to!")
		return
	}

	// Get the map
	nomod := osuapi.Mods(0)
	switch submatches[3] {
	case "s":
		beatmap = osutools.BeatmapParse(submatches[4], "set", &nomod)
	case "b":
		beatmap = osutools.BeatmapParse(submatches[4], "map", &nomod)
	case "beatmaps":
		beatmap = osutools.BeatmapParse(submatches[4], "map", &nomod)
	case "beatmapsets":
		if len(submatches[7]) > 0 {
			beatmap = osutools.BeatmapParse(submatches[7], "map", &nomod)
		} else {
			beatmap = osutools.BeatmapParse(submatches[4], "set", &nomod)
		}
	}
	if beatmap.BeatmapID == 0 {
		s.ChannelMessageSend(m.ChannelID, "No map to compare to!")
		return
	} else if beatmap.Approved < 1 {
		s.ChannelMessageSend(m.ChannelID, "The map `"+beatmap.Artist+" - "+beatmap.Title+"` does not have a leaderboard!")
		return
	}
	username = strings.TrimSpace(strings.Replace(username, submatches[0], "", -1))

	// Get user
	var user osuapi.User
	for _, player := range cache {
		if username != "" {
			if username == player.Osu.Username {
				user = player.Osu
				break
			}
		} else if m.Author.ID == player.Discord.ID && player.Osu.Username != "" {
			user = player.Osu
			break
		}
	}

	// Check if user even exists
	if user.UserID == 0 {
		if username == "" {
			s.ChannelMessageSend(m.ChannelID, "No user mentioned in message/linked to your account! Please use `set` or `link` to link an osu! account to you, or name a user to obtain their recent score of!")
		}
		test, err := OsuAPI.GetUser(osuapi.GetUserOpts{
			Username: username,
		})
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, "User "+username+" may not exist! Are you sure you replaced spaces with `_`?")
			return
		}
		user = *test
	}

	// API call
	scoreOpts := osuapi.GetScoresOpts{
		BeatmapID: beatmap.BeatmapID,
		UserID:    user.UserID,
		Limit:     100,
	}
	scores, err := OsuAPI.GetScores(scoreOpts)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "The osu! API just owned me. Please try again!")
		return
	}
	if len(scores) == 0 {
		if username != "" {
			s.ChannelMessageSend(m.ChannelID, username+" hasn't set a score on this!")
		} else {
			s.ChannelMessageSend(m.ChannelID, "You haven't set a score on this with any mod combination!")
		}
		return
	}

	// Mod filter
	if mods != "NM" {
		for i := 0; i < len(scores); i++ {
			if (strict && scores[i].Mods != parsedMods) || (!strict && ((parsedMods == 0 && scores[i].Mods != 0) || scores[i].Mods&parsedMods != parsedMods)) {
				scores = append(scores[:i], scores[i+1:]...)
				i--
			}
		}
		if len(scores) == 0 {
			s.ChannelMessageSend(m.ChannelID, "No scores with the mod combination **"+mods+"** exist!")
			return
		}
	}
	topScore := scores[0]

	// Sort by PP
	sort.Slice(scores, func(i, j int) bool {
		return scores[i].PP > scores[j].PP
	})

	// Get the beatmap but with mods applied if not all
	if !allRegex.MatchString(m.Content) {
		diffMods := 338 & scores[0].Mods
		if diffMods&256 != 0 && diffMods&64 != 0 { // Remove DTHT
			diffMods -= 320
		}
		if diffMods&2 != 0 && diffMods&16 != 0 { // Remove EZHR
			diffMods -= 18
		}
		beatmap = osutools.BeatmapParse(strconv.Itoa(beatmap.BeatmapID), "map", &diffMods)
	}

	// Create embed
	// Assign timing variables for map specs
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
	sr := "**SR:** " + strconv.FormatFloat(beatmap.DifficultyRating, 'f', 2, 64) + " **Aim:** " + strconv.FormatFloat(beatmap.DifficultyAim, 'f', 2, 64) + " **Speed:** " + strconv.FormatFloat(beatmap.DifficultySpeed, 'f', 2, 64)
	length := "**Length:** " + fmt.Sprint(totalMinutes) + ":" + fmt.Sprint(totalSeconds) + " (" + fmt.Sprint(hitMinutes) + ":" + fmt.Sprint(hitSeconds) + ") "
	bpm := "**BPM:** " + fmt.Sprint(beatmap.BPM) + " "
	mapStats := "**CS:** " + strconv.FormatFloat(beatmap.CircleSize, 'f', 1, 64) + " **AR:** " + strconv.FormatFloat(beatmap.ApproachRate, 'f', 1, 64) + " **OD:** " + strconv.FormatFloat(beatmap.OverallDifficulty, 'f', 1, 64) + " **HP:** " + strconv.FormatFloat(beatmap.HPDrain, 'f', 1, 64)
	mapObjs := "**Circles:** " + strconv.Itoa(beatmap.Circles) + " **Sliders:** " + strconv.Itoa(beatmap.Sliders) + " **Spinners:** " + strconv.Itoa(beatmap.Spinners)
	Color := osutools.ModeColour(beatmap.Mode)

	embed := &discordgo.MessageEmbed{
		Color: Color,
		Author: &discordgo.MessageEmbedAuthor{
			URL:     "https://osu.ppy.sh/users/" + strconv.Itoa(user.UserID),
			Name:    user.Username,
			IconURL: "https://a.ppy.sh/" + strconv.Itoa(user.UserID) + "?" + strconv.Itoa(rand.Int()) + ".jpeg",
		},
		Description: sr + "\n" +
			length + bpm + "\n" +
			mapStats + "\n" +
			mapObjs + "\n\n",
		Title: beatmap.Artist + " - " + beatmap.Title + " [" + beatmap.DiffName + "] by " + strings.Replace(beatmap.Creator, "_", `\_`, -1),
		URL:   "https://osu.ppy.sh/beatmaps/" + strconv.Itoa(beatmap.BeatmapID),
		Thumbnail: &discordgo.MessageEmbedThumbnail{
			URL: "https://b.ppy.sh/thumb/" + strconv.Itoa(beatmap.BeatmapSetID) + "l.jpg",
		},
	}
	if strings.ToLower(beatmap.Title) == "crab rave" {
		embed.Image = &discordgo.MessageEmbedImage{
			URL: "https://cdn.discordapp.com/emojis/510169818893385729.gif",
		}
	}
	for i := 0; i < len(scores); i++ {
		score := scores[i]

		// Get time since play
		timeParse, _ := time.Parse("2006-01-02 15:04:05", score.Date.String())
		time := tools.TimeSince(timeParse)

		// Assign values
		mods = score.Mods.String()
		accCalc := (50.0*float64(score.Count50) + 100.0*float64(score.Count100) + 300.0*float64(score.Count300)) / (300.0 * float64(score.CountMiss+score.Count50+score.Count100+score.Count300)) * 100.0
		scorePrint := " **" + tools.Comma(score.Score.Score) + "** "
		acc := "** " + strconv.FormatFloat(accCalc, 'f', 2, 64) + "%** "
		hits := "**Hits:** [" + strconv.Itoa(score.Count300) + "/" + strconv.Itoa(score.Count100) + "/" + strconv.Itoa(score.Count50) + "/" + strconv.Itoa(score.CountMiss) + "]"

		replay := ""
		if score.Replay {
			replay = "| [**Replay**](https://osu.ppy.sh/scores/osu/" + strconv.FormatInt(score.ScoreID, 10) + "/download)"
		}

		if mods == "" {
			mods = "NM"
		}

		if strings.Contains(mods, "DTNC") {
			mods = strings.Replace(mods, "DTNC", "NC", 1)
		}

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

		mapCompletion := ""
		if i == 0 { // Only matters for the top pp score Lol
			orderedScores, err := OsuAPI.GetUserBest(osuapi.GetUserScoresOpts{
				Username: user.Username,
				Limit:    100,
			})
			if err != nil {
				s.ChannelMessageSend(m.ChannelID, "The osu! API just owned me. Please try again!")
				return
			}
			for i, orderedScore := range orderedScores {
				if score.Score.Score == orderedScore.Score.Score {
					mapCompletion += "**#" + strconv.Itoa(i+1) + "** in top performances! \n"
					break
				}
			}
		}
		if topScore.Score.Score == score.Score.Score { // Only matters for the top score Lol
			mapScores, err := OsuAPI.GetScores(osuapi.GetScoresOpts{
				BeatmapID: beatmap.BeatmapID,
				Limit:     100,
			})
			if err != nil {
				s.ChannelMessageSend(m.ChannelID, "The osu! API just owned me. Please try again!")
				return
			}
			for i, mapScore := range mapScores {
				if score.UserID == mapScore.UserID && score.Score.Score == mapScore.Score.Score {
					mapCompletion += "**#" + strconv.Itoa(i+1) + "** on leaderboard! \n"
					break
				}
			}
		}

		// Get pp values
		var pp string
		totalObjs := beatmap.Circles + beatmap.Sliders + beatmap.Spinners
		if score.Score.FullCombo { // If play was a perfect combo
			pp = "**" + strconv.FormatFloat(score.PP, 'f', 2, 64) + "pp**/" + strconv.FormatFloat(score.PP, 'f', 2, 64) + "pp "
		} else { // If map was finished, but play was not a perfect combo
			ppValues := make(chan string, 1)
			accCalcNoMiss := (50.0*float64(score.Count50) + 100.0*float64(score.Count100) + 300.0*float64(totalObjs-score.Count50-score.Count100)) / (300.0 * float64(totalObjs)) * 100.0
			go osutools.PPCalc(beatmap, accCalcNoMiss, "", "", mods, ppValues)
			pp = "**" + strconv.FormatFloat(score.PP, 'f', 2, 64) + "pp**/" + <-ppValues + "pp "
		}
		mods = " **+" + mods + "** "

		g, _ := s.Guild(config.Conf.Server)
		scoreRank := ""
		for _, emoji := range g.Emojis {
			if emoji.Name == score.Rank+"_" {
				scoreRank = emoji.MessageFormat()
			}
		}

		if !allRegex.MatchString(m.Content) || len(scores) == 1 {
			embed.Description += scoreRank + scorePrint + mods + combo + acc + replay + "\n" +
				mapCompletion + "\n" +
				pp + hits + "\n\n"
			embed.Footer = &discordgo.MessageEmbedFooter{
				Text: time,
			}
			message, err := s.ChannelMessageSendEmbed(m.ChannelID, embed)
			if scorePostRegex.MatchString(m.Content) && err == nil {
				ScorePost(s, &discordgo.MessageCreate{message}, cache, "")
			}
			return
		}
		embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
			Name: "#" + strconv.Itoa(i+1) + " | " + time,
			Value: scoreRank + scorePrint + mods + combo + acc + replay + "\n" +
				mapCompletion +
				pp + hits + "\n\n",
		})
		if (i+1)%25 == 0 {
			s.ChannelMessageSendEmbed(m.ChannelID, embed)
			embed.Fields = []*discordgo.MessageEmbedField{}
		}
	}
	if len(scores)%25 != 0 {
		s.ChannelMessageSendEmbed(m.ChannelID, embed)
	}
}
