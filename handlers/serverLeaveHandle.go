package handlers

import (
	"strconv"

	"github.com/bwmarrin/discordgo"
)

// ServerLeave is to send a message when the bot leaves a server
func ServerLeave(s *discordgo.Session, g *discordgo.GuildDelete) {
	s.UpdateStatus(0, strconv.Itoa(len(s.State.Guilds))+" servers")
}
