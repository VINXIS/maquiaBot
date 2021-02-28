package osucommands

import (
	"encoding/json"
	"io/ioutil"
	"math/rand"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	osuapi "maquiaBot/osu-api"
	osutools "maquiaBot/osu-tools"

	"github.com/bwmarrin/discordgo"
)

// SimilarObject holds information for successful search attemps from https://github.com/Xarib/OsuMapMatcher
type SimilarObject struct {
	DifficultyName string  `json:"DifficultyName"`
	Artist         string  `json:"Artist"`
	Title          string  `json:"Title"`
	Mapper         string  `json:"Mapper"`
	Link           string  `json:"MapLink"`
	KDistance      float64 `json:"KDistance"`
}

var lastReq = time.Now()

// Similar finds similar maps using https://github.com/Xarib/OsuMapMatcher
func Similar(s *discordgo.Session, m *discordgo.MessageCreate) {
	mapRegex, _ := regexp.Compile(`(?i)(https:\/\/)?(osu|old)\.ppy\.sh\/(s|b|beatmaps|beatmapsets)\/(\d+)(#(osu|taiko|fruits|mania)\/(\d+))?`)

	// Check last req time
	currentTime := time.Now()
	if currentTime.Sub(lastReq) < 333*time.Millisecond {
		time.Sleep(333 * time.Millisecond)
	}

	// Check if it should be listed or not
	noList := true
	if strings.Contains(strings.ToLower(m.Content), "-list") {
		noList = false
	}

	// Get the map
	var beatmap osuapi.Beatmap
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

	// Get the map
	nomod := osuapi.Mods(0)
	switch submatches[3] {
	case "s":
		beatmap = osutools.BeatmapParse(submatches[4], "set", &nomod)
	case "b":
		beatmap = osutools.BeatmapParse(submatches[4], "map", &nomod)
	case "beatmaps":
		beatmap = osutools.BeatmapParse(submatches[4], "map", &nomod)
	case "beatmapsets":
		if len(submatches[7]) > 0 {
			beatmap = osutools.BeatmapParse(submatches[7], "map", &nomod)
		} else {
			beatmap = osutools.BeatmapParse(submatches[4], "set", &nomod)
		}
	}
	if beatmap.BeatmapID == 0 {
		s.ChannelMessageSend(m.ChannelID, "No map to compare to!")
		return
	} else if beatmap.Approved < 1 {
		s.ChannelMessageSend(m.ChannelID, "The map `"+beatmap.Artist+" - "+beatmap.Title+"` does not have a leaderboard!")
		return
	}

	// Search
	res, err := http.Get("https://omm.xarib.ch/api/knn/search?id=" + strconv.Itoa(beatmap.BeatmapID))
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "Erorr in requesting search results from https://omm.xarib.ch/api/knn/search?id="+strconv.Itoa(beatmap.BeatmapID))
		return
	}
	defer res.Body.Close()
	byteArray, err := ioutil.ReadAll(res.Body)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "Error in formatting search results into bytes from https://omm.xarib.ch/api/knn/search?id="+strconv.Itoa(beatmap.BeatmapID))
		return
	}

	// Format result
	var searchResult []SimilarObject
	err = json.Unmarshal(byteArray, &searchResult)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "Error in formatting search results into readable JSON from https://omm.xarib.ch/api/knn/search?id="+strconv.Itoa(beatmap.BeatmapID)+"\nPlease note that this feature only works for ranked/approved/loved beatmaps that are not **recently** ranked/approved/loved as well.")
		return
	}

	if !noList {
		message := "Results from https://omm.xarib.ch/api/knn/search?id=" + strconv.Itoa(beatmap.BeatmapID) + ":\n"
		for _, result := range searchResult {
			message += "**" + result.Artist + " - " + result.Title + "** hosted by **" + result.Mapper + " [" + result.DifficultyName + "]**: <" + result.Link + "> (KDistance: " + strconv.FormatFloat(result.KDistance, 'f', 2, 64) + ")\n"
		}
		go s.ChannelMessageSend(m.ChannelID, message)
	} else {
		mapNum := rand.Intn(len(searchResult))
		go BeatmapMessage(s, &discordgo.MessageCreate{
			&discordgo.Message{
				ChannelID: m.ChannelID,
				Content:   searchResult[mapNum].Link,
			},
		}, "Result from https://omm.xarib.ch/api/knn/search?id="+strconv.Itoa(beatmap.BeatmapID)+"\n(KDistance: "+strconv.FormatFloat(searchResult[mapNum].KDistance, 'f', 2, 64)+")", mapRegex)
	}
	lastReq = time.Now()
}
