package osucommands

import (
	"math"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
	osuapi "maquiaBot/osu-api"
	structs "maquiaBot/structs"
)

// PPAdd calculates final pp after obtaining given pp score
func PPAdd(s *discordgo.Session, m *discordgo.MessageCreate, cache []structs.PlayerData) {
	ppRegex, _ := regexp.Compile(`(?i)(ppadd|addpp)\s+(.+)`)

	if !ppRegex.MatchString(m.Content) {
		s.ChannelMessageSend(m.ChannelID, "No pp given!")
		return
	}

	var pp float64
	var err error
	username := ppRegex.FindStringSubmatch(m.Content)[2]

	// Check if only pp was given, or if a user was also given as well
	ppSplit := strings.Split(username, " ")
	for _, txt := range ppSplit {
		if i, err := strconv.ParseFloat(strings.Replace(txt, "pp", "", -1), 64); err == nil {
			username = strings.TrimSpace(strings.Replace(username, txt, "", 1))
			pp = i
			break
		}
	}

	if username == "" {
		for _, player := range cache {
			if m.Author.ID == player.Discord.ID && player.Osu.Username != "" {
				username = player.Osu.Username
				break
			}
		}
		if username == "" {
			s.ChannelMessageSend(m.ChannelID, "Could not find any osu! account linked for "+m.Author.Mention()+" ! Please use `set` or `link` to link an osu! account to you!")
			return
		}
	}

	if math.IsNaN(pp) || math.IsInf(pp, 0) || pp < 0 {
		s.ChannelMessageSend(m.ChannelID, "haha ur so funny")
		return
	}

	// Get user and their best scores
	user, err := OsuAPI.GetUser(osuapi.GetUserOpts{
		Username: username,
	})
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "Either the API just owned me, or the osu! user **"+username+"** may not exist! Check if the osu! user exists or try again later.")
		return
	}
	scores, err := OsuAPI.GetUserBest(osuapi.GetUserScoresOpts{
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
	placement := 101
	for i, ppVal := range ppOnly {
		if i == 100 {
			break
		}
		if ppVal == pp {
			placement = i + 1
		}
		newPP += ppVal * math.Pow(0.95, float64(i))
	}

	// The result
	ppChange := newPP - originalPP
	totalPP := user.PP + ppChange
	text := "**" + user.Username + "**: " + strconv.FormatFloat(user.PP, 'f', 2, 64) + " -> " + strconv.FormatFloat(totalPP, 'f', 2, 64) + " (+" + strconv.FormatFloat(ppChange, 'f', 2, 64) + "pp)"
	if placement <= 100 {
		text += " | **#" + strconv.Itoa(placement) + "** in top performances!"
	}
	s.ChannelMessageSend(m.ChannelID, text)
}
