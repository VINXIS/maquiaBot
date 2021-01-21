package gencommands

import (
	"maquiaBot/tools"
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
)

// Ntow changes numbers to words
func Ntow(s *discordgo.Session, m *discordgo.MessageCreate) {
	parts := strings.Split(m.Content, " ")[1:]
	converted := ""
	for _, part := range parts {
		if num, err := strconv.Atoi(part); err == nil {
			converted = tools.Ntow(int64(num))
			break
		}
	}

	if converted == "" {
		s.ChannelMessageSend(m.ChannelID, "No number found! Reminder that the largest number allowed is `9223372036854775807`")
	} else {
		s.ChannelMessageSend(m.ChannelID, converted)
	}
}

// Wton changes words to numbers
func Wton(s *discordgo.Session, m *discordgo.MessageCreate) {

}
