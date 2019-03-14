package osucommands

import (
	"fmt"
	"math"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	osutools "../../osu-functions"
	structs "../../structs"
	tools "../../tools"
	"github.com/bwmarrin/discordgo"
	"github.com/thehowl/go-osuapi"
)

// Compare compares finds a score from the current user on the previous map linked by the bot
func Compare(s *discordgo.Session, m *discordgo.MessageCreate, args []string, osuAPI *osuapi.Client, cache []structs.PlayerData, serverPrefix string, mapCache []structs.MapData) {
	userArg := ""
	if len(args) > 1 {
		userArg = args[1]
	}

	// Get prev messages
	messages, err := s.ChannelMessages(m.ChannelID, 100, "", "", "")
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
	var mapID int
	foundMap := false
	for _, msg := range messages {
		msgTime, err := time.Parse(time.RFC3339, string(msg.Timestamp))
		if msg.Author.ID == s.State.User.ID && time.Since(msgTime) < time.Hour {
			if msg.ID != (discordgo.Message{}).ID && len(msg.Embeds) > 0 {
				if beatmapRegex.MatchString(msg.Embeds[0].URL) {
					foundMap = true
					mapID, err = strconv.Atoi(beatmapRegex.FindStringSubmatch(msg.Embeds[0].URL)[1])
					tools.ErrRead(err)
					break
				} else if beatmapRegex.MatchString(msg.Embeds[0].Author.URL) {
					foundMap = true
					mapID, err = strconv.Atoi(beatmapRegex.FindStringSubmatch(msg.Embeds[0].Author.URL)[1])
					tools.ErrRead(err)
					break
				}
			}
		}
	}
	if !foundMap {
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
		test, err := osuAPI.GetUser(osuapi.GetUserOpts{
			Username: userArg,
		})
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, "User "+userArg+" doesn't exist! Are you sure you replaced spaces with `_`?")
			return
		}
		user = *test
	}

	// API call
	scores, err := osuAPI.GetScores(osuapi.GetScoresOpts{
		BeatmapID: mapID,
		UserID:    user.UserID,
	})
	beatmaps, err := osuAPI.GetBeatmaps(osuapi.GetBeatmapsOpts{
		BeatmapID: mapID,
	})
	tools.ErrRead(err)
	beatmap := beatmaps[0]

	if len(scores) == 0 {
		if userArg != "" {
			s.ChannelMessageSend(m.ChannelID, userArg+" hasn't set a score on this!")
			return
		}
		s.ChannelMessageSend(m.ChannelID, "You haven't set a score on this!")
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
	mods := "NM"
	if score.Mods != 0 {
		mods = score.Mods.String()
	}
	accCalc := (50.0*float64(score.Count50) + 100.0*float64(score.Count100) + 300.0*float64(score.Count300)) / (300.0 * float64(score.CountMiss+score.Count50+score.Count100+score.Count300)) * 100.0
	Color := osutools.ModeColour(osuapi.ModeOsu)
	sr, _, _, _, _, _ := osutools.BeatmapCache(mods, beatmap, mapCache)
	length := "**Length:** " + fmt.Sprint(totalMinutes) + ":" + fmt.Sprint(totalSeconds) + " (" + fmt.Sprint(hitMinutes) + ":" + fmt.Sprint(hitSeconds) + ") "
	bpm := "**BPM:** " + fmt.Sprint(beatmap.BPM) + " "
	scorePrint := " **" + tools.Comma(score.Score.Score) + "** "
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

	// Get pp values
	var pp string
	if score.Score.FullCombo { // If play was a perfect combo
		pp = "**" + strconv.FormatFloat(score.PP, 'f', 0, 64) + "pp**/" + strconv.FormatFloat(score.PP, 'f', 0, 64) + "pp "
	} else { // If map was finished, but play was not a perfect combo
		ppValues := make(chan string, 1)
		accCalcNoMiss := (50.0*float64(score.Count50) + 100.0*float64(score.Count100) + 300.0*float64(score.Count300+score.CountMiss)) / (300.0 * float64(score.Count50+score.Count100+score.Count300)) * 100.0
		go osutools.PPCalc(beatmap, accCalcNoMiss, "", "", score.Mods.String(), ppValues)
		pp = "**" + strconv.FormatFloat(score.PP, 'f', 0, 64) + "pp**/" + <-ppValues + "pp "
	}
	acc := "**Acc:** " + strconv.FormatFloat(accCalc, 'f', 2, 64) + "% "

	hits := "**Hits:** [" + strconv.Itoa(score.Count300) + "/" + strconv.Itoa(score.Count100) + "/" + strconv.Itoa(score.Count50) + "/" + strconv.Itoa(score.CountMiss) + "]"

	s.ChannelMessageSendEmbed(m.ChannelID, &discordgo.MessageEmbed{
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
		Description: sr + length + bpm + "\n\n" +
			scorePrint + mods + combo + "\n\n" +
			pp + acc + hits + "\n\n",
		Footer: &discordgo.MessageEmbedFooter{
			Text: time,
		},
	})
	return
}
