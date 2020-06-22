package gencommands

import (
	"net/http"

	"github.com/bwmarrin/discordgo"
	config "maquiaBot/config"
)

// Cheers sends a cheers message
func Cheers(s *discordgo.Session, m *discordgo.MessageCreate) {
	response, err := http.Get(config.Conf.Cheers)
	if err != nil {
		return
	}

	message := &discordgo.MessageSend{
		File: &discordgo.File{
			Name:   "cheers.mp4",
			Reader: response.Body,
		},
	}
	s.ChannelMessageSendComplex(m.ChannelID, message)
	response.Body.Close()
}
