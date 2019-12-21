package osucommands

import (
	"fmt"
	"io/ioutil"
	"math"
	"math/rand"
	"net/http"
	"regexp"
	"sort"
	"strconv"
	"strings"

	config "../../config"
	osuapi "../../osu-api"
	osutools "../../osu-functions"
	structs "../../structs"
	tools "../../tools"
	"github.com/bwmarrin/discordgo"
)

// ReplayMessage posts replay information fopr a given replay
func ReplayMessage(s *discordgo.Session, m *discordgo.MessageCreate, linkRegex *regexp.Regexp, mapCache []structs.MapData) {
	// Get URL
	url := ""
	if len(m.Attachments) > 0 {
		url = m.Attachments[0].URL
	} else if linkRegex.MatchString(m.Content) {
		url = linkRegex.FindStringSubmatch(m.Content)[0]
	} else {
		return
	}

	if !strings.HasSuffix(url, ".osr") {
		return
	}

	// Get byte array
	res, err := http.Get(url)
	if err != nil {
		return
	}
	defer res.Body.Close()

	replayInfo, err := ioutil.ReadAll(res.Body)
	if err != nil || len(replayInfo) <= 81 {
		return
	}

	// Parse replay data
	replay := structs.ReplayData{
		Data: replayInfo,
	}
	replay.ParseReplay(OsuAPI)
	if replay.Beatmap.BeatmapID != 0 {
		osutools.BeatmapParse(strconv.Itoa(replay.Beatmap.BeatmapID), "map")
	}

	// Get time since play
	time := tools.TimeSince(replay.Time)

	// Assign timing variables for map specs
	totalMinutes := math.Floor(float64(replay.Beatmap.TotalLength / 60))
	totalSeconds := fmt.Sprint(math.Mod(float64(replay.Beatmap.TotalLength), float64(60)))
	if len(totalSeconds) == 1 {
		totalSeconds = "0" + totalSeconds
	}
	hitMinutes := math.Floor(float64(replay.Beatmap.HitLength / 60))
	hitSeconds := fmt.Sprint(math.Mod(float64(replay.Beatmap.HitLength), float64(60)))
	if len(hitSeconds) == 1 {
		hitSeconds = "0" + hitSeconds
	}

	// Assign values
	mods := replay.Score.Mods.String()
	accCalc := (50.0*float64(replay.Score.Count50) + 100.0*float64(replay.Score.Count100) + 300.0*float64(replay.Score.Count300)) / (300.0 * float64(replay.Score.CountMiss+replay.Score.Count50+replay.Score.Count100+replay.Score.Count300)) * 100.0
	Color := osutools.ModeColour(replay.Beatmap.Mode)
	sr := "**SR:** " + strconv.FormatFloat(replay.Beatmap.DifficultyRating, 'f', 2, 64) + " "
	length := "**Length:** " + fmt.Sprint(totalMinutes) + ":" + fmt.Sprint(totalSeconds) + " (" + fmt.Sprint(hitMinutes) + ":" + fmt.Sprint(hitSeconds) + ") "
	bpm := "**BPM:** " + fmt.Sprint(replay.Beatmap.BPM) + " "
	scorePrint := " **" + tools.Comma(replay.Score.Score) + "** "
	mapStats := "**CS:** " + strconv.FormatFloat(replay.Beatmap.CircleSize, 'f', 1, 64) + " **AR:** " + strconv.FormatFloat(replay.Beatmap.ApproachRate, 'f', 1, 64) + " **OD:** " + strconv.FormatFloat(replay.Beatmap.OverallDifficulty, 'f', 1, 64) + " **HP:** " + strconv.FormatFloat(replay.Beatmap.HPDrain, 'f', 1, 64)
	mapObjs := "**Circles:** " + strconv.Itoa(replay.Beatmap.Circles) + " **Sliders:** " + strconv.Itoa(replay.Beatmap.Sliders) + " **Spinners:** " + strconv.Itoa(replay.Beatmap.Spinners)
	acc := "** " + strconv.FormatFloat(accCalc, 'f', 2, 64) + "%** "
	hits := "**Hits:** [" + strconv.Itoa(replay.Score.Count300) + "/" + strconv.Itoa(replay.Score.Count100) + "/" + strconv.Itoa(replay.Score.Count50) + "/" + strconv.Itoa(replay.Score.CountMiss) + "]"
	status := "**Rank Status:** " + strings.Title(replay.Beatmap.Approved.String())
	if mods == "" {
		mods = "NM"
	}

	if strings.Contains(mods, "DTNC") {
		mods = strings.Replace(mods, "DTNC", "NC", 1)
	}

	var combo string
	if replay.Score.MaxCombo == replay.Beatmap.MaxCombo {
		if accCalc == 100.0 {
			combo = " **SS** "
		} else {
			combo = " **FC** "
		}
	} else if replay.Beatmap.MaxCombo != 0 {
		combo = " **x" + strconv.Itoa(replay.Score.MaxCombo) + "**/" + strconv.Itoa(replay.Beatmap.MaxCombo) + " "
	}

	mapCompletion := ""
	scoreRank := replay.Score.Rank
	if replay.Beatmap.Approved > 0 {
		orderedScores, err := OsuAPI.GetUserBest(osuapi.GetUserScoresOpts{
			Username: replay.Player.Username,
			Limit:    100,
		})
		if err == nil {
			for i, orderedScore := range orderedScores {
				if replay.Score.Score == orderedScore.Score.Score {
					replay.Score.Rank = orderedScore.Score.Rank
					mapCompletion += "**#" + strconv.Itoa(i+1) + "** in top performances! \n"
					break
				}
			}
		}
		mapScores, err := OsuAPI.GetScores(osuapi.GetScoresOpts{
			BeatmapID: replay.Beatmap.BeatmapID,
			Limit:     100,
		})
		if err == nil {
			for i, mapScore := range mapScores {
				if replay.Player.UserID == mapScore.UserID && replay.Score.Score == mapScore.Score.Score {
					replay.Score.Rank = mapScore.Rank
					mapCompletion += "**#" + strconv.Itoa(i+1) + "** on leaderboard! \n"
					break
				}
			}
		}
	}

	g, _ := s.Guild(config.Conf.Server)
		tools.ErrRead(err)
		for _, emoji := range g.Emojis {
			if emoji.Name == replay.Score.Rank+"_" {
				scoreRank = emoji.MessageFormat()
			}
		}

	ppValues := make(chan string, 2)
	var ppValueArray [2]string
	totalObjs := replay.Beatmap.Circles + replay.Beatmap.Sliders + replay.Beatmap.Spinners
	accCalcNoMiss := (50.0*float64(replay.Score.Count50) + 100.0*float64(replay.Score.Count100) + 300.0*float64(totalObjs-replay.Score.Count50-replay.Score.Count100)) / (300.0 * float64(totalObjs)) * 100.0
	go osutools.PPCalc(replay.Beatmap, accCalcNoMiss, "", "", mods, ppValues)
	go osutools.PPCalc(replay.Beatmap, accCalc, strconv.Itoa(replay.Score.MaxCombo), strconv.Itoa(replay.Score.CountMiss), mods, ppValues)
	for v := 0; v < 2; v++ {
		ppValueArray[v] = <-ppValues
	}
	sort.Slice(ppValueArray[:], func(i, j int) bool {
		pp1, _ := strconv.ParseFloat(ppValueArray[i], 64)
		pp2, _ := strconv.ParseFloat(ppValueArray[j], 64)
		return pp1 > pp2
	})
	pp := "**" + ppValueArray[1] + "pp**/" + ppValueArray[0] + "pp "
	mods = " **+" + mods + "** "

	// Create embed
	var embed = &discordgo.MessageEmbed{}
	if replay.Beatmap.BeatmapID == 0 {
		embed = &discordgo.MessageEmbed{
			Color: Color,
			Author: &discordgo.MessageEmbedAuthor{
				URL:     "https://osu.ppy.sh/users/" + strconv.Itoa(replay.Player.UserID),
				Name:    replay.Player.Username,
				IconURL: "https://a.ppy.sh/" + strconv.Itoa(replay.Player.UserID) + "?" + strconv.Itoa(rand.Int()) + ".jpeg",
			},
			Title: "Unknown / Unsubmitted map",
			Description: scorePrint + mods + combo + acc + scoreRank + "\n" +
				mapCompletion + "\n" +
				hits + "\n\n",
			Footer: &discordgo.MessageEmbedFooter{
				Text: time,
			},
		}
	} else {
		embed = &discordgo.MessageEmbed{
			Color: Color,
			Author: &discordgo.MessageEmbedAuthor{
				URL:     "https://osu.ppy.sh/users/" + strconv.Itoa(replay.Player.UserID),
				Name:    replay.Player.Username,
				IconURL: "https://a.ppy.sh/" + strconv.Itoa(replay.Player.UserID) + "?" + strconv.Itoa(rand.Int()) + ".jpeg",
			},
			Title: replay.Beatmap.Artist + " - " + replay.Beatmap.Title + " [" + replay.Beatmap.DiffName + "] by " + replay.Beatmap.Creator,
			URL:   "https://osu.ppy.sh/beatmaps/" + strconv.Itoa(replay.Beatmap.BeatmapID),
			Thumbnail: &discordgo.MessageEmbedThumbnail{
				URL: "https://b.ppy.sh/thumb/" + strconv.Itoa(replay.Beatmap.BeatmapSetID) + "l.jpg",
			},
			Description: sr + length + bpm + "\n" +
				mapStats + "\n" +
				mapObjs + "\n" + 
				status + "\n\n" +
				scorePrint + mods + combo + acc + scoreRank + "\n" +
				mapCompletion + "\n" +
				pp + hits + "\n\n",
			Footer: &discordgo.MessageEmbedFooter{
				Text: time,
			},
		}
		if strings.ToLower(replay.Beatmap.Title) == "crab rave" {
			embed.Image = &discordgo.MessageEmbedImage{
				URL: "https://cdn.discordapp.com/emojis/510169818893385729.gif",
			}
		}
	}
	if replay.Player.UserID == 0 {
		if replay.Score.Mods&osuapi.ModAutoplay != 0 {
			embed.Author = &discordgo.MessageEmbedAuthor{
				Name:    "osu!",
				IconURL: "https://osu.ppy.sh/images/layout/avatar-guest.png",
			}
		} else {
			embed.Author = &discordgo.MessageEmbedAuthor{
				Name:    "Unknown player",
				IconURL: "https://osu.ppy.sh/images/layout/avatar-guest.png",
			}
		}
	}
	s.ChannelMessageSendEmbed(m.ChannelID, embed)
}
