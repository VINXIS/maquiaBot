package gencommands

import (
	"time"

	"github.com/bwmarrin/discordgo"
)

// Ping checks how long it takes to reply to a message
func Ping(s *discordgo.Session, m *discordgo.MessageCreate) {
	messageTime, _ := m.Timestamp.Parse()
	now := time.Now().UTC()
	timeSince := now.Sub(messageTime).String()
	s.ChannelMessageSend(m.ChannelID, "Message timestamp: "+messageTime.Format(time.RFC822)+"\nObtained at: "+now.Format(time.RFC822)+"\nTime elapsed: "+timeSince)
}
