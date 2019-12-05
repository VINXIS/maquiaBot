package osucommands

import (
	"math"
	"math/rand"
	"regexp"
	"strconv"
	"time"

	osuapi "../../osu-api"
	structs "../../structs"
	"github.com/bwmarrin/discordgo"
)

// BPM gives a player's BPM of the day
func BPM(s *discordgo.Session, m *discordgo.MessageCreate, osuAPI *osuapi.Client, cache []structs.PlayerData) {
	bpmregex, _ := regexp.Compile(`bpm\s+(.+)`)
	username := ""
	if bpmregex.MatchString(m.Content) {
		username = bpmregex.FindStringSubmatch(m.Content)[1]
	}

	// Check for user
	if username == "" {
		for _, cachePlayer := range cache {
			if m.Author.ID == cachePlayer.Discord.ID {
				if cachePlayer.Osu.Username == "" {
					s.ChannelMessageSend(m.ChannelID, "No user linked to your discord account! Use `link` or `set` to link your account!")
					return
				}
				username = cachePlayer.Osu.Username
				break
			}
		}
	}
	if username == "" {
		s.ChannelMessageSend(m.ChannelID, "No user found for you! Use `link` or `set` to link your account!")
		return
	}

	// Get top scores
	player, err := osuAPI.GetUser(osuapi.GetUserOpts{
		Username: username,
	})
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "User: **"+username+"** may not exist!")
		return
	}

	orderedScores, err := osuAPI.GetUserBest(osuapi.GetUserScoresOpts{
		Username: username,
		Limit:    100,
	})
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "An error occured in trying to get **"+username+"'s** top scores!")
		return
	} else if len(orderedScores) == 0 {
		s.ChannelMessageSend(m.ChannelID, username+" has no top scores!")
		return
	}

	// Obtain average and stddev
	var averageBPM float64
	var mapBPMs []float64
	for _, score := range orderedScores {
		beatmap, err := osuAPI.GetBeatmaps(osuapi.GetBeatmapsOpts{
			BeatmapID: score.BeatmapID,
		})
		if err != nil {
			continue
		}
		if score.Mods&osuapi.ModDoubleTime != 0 {
			beatmap[0].BPM *= 1.5
		} else if score.Mods&osuapi.ModHalfTime != 0 {
			beatmap[0].BPM *= 0.75
		}
		averageBPM += beatmap[0].BPM
		mapBPMs = append(mapBPMs, beatmap[0].BPM)
	}

	averageBPM /= float64(len(mapBPMs))
	var stddevBPM float64
	for _, BPM := range mapBPMs {
		stddevBPM += math.Pow(BPM-averageBPM, 2)
	}
	stddevBPM = math.Sqrt(stddevBPM / float64(len(mapBPMs)-1))

	// Create randomizer based on osu! ID and date
	year, month, day := time.Now().Date()
	random := rand.New(rand.NewSource(int64(player.UserID + day + int(month) + year)))

	todayBPM := math.Min(math.Max(50, random.NormFloat64()*stddevBPM+averageBPM), 400)

	s.ChannelMessageSend(m.ChannelID, "BPM of the day for **"+username+":** "+strconv.FormatFloat(todayBPM, 'f', 0, 64))
}
