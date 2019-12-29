package osucommands

import (
	"bytes"
	"fmt"
	"io/ioutil"
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

// Recent gets the most recent score done/nth score done
func Recent(s *discordgo.Session, m *discordgo.MessageCreate, option string, cache []structs.PlayerData) {
	recentRegex, _ := regexp.Compile(`(r|recent|rs|rb|recentb|recentbest)\s+(.+)`)
	modRegex, _ := regexp.Compile(`-m\s*(\S+)`)
	strictRegex, _ := regexp.Compile(`-nostrict`)
	scorePostRegex, _ := regexp.Compile(`-sp`)
	genOSR, _ := regexp.Compile(`-osr`)

	username := ""
	mods := ""
	index := 1
	strict := true

	// Obtain index, mods, strict, and username
	if recentRegex.MatchString(m.Content) {
		username = recentRegex.FindStringSubmatch(m.Content)[2]
		if modRegex.MatchString(username) {
			mods = strings.ToUpper(modRegex.FindStringSubmatch(username)[1])
			if strings.Contains(mods, "NC") && !strings.Contains(mods, "DT") {
				mods += "DT"
			}
			username = strings.TrimSpace(strings.Replace(username, modRegex.FindStringSubmatch(username)[0], "", 1))
		}
		usernameSplit := strings.Split(username, " ")
		for _, txt := range usernameSplit {
			if i, err := strconv.Atoi(txt); err == nil && i > 0 && i <= 100 {
				username = strings.TrimSpace(strings.Replace(username, txt, "", 1))
				index = i
				break
			}
		}
		if strictRegex.MatchString(m.Content) {
			strict = false
			username = strings.TrimSpace(strings.Replace(username, strictRegex.FindStringSubmatch(m.Content)[0], "", 1))
		}
		if scorePostRegex.MatchString(m.Content) {
			username = strings.TrimSpace(strings.Replace(username, scorePostRegex.FindStringSubmatch(m.Content)[0], "", 1))
		}
		if genOSR.MatchString(m.Content) {
			username = strings.TrimSpace(strings.Replace(username, genOSR.FindStringSubmatch(m.Content)[0], "", 1))
		}
	}

	// Get message author's osu! user if no user was specified
	if username == "" {
		for _, player := range cache {
			if m.Author.ID == player.Discord.ID && player.Osu.Username != "" {
				username = player.Osu.Username
				break
			}
		}
		if username == "" {
			s.ChannelMessageSend(m.ChannelID, "Could not find any osu! account linked for "+m.Author.Mention()+" ! Please use `set` or `link` to link an osu! account to you!")
			return
		}
	}

	// Get user
	userP, err := OsuAPI.GetUser(osuapi.GetUserOpts{
		Username: username,
	})
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "User **"+username+"** may not exist!")
		return
	}

	// Run api call for user best/recent
	var scoreList []osuapi.GUSScore
	switch option {
	case "recent":
		scoreList, err = OsuAPI.GetUserRecent(osuapi.GetUserScoresOpts{
			Username: username,
			Limit:    50,
		})
		tools.ErrRead(err)
		if len(scoreList) == 0 {
			s.ChannelMessageSend(m.ChannelID, username+" has not played recently!")
			return
		}
	case "best":
		scoreList, err = OsuAPI.GetUserBest(osuapi.GetUserScoresOpts{
			Username: username,
			Limit:    100,
		})
		tools.ErrRead(err)
		if len(scoreList) == 0 {
			s.ChannelMessageSend(m.ChannelID, username+" has no top scores!")
			return
		}
	}

	// Sort scores by date and get score
	sort.Slice(scoreList, func(i, j int) bool {
		time1, err := time.Parse("2006-01-02 15:04:05", scoreList[i].Date.String())
		tools.ErrRead(err)
		time2, err := time.Parse("2006-01-02 15:04:05", scoreList[j].Date.String())
		tools.ErrRead(err)

		return time1.Unix() > time2.Unix()
	})

	// Mod filter
	if mods != "" {
		parsedMods := osuapi.ParseMods(mods)
		for i := 0; i < len(scoreList); i++ {
			if (strict && scoreList[i].Mods != parsedMods) || (!strict && ((parsedMods == 0 && scoreList[i].Mods != 0) || scoreList[i].Mods&parsedMods != parsedMods)) {
				scoreList = append(scoreList[:i], scoreList[i+1:]...)
				i--
			}
		}
		if len(scoreList) == 0 {
			s.ChannelMessageSend(m.ChannelID, "No scores with the mod combination **"+mods+"** exist!")
			return
		}
	}

	warning := ""
	if index > len(scoreList) {
		index = len(scoreList)
		warning = "Defaulted to max: " + strconv.Itoa(len(scoreList))
	}
	score := scoreList[index-1]

	// Get beatmap, acc, and mods
	diffMods := osuapi.Mods(338) & score.Mods
	beatmap := osutools.BeatmapParse(strconv.Itoa(score.BeatmapID), "map", &diffMods)
	accCalc := (50.0*float64(score.Count50) + 100.0*float64(score.Count100) + 300.0*float64(score.Count300)) / (300.0 * float64(score.CountMiss+score.Count50+score.Count100+score.Count300)) * 100.0
	mods = score.Mods.String()

	// Count number of tries
	try := 0
	for i := index - 1; i < len(scoreList); i++ {
		if scoreList[i].BeatmapID == score.BeatmapID {
			try++
		} else {
			break
		}
	}

	// Count number of objs
	objCount := beatmap.Circles + beatmap.Sliders + beatmap.Spinners
	playObjCount := score.CountMiss + score.Count100 + score.Count300 + score.Count50

	// Get time since play
	timeParse, _ := time.Parse("2006-01-02 15:04:05", score.Date.String())

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

	// Assign misc variables
	Color := osutools.ModeColour(beatmap.Mode)
	sr := "**SR:** " + strconv.FormatFloat(beatmap.DifficultyRating, 'f', 2, 64) + " **Aim:** " + strconv.FormatFloat(beatmap.DifficultyAim, 'f', 2, 64) + " **Speed:** " + strconv.FormatFloat(beatmap.DifficultySpeed, 'f', 2, 64)
	length := "**Length:** " + fmt.Sprint(totalMinutes) + ":" + fmt.Sprint(totalSeconds) + " (" + fmt.Sprint(hitMinutes) + ":" + fmt.Sprint(hitSeconds) + ") "
	bpm := "**BPM:** " + fmt.Sprint(beatmap.BPM) + " "
	mapStats := "**CS:** " + strconv.FormatFloat(beatmap.CircleSize, 'f', 1, 64) + " **AR:** " + strconv.FormatFloat(beatmap.ApproachRate, 'f', 1, 64) + " **OD:** " + strconv.FormatFloat(beatmap.OverallDifficulty, 'f', 1, 64) + " **HP:** " + strconv.FormatFloat(beatmap.HPDrain, 'f', 1, 64)
	mapObjs := "**Circles:** " + strconv.Itoa(beatmap.Circles) + " **Sliders:** " + strconv.Itoa(beatmap.Sliders) + " **Spinners:** " + strconv.Itoa(beatmap.Spinners)
	scorePrint := " **" + tools.Comma(score.Score.Score) + "** "
	var combo string
	var mapCompletion string

	if strings.Contains(mods, "DTNC") {
		mods = strings.Replace(mods, "DTNC", "NC", 1)
	}

	if score.MaxCombo == beatmap.MaxCombo {
		if accCalc == 100.0 {
			combo = " **SS** "
		} else {
			combo = " **FC** "
		}
	} else {
		combo = " **" + strconv.Itoa(score.MaxCombo) + "**/" + strconv.Itoa(beatmap.MaxCombo) + "x "
	}

	if objCount != playObjCount {
		completed := float64(playObjCount) / float64(objCount) * 100.0
		mapCompletion = "**" + strconv.FormatFloat(completed, 'f', 2, 64) + "%** completed \n"
	} else {
		orderedScores, err := OsuAPI.GetUserBest(osuapi.GetUserScoresOpts{
			Username: userP.Username,
			Limit:    100,
		})
		tools.ErrRead(err)
		for i, orderedScore := range orderedScores {
			if score.Score.Score == orderedScore.Score.Score {
				mapCompletion += "**#" + strconv.Itoa(i+1) + "** in top performances! \n"
				break
			}
		}
		mapScores, err := OsuAPI.GetScores(osuapi.GetScoresOpts{
			BeatmapID: beatmap.BeatmapID,
			Limit:     100,
		})
		tools.ErrRead(err)
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
	if score.PP == 0 { // If map was not finished
		ppValues := make(chan string, 2)
		var ppValueArray [2]string
		accCalcNoMiss := (50.0*float64(score.Count50) + 100.0*float64(score.Count100) + 300.0*float64(totalObjs-score.Count50-score.Count100)) / (300.0 * float64(totalObjs)) * 100.0
		go osutools.PPCalc(beatmap, accCalcNoMiss, "", "", mods, ppValues)
		go osutools.PPCalc(beatmap, accCalc, strconv.Itoa(score.MaxCombo), strconv.Itoa(score.CountMiss), mods, ppValues)
		for v := 0; v < 2; v++ {
			ppValueArray[v] = <-ppValues
		}
		sort.Slice(ppValueArray[:], func(i, j int) bool {
			pp1, _ := strconv.ParseFloat(ppValueArray[i], 64)
			pp2, _ := strconv.ParseFloat(ppValueArray[j], 64)
			return pp1 > pp2
		})
		if objCount != playObjCount {
			pp = "~~**" + ppValueArray[1] + "pp**~~/" + ppValueArray[0] + "pp "
		} else {
			pp = "**" + ppValueArray[1] + "pp**/" + ppValueArray[0] + "pp "
		}
	} else if score.Score.FullCombo { // If play was a perfect combo
		pp = "**" + strconv.FormatFloat(score.PP, 'f', 2, 64) + "pp**/" + strconv.FormatFloat(score.PP, 'f', 2, 64) + "pp "
	} else { // If map was finished, but play was not a perfect combo
		ppValues := make(chan string, 1)
		accCalcNoMiss := (50.0*float64(score.Count50) + 100.0*float64(score.Count100) + 300.0*float64(totalObjs-score.Count50-score.Count100)) / (300.0 * float64(totalObjs)) * 100.0
		go osutools.PPCalc(beatmap, accCalcNoMiss, "", "", mods, ppValues)
		pp = "**" + strconv.FormatFloat(score.PP, 'f', 2, 64) + "pp**/" + <-ppValues + "pp "
	}
	acc := "** " + strconv.FormatFloat(accCalc, 'f', 2, 64) + "%** "
	hits := "**Hits:** [" + strconv.Itoa(score.Count300) + "/" + strconv.Itoa(score.Count100) + "/" + strconv.Itoa(score.Count50) + "/" + strconv.Itoa(score.CountMiss) + "]"
	mods = " **+" + mods + "** "

	replay := ""
	if score.Replay {
		replay = "| [**Replay**](https://osu.ppy.sh/scores/osu/" + strconv.FormatInt(score.ScoreID, 10) + "/download)"
		reader, _ := OsuAPI.GetReplay(osuapi.GetReplayOpts{
			Username:  userP.Username,
			Mode:      beatmap.Mode,
			BeatmapID: beatmap.BeatmapID,
			Mods: &score.Mods,
		})
		buf := new(bytes.Buffer)
		buf.ReadFrom(reader)
		replayData := structs.ReplayData{
			Mode:    beatmap.Mode,
			Beatmap: beatmap,
			Score:   score.Score,
			Player:  *userP,
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
		if genOSR.MatchString(m.Content) && m.Author.ID == config.Conf.BotHoster.UserID {
			fileContent := replayData.CreateOSR()
			ioutil.WriteFile("./"+userP.Username+strconv.Itoa(beatmap.BeatmapID)+strconv.Itoa(int(score.Mods))+".osr", fileContent, 0644)
		}
	} else if option == "recent" && objCount == playObjCount {
		replayScore, _ := OsuAPI.GetScores(osuapi.GetScoresOpts{
			BeatmapID: beatmap.BeatmapID,
			UserID:    userP.UserID,
			Mods:      &score.Mods,
		})
		if replayScore[0].Replay && replayScore[0].Score.Score == score.Score.Score {
			replay = "| [**Replay**](https://osu.ppy.sh/scores/osu/" + strconv.FormatInt(replayScore[0].ScoreID, 10) + "/download)"
			reader, _ := OsuAPI.GetReplay(osuapi.GetReplayOpts{
				Username:  userP.Username,
				Mode:      beatmap.Mode,
				BeatmapID: beatmap.BeatmapID,
				Mods: &score.Mods,
			})
			buf := new(bytes.Buffer)
			buf.ReadFrom(reader)
			replayData := structs.ReplayData{
				Mode:    beatmap.Mode,
				Beatmap: beatmap,
				Score:   replayScore[0].Score,
				Player:  *userP,
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
	}

	score.Rank = strings.Replace(score.Rank, "X", "SS", -1)
	g, _ := s.Guild(config.Conf.Server)
	scoreRank := ""
	for _, emoji := range g.Emojis {
		if emoji.Name == score.Rank+"_" {
			scoreRank = emoji.MessageFormat()
		}
	}

	// Create embed
	embed := &discordgo.MessageEmbed{
		Color: Color,
		Author: &discordgo.MessageEmbedAuthor{
			URL:     "https://osu.ppy.sh/users/" + strconv.Itoa(userP.UserID),
			Name:    userP.Username,
			IconURL: "https://a.ppy.sh/" + strconv.Itoa(userP.UserID) + "?" + strconv.Itoa(rand.Int()) + ".jpeg",
		},
		Title: beatmap.Artist + " - " + beatmap.Title + " [" + beatmap.DiffName + "] by " + strings.Replace(beatmap.Creator, "_", `\_`, -1),
		URL:   "https://osu.ppy.sh/beatmaps/" + strconv.Itoa(beatmap.BeatmapID),
		Description: sr + "\n" +
			length + bpm + "\n" +
			mapStats + "\n" +
			mapObjs + "\n\n" +
			scoreRank + scorePrint + mods + combo + acc + replay + "\n" +
			mapCompletion + "\n" +
			pp + hits + "\n\n",
		Thumbnail: &discordgo.MessageEmbedThumbnail{
			URL: "https://b.ppy.sh/thumb/" + strconv.Itoa(beatmap.BeatmapSetID) + "l.jpg",
		},
	}
	if strings.ToLower(beatmap.Title) == "crab rave" {
		embed.Image = &discordgo.MessageEmbedImage{
			URL: "https://cdn.discordapp.com/emojis/510169818893385729.gif",
		}
	}
	if option == "best" {
		embed.Footer = &discordgo.MessageEmbedFooter{
			Text: tools.TimeSince(timeParse),
		}
	} else if option == "recent" {
		embed.Footer = &discordgo.MessageEmbedFooter{
			Text: "Try #" + strconv.Itoa(try) + " | " + tools.TimeSince(timeParse),
		}
	}
	message, err := s.ChannelMessageSendComplex(m.ChannelID, &discordgo.MessageSend{
		Content: warning,
		Embed:   embed,
	})
	if scorePostRegex.MatchString(m.Content) && err == nil {
		if option == "best" {
			ScorePost(s, &discordgo.MessageCreate{message}, cache, "recentBest")
		} else if option == "recent" {
			ScorePost(s, &discordgo.MessageCreate{message}, cache, "recent")
		}
	}
}
