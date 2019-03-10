package osucommands

import (
	"fmt"
	"math"
	"sort"
	"strconv"
	"time"

	osutools "../../osu-functions"
	structs "../../structs"
	tools "../../tools"
	"github.com/bwmarrin/discordgo"
	"github.com/thehowl/go-osuapi"
)

// Recent gets the most recent score done/nth score done
func Recent(s *discordgo.Session, m *discordgo.MessageCreate, args []string, osuAPI *osuapi.Client, cache []structs.PlayerData, option string, mapCache []structs.MapData) {
	emptyUser := osuapi.User{}
	index := 1
	username := ""

	if len(args) > 1 {
		check, err := strconv.Atoi(args[1])
		if err != nil {
			username = args[1]
			if len(args) > 2 {
				check, err = strconv.Atoi(args[2])
				if err != nil {
					s.ChannelMessageSend(m.ChannelID, "Please use _ for username spaces!")
					return
				}
			} else {
				check = 1
			}
		}
		index = check
	}

	for _, player := range cache {
		if username == "" || (m.Author.ID == player.Discord.ID && player.Osu.Username != emptyUser.Username) {
			user := osuapi.User{}
			if username == "" {
				username = player.Osu.Username
				user = player.Osu
			} else {
				userP, err := osuAPI.GetUser(osuapi.GetUserOpts{
					Username: username,
				})
				if err != nil {
					s.ChannelMessageSend(m.ChannelID, "User "+username+" not found!")
					return
				}
				user = *userP
				go osutools.PlayerCache(user, cache)
			}
			score := osuapi.GUSScore{}
			scoreList := []osuapi.GUSScore{}

			// Run api call for user best/recent
			if option == "recent" {
				scores, err := osuAPI.GetUserRecent(osuapi.GetUserScoresOpts{
					Username: username,
					Limit:    50,
				})
				tools.ErrRead(err)
				if len(scores) == 0 {
					s.ChannelMessageSend(m.ChannelID, username+" has not played recently!")
					return
				}

				scoreList = scores
			} else if option == "best" {
				scores, err := osuAPI.GetUserBest(osuapi.GetUserScoresOpts{
					Username: username,
					Limit:    100,
				})
				tools.ErrRead(err)
				if len(scores) == 0 {
					s.ChannelMessageSend(m.ChannelID, username+" has not played recently!")
					return
				}

				scoreList = scores
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

			if len(scoreList) < index {
				index = len(scoreList)
			}
			score = scoreList[index-1]

			// Get beatmap, and mods
			beatmap := osutools.BeatmapParse(strconv.Itoa(score.BeatmapID), "map", osuAPI)

			mods := "NM"
			if score.Mods != 0 {
				mods = score.Mods.String()
			}

			// Get pp values
			ppValues := make(chan string, 2)
			var ppValueArray [2]string
			accCalcNoMiss := (50.0*float64(score.Count50) + 100.0*float64(score.Count100) + 300.0*float64(score.Count300+score.CountMiss)) / (300.0 * float64(score.Count50+score.Count100+score.Count300)) * 100.0
			accCalc := (50.0*float64(score.Count50) + 100.0*float64(score.Count100) + 300.0*float64(score.Count300)) / (300.0 * float64(score.CountMiss+score.Count50+score.Count100+score.Count300)) * 100.0
			go osutools.PPCalc(beatmap, accCalcNoMiss, "", "", score.Mods.String(), ppValues)
			go osutools.PPCalc(beatmap, accCalc, strconv.Itoa(score.MaxCombo), strconv.Itoa(score.CountMiss), score.Mods.String(), ppValues)
			for v := 0; v < 2; v++ {
				ppValueArray[v] = <-ppValues
			}
			sort.Slice(ppValueArray[:], func(i, j int) bool {
				return ppValueArray[i] > ppValueArray[j]
			})

			// Count number of tries
			try := 1
			for i := index; i < len(scoreList); i++ {
				if scoreList[i].BeatmapID == score.BeatmapID {
					try++
				} else {
					break
				}
			}

			// Assign variables for map specs
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

			Color := osutools.ModeColour(osuapi.ModeOsu)
			sr, _, _, _, _, _ := osutools.BeatmapCache(mods, beatmap, mapCache)
			length := "**Length:** " + fmt.Sprint(totalMinutes) + ":" + fmt.Sprint(totalSeconds) + " (" + fmt.Sprint(hitMinutes) + ":" + fmt.Sprint(hitSeconds) + ") "
			bpm := "**BPM:** " + fmt.Sprint(beatmap.BPM) + " "

			mods = "**Mods:** " + mods + " "
			scorePrint := "**Score:** " + strconv.FormatInt(score.Score.Score, 10) + " "
			combo := "**Combo:** " + strconv.Itoa(score.MaxCombo) + "/" + strconv.Itoa(beatmap.MaxCombo) + "x"

			pp := "**PP:** " + ppValueArray[1] + " / " + ppValueArray[0] + " "
			acc := "**Acc:** " + strconv.FormatFloat(accCalc, 'f', 2, 64) + "%"

			if option == "best" {
				s.ChannelMessageSendEmbed(m.ChannelID, &discordgo.MessageEmbed{
					Color: Color,
					Author: &discordgo.MessageEmbedAuthor{
						URL:     "https://osu.ppy.sh/users/" + strconv.Itoa(user.UserID),
						Name:    user.Username,
						IconURL: "https://a.ppy.sh/" + strconv.Itoa(user.UserID),
					},
					Title: beatmap.Artist + " - " + beatmap.Title + " [" + beatmap.DiffName + "]",
					URL:   "https://osu.ppy.sh/beatmaps/" + strconv.Itoa(beatmap.BeatmapID),
					Thumbnail: &discordgo.MessageEmbedThumbnail{
						URL: "https://b.ppy.sh/thumb/" + strconv.Itoa(beatmap.BeatmapSetID) + "l.jpg",
					},
					Description: sr + length + bpm + "\n" +
						mods + scorePrint + combo + "\n" +
						pp + acc,
				})
				return
			}
			message, err := s.ChannelMessageSend(m.ChannelID, "Try #"+strconv.Itoa(try))
			tools.ErrRead(err)
			s.ChannelMessageEditEmbed(m.ChannelID, message.ID, &discordgo.MessageEmbed{
				Color: Color,
				Author: &discordgo.MessageEmbedAuthor{
					URL:     "https://osu.ppy.sh/users/" + strconv.Itoa(user.UserID),
					Name:    user.Username,
					IconURL: "https://a.ppy.sh/" + strconv.Itoa(user.UserID),
				},
				Title: beatmap.Artist + " - " + beatmap.Title + " [" + beatmap.DiffName + "]",
				URL:   "https://osu.ppy.sh/beatmaps/" + strconv.Itoa(beatmap.BeatmapID),
				Thumbnail: &discordgo.MessageEmbedThumbnail{
					URL: "https://b.ppy.sh/thumb/" + strconv.Itoa(beatmap.BeatmapSetID) + "l.jpg",
				},
				Description: sr + length + bpm + "\n" +
					mods + scorePrint + combo + "\n" +
					pp + acc,
			})
			return
		}
	}

	s.ChannelMessageSend(m.ChannelID, "Could not find any osu! account linked for "+m.Author.Mention()+" !")
	return
}
