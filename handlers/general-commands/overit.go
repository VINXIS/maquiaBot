package gencommands

import (
	"net/http"

	config "../../config"
	"github.com/bwmarrin/discordgo"
)

// OverIt sends an over it image
func OverIt(s *discordgo.Session, m *discordgo.MessageCreate) {
	response, err := http.Get(config.Conf.OverIt)
	if err != nil {
		return
	}

	message := &discordgo.MessageSend{
		File: &discordgo.File{
			Name:   "overit.png",
			Reader: response.Body,
		},
	}
	s.ChannelMessageSendComplex(m.ChannelID, message)
	response.Body.Close()
}
