package handlers

import (
	tools "../tools"
	"github.com/bwmarrin/discordgo"
)

// ServerJoin is to send a message when the bot joins a server
func ServerJoin(s *discordgo.Session, g *discordgo.GuildCreate) {
	dm, err := s.UserChannelCreate(g.OwnerID)
	tools.ErrRead(err)

	s.ChannelMessageSend(dm.ID, "Hello server owner!")
	return
}
