package osutools

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"math/rand"
	"strconv"
	"strings"
	"time"

	osuapi "../osu-api"
	structs "../structs"
	tools "../tools"
	"github.com/bwmarrin/discordgo"
)

// TrackMapperPost tracks mappers and posts their new maps into the channels
func TrackMapperPost(s *discordgo.Session) {
	startTime := time.Now()
	ticker := time.NewTicker(10 * time.Second)
	for {
		select {
		case <-ticker.C:
			// Obtain mapper data
			var mapperData []structs.MapperData
			f, err := ioutil.ReadFile("./data/osuData/mapperData.json")
			tools.ErrRead(err)
			_ = json.Unmarshal(f, &mapperData)

			for i := 0; i < len(mapperData); i++ {
				user, err := OsuAPI.GetUser(osuapi.GetUserOpts{
					UserID: mapperData[i].Mapper.UserID,
				})
				if err != nil {
					mapperData = append(mapperData[:i], mapperData[i+1:]...)
					i--
					continue
				}
				beatmaps, err := OsuAPI.GetBeatmaps(osuapi.GetBeatmapsOpts{
					UserID: user.UserID,
				})

				setsChecked := make(map[int]bool)
				for _, beatmap := range beatmaps {
					// see if set was already checked
					_, t := setsChecked[beatmap.BeatmapSetID]
					if t {
						continue
					}
					setsChecked[beatmap.BeatmapSetID] = true

					// see if set satisfies approvalstatus change/submission requirements
					osuMap := BeatmapParse(strconv.Itoa(beatmap.BeatmapSetID), "set")
					var targetMap osuapi.Beatmap
					var approvals []osuapi.ApprovedStatus // Store all approvedstatuses because sometimes the osu!api gives sets with multiple approved status
					for _, dataMap := range mapperData[i].Beatmaps {
						if dataMap.BeatmapSetID == osuMap.BeatmapSetID {
							targetMap = dataMap
							approvals = append(approvals, targetMap.Approved)
						}
					}
					for _, approval := range approvals {
						if osuMap.Approved == approval {
							targetMap.BeatmapSetID = 0
							break
						}
					}
					submitDate, _ := time.Parse("2006-01-02 15:04:05", osuMap.SubmitDate.String())
					if targetMap.BeatmapSetID == 0 && startTime.After(submitDate) {
						continue
					}

					// Assign embed colour for different modes
					Color := ModeColour(osuMap.Mode)

					// Obtain whole set
					set, _ := OsuAPI.GetBeatmaps(osuapi.GetBeatmapsOpts{
						BeatmapSetID: osuMap.BeatmapSetID,
					})

					// Assign variables for map specs
					totalMinutes := math.Floor(float64(osuMap.TotalLength / 60))
					totalSeconds := fmt.Sprint(math.Mod(float64(osuMap.TotalLength), float64(60)))
					if len(totalSeconds) == 1 {
						totalSeconds = "0" + totalSeconds
					}
					hitMinutes := math.Floor(float64(osuMap.HitLength / 60))
					hitSeconds := fmt.Sprint(math.Mod(float64(osuMap.HitLength), float64(60)))
					if len(hitSeconds) == 1 {
						hitSeconds = "0" + hitSeconds
					}

					length := "**Length:** " + fmt.Sprint(totalMinutes) + ":" + totalSeconds + " (" + fmt.Sprint(hitMinutes) + ":" + hitSeconds + ") "
					bpm := "**BPM:** " + fmt.Sprint(osuMap.BPM) + " "
					combo := "**FC:** " + strconv.Itoa(osuMap.MaxCombo) + "x"
					mapStats := "**CS:** " + strconv.FormatFloat(osuMap.CircleSize, 'f', 1, 64) + " **AR:** " + strconv.FormatFloat(osuMap.ApproachRate, 'f', 1, 64) + " **OD:** " + strconv.FormatFloat(osuMap.OverallDifficulty, 'f', 1, 64) + " **HP:** " + strconv.FormatFloat(osuMap.HPDrain, 'f', 1, 64)
					mapObjs := "**Circles:** " + strconv.Itoa(osuMap.Circles) + " **Sliders:** " + strconv.Itoa(osuMap.Sliders) + " **Spinners:** " + strconv.Itoa(osuMap.Spinners)

					status := "**Rank Status:** " + targetMap.Approved.String() + " -> **" + osuMap.Approved.String() + "**"

					download := "**Download:** [osz link](https://osu.ppy.sh/beatmapsets/" + strconv.Itoa(osuMap.BeatmapSetID) + "/download)" + " | <osu://dl/" + strconv.Itoa(osuMap.BeatmapSetID) + ">"
					var diffs string
					if len(beatmaps) == 1 {
						diffs = "**1** difficulty <:ahFuck:550808614202245131>"
					} else {
						diffs = "**" + strconv.Itoa(len(set)) + "** difficulties <:ahFuck:550808614202245131>"
					}

					// Calculate SR and PP
					values := BeatmapCalc("NM", "", "", "", osuMap)

					// Create embed
					embed := &discordgo.MessageEmbed{
						Author: &discordgo.MessageEmbedAuthor{
							URL:     "https://osu.ppy.sh/beatmaps/" + strconv.Itoa(osuMap.BeatmapID),
							Name:    osuMap.Artist + " - " + osuMap.Title + " by " + osuMap.Creator,
							IconURL: "https://a.ppy.sh/" + strconv.Itoa(osuMap.CreatorID) + "?" + strconv.Itoa(rand.Int()) + ".jpeg",
						},
						Color: Color,
						Description: values[0] + length + bpm + combo + "\n" +
							mapStats + "\n" +
							mapObjs + "\n" +
							status + "\n" +
							download + "\n" +
							diffs + "\n" + "\n" +
							"**[" + osuMap.DiffName + "]** with mods: **NM**\n" +
							//aimRating + speedRating + totalRating + "\n" + TODO: Make SR calc work
							values[1] + values[2] + values[3] + values[4] + values[5],
						Thumbnail: &discordgo.MessageEmbedThumbnail{
							URL: "https://b.ppy.sh/thumb/" + strconv.Itoa(beatmap.BeatmapSetID) + "l.jpg",
						},
					}
					if strings.ToLower(osuMap.Title) == "crab rave" {
						embed.Image = &discordgo.MessageEmbedImage{
							URL: "https://cdn.discordapp.com/emojis/510169818893385729.gif",
						}
					}
					text := ""
					if startTime.Before(submitDate) {
						text = "New map by " + user.Username + "!"
					} else {
						text = osuMap.Artist + " - " + osuMap.Title + " by " + osuMap.Creator + " has changed from **" + targetMap.Approved.String() + "** to **" + osuMap.Approved.String() + "!**"
					}
					for j := 0; j < len(mapperData[i].Channels); j++ {
						_, err = s.ChannelMessageSendComplex(mapperData[i].Channels[j], &discordgo.MessageSend{
							Content: text,
							Embed:   embed,
						})
						if err != nil {
							fmt.Println("TrackMapperPost err: " + err.Error())
							mapperData[j].Channels = append(mapperData[i].Channels[:j], mapperData[i].Channels[j+1:]...)
							j--
						}
					}
					break
				}
				mapperData[i].Mapper = *user
				mapperData[i].Beatmaps = beatmaps
			}

			// Save mapper data
			jsonCache, err := json.Marshal(mapperData)
			tools.ErrRead(err)

			err = ioutil.WriteFile("./data/osuData/mapperData.json", jsonCache, 0644)
			tools.ErrRead(err)

			startTime = time.Now()
		}
	}
}
