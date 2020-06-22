package gencommands

import (
	"net/http"

	"github.com/bwmarrin/discordgo"
	config "maquiaBot/config"
)

// Late sends a video related to someone saying late.
func Late(s *discordgo.Session, m *discordgo.MessageCreate) {
	response, err := http.Get(config.Conf.Late)
	if err != nil {
		return
	}

	message := &discordgo.MessageSend{
		File: &discordgo.File{
			Name:   "late.mp4",
			Reader: response.Body,
		},
	}
	s.ChannelMessageSendComplex(m.ChannelID, message)
	response.Body.Close()
}
