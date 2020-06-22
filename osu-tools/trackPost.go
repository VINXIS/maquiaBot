package osutools

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"math/rand"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	config "maquiaBot/config"
	osuapi "maquiaBot/osu-api"
	structs "maquiaBot/structs"
	tools "maquiaBot/tools"
)

// TrackPost posts scores for users tracked for that channel
func TrackPost(channel discordgo.Channel, s *discordgo.Session) {

	startTime := time.Now()

	ticker := time.NewTicker(20 * time.Second)
	for {
		select {
		case <-ticker.C:

			// Get channel data
			ch, newCh := tools.GetChannel(channel, s)

			if !ch.Tracking || len(ch.Users) == 0 || newCh {
				return
			}

			for l, user := range ch.Users {
				recentScores, err := OsuAPI.GetUserRecent(osuapi.GetUserScoresOpts{
					UserID: user.UserID,
					Limit:  100,
				})
				if err != nil {
					continue
				}

				for index, score := range recentScores {
					scoreTime, _ := time.Parse("2006-01-02 15:04:05", score.Date.String())
					if score.Rank != "F" && startTime.Before(scoreTime) {

						// Get mods
						mods := "NM"
						if score.Mods != 0 {
							mods = score.Mods.String()
						}

						// Save beatmap
						diffMods := osuapi.Mods(338) & score.Mods
						beatmap := BeatmapParse(strconv.Itoa(score.BeatmapID), "map", &diffMods)

						// Assign timing variables for values below
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
						accCalc := 100.0 * float64(score.Count50+2*score.Count100+6*score.Count300) / float64(6*(score.CountMiss+score.Count50+score.Count100+score.Count300))
						Color := ModeColour(osuapi.ModeOsu)
						sr := "**SR:** " + strconv.FormatFloat(beatmap.DifficultyRating, 'f', 2, 64)
						if beatmap.Mode == osuapi.ModeOsu {
							sr += " **Aim:** " + strconv.FormatFloat(beatmap.DifficultyAim, 'f', 2, 64) + " **Speed:** " + strconv.FormatFloat(beatmap.DifficultySpeed, 'f', 2, 64)
						}
						length := "**Length:** " + fmt.Sprint(totalMinutes) + ":" + fmt.Sprint(totalSeconds) + " (" + fmt.Sprint(hitMinutes) + ":" + fmt.Sprint(hitSeconds) + ") "
						bpm := "**BPM:** " + fmt.Sprint(beatmap.BPM) + " "
						scorePrint := " **" + tools.Comma(score.Score.Score) + "** "
						mapStats := "**CS:** " + strconv.FormatFloat(beatmap.CircleSize, 'f', 1, 64) + " **AR:** " + strconv.FormatFloat(beatmap.ApproachRate, 'f', 1, 64) + " **OD:** " + strconv.FormatFloat(beatmap.OverallDifficulty, 'f', 1, 64) + " **HP:** " + strconv.FormatFloat(beatmap.HPDrain, 'f', 1, 64)
						mapObjs := "**Circles:** " + strconv.Itoa(beatmap.Circles) + " **Sliders:** " + strconv.Itoa(beatmap.Sliders) + " **Spinners:** " + strconv.Itoa(beatmap.Spinners)
						timeText := tools.TimeSince(scoreTime)
						var combo string

						// Remove DT for NC scores
						if strings.Contains(mods, "DTNC") {
							mods = strings.Replace(mods, "DTNC", "NC", 1)
						}
						mods = " **+" + mods + "** "

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
						topNum := 101
						leaderboardNum := 101
						orderedScores, err := OsuAPI.GetUserBest(osuapi.GetUserScoresOpts{
							UserID: user.UserID,
							Limit:  100,
						})
						tools.ErrRead(s, err)
						for i, orderedScore := range orderedScores {
							if score.Score.Score == orderedScore.Score.Score {
								topNum = i + 1
								mapCompletion += "**#" + strconv.Itoa(topNum) + "** in top performances! \n"
								break
							}
						}
						mapScores, err := OsuAPI.GetScores(osuapi.GetScoresOpts{
							BeatmapID: beatmap.BeatmapID,
							Limit:     100,
						})
						tools.ErrRead(s, err)
						for i, mapScore := range mapScores {
							if score.UserID == mapScore.UserID && score.Score.Score == mapScore.Score.Score {
								topNum = i + 1
								mapCompletion += "**#" + strconv.Itoa(topNum) + "** on leaderboard! \n"
								break
							}
						}

						// Get pp values
						var pp string
						totalObjs := beatmap.Circles + beatmap.Sliders + beatmap.Spinners
						if score.PP == 0 { // If map was not finished
							ppValues := make(chan string, 2)
							var ppValueArray [2]string
							go PPCalc(beatmap, osuapi.Score{
								MaxCombo: beatmap.MaxCombo,
								Count50:  score.Count50,
								Count100: score.Count100,
								Count300: totalObjs - score.Count50 - score.Count100,
								Mods:     score.Mods,
							}, ppValues)
							go PPCalc(beatmap, score.Score, ppValues)
							for v := 0; v < 2; v++ {
								ppValueArray[v] = <-ppValues
							}
							sort.Slice(ppValueArray[:], func(i, j int) bool {
								pp1, _ := strconv.ParseFloat(ppValueArray[i], 64)
								pp2, _ := strconv.ParseFloat(ppValueArray[j], 64)
								return pp1 > pp2
							})
							score.PP, _ = strconv.ParseFloat(ppValueArray[1], 64)
							pp = "**" + ppValueArray[1] + "pp**/" + ppValueArray[0] + "pp "
						} else if score.Score.FullCombo { // If play was a perfect combo
							pp = "**" + strconv.FormatFloat(score.PP, 'f', 0, 64) + "pp**/" + strconv.FormatFloat(score.PP, 'f', 0, 64) + "pp "
						} else { // If map was finished, but play was not a perfect combo
							ppValues := make(chan string, 1)
							go PPCalc(beatmap, osuapi.Score{
								MaxCombo: beatmap.MaxCombo,
								Count50:  score.Count50,
								Count100: score.Count100,
								Count300: totalObjs - score.Count50 - score.Count100,
								Mods:     score.Mods,
							}, ppValues)
							pp = "**" + strconv.FormatFloat(score.PP, 'f', 0, 64) + "pp**/" + <-ppValues + "pp "
						}

						// Check params
						rankCheck := ch.Ranked && (beatmap.Approved == osuapi.StatusRanked || beatmap.Approved == osuapi.StatusApproved)
						qualCheck := ch.Qualified && beatmap.Approved == osuapi.StatusQualified
						loveCheck := ch.Loved && beatmap.Approved == osuapi.StatusLoved
						mapCheck := rankCheck || qualCheck || loveCheck

						// score checking is a bit more complicated than map checking since we should be ignoring parameters if they were not technically given
						ppCheck := score.PP >= ch.PPReq
						leaderboardCheck := leaderboardNum <= ch.LeaderboardReq
						topCheck := topNum <= ch.TopReq
						scoreCheck := false
						if ch.PPReq != -1 {
							scoreCheck = scoreCheck || ppCheck
						}
						if ch.LeaderboardReq != 101 {
							scoreCheck = scoreCheck || leaderboardCheck
						}
						if ch.TopReq != 101 {
							scoreCheck = scoreCheck || topCheck
						}
						// If no requirements were given for all 3 areas
						if ch.PPReq == -1 && ch.LeaderboardReq == 101 && ch.TopReq == 101 {
							scoreCheck = true
						}

						checkPass := mapCheck && scoreCheck
						if checkPass {
							// Count number of tries
							try := 0
							tryInc := true
							for i := index; i < len(recentScores); i++ {
								if tryInc && recentScores[i].BeatmapID == score.BeatmapID {
									try++
								} else {
									tryInc = false
								}
							}
							acc := "** " + strconv.FormatFloat(accCalc, 'f', 2, 64) + "%** "
							hits := "**Hits:** [" + strconv.Itoa(score.Count300) + "/" + strconv.Itoa(score.Count100) + "/" + strconv.Itoa(score.Count50) + "/" + strconv.Itoa(score.CountMiss) + "]"

							replay := ""
							replayScore, _ := OsuAPI.GetScores(osuapi.GetScoresOpts{
								BeatmapID: beatmap.BeatmapID,
								UserID:    user.UserID,
								Mods:      &score.Mods,
							})
							if replayScore[0].Replay && replayScore[0].Score.Score == score.Score.Score {
								replay = "| [**Replay**](https://osu.ppy.sh/scores/osu/" + strconv.FormatInt(replayScore[0].ScoreID, 10) + "/download)"
								reader, _ := OsuAPI.GetReplay(osuapi.GetReplayOpts{
									UserID:    user.UserID,
									Mode:      beatmap.Mode,
									BeatmapID: beatmap.BeatmapID,
									Mods:      &score.Mods,
								})
								buf := new(bytes.Buffer)
								buf.ReadFrom(reader)
								replayData := structs.ReplayData{
									Mode:    beatmap.Mode,
									Beatmap: beatmap,
									Score:   replayScore[0].Score,
									Player:  user,
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
									URL:     "https://osu.ppy.sh/users/" + strconv.Itoa(user.UserID),
									IconURL: "https://a.ppy.sh/" + strconv.Itoa(user.UserID) + "?" + strconv.Itoa(rand.Int()) + ".jpeg",
								},
								Title: beatmap.Artist + " - " + beatmap.Title + " [" + beatmap.DiffName + "] by " + strings.Replace(beatmap.Creator, "_", `\_`, -1),
								URL:   "https://osu.ppy.sh/beatmaps/" + strconv.Itoa(beatmap.BeatmapID),
								Thumbnail: &discordgo.MessageEmbedThumbnail{
									URL: "https://b.ppy.sh/thumb/" + strconv.Itoa(beatmap.BeatmapSetID) + "l.jpg",
								},
								Description: sr + "\n" +
									length + bpm + "\n" +
									mapStats + "\n" +
									mapObjs + "\n\n" +
									scoreRank + scorePrint + mods + combo + acc + replay + "\n" +
									mapCompletion + "\n" +
									pp + hits + "\n\n",
								Footer: &discordgo.MessageEmbedFooter{
									Text: "Try #" + strconv.Itoa(try) + " | " + timeText,
								},
							}

							recentUser, err := OsuAPI.GetUser(osuapi.GetUserOpts{
								UserID: user.UserID,
							})
							tools.ErrRead(s, err)

							if recentUser.PP != user.PP {
								embed.Author.Name = user.Username + " " + strconv.FormatFloat(user.PP, 'f', 2, 64) + " -> " + strconv.FormatFloat(recentUser.PP, 'f', 2, 64)
								if recentUser.PP-user.PP > 0 {
									embed.Author.Name += " (+" + strconv.FormatFloat(recentUser.PP-user.PP, 'f', 2, 64) + "pp)"
								} else {
									embed.Author.Name += " (" + strconv.FormatFloat(recentUser.PP-user.PP, 'f', 2, 64) + "pp)"
								}
							} else {
								embed.Author.Name = user.Username
							}
							if strings.ToLower(beatmap.Title) == "crab rave" {
								embed.Image = &discordgo.MessageEmbedImage{
									URL: "https://cdn.discordapp.com/emojis/510169818893385729.gif",
								}
							}
							_, err = s.ChannelMessageSendEmbed(ch.Channel.ID, embed)
							if err != nil {
								log.Println("TrackPost err: " + err.Error())
								log.Println(ch)
							}
							ch.Users[l] = *recentUser
						}
					}
				}
			}

			// Write data to JSON
			jsonCache, err := json.Marshal(ch)
			tools.ErrRead(s, err)

			err = ioutil.WriteFile("./data/channelData/"+ch.Channel.ID+".json", jsonCache, 0644)
			tools.ErrRead(s, err)

			startTime = time.Now()
		}
	}
}
