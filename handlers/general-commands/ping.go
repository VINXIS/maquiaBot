package gencommands

import (
	"time"

	"github.com/bwmarrin/discordgo"
)

// Ping checks how long it takes to reply to a message
func Ping(s *discordgo.Session, m *discordgo.MessageCreate) {
	messageTime := time.Now().UTC()
	msg, err := s.ChannelMessageSend(m.ChannelID, "This Message's timestamp: "+messageTime.Format(time.RFC3339Nano))
	if err != nil {
		return
	}
	editTime := time.Now().UTC()
	s.ChannelMessageEdit(m.ChannelID, msg.ID, "This Message's timestamp: "+messageTime.Format(time.RFC3339Nano)+"\nEdit time: "+editTime.Format(time.RFC3339Nano)+"\nTime elapsed: "+editTime.Sub(messageTime).String())
}
