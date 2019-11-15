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

// Recent gets the most recent score done/nth score done
func Recent(s *discordgo.Session, m *discordgo.MessageCreate, osuAPI *osuapi.Client, cache []structs.PlayerData, option string, mapCache []structs.MapData) {
	index := 1
	username := ""
	var err error
	recentRegex, _ := regexp.Compile(`(.+)(r|recent|rs|rb|recentb|recentbest)\s+(.+)`)

	// Obtain index and username
	if recentRegex.MatchString(m.Content) {
		args := strings.Split(recentRegex.FindStringSubmatch(m.Content)[3], " ")
		if len(args) > 1 {
			if index, err = strconv.Atoi(args[len(args)-1]); err == nil {
				username = strings.Join(args[:len(args)-1], " ")
			} else {
				index = 1
				username = recentRegex.FindStringSubmatch(m.Content)[3]
			}
		} else {
			index, err = strconv.Atoi(args[0])
			if err != nil {
				index = 1
				username = args[0]
			}
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
	userP, err := osuAPI.GetUser(osuapi.GetUserOpts{
		Username: username,
	})
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "User **"+username+"** may not exist!")
		return
	}
	go osutools.PlayerCache(*userP, cache)
	scoreList := []osuapi.GUSScore{}

	// Run api call for user best/recent
	if option == "recent" {
		scoreList, err = osuAPI.GetUserRecent(osuapi.GetUserScoresOpts{
			Username: username,
			Limit:    50,
		})
		tools.ErrRead(err)
		if len(scoreList) == 0 {
			s.ChannelMessageSend(m.ChannelID, username+" has not played recently!")
			return
		}
	} else if option == "best" {
		scoreList, err = osuAPI.GetUserBest(osuapi.GetUserScoresOpts{
			Username: username,
			Limit:    100,
		})
		tools.ErrRead(err)
		if len(scoreList) == 0 {
			s.ChannelMessageSend(m.ChannelID, username+" has no top scores!")
			return
		}
	} else {
		s.ChannelMessageSend(m.ChannelID, "Oops! Something went wrong! I was somehow not given recent or best as an option!")
		return
	}

	// Sort scores by date and get score
	sort.Slice(scoreList, func(i, j int) bool {
		time1, err := time.Parse("2006-01-02 15:04:05", scoreList[i].Date.String())
		tools.ErrRead(err)
		time2, err := time.Parse("2006-01-02 15:04:05", scoreList[j].Date.String())
		tools.ErrRead(err)

		return time1.Unix() > time2.Unix()
	})

	warning := ""
	if len(scoreList) < index {
		index = len(scoreList)
		warning = "Defaulted to max: " + strconv.Itoa(len(scoreList))
	}
	score := scoreList[index-1]

	// Get beatmap, acc, and mods
	beatmap := osutools.BeatmapParse(strconv.Itoa(score.BeatmapID), "map", score.Mods, osuAPI)
	accCalc := (50.0*float64(score.Count50) + 100.0*float64(score.Count100) + 300.0*float64(score.Count300)) / (300.0 * float64(score.CountMiss+score.Count50+score.Count100+score.Count300)) * 100.0
	mods := "NM"
	if score.Mods != 0 {
		mods = score.Mods.String()
	}

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
	objCount := osutools.CountObjs(beatmap)
	playObjCount := score.CountMiss + score.Count100 + score.Count300 + score.Count50

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

	// Assign misc variables
	Color := osutools.ModeColour(osuapi.ModeOsu)
	sr, _, _, _, _, _ := osutools.BeatmapCache(mods, beatmap, mapCache)
	length := "**Length:** " + fmt.Sprint(totalMinutes) + ":" + fmt.Sprint(totalSeconds) + " (" + fmt.Sprint(hitMinutes) + ":" + fmt.Sprint(hitSeconds) + ") "
	bpm := "**BPM:** " + fmt.Sprint(beatmap.BPM) + " "
	mapStats := "**CS:** " + strconv.FormatFloat(beatmap.CircleSize, 'f', 1, 64) + " **AR:** " + strconv.FormatFloat(beatmap.ApproachRate, 'f', 1, 64) + " **OD:** " + strconv.FormatFloat(beatmap.OverallDifficulty, 'f', 1, 64) + " **HP:** " + strconv.FormatFloat(beatmap.HPDrain, 'f', 1, 64)
	mapObjs := "**Circles:** " + strconv.Itoa(beatmap.Circles) + " **Sliders:** " + strconv.Itoa(beatmap.Sliders) + " **Spinners:** " + strconv.Itoa(beatmap.Spinners)
	scorePrint := " **" + tools.Comma(score.Score.Score) + "** "
	var combo string
	var mapCompletion string
	var mapCompletion2 string

	if strings.Contains(mods, "DTNC") {
		mods = strings.Replace(mods, "DTNC", "NC", 1)
	}
	scoreMods := mods
	mods = " **+" + mods + "** "

	if score.MaxCombo == beatmap.MaxCombo {
		if accCalc == 100.0 {
			combo = " **SS** "
		} else {
			combo = " **FC** "
		}
	} else {
		combo = " **x" + strconv.Itoa(score.MaxCombo) + "**/" + strconv.Itoa(beatmap.MaxCombo) + " "
	}

	if objCount != playObjCount {
		completed := float64(playObjCount) / float64(objCount) * 100.0
		mapCompletion = "**" + strconv.FormatFloat(completed, 'f', 2, 64) + "%** completed \n"
		mapCompletion2 = ""
	} else {
		if option == "best" {
			sort.Slice(scoreList, func(i, j int) bool {
				return scoreList[i].PP > scoreList[j].PP
			})
			for i, bestScore := range scoreList {
				if score.BeatmapID == bestScore.BeatmapID {
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
		} else if option == "recent" {
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
						Username: username,
						Limit:    100,
					})
					if err != nil {
						s.ChannelMessageSend(m.ChannelID, "The osu! API just owned me. Please try again!")
						return
					}
					if len(scores) == 0 {
						s.ChannelMessageSend(m.ChannelID, username+" has no top scores!")
						return
					}

					for j, bestScore := range scores {
						if score.BeatmapID == bestScore.BeatmapID {
							mapCompletion2 = "**#" + strconv.Itoa(j+1) + "** in top performances! \n"
							break
						}
					}
					break
				}
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
		go osutools.PPCalc(beatmap, accCalcNoMiss, "", "", scoreMods, ppValues)
		go osutools.PPCalc(beatmap, accCalc, strconv.Itoa(score.MaxCombo), strconv.Itoa(score.CountMiss), scoreMods, ppValues)
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
			URL:     "https://osu.ppy.sh/users/" + strconv.Itoa(userP.UserID),
			Name:    userP.Username,
			IconURL: "https://a.ppy.sh/" + strconv.Itoa(userP.UserID) + "?" + strconv.Itoa(rand.Int()) + ".jpeg",
		},
		Title: beatmap.Artist + " - " + beatmap.Title + " [" + beatmap.DiffName + "] by " + beatmap.Creator,
		URL:   "https://osu.ppy.sh/beatmaps/" + strconv.Itoa(beatmap.BeatmapID),
		Description: sr + length + bpm + "\n" +
			mapStats + "\n" +
			mapObjs + "\n\n" +
			scorePrint + mods + combo + acc + scoreRank + "\n" +
			mapCompletion + mapCompletion2 + "\n" +
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
			Text: time,
		}
	} else if option == "recent" {
		embed.Footer = &discordgo.MessageEmbedFooter{
			Text: "Try #" + strconv.Itoa(try) + " | " + time,
		}
	}
	s.ChannelMessageSendComplex(m.ChannelID, &discordgo.MessageSend{
		Content: warning,
		Embed:   embed,
	})
	return
}
