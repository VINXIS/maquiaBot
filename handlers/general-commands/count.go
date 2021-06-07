package gencommands

import (
	"regexp"
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
)

// CharCount counts the number of characters in text
func CharCount(s *discordgo.Session, m *discordgo.MessageCreate) {
	charRegex, _ := regexp.Compile(`(?i)(c(har)?)?count\s+(.+)`)

	if !charRegex.MatchString(m.Content) {
		s.ChannelMessageSend(m.ChannelID, "No text provided to char count for!")
		return
	}

	text := charRegex.FindStringSubmatch(m.Content)[3]
	s.ChannelMessageSend(m.ChannelID, strconv.Itoa(len(text)))
}

// WordCount counts the number of words in text
func WordCount(s *discordgo.Session, m *discordgo.MessageCreate) {
	wordRegex, _ := regexp.Compile(`(?i)w(ord)?count\s+(.+)`)

	if !wordRegex.MatchString(m.Content) {
		s.ChannelMessageSend(m.ChannelID, "No text provided to word count for!")
		return
	}

	text := wordRegex.FindStringSubmatch(m.Content)[2]
	s.ChannelMessageSend(m.ChannelID, strconv.Itoa(len(strings.Fields(text))))
}
