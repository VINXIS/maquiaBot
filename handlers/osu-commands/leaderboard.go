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

	config "../../config"
	osuapi "../../osu-api"
	osutools "../../osu-functions"
	structs "../../structs"
	tools "../../tools"
	"github.com/bwmarrin/discordgo"
)

// Leaderboard gives you the leaderboard for a map
func Leaderboard(s *discordgo.Session, m *discordgo.MessageCreate, regex *regexp.Regexp, cache []structs.PlayerData) {
	mapRegex, _ := regexp.Compile(`(https:\/\/)?(osu|old)\.ppy\.sh\/(s|b|beatmaps|beatmapsets)\/(\d+)(#(osu|taiko|fruits|mania)\/(\d+))?`)
	modRegex, _ := regexp.Compile(`-m\s*(\S+)`)
	showRegex, _ := regexp.Compile(`-n\s*(\d+)`)
	var beatmap osuapi.Beatmap
	var err error

	mods := ""
	parsedMods := osuapi.Mods(0)
	if modRegex.MatchString(m.Content) {
		mods = strings.ToUpper(modRegex.FindStringSubmatch(m.Content)[1])
		if strings.Contains(mods, "NC") && !strings.Contains(mods, "DT") {
			mods += "DT"
		}
		parsedMods = osuapi.ParseMods(mods)
	}

	count := 5
	if showRegex.MatchString(m.Content) {
		count, _ = strconv.Atoi(showRegex.FindStringSubmatch(m.Content)[1])
		if count > 100 {
			count = 100
		}
	}

	// Get the map
	var submatches []string
	if mapRegex.MatchString(m.Content) {
		submatches = mapRegex.FindStringSubmatch(m.Content)
	} else {
		// Get prev messages
		messages, err := s.ChannelMessages(m.ChannelID, -1, "", "", "")
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, "No map to compare to!")
			return
		}

		// Look for a valid beatmap ID
		for _, msg := range messages {
			if len(msg.Embeds) > 0 && msg.Embeds[0].Author != nil {
				if mapRegex.MatchString(msg.Embeds[0].URL) {
					submatches = mapRegex.FindStringSubmatch(msg.Embeds[0].URL)
					break
				} else if mapRegex.MatchString(msg.Embeds[0].Author.URL) {
					submatches = mapRegex.FindStringSubmatch(msg.Embeds[0].Author.URL)
					break
				}
			} else if mapRegex.MatchString(msg.Content) {
				submatches = mapRegex.FindStringSubmatch(msg.Content)
				break
			}
		}
	}

	// Check if found
	if len(submatches) == 0 {
		s.ChannelMessageSend(m.ChannelID, "No map to compare to!")
		return
	}

	// Get the map
	diffMods := 338 & parsedMods
	if diffMods&256 != 0 && diffMods&64 != 0 { // Remove DTHT
		diffMods -= 320
	}
	if diffMods&2 != 0 && diffMods&16 != 0 { // Remove EZHR
		diffMods -= 18
	}
	switch submatches[3] {
	case "s":
		beatmap = osutools.BeatmapParse(submatches[4], "set", &diffMods)
	case "b":
		beatmap = osutools.BeatmapParse(submatches[4], "map", &diffMods)
	case "beatmaps":
		beatmap = osutools.BeatmapParse(submatches[4], "map", &diffMods)
	case "beatmapsets":
		if len(submatches[7]) > 0 {
			beatmap = osutools.BeatmapParse(submatches[7], "map", &diffMods)
		} else {
			beatmap = osutools.BeatmapParse(submatches[4], "set", &diffMods)
		}
	}
	if beatmap.BeatmapID == 0 {
		s.ChannelMessageSend(m.ChannelID, "No map to compare to!")
		return
	} else if beatmap.Approved < 1 {
		s.ChannelMessageSend(m.ChannelID, "The map `"+beatmap.Artist+" - "+beatmap.Title+"` does not have a leaderboard!")
		return
	}

	// API call(s)
	var scores []osuapi.GSScore
	scoreOpts := osuapi.GetScoresOpts{
		BeatmapID: beatmap.BeatmapID,
		Limit:     100,
	}
	if mods != "" {
		scoreOpts.Mods = &parsedMods
	}
	if strings.Contains(m.Content, "-s") {
		members, err := s.GuildMembers(m.GuildID, "", 1000)
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, "This is not a server!")
			return
		}
		trueCache := []structs.PlayerData{}

		for _, player := range cache {
			for _, member := range members {
				if player.Discord.ID == member.User.ID && player.Osu.Username != "" {
					scoreOpts.UserID = player.Osu.UserID
					userScores, err := OsuAPI.GetScores(scoreOpts)
					if err == nil {
						scores = append(scores, userScores...)
					}
					break
				}
			}
		}

		cache = trueCache
	} else {
		scores, err = OsuAPI.GetScores(scoreOpts)
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, "The osu! API just owned me. Please try again!")
			return
		}
	}
	if len(scores) == 0 {
		s.ChannelMessageSend(m.ChannelID, "There are no scores on this map currently!")
		return
	}
	if count > len(scores) {
		count = len(scores)
	}

	// Sort by Score
	sort.Slice(scores, func(i, j int) bool {
		return scores[i].Score.Score > scores[j].Score.Score
	})

	// Assign variables for map specs
	Color := osutools.ModeColour(beatmap.Mode)
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

	sr := "**SR:** " + strconv.FormatFloat(beatmap.DifficultyRating, 'f', 2, 64) + " **Aim:** " + strconv.FormatFloat(beatmap.DifficultyAim, 'f', 2, 64) + " **Speed:** " + strconv.FormatFloat(beatmap.DifficultySpeed, 'f', 2, 64)
	length := "**Length:** " + fmt.Sprint(totalMinutes) + ":" + totalSeconds + " (" + fmt.Sprint(hitMinutes) + ":" + hitSeconds + ") "
	bpm := "**BPM:** " + fmt.Sprint(beatmap.BPM) + " "
	combo := "**FC:** " + strconv.Itoa(beatmap.MaxCombo) + "x"
	mapStats := "**CS:** " + strconv.FormatFloat(beatmap.CircleSize, 'f', 1, 64) + " **AR:** " + strconv.FormatFloat(beatmap.ApproachRate, 'f', 1, 64) + " **OD:** " + strconv.FormatFloat(beatmap.OverallDifficulty, 'f', 1, 64) + " **HP:** " + strconv.FormatFloat(beatmap.HPDrain, 'f', 1, 64)
	mapObjs := "**Circles:** " + strconv.Itoa(beatmap.Circles) + " **Sliders:** " + strconv.Itoa(beatmap.Sliders) + " **Spinners:** " + strconv.Itoa(beatmap.Spinners)

	status := "**Rank Status:** " + strings.Title(beatmap.Approved.String())

	download := "**Download:** [osz link](https://osu.ppy.sh/beatmapsets/" + strconv.Itoa(beatmap.BeatmapSetID) + "/download)" + " | <osu://dl/" + strconv.Itoa(beatmap.BeatmapSetID) + ">"

	// Create embed
	embed := &discordgo.MessageEmbed{
		Author: &discordgo.MessageEmbedAuthor{
			URL:     "https://osu.ppy.sh/beatmaps/" + strconv.Itoa(beatmap.BeatmapID),
			Name:    beatmap.Artist + " - " + beatmap.Title + " [" + beatmap.DiffName + "] by " + beatmap.Creator,
			IconURL: "https://a.ppy.sh/" + strconv.Itoa(beatmap.CreatorID) + "?" + strconv.Itoa(rand.Int()) + ".jpeg",
		},
		Color: Color,
		Description: sr + "\n" + 
			length + bpm + combo + "\n" +
			mapStats + "\n" +
			mapObjs + "\n" +
			status + "\n" +
			download,
		Thumbnail: &discordgo.MessageEmbedThumbnail{
			URL: "https://b.ppy.sh/thumb/" + strconv.Itoa(beatmap.BeatmapSetID) + "l.jpg",
		},
	}

	g, _ := s.Guild(config.Conf.Server)
	i := count
	if i > 25 {
		i = 25
	}
	initial := i
	for {
		for j, score := range scores[i-initial : i] {
			scoreRank := ""
			for _, emoji := range g.Emojis {
				if emoji.Name == score.Rank+"_" {
					scoreRank = emoji.MessageFormat()
				}
			}
			scorePrint := " **" + tools.Comma(score.Score.Score) + "** "
			scoreMods := score.Mods.String()
			if strings.Contains(scoreMods, "DTNC") {
				scoreMods = strings.Replace(scoreMods, "DTNC", "NC", -1)
			}
			accCalc := (50.0*float64(score.Count50) + 100.0*float64(score.Count100) + 300.0*float64(score.Count300)) / (300.0 * float64(score.CountMiss+score.Count50+score.Count100+score.Count300)) * 100.0
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
			acc := "** " + strconv.FormatFloat(accCalc, 'f', 2, 64) + "%** "
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
			hits := "[" + strconv.Itoa(score.Count300) + "/" + strconv.Itoa(score.Count100) + "/" + strconv.Itoa(score.Count50) + "/" + strconv.Itoa(score.CountMiss) + "]"
			timeParse, _ := time.Parse("2006-01-02 15:04:05", score.Date.String())
			time := tools.TimeSince(timeParse)
			scoreMods = " **+" + scoreMods + "** "

			embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
				Name: "#" + strconv.Itoa(i-initial+j+1) + " **" + score.Username + "** (" + strconv.Itoa(score.UserID) + ")",
				Value: scoreRank + scorePrint + scoreMods + combo + acc + "\n" +
					pp + hits + "\n" +
					time,
			})
		}
		s.ChannelMessageSendEmbed(m.ChannelID, embed)
		if i == count {
			break
		}
		i += 25
		initial = 25
		if i > count {
			initial = count - i + 25
			i = count
		}
		embed = &discordgo.MessageEmbed{Color: Color}
	}
}
