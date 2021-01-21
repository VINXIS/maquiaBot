package admincommands

import (
	"bytes"
	"encoding/json"
	tools "maquiaBot/tools"

	"github.com/bwmarrin/discordgo"
)

// DownloadChannel lets admins download any data stored for the channel
func DownloadChannel(s *discordgo.Session, m *discordgo.MessageCreate) {
	channel, err := s.Channel(m.ChannelID)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "This is not an allowed channel!")
		return
	}

	server, err := s.Guild(m.GuildID)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "This is not a server!")
		return
	}

	if !tools.AdminCheck(s, m, *server) {
		s.ChannelMessageSend(m.ChannelID, "You must be an admin, server manager, or server owner!")
		return
	}

	// Obtain channel data
	channelData, new := tools.GetChannel(*channel, s)

	if new {
		s.ChannelMessageSend(m.ChannelID, "There is currently no data stored for this channel.")
		return
	}

	b, err := json.Marshal(channelData)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "Error in parsing the channel data. Please contact `@vinxis1` on twitter or `VINXIS#1000` on discord about this!")
		return
	}
	s.ChannelMessageSendComplex(m.ChannelID, &discordgo.MessageSend{
		Content: "Here is the data stored for the channel.",
		File: &discordgo.File{
			Name:   channel.Name + ".json",
			Reader: bytes.NewReader(b),
		},
	})
}

// DownloadServer lets admins download any data stored for the server
func DownloadServer(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Check if server exists
	server, err := s.Guild(m.GuildID)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "This is not a server!")
		return
	}

	if !tools.AdminCheck(s, m, *server) {
		s.ChannelMessageSend(m.ChannelID, "You must be an admin, server manager, or server owner!")
		return
	}

	// Obtain server data
	serverData, new := tools.GetServer(*server, s)
	if new {
		s.ChannelMessageSend(m.ChannelID, "There is currently no data stored for this server.")
		return
	}

	b, err := json.Marshal(serverData)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "Error in parsing the server data. Please contact `@vinxis1` on twitter or `VINXIS#1000` on discord about this!")
		return
	}
	s.ChannelMessageSendComplex(m.ChannelID, &discordgo.MessageSend{
		Content: "Here is the data stored for the server.",
		File: &discordgo.File{
			Name:   server.Name + ".json",
			Reader: bytes.NewReader(b),
		},
	})
}
