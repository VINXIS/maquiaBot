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

	osuapi "../../osu-api"
	osutools "../../osu-functions"
	structs "../../structs"
	tools "../../tools"
	"github.com/bwmarrin/discordgo"
)

// Compare compares finds a score from the current user on the previous map linked by the bot
func Compare(s *discordgo.Session, m *discordgo.MessageCreate, args []string, osuAPI *osuapi.Client, cache []structs.PlayerData, serverPrefix string, mapCache []structs.MapData) {
	mapRegex, _ := regexp.Compile(`(https:\/\/)?(osu|old)\.ppy\.sh\/(s|b|beatmaps|beatmapsets)\/(\d+)(#(osu|taiko|fruits|mania)\/(\d+))?`)
	modRegex, _ := regexp.Compile(`-m\s*(\S+)`)
	compareRegex, _ := regexp.Compile(`(c|compare)\s*(.+)?`)
	strictRegex, _ := regexp.Compile(`-nostrict`)
	var beatmap osuapi.Beatmap

	// Obtain username and mods
	username := ""
	mods := ""
	strict := true
	if compareRegex.MatchString(m.Content) {
		username = compareRegex.FindStringSubmatch(m.Content)[2]
		if modRegex.MatchString(username) {
			mods = strings.ToUpper(modRegex.FindStringSubmatch(username)[1])
			if strings.Contains(mods, "NC") && !strings.Contains(mods, "DT") {
				mods += "DT"
			}
			username = strings.TrimSpace(strings.Replace(username, modRegex.FindStringSubmatch(username)[0], "", 1))
		}
		if strictRegex.MatchString(m.Content) {
			strict = false
			username = strings.TrimSpace(strings.Replace(username, strictRegex.FindStringSubmatch(m.Content)[0], "", 1))
		}
	}

	// Get the map
	if mapRegex.MatchString(m.Content) {
		submatches := mapRegex.FindStringSubmatch(m.Content)
		switch submatches[3] {
		case "s":
			beatmap = osutools.BeatmapParse(submatches[4], "set", osuapi.ParseMods(mods), osuAPI)
		case "b":
			beatmap = osutools.BeatmapParse(submatches[4], "map", osuapi.ParseMods(mods), osuAPI)
		case "beatmaps":
			beatmap = osutools.BeatmapParse(submatches[4], "map", osuapi.ParseMods(mods), osuAPI)
		case "beatmapsets":
			if len(submatches[7]) > 0 {
				beatmap = osutools.BeatmapParse(submatches[7], "map", osuapi.ParseMods(mods), osuAPI)
			} else {
				beatmap = osutools.BeatmapParse(submatches[4], "set", osuapi.ParseMods(mods), osuAPI)
			}
		}
		if beatmap.BeatmapID == 0 {
			s.ChannelMessageSend(m.ChannelID, "Map does not exist!")
			return
		} else if beatmap.Approved < 1 {
			s.ChannelMessageSend(m.ChannelID, "The map `"+beatmap.Artist+" - "+beatmap.Title+"` does not have a leaderboard!")
			return
		}
		username = strings.TrimSpace(strings.Replace(username, submatches[0], "", -1))
	}

	// Check if message linked a map or not, if not then check previous messages instead
	if beatmap.BeatmapID == 0 {
		// Get prev messages
		messages, err := s.ChannelMessages(m.ChannelID, -1, "", "", "")
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, "No map to compare to!")
			return
		}

		// Sort by date
		sort.Slice(messages, func(i, j int) bool {
			time1, _ := messages[i].Timestamp.Parse()
			time2, _ := messages[j].Timestamp.Parse()
			return time1.After(time2)
		})
		beatmapRegex, _ := regexp.Compile(`https://osu.ppy.sh/beatmaps/(\d+)`)
		mapID := 0
		found := false

		// Look for a valid beatmap ID
		for _, msg := range messages {
			msgTime, _ := msg.Timestamp.Parse()
			if msg.Author.ID == s.State.User.ID && time.Since(msgTime) < time.Hour {
				if msg.ID != (discordgo.Message{}).ID && len(msg.Embeds) > 0 && msg.Embeds[0].Author != nil {
					if beatmapRegex.MatchString(msg.Embeds[0].URL) {
						mapID, _ = strconv.Atoi(beatmapRegex.FindStringSubmatch(msg.Embeds[0].URL)[1])
						found = true
						break
					} else if beatmapRegex.MatchString(msg.Embeds[0].Author.URL) {
						mapID, _ = strconv.Atoi(beatmapRegex.FindStringSubmatch(msg.Embeds[0].Author.URL)[1])
						found = true
						break
					}
				}
			}
		}

		// Check if found
		if found == false {
			s.ChannelMessageSend(m.ChannelID, "No map to compare to!")
			return
		}

		// Get the map
		beatmap = osutools.BeatmapParse(strconv.Itoa(mapID), "map", osuapi.ParseMods(mods), osuAPI)
		if beatmap.BeatmapID == 0 {
			s.ChannelMessageSend(m.ChannelID, "Map does not exist!")
			return
		} else if beatmap.Approved < 1 {
			s.ChannelMessageSend(m.ChannelID, "The map `"+beatmap.Artist+" - "+beatmap.Title+"` does not have a leaderboard!")
			return
		}
	}

	if beatmap.BeatmapID == 0 {
		s.ChannelMessageSend(m.ChannelID, "No map to compare to!")
		return
	}

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
		test, err := osuAPI.GetUser(osuapi.GetUserOpts{
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
	}
	scores, err := osuAPI.GetScores(scoreOpts)
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
	if mods != "" {
		parsedMods := osuapi.ParseMods(mods)
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

	// Sort by PP
	sort.Slice(scores, func(i, j int) bool {
		return scores[i].PP > scores[j].PP
	})

	score := scores[0]

	// Get time since play
	timeParse, _ := time.Parse("2006-01-02 15:04:05", score.Date.String())
	time := tools.TimeSince(timeParse)

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

	// Assign values
	mods = score.Mods.String()
	accCalc := (50.0*float64(score.Count50) + 100.0*float64(score.Count100) + 300.0*float64(score.Count300)) / (300.0 * float64(score.CountMiss+score.Count50+score.Count100+score.Count300)) * 100.0
	Color := osutools.ModeColour(beatmap.Mode)
	sr, _, _, _, _, _ := osutools.BeatmapCache(mods, beatmap, mapCache)
	length := "**Length:** " + fmt.Sprint(totalMinutes) + ":" + fmt.Sprint(totalSeconds) + " (" + fmt.Sprint(hitMinutes) + ":" + fmt.Sprint(hitSeconds) + ") "
	bpm := "**BPM:** " + fmt.Sprint(beatmap.BPM) + " "
	scorePrint := " **" + tools.Comma(score.Score.Score) + "** "
	mapStats := "**CS:** " + strconv.FormatFloat(beatmap.CircleSize, 'f', 1, 64) + " **AR:** " + strconv.FormatFloat(beatmap.ApproachRate, 'f', 1, 64) + " **OD:** " + strconv.FormatFloat(beatmap.OverallDifficulty, 'f', 1, 64) + " **HP:** " + strconv.FormatFloat(beatmap.HPDrain, 'f', 1, 64)
	mapObjs := "**Circles:** " + strconv.Itoa(beatmap.Circles) + " **Sliders:** " + strconv.Itoa(beatmap.Sliders) + " **Spinners:** " + strconv.Itoa(beatmap.Spinners)
	acc := "** " + strconv.FormatFloat(accCalc, 'f', 2, 64) + "%** "
	hits := "**Hits:** [" + strconv.Itoa(score.Count300) + "/" + strconv.Itoa(score.Count100) + "/" + strconv.Itoa(score.Count50) + "/" + strconv.Itoa(score.CountMiss) + "]"

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
		combo = " **x" + strconv.Itoa(score.MaxCombo) + "**/" + strconv.Itoa(beatmap.MaxCombo) + " "
	}

	mapCompletion := ""
	orderedScores, err := osuAPI.GetUserBest(osuapi.GetUserScoresOpts{
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
	mapScores, err := osuAPI.GetScores(osuapi.GetScoresOpts{
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

	g, _ := s.Guild("556243477084635170")
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
			URL:     "https://osu.ppy.sh/users/" + strconv.Itoa(user.UserID),
			Name:    user.Username,
			IconURL: "https://a.ppy.sh/" + strconv.Itoa(user.UserID) + "?" + strconv.Itoa(rand.Int()) + ".jpeg",
		},
		Title: beatmap.Artist + " - " + beatmap.Title + " [" + beatmap.DiffName + "] by " + beatmap.Creator,
		URL:   "https://osu.ppy.sh/beatmaps/" + strconv.Itoa(beatmap.BeatmapID),
		Thumbnail: &discordgo.MessageEmbedThumbnail{
			URL: "https://b.ppy.sh/thumb/" + strconv.Itoa(beatmap.BeatmapSetID) + "l.jpg",
		},
		Description: sr + length + bpm + "\n" +
			mapStats + "\n" +
			mapObjs + "\n\n" +
			scorePrint + mods + combo + acc + scoreRank + "\n" +
			mapCompletion + "\n" +
			pp + hits + "\n\n",
		Footer: &discordgo.MessageEmbedFooter{
			Text: time,
		},
	}
	if strings.ToLower(beatmap.Title) == "crab rave" {
		embed.Image = &discordgo.MessageEmbedImage{
			URL: "https://cdn.discordapp.com/emojis/510169818893385729.gif",
		}
	}
	s.ChannelMessageSendEmbed(m.ChannelID, embed)
}
