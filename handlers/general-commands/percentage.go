package gencommands

import (
	"math"
	"math/rand"
	"regexp"
	"strconv"
	"strings"
	"time"

	tools "../../tools"
	"github.com/bwmarrin/discordgo"
)

// Percentage gives a random percentage alongside bars
func Percentage(s *discordgo.Session, m *discordgo.MessageCreate) {
	percentRegex, _ := regexp.Compile(`(p|percentage|per|percent)\s+(.+)?`)

	var text string
	if percentRegex.MatchString(m.Content) {
		text = "Percentage for `" + strings.ReplaceAll(percentRegex.FindStringSubmatch(m.Content)[2], "`", "") + "`:\n"
	}
	authorid, _ := strconv.Atoi(m.Author.ID)
	skillRang := rand.New(rand.NewSource(int64(authorid) + time.Now().UnixNano()))
	percent := math.Max(0, math.Min(100, skillRang.NormFloat64()*12.5+50))
	bar := tools.BarCreation(percent / 100)

	s.ChannelMessageSend(m.ChannelID, text+"```\n"+bar+" "+strconv.FormatFloat(percent, 'f', 0, 64)+"%```")
}
