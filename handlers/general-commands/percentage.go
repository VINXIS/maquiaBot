package gencommands

import (
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
	textLength := 0
	if percentRegex.MatchString(m.Content) {
		text = "Percentage for `" + strings.ReplaceAll(percentRegex.FindStringSubmatch(m.Content)[2], "`", "") + "`:\n"
		textLength = len(strings.ReplaceAll(percentRegex.FindStringSubmatch(m.Content)[2], "`", ""))
	}
	authorid, _ := strconv.Atoi(m.Author.ID)
	skillRang := rand.New(rand.NewSource(int64(authorid+textLength) + time.Now().UnixNano()))
	percent := skillRang.Int63n(100) + 1
	bar := tools.BarCreation(float64(percent) / 100)

	s.ChannelMessageSend(m.ChannelID, text+"```\n"+bar+" "+strconv.FormatInt(percent, 10)+"%```")
}
