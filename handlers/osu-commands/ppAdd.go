package osucommands

import (
	"math"
	"regexp"
	"sort"
	"strconv"
	"strings"

	osuapi "../../osu-api"
	structs "../../structs"
	"github.com/bwmarrin/discordgo"
)

// PPAdd calculates final pp after obtaining given pp score
func PPAdd(s *discordgo.Session, m *discordgo.MessageCreate, osuAPI *osuapi.Client, cache []structs.PlayerData) {
	ppRegex, _ := regexp.Compile(`(.+)ppadd\s*(.+)`)
	userAndpp := ppRegex.FindStringSubmatch(m.Content)

	if len(userAndpp) == 0 {
		s.ChannelMessageSend(m.ChannelID, "No pp given!")
		return
	}

	var pp float64
	var username string
	var err error

	// Check if only pp was given, or if a user was also given as well
	args := strings.Split(userAndpp[2], " ")
	if pp, err = strconv.ParseFloat(args[0], 64); err == nil && len(args) == 1 { // No user given, use message author user
		for _, player := range cache {
			if m.Author.ID == player.Discord.ID && player.Osu.Username != "" {
				username = player.Osu.Username
				break
			}
		}
		if username == "" {
			s.ChannelMessageSend(m.ChannelID, "Could not find any osu! account linked for "+m.Author.Mention()+" !")
			return
		}
	} else if pp, err = strconv.ParseFloat(args[len(args)-1], 64); err == nil { // User given
		username = strings.Join(args[:len(args)-1], " ")
	} else { // Nothing given at all
		s.ChannelMessageSend(m.ChannelID, "No pp given!")
		return
	}
	if math.IsNaN(pp) || math.IsInf(pp, 0) || pp < 0 {
		s.ChannelMessageSend(m.ChannelID, "haha ur so funny")
		return
	}

	// Get user and their best scores
	user, err := osuAPI.GetUser(osuapi.GetUserOpts{
		Username: username,
	})
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "Either the API just owned me, or the osu! user **"+username+"** may not exist! Check if the osu! user exists or try again later.")
		return
	}
	scores, err := osuAPI.GetUserBest(osuapi.GetUserScoresOpts{
		Username: user.Username,
		Limit:    100,
	})
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "The API just owned me try again")
		return
	}

	// Get the original pp values, and the new pp values
	ppOnly := []float64{pp}
	var originalPP, newPP float64
	for i, score := range scores {
		ppOnly = append(ppOnly, score.Score.PP)
		originalPP += score.Score.PP * math.Pow(0.95, float64(i))
	}
	sort.Slice(ppOnly, func(i, j int) bool { return ppOnly[i] > ppOnly[j] })
	for i, ppVal := range ppOnly {
		if i == 100 {
			break
		}
		newPP += ppVal * math.Pow(0.95, float64(i))
	}

	// The result
	ppChange := newPP - originalPP
	totalPP := user.PP + ppChange
	s.ChannelMessageSend(m.ChannelID, "**"+user.Username+"**: "+strconv.FormatFloat(user.PP, 'f', 2, 64)+" -> "+strconv.FormatFloat(totalPP, 'f', 2, 64)+" (+"+strconv.FormatFloat(ppChange, 'f', 2, 64)+"pp)")
}
