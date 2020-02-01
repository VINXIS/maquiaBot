package tools

import (
	"log"

	"github.com/bwmarrin/discordgo"
)

// CommandLog logs commands inputted
func CommandLog(s *discordgo.Session, m *discordgo.MessageCreate, command string) {
	channel, err := s.Channel(m.ChannelID)
	logText := m.Author.Username + " has used the " + command + " command"
	if err == nil {
		logText += " in #" + channel.Name
	}

	// Check if server or not
	server, err := s.Guild(m.GuildID)
	if err != nil {
		log.Println(logText + ".\nText: " + m.Content)
		return
	}
	log.Println(logText + ", which is in the " + server.Name + " server.\nText: " + m.Content)
}
