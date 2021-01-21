package admincommands

import (
	"maquiaBot/tools"

	"github.com/bwmarrin/discordgo"
)

// RemoveChannel lets admins remove channel data
func RemoveChannel(s *discordgo.Session, m *discordgo.MessageCreate) {
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
	_, new := tools.GetChannel(*channel, s)
	if new {
		s.ChannelMessageSend(m.ChannelID, "There is currently no data stored for this channel.")
		return
	}

	tools.DeleteFile("./data/channelData/" + m.ChannelID + ".json")
	s.ChannelMessageSend(m.ChannelID, "Removed data for this channel!")
}

// RemoveServer lets admins remove server data
func RemoveServer(s *discordgo.Session, m *discordgo.MessageCreate) {
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
	_, new := tools.GetServer(*server, s)
	if new {
		s.ChannelMessageSend(m.ChannelID, "There is currently no data stored for this server.")
		return
	}

	tools.DeleteFile("./data/serverData/" + m.GuildID + ".json")
	s.ChannelMessageSend(m.ChannelID, "Removed data for this server!")
}
