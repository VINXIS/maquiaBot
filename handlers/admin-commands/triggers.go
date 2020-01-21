package admincommands

import (
	"fmt"

	tools "../../tools"
	"github.com/bwmarrin/discordgo"
)

// Trigger adds / removes triggers
func Trigger(s *discordgo.Session, m *discordgo.MessageCreate) {
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
	serverData := tools.GetServer(*server)
	fmt.Println(serverData)
}
