package osutools

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	structs "../structs"
	tools "../tools"
	"github.com/bwmarrin/discordgo"
	"github.com/thehowl/go-osuapi"
)

// TrackPost posts scores tracked for that channel
func TrackPost(channel string, s *discordgo.Session, mapCache []structs.MapData) {
	osuAPI := osuapi.NewClient(os.Getenv("OSU_API"))

	startTime := time.Now()

	ticker := time.NewTicker(2 * time.Second)
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

			if !ch.Tracking {
				return
			}

			if len(ch.Users) != 0 {
				for _, user := range ch.Users {
					recentScores, err := osuAPI.GetUserRecent(osuapi.GetUserScoresOpts{
						UserID: user.UserID,
						Limit:  100,
					})
					tools.ErrRead(err)

					for _, score := range recentScores {
						scoreTime, err := time.Parse("2006-01-02 15:04:05", score.Date.String())
						tools.ErrRead(err)
						if score.Rank != "F" && startTime.Before(scoreTime) {
							beatmap := BeatmapParse(strconv.Itoa(score.BeatmapID), "map", osuAPI)
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
							mods := "NM"
							if score.Mods != 0 {
								mods = score.Mods.String()
							}
							accCalc := (50.0*float64(score.Count50) + 100.0*float64(score.Count100) + 300.0*float64(score.Count300)) / (300.0 * float64(score.CountMiss+score.Count50+score.Count100+score.Count300)) * 100.0
							Color := ModeColour(osuapi.ModeOsu)
							sr, _, _, _, _, _ := BeatmapCache(mods, beatmap, mapCache)
							length := "**Length:** " + fmt.Sprint(totalMinutes) + ":" + fmt.Sprint(totalSeconds) + " (" + fmt.Sprint(hitMinutes) + ":" + fmt.Sprint(hitSeconds) + ") "
							bpm := "**BPM:** " + fmt.Sprint(beatmap.BPM) + " "
							scorePrint := " **" + tools.Comma(score.Score.Score) + "** "
							mapStats := "**CS:** " + strconv.FormatFloat(beatmap.CircleSize, 'f', 1, 64) + " **AR:** " + strconv.FormatFloat(beatmap.ApproachRate, 'f', 1, 64) + " **OD:** " + strconv.FormatFloat(beatmap.OverallDifficulty, 'f', 1, 64) + " **HP:** " + strconv.FormatFloat(beatmap.HPDrain, 'f', 1, 64)
							var combo string

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
									if score.UserID == mapScore.UserID {
										playNum = i + 1
										mapCompletion = "**#" + strconv.Itoa(i+1) + "** on leaderboard! \n"
									}
								}
							}

							// Get pp values
							var pp string
							var ppInt int
							if score.PP == 0 { // If map was not finished
								ppValues := make(chan string, 2)
								var ppValueArray [2]string
								accCalcNoMiss := (50.0*float64(score.Count50) + 100.0*float64(score.Count100) + 300.0*float64(score.Count300+score.CountMiss)) / (300.0 * float64(score.Count50+score.Count100+score.Count300)) * 100.0
								go PPCalc(beatmap, accCalcNoMiss, "", "", score.Mods.String(), ppValues)
								go PPCalc(beatmap, accCalc, strconv.Itoa(score.MaxCombo), strconv.Itoa(score.CountMiss), score.Mods.String(), ppValues)
								for v := 0; v < 2; v++ {
									ppValueArray[v] = <-ppValues
								}
								sort.Slice(ppValueArray[:], func(i, j int) bool {
									pp1, _ := strconv.Atoi(ppValueArray[i])
									pp2, _ := strconv.Atoi(ppValueArray[j])
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
								accCalcNoMiss := (50.0*float64(score.Count50) + 100.0*float64(score.Count100) + 300.0*float64(score.Count300+score.CountMiss)) / (300.0 * float64(score.Count50+score.Count100+score.Count300)) * 100.0
								go PPCalc(beatmap, accCalcNoMiss, "", "", score.Mods.String(), ppValues)
								pp = "**" + strconv.FormatFloat(score.PP, 'f', 0, 64) + "pp**/" + <-ppValues + "pp "
							}

							// Check params
							check1 := (ch.PPLimit != 0 && ch.TopPlay != 100 && ppInt < ch.PPLimit && playNum > ch.TopPlay)
							check2 := (ch.PPLimit == 0 && playNum > ch.TopPlay)
							check3 := (ch.TopPlay == 100 && ppInt < ch.PPLimit)
							if !check1 && !check2 && !check3 {
								acc := "** " + strconv.FormatFloat(accCalc, 'f', 2, 64) + "%** "
								hits := "**Hits:** [" + strconv.Itoa(score.Count300) + "/" + strconv.Itoa(score.Count100) + "/" + strconv.Itoa(score.Count50) + "/" + strconv.Itoa(score.CountMiss) + "]"

								// Create embed
								embed := &discordgo.MessageEmbed{
									Color: Color,
									Author: &discordgo.MessageEmbedAuthor{
										URL:     "https://osu.ppy.sh/users/" + strconv.Itoa(user.UserID),
										Name:    user.Username,
										IconURL: "https://a.ppy.sh/" + strconv.Itoa(user.UserID),
									},
									Title: beatmap.Artist + " - " + beatmap.Title + " [" + beatmap.DiffName + "] by " + beatmap.Creator,
									URL:   "https://osu.ppy.sh/beatmaps/" + strconv.Itoa(beatmap.BeatmapID),
									Thumbnail: &discordgo.MessageEmbedThumbnail{
										URL: "https://b.ppy.sh/thumb/" + strconv.Itoa(beatmap.BeatmapSetID) + "l.jpg",
									},
									Description: sr + length + bpm + "\n" +
										mapStats + "\n\n" +
										scorePrint + mods + combo + acc + "\n" +
										mapCompletion + "\n" +
										pp + hits + "\n\n",
									Footer: &discordgo.MessageEmbedFooter{
										Text: timeText,
									},
								}
								if strings.ToLower(beatmap.Title) == "crab rave" {
									embed.Image = &discordgo.MessageEmbedImage{
										URL: "https://cdn.discordapp.com/emojis/510169818893385729.gif",
									}
								}
								s.ChannelMessageSendEmbed(ch.Channel.ID, embed)
							}
						}
					}
				}
			}
			startTime = time.Now()
		}
	}
}
