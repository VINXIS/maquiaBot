package gencommands

import (
	"regexp"
	"strconv"

	"github.com/bwmarrin/discordgo"
	tools "maquiaBot/tools"
)

// Levenshtein gives the Levenshtein distance of 2 messages
func Levenshtein(s *discordgo.Session, m *discordgo.MessageCreate) {
	messageRegex, _ := regexp.Compile(`(?i)l(even(shtein)?)?\s+(.+)\s+(.+)`)

	if !messageRegex.MatchString(m.Content) {
		s.ChannelMessageSend(m.ChannelID, "Please give 2 words to compare!")
		return
	}
	res := messageRegex.FindStringSubmatch(m.Content)
	word1 := res[3]
	word2 := res[4]

	value := tools.Levenshtein(word1, word2)

	s.ChannelMessageSend(m.ChannelID, "The levenshtein value of `"+word1+"` and `"+word2+"` is "+strconv.FormatFloat(value, 'f', 2, 64))
}
