package paginationcommands

import (
	"maquiaBot/structs"

	"github.com/bwmarrin/discordgo"
)

// Triggers handles the pagination of a list of triggers
func Triggers(s *discordgo.Session, r *discordgo.MessageReactionAdd, msg *discordgo.Message, serverData structs.ServerData, num, numend int) (*discordgo.MessageEmbed, bool) {
	embed := &discordgo.MessageEmbed{}
	return embed, true
}
