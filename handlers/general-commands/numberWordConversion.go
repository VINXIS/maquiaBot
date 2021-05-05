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
		if num, err := strconv.ParseFloat(part, 64); err == nil && num < 99999999999999999999999999999999999999999999999999999999999999999 {
			converted = tools.Ntow(num)
			break
		}
	}

	if converted == "" {
		s.ChannelMessageSend(m.ChannelID, "No number found! Reminder that the largest number allowed is (technically) `99999999999999999999999999999999999999999999999999999999999999999`")
	} else {
		s.ChannelMessageSend(m.ChannelID, converted)
	}
}

// Wton changes words to numbers
func Wton(s *discordgo.Session, m *discordgo.MessageCreate) {

}
