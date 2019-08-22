package osucommands

import (
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"

	osutools "../../osu-functions"
	structs "../../structs"
	tools "../../tools"
	"github.com/bwmarrin/discordgo"
	"github.com/thehowl/go-osuapi"
)

// Top gets the nth top pp score
func Top(s *discordgo.Session, m *discordgo.MessageCreate, args []string, osuAPI *osuapi.Client, cache []structs.PlayerData, serverPrefix string, mapCache []structs.MapData) {
	emptyUser := osuapi.User{}
	username := ""
	mods := ""
	index := 1
	var APIMods osuapi.Mods
	var emptyAPI osuapi.Mods

	if len(args) > 1 {
		if args[0] == serverPrefix+"osu" && len(args) > 2 {
			check, err := strconv.Atoi(args[2])
			if err == nil {
				index = check
				if len(args) > 3 {
					if args[3] == "-m" && len(args) > 4 {
						mods = strings.ToUpper(args[4])
					} else if len(args[3])%2 == 0 && osuapi.ParseMods(strings.ToUpper(args[3])).String() == strings.ToUpper(args[3]) {
						mods = strings.ToUpper(args[3])
					}
				}
			} else if args[2] == "-m" && len(args) > 3 {
				mods = strings.ToUpper(args[3])
				if len(args) > 4 {
					check, err = strconv.Atoi(args[4])
					if err == nil {
						index = check
					}
				}
			} else if len(args[2])%2 == 0 && osuapi.ParseMods(strings.ToUpper(args[2])).String() == strings.ToUpper(args[2]) {
				mods = strings.ToUpper(args[2])
				if len(args) > 3 {
					check, err = strconv.Atoi(args[3])
					if err == nil {
						index = check
					}
				}
			} else {
				username = args[2]
				if len(args) > 3 {
					check, err := strconv.Atoi(args[3])
					if err == nil {
						index = check
						if len(args) > 4 {
							if args[4] == "-m" && len(args) > 5 {
								mods = strings.ToUpper(args[5])
							} else {
								mods = strings.ToUpper(args[4])
							}
						}
					} else {
						if args[3] == "-m" && len(args) > 4 {
							mods = strings.ToUpper(args[4])
							if len(args) > 5 {
								check, err = strconv.Atoi(args[5])
								if err == nil {
									index = check
								}
							}
						} else {
							mods = strings.ToUpper(args[3])
							if len(args) > 4 {
								check, err = strconv.Atoi(args[4])
								if err == nil {
									index = check
								}
							}
						}
					}
				}
			}
		} else {
			check, err := strconv.Atoi(args[1])
			if err == nil {
				index = check
				if len(args) > 2 {
					if args[2] == "-m" && len(args) > 3 {
						mods = strings.ToUpper(args[3])
					} else if len(args[2])%2 == 0 && osuapi.ParseMods(strings.ToUpper(args[2])).String() == strings.ToUpper(args[2]) {
						mods = strings.ToUpper(args[2])
					}
				}
			} else if args[1] == "-m" && len(args) > 2 {
				mods = strings.ToUpper(args[2])
				if len(args) > 3 {
					check, err = strconv.Atoi(args[3])
					if err == nil {
						index = check
					}
				}
			} else if len(args[1])%2 == 0 && osuapi.ParseMods(strings.ToUpper(args[1])).String() == strings.ToUpper(args[1]) {
				mods = strings.ToUpper(args[1])
				if len(args) > 2 {
					check, err = strconv.Atoi(args[2])
					if err == nil {
						index = check
					}
				}
			} else {
				username = args[1]
				if len(args) > 2 {
					check, err := strconv.Atoi(args[2])
					if err == nil {
						index = check
						if len(args) > 3 {
							if args[3] == "-m" && len(args) > 4 {
								mods = strings.ToUpper(args[4])
							} else {
								mods = strings.ToUpper(args[3])
							}
						}
					} else {
						if args[2] == "-m" && len(args) > 3 {
							mods = strings.ToUpper(args[3])
							if len(args) > 4 {
								check, err = strconv.Atoi(args[4])
								if err == nil {
									index = check
								}
							}
						} else {
							mods = strings.ToUpper(args[2])
							if len(args) > 3 {
								check, err = strconv.Atoi(args[3])
								if err == nil {
									index = check
								}
							}
						}
					}
				}
			}
		}
	}

	if mods != "" {
		if len(mods)%2 == 0 && len(osuapi.ParseMods(mods).String()) > 0 {
			if strings.Contains(mods, "NC") && !strings.Contains(mods, "DTNC") {
				mods = strings.Replace(mods, "NC", "DTNC", 1)
			}
			APIMods = osuapi.ParseMods(mods)
			mods = APIMods.String()
		} else {
			APIMods = 0
			mods = "NM"
		}
	}

	for _, player := range cache {
		if username != "" || (m.Author.ID == player.Discord.ID && player.Osu.Username != emptyUser.Username) {
			// Check for user
			user := osuapi.User{}
			if username == "" {
				username = player.Osu.Username
				user = player.Osu
			} else {
				userP, err := osuAPI.GetUser(osuapi.GetUserOpts{
					Username: username,
				})
				if err != nil {
					s.ChannelMessageSend(m.ChannelID, "User "+username+" may not exist!")
					return
				}
				user = *userP
				go osutools.PlayerCache(user, cache)
			}
			score := osuapi.GUSScore{}

			// Get best scores
			scoreList, err := osuAPI.GetUserBest(osuapi.GetUserScoresOpts{
				Username: username,
				Limit:    100,
			})
			if err != nil {
				s.ChannelMessageSend(m.ChannelID, "The osu! API just owned me. Please try again!")
				return
			}
			if len(scoreList) == 0 {
				s.ChannelMessageSend(m.ChannelID, username+" has no top scores!")
				return
			}

			warning := ""
			if len(scoreList) < index {
				index = len(scoreList)
				warning = "Defaulted to max: " + strconv.Itoa(len(scoreList))
			}
			if APIMods != emptyAPI {
				num := 1
				for _, scor := range scoreList {
					if scor.Mods == APIMods {
						if num != index {
							num++
						} else if num == index {
							score = scor
							num++
						}
					}
				}
			} else {
				score = scoreList[index-1]
			}

			if score == (osuapi.GUSScore{}) {
				if strings.Contains(mods, "DTNC") {
					mods = strings.Replace(mods, "DTNC", "NC", 1)
				}
				s.ChannelMessageSend(m.ChannelID, "No score found in your best performances with the mods: **"+mods+"**")
				return
			}

			// Get beatmap, acc, and mods
			beatmap := osutools.BeatmapParse(strconv.Itoa(score.BeatmapID), "map", osuAPI)
			accCalc := (50.0*float64(score.Count50) + 100.0*float64(score.Count100) + 300.0*float64(score.Count300)) / (300.0 * float64(score.CountMiss+score.Count50+score.Count100+score.Count300)) * 100.0
			mods = score.Mods.String()
			if mods == "" {
				mods = "NM"
			}

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

			// Get pp values
			var pp string
			if score.Score.FullCombo { // If play was a perfect combo
				pp = "**" + strconv.FormatFloat(score.PP, 'f', 0, 64) + "pp**/" + strconv.FormatFloat(score.PP, 'f', 0, 64) + "pp "
			} else { // If play wasn't a perfect combo
				ppValues := make(chan string, 1)
				accCalcNoMiss := (50.0*float64(score.Count50) + 100.0*float64(score.Count100) + 300.0*float64(score.Count300+score.CountMiss)) / (300.0 * float64(score.Count50+score.Count100+score.Count300+score.CountMiss)) * 100.0
				go osutools.PPCalc(beatmap, accCalcNoMiss, "", "", scoreMods, ppValues)
				pp = "**" + strconv.FormatFloat(score.PP, 'f', 0, 64) + "pp**/" + <-ppValues + "pp "
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
					IconURL: "https://a.ppy.sh/" + strconv.Itoa(user.UserID),
				},
				Title: beatmap.Artist + " - " + beatmap.Title + " [" + beatmap.DiffName + "] by " + beatmap.Creator,
				URL:   "https://osu.ppy.sh/beatmaps/" + strconv.Itoa(beatmap.BeatmapID),
				Description: sr + length + bpm + "\n" +
					mapStats + "\n\n" +
					scorePrint + mods + combo + acc + scoreRank + "\n" +
					mapCompletion + mapCompletion2 + "\n" +
					pp + hits + "\n\n",
				Thumbnail: &discordgo.MessageEmbedThumbnail{
					URL: "https://b.ppy.sh/thumb/" + strconv.Itoa(beatmap.BeatmapSetID) + "l.jpg",
				},
				Footer: &discordgo.MessageEmbedFooter{
					Text: time,
				},
			}
			if strings.ToLower(beatmap.Title) == "crab rave" {
				embed.Image = &discordgo.MessageEmbedImage{
					URL: "https://cdn.discordapp.com/emojis/510169818893385729.gif",
				}
			}
			s.ChannelMessageSendComplex(m.ChannelID, &discordgo.MessageSend{
				Content: warning,
				Embed:   embed,
			})
			return
		}
	}

	s.ChannelMessageSend(m.ChannelID, "Could not find any osu! account linked for "+m.Author.Mention()+" !")
	return
}
