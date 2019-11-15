package osutools

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"math/rand"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	osuapi "../osu-api"
	structs "../structs"
	tools "../tools"
	"github.com/bwmarrin/discordgo"
)

// TrackPost posts scores for users tracked for that channel
func TrackPost(channel string, s *discordgo.Session, mapCache []structs.MapData) {
	osuAPI := osuapi.NewClient(os.Getenv("OSU_API"))

	startTime := time.Now()

	ticker := time.NewTicker(20 * time.Second)
	for {
		select {
		case <-ticker.C:

			// Get channel data
			ch := structs.ChannelData{}
			f, err := ioutil.ReadFile("./" + channel)
			if err != nil {
				return
			}
			_ = json.Unmarshal(f, &ch)

			if !ch.Tracking || len(ch.Users) == 0 {
				return
			}

			for l, user := range ch.Users {
				recentScores, err := osuAPI.GetUserRecent(osuapi.GetUserScoresOpts{
					UserID: user.UserID,
					Limit:  100,
				})
				if err != nil {
					fmt.Println(err)
					continue
				}

				for index, score := range recentScores {
					scoreTime, err := time.Parse("2006-01-02 15:04:05", score.Date.String())
					tools.ErrRead(err)
					if score.Rank != "F" && startTime.Before(scoreTime) {
						mods := "NM"
						if score.Mods != 0 {
							mods = score.Mods.String()
						}

						beatmap := BeatmapParse(strconv.Itoa(score.BeatmapID), "map", score.Mods, osuAPI)
						timeText := tools.TimeSince(scoreTime)

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
						accCalc := (50.0*float64(score.Count50) + 100.0*float64(score.Count100) + 300.0*float64(score.Count300)) / (300.0 * float64(score.CountMiss+score.Count50+score.Count100+score.Count300)) * 100.0
						Color := ModeColour(osuapi.ModeOsu)
						sr, _, _, _, _, _ := BeatmapCache(mods, beatmap, mapCache)
						length := "**Length:** " + fmt.Sprint(totalMinutes) + ":" + fmt.Sprint(totalSeconds) + " (" + fmt.Sprint(hitMinutes) + ":" + fmt.Sprint(hitSeconds) + ") "
						bpm := "**BPM:** " + fmt.Sprint(beatmap.BPM) + " "
						scorePrint := " **" + tools.Comma(score.Score.Score) + "** "
						mapStats := "**CS:** " + strconv.FormatFloat(beatmap.CircleSize, 'f', 1, 64) + " **AR:** " + strconv.FormatFloat(beatmap.ApproachRate, 'f', 1, 64) + " **OD:** " + strconv.FormatFloat(beatmap.OverallDifficulty, 'f', 1, 64) + " **HP:** " + strconv.FormatFloat(beatmap.HPDrain, 'f', 1, 64)
						mapObjs := "**Circles:** " + strconv.Itoa(beatmap.Circles) + " **Sliders:** " + strconv.Itoa(beatmap.Sliders) + " **Spinners:** " + strconv.Itoa(beatmap.Spinners)
						var combo string

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

						mapCompletion := ""
						playNum := 101
						orderedScores, err := osuAPI.GetUserBest(osuapi.GetUserScoresOpts{
							Username: user.Username,
							Limit:    100,
						})
						tools.ErrRead(err)
						for i, orderedScore := range orderedScores {
							if score.Score.Score == orderedScore.Score.Score {
								playNum = i + 1
								mapCompletion = "**#" + strconv.Itoa(i+1) + "** in top performances! \n"
							}
						}
						if mapCompletion == "" {
							mapScores, err := osuAPI.GetScores(osuapi.GetScoresOpts{
								BeatmapID: beatmap.BeatmapID,
								Limit:     100,
							})
							tools.ErrRead(err)
							for i, mapScore := range mapScores {
								if score.UserID == mapScore.UserID && score.Score.Score == mapScore.Score.Score {
									playNum = i + 1
									mapCompletion = "**#" + strconv.Itoa(i+1) + "** on leaderboard! \n"
								}
							}
						}

						// Get pp values
						var pp string
						var ppInt int
						totalObjs := beatmap.Circles + beatmap.Sliders + beatmap.Spinners
						if score.PP == 0 { // If map was not finished
							ppValues := make(chan string, 2)
							var ppValueArray [2]string
							accCalcNoMiss := (50.0*float64(score.Count50) + 100.0*float64(score.Count100) + 300.0*float64(totalObjs-score.Count50-score.Count100)) / (300.0 * float64(totalObjs)) * 100.0
							go PPCalc(beatmap, accCalcNoMiss, "", "", scoreMods, ppValues)
							go PPCalc(beatmap, accCalc, strconv.Itoa(score.MaxCombo), strconv.Itoa(score.CountMiss), scoreMods, ppValues)
							for v := 0; v < 2; v++ {
								ppValueArray[v] = <-ppValues
							}
							sort.Slice(ppValueArray[:], func(i, j int) bool {
								pp1, _ := strconv.ParseFloat(ppValueArray[i], 64)
								pp2, _ := strconv.ParseFloat(ppValueArray[j], 64)
								return pp1 > pp2
							})
							ppInt, _ = strconv.Atoi(ppValueArray[1])
							pp = "**" + ppValueArray[1] + "pp**/" + ppValueArray[0] + "pp "
						} else if score.Score.FullCombo { // If play was a perfect combo
							ppInt = int(score.PP)
							pp = "**" + strconv.FormatFloat(score.PP, 'f', 0, 64) + "pp**/" + strconv.FormatFloat(score.PP, 'f', 0, 64) + "pp "
						} else { // If map was finished, but play was not a perfect combo
							ppInt = int(score.PP)
							ppValues := make(chan string, 1)
							accCalcNoMiss := (50.0*float64(score.Count50) + 100.0*float64(score.Count100) + 300.0*float64(totalObjs-score.Count50-score.Count100)) / (300.0 * float64(totalObjs)) * 100.0
							go PPCalc(beatmap, accCalcNoMiss, "", "", scoreMods, ppValues)
							pp = "**" + strconv.FormatFloat(score.PP, 'f', 0, 64) + "pp**/" + <-ppValues + "pp "
						}

						// Check params
						check1 := (ch.PPLimit != 0 && ch.TopPlay != 100 && ppInt < ch.PPLimit && playNum > ch.TopPlay)
						check2 := (ch.PPLimit == 0 && playNum > ch.TopPlay)
						check3 := (ch.TopPlay == 100 && ppInt < ch.PPLimit)
						if !check1 && !check2 && !check3 {
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
									Text: "Try #" + strconv.Itoa(try) + " | " + timeText,
								},
							}

							recentUser, err := osuAPI.GetUser(osuapi.GetUserOpts{
								UserID: user.UserID,
							})
							tools.ErrRead(err)

							if recentUser.PP != user.PP {
								embed.Author.Name = user.Username + " " + strconv.FormatFloat(user.PP, 'f', 2, 64) + " -> " + strconv.FormatFloat(recentUser.PP, 'f', 2, 64) + " (+" + strconv.FormatFloat(recentUser.PP-user.PP, 'f', 2, 64) + "pp)"
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
								fmt.Println(err)
								fmt.Println(ch)
							}
							ch.Users[l] = *recentUser
						}
					}
				}
			}

			jsonData, err := json.Marshal(ch)
			err = ioutil.WriteFile("./"+channel, jsonData, 0644)
			tools.ErrRead(err)

			startTime = time.Now()
		}
	}
}
