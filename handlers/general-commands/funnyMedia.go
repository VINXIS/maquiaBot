package gencommands

import (
	"net/http"

	config "maquiaBot/config"

	"github.com/bwmarrin/discordgo"
)

// FunnyMedia sends specific funny images/videos
func FunnyMedia(s *discordgo.Session, m *discordgo.MessageCreate, t string) {
	var response *http.Response
	var err error
	switch t {
	case "cheers":
		response, err = http.Get(config.Conf.Cheers)
	case "crab":
		response, err = http.Get(config.Conf.Crab)
	case "late":
		response, err = http.Get(config.Conf.Late)
	case "over":
		response, err = http.Get(config.Conf.OverIt)
	case "idea":
		s.ChannelMessageSend(m.ChannelID, "https://www.youtube.com/watch?v=aAxjVu3iZps")
		return
	}
	if err != nil {
		return
	}

	message := &discordgo.MessageSend{
		File: &discordgo.File{
			Name:   "funny.mp4",
			Reader: response.Body,
		},
	}
	if t == "crab" {
		message.File.Name = "funny.gif"
	} else if t == "over" {
		message.File.Name = "funny.png"
	}
	s.ChannelMessageSendComplex(m.ChannelID, message)
	response.Body.Close()
}
