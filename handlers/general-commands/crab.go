package gencommands

import (
	"net/http"

	"github.com/bwmarrin/discordgo"
)

// Crab sends a crab rave image
func Crab(s *discordgo.Session, m *discordgo.MessageCreate) {
	response, err := http.Get("https://cdn.discordapp.com/emojis/510169818893385729.gif")
	if err != nil {
		return
	}

	message := &discordgo.MessageSend{
		File: &discordgo.File{
			Name:   "crab.gif",
			Reader: response.Body,
		},
	}
	s.ChannelMessageSendComplex(m.ChannelID, message)
	response.Body.Close()
}
