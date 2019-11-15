package osucommands

import (
	"regexp"

	"github.com/bwmarrin/discordgo"
)

// TimestampMessage parses osu! timestamps
func TimestampMessage(s *discordgo.Session, m *discordgo.MessageCreate, regex *regexp.Regexp) {
	args := regex.FindAllStringSubmatch(m.Content, -1)
	msg := ""
	for _, timestamp := range args {
		msg = msg + "<osu://edit/" + timestamp[1] + ":" + timestamp[2] + ":" + timestamp[3]
		if timestamp[4] != "" {
			msg = msg + "-" + timestamp[4]
		}
		msg = msg + ">\n"
	}
	s.ChannelMessageSend(m.ChannelID, msg)
}
