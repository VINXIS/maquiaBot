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
	messageRegex, _ := regexp.Compile(`(https:\/\/)?(osu|old)\.ppy\.sh\/(s|b|beatmaps|beatmapsets)\/(\d+)(#(osu|taiko|fruits|mania)\/(\d+))?`)
	modRegex, _ := regexp.Compile(`-m\s*(\S+)`)
	compareRegex, _ := regexp.Compile(`(c|compare)\s*(.+)?`)

	// Obtain username, mods, and map id
	userArg := ""
	mods := ""
	beatmaps := []osuapi.Beatmap{}
	beatmap := osuapi.Beatmap{}
	var err error
	var APIMods osuapi.Mods
	if compareRegex.MatchString(m.Content) {
		userArg = compareRegex.FindStringSubmatch(m.Content)[2]
	}

	if modRegex.MatchString(m.Content) {
		modCheck := modRegex.FindStringSubmatch(m.Content)
		modsRaw := strings.ToUpper(modCheck[1])
		if len(modsRaw)%2 == 0 && len(osuapi.ParseMods(modsRaw).String()) > 0 {
			APIMods = osuapi.ParseMods(modsRaw)
			mods = APIMods.String()
		}
		userArg = strings.TrimSpace(strings.Replace(userArg, modCheck[0], "", -1))
	}

	if messageRegex.MatchString(m.Content) {
		submatches := messageRegex.FindStringSubmatch(m.Content)
		mapID := 0
		setID := 0
		switch submatches[3] {
		case "s":
			setID, _ = strconv.Atoi(submatches[4])
		case "b", "beatmaps":
			mapID, _ = strconv.Atoi(submatches[4])
		case "beatmapsets":
			if len(submatches[7]) > 0 {
				mapID, _ = strconv.Atoi(submatches[7])
			} else {
				setID, _ = strconv.Atoi(submatches[4])
			}
		}
		if setID != 0 {
			beatmaps, err = osuAPI.GetBeatmaps(osuapi.GetBeatmapsOpts{
				BeatmapSetID: setID,
			})
			if err != nil {
				s.ChannelMessageSend(m.ChannelID, "The osu! API just owned me. Please try again!")
				return
			}
			if len(beatmaps) == 0 {
				s.ChannelMessageSend(m.ChannelID, "This set does not exist!")
				return
			}
			sort.Slice(beatmaps, func(i, j int) bool {
				return beatmaps[i].DifficultyRating > beatmaps[j].DifficultyRating
			})
			beatmap = beatmaps[0]
			if beatmap.Approved == osuapi.StatusPending || beatmap.Approved == osuapi.StatusGraveyard || beatmap.Approved == osuapi.StatusWIP {
				s.ChannelMessageSend(m.ChannelID, "The map `"+beatmap.Artist+" - "+beatmap.Title+"` does not have a leaderboard!")
				return
			}
		} else {
			beatmaps, err = osuAPI.GetBeatmaps(osuapi.GetBeatmapsOpts{
				BeatmapID: mapID,
			})
			if err != nil {
				s.ChannelMessageSend(m.ChannelID, "The osu! API just owned me. Please try again!")
				return
			}
			if len(beatmaps) == 0 {
				s.ChannelMessageSend(m.ChannelID, "This map does not exist!")
				return
			}
			beatmap = beatmaps[0]
			if beatmap.Approved == osuapi.StatusPending || beatmap.Approved == osuapi.StatusGraveyard || beatmap.Approved == osuapi.StatusWIP {
				s.ChannelMessageSend(m.ChannelID, "The map `"+beatmap.Artist+" - "+beatmap.Title+"` does not have a leaderboard!")
				return
			}
		}
		userArg = strings.TrimSpace(strings.Replace(userArg, submatches[0], "", -1))
	}

	// Get prev messages
	messages, err := s.ChannelMessages(m.ChannelID, -1, "", "", "")
	tools.ErrRead(err)

	// Sort by date
	sort.Slice(messages, func(i, j int) bool {
		time1, err := time.Parse(time.RFC3339, string(messages[i].Timestamp))
		tools.ErrRead(err)
		time2, err := time.Parse(time.RFC3339, string(messages[j].Timestamp))
		tools.ErrRead(err)
		return time1.After(time2)
	})

	// Check if message linked a map or not
	beatmapRegex, _ := regexp.Compile(`https://osu.ppy.sh/beatmaps/(\d+)`)
	if beatmap.BeatmapID == 0 {
		mapID := 0
		found := false
		for _, msg := range messages {
			msgTime, err := time.Parse(time.RFC3339, string(msg.Timestamp))
			if msg.Author.ID == s.State.User.ID && time.Since(msgTime) < time.Hour {
				if msg.ID != (discordgo.Message{}).ID && len(msg.Embeds) > 0 && msg.Embeds[0].Author != nil {
					if beatmapRegex.MatchString(msg.Embeds[0].URL) {
						mapID, err = strconv.Atoi(beatmapRegex.FindStringSubmatch(msg.Embeds[0].URL)[1])
						tools.ErrRead(err)
						found = true
						break
					} else if beatmapRegex.MatchString(msg.Embeds[0].Author.URL) {
						mapID, err = strconv.Atoi(beatmapRegex.FindStringSubmatch(msg.Embeds[0].Author.URL)[1])
						tools.ErrRead(err)
						found = true
						break
					}
				}
			}
		}
		if found == false {
			s.ChannelMessageSend(m.ChannelID, "No map to compare to!")
			return
		}
		beatmaps, err = osuAPI.GetBeatmaps(osuapi.GetBeatmapsOpts{
			BeatmapID: mapID,
		})
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, "The osu! API just owned me. Please try again!")
			return
		}
		if len(beatmaps) == 0 {
			s.ChannelMessageSend(m.ChannelID, "This map does not exist!")
			return
		}
		beatmap = beatmaps[0]
		if beatmap.Approved == osuapi.StatusPending || beatmap.Approved == osuapi.StatusGraveyard || beatmap.Approved == osuapi.StatusWIP {
			s.ChannelMessageSend(m.ChannelID, "The map `"+beatmap.Artist+" - "+beatmap.Title+"` does not have a leaderboard!")
			return
		}
	}

	if beatmap.BeatmapID == 0 {
		s.ChannelMessageSend(m.ChannelID, "No map to compare to!")
		return
	}

	// Get user
	user := osuapi.User{}
	for _, player := range cache {
		if userArg != "" {
			if userArg == user.Username {
				user = player.Osu
				break
			}
		} else if m.Author.ID == player.Discord.ID && player.Osu.Username != (osuapi.User{}).Username {
			user = player.Osu
			break
		}
	}
	if user.UserID == (osuapi.User{}).UserID {
		// Check if user even exists
		if userArg == "" {
			s.ChannelMessageSend(m.ChannelID, "No user mentioned in message/linked to your account! Please use `set` or `link` to link an osu! account to you, or name a user to obtain their recent score of!")
		}
		test, err := osuAPI.GetUser(osuapi.GetUserOpts{
			Username: userArg,
		})
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, "User "+userArg+" may not exist! Are you sure you replaced spaces with `_`?")
			return
		}
		user = *test
	}

	// API call
	scoreOpts := osuapi.GetScoresOpts{
		BeatmapID: beatmap.BeatmapID,
		UserID:    user.UserID,
	}
	if mods != "" {
		scoreOpts.Mods = &APIMods
	}

	scores, err := osuAPI.GetScores(scoreOpts)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "The osu! API just owned me. Please try again!")
		return
	}

	if strings.Contains(mods, "DTNC") {
		mods = strings.Replace(mods, "DTNC", "NC", 1)
	}

	if len(scores) == 0 {
		if userArg != "" {
			s.ChannelMessageSend(m.ChannelID, userArg+" hasn't set a score on this!")
		} else if mods != "" {
			if strings.Contains(mods, "DT") {
				s.ChannelMessageSend(m.ChannelID, "You haven't set a score on this with the mods: **"+mods+"**! Are you sure you didn't play with **NC** like the dumbass weaboo u are?")
			} else if strings.Contains(mods, "NC") {
				s.ChannelMessageSend(m.ChannelID, "You haven't set a score on this with the mods: **"+mods+"**! Are you sure you didn't play with **DT** like a sane person?")
			} else {
				s.ChannelMessageSend(m.ChannelID, "You haven't set a score on this with the mods: **"+mods+"**!")
			}
		} else {
			s.ChannelMessageSend(m.ChannelID, "You haven't set a score on this with any mod combination!")
		}
		return
	}

	// Sort by PP
	sort.Slice(scores, func(i, j int) bool {
		return scores[i].PP > scores[j].PP
	})

	score := scores[0]

	// Get time since play
	timeParse, err := time.Parse("2006-01-02 15:04:05", score.Date.String())
	tools.ErrRead(err)
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
	if score.Mods != 0 {
		mods = score.Mods.String()
	}
	accCalc := (50.0*float64(score.Count50) + 100.0*float64(score.Count100) + 300.0*float64(score.Count300)) / (300.0 * float64(score.CountMiss+score.Count50+score.Count100+score.Count300)) * 100.0
	Color := osutools.ModeColour(osuapi.ModeOsu)
	sr, _, _, _, _, _ := osutools.BeatmapCache(mods, beatmap, mapCache)
	length := "**Length:** " + fmt.Sprint(totalMinutes) + ":" + fmt.Sprint(totalSeconds) + " (" + fmt.Sprint(hitMinutes) + ":" + fmt.Sprint(hitSeconds) + ") "
	bpm := "**BPM:** " + fmt.Sprint(beatmap.BPM) + " "
	scorePrint := " **" + tools.Comma(score.Score.Score) + "** "
	mapStats := "**CS:** " + strconv.FormatFloat(beatmap.CircleSize, 'f', 1, 64) + " **AR:** " + strconv.FormatFloat(beatmap.ApproachRate, 'f', 1, 64) + " **OD:** " + strconv.FormatFloat(beatmap.OverallDifficulty, 'f', 1, 64) + " **HP:** " + strconv.FormatFloat(beatmap.HPDrain, 'f', 1, 64)
	mapObjs := "**Circles:** " + strconv.Itoa(beatmap.Circles) + " **Sliders:** " + strconv.Itoa(beatmap.Sliders) + " **Spinners:** " + strconv.Itoa(beatmap.Spinners)

	if mods == "" {
		mods = "NM"
	}

	if strings.Contains(mods, "DTNC") {
		mods = strings.Replace(mods, "DTNC", "NC", 1)
	}
	scoreMods := mods
	mods = " **+" + mods + "** "

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
	mapCompletion2 := ""
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
			mapCompletion = "**#" + strconv.Itoa(i+1) + "** in top performances! \n"
			mapScores, err := osuAPI.GetScores(osuapi.GetScoresOpts{
				BeatmapID: beatmap.BeatmapID,
				Limit:     100,
			})
			if err != nil {
				s.ChannelMessageSend(m.ChannelID, "The osu! API just owned me. Please try again!")
				return
			}
			for j, mapScore := range mapScores {
				if score.UserID == mapScore.UserID && score.Score.Score == mapScore.Score.Score {
					mapCompletion2 = "**#" + strconv.Itoa(j+1) + "** on leaderboard! \n"
					break
				}
			}
			break
		}
	}
	if mapCompletion == "" {
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
				mapCompletion = "**#" + strconv.Itoa(i+1) + "** on leaderboard! \n"
				scores, err := osuAPI.GetUserBest(osuapi.GetUserScoresOpts{
					Username: user.Username,
					Limit:    100,
				})
				if err != nil {
					s.ChannelMessageSend(m.ChannelID, "The osu! API just owned me. Please try again!")
					return
				}
				if len(scores) == 0 {
					s.ChannelMessageSend(m.ChannelID, user.Username+" has no top scores!")
					return
				}

				for j, bestScore := range scores {
					if score.UserID == bestScore.UserID && score.Score.Score == bestScore.Score.Score {
						mapCompletion2 = "**#" + strconv.Itoa(j+1) + "** in top performances! \n"
						break
					}
				}
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
		go osutools.PPCalc(beatmap, accCalcNoMiss, "", "", scoreMods, ppValues)
		pp = "**" + strconv.FormatFloat(score.PP, 'f', 2, 64) + "pp**/" + <-ppValues + "pp "
	}
	acc := "** " + strconv.FormatFloat(accCalc, 'f', 2, 64) + "%** "
	hits := "**Hits:** [" + strconv.Itoa(score.Count300) + "/" + strconv.Itoa(score.Count100) + "/" + strconv.Itoa(score.Count50) + "/" + strconv.Itoa(score.CountMiss) + "]"

	g, err := s.Guild("556243477084635170")
	tools.ErrRead(err)
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
			mapCompletion + mapCompletion2 + "\n" +
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
	return
}
