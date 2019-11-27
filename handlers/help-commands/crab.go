package helpcommands

import (
	"github.com/bwmarrin/discordgo"
)

// Crab explains the crab functionality
func Crab(embed *discordgo.MessageEmbed) *discordgo.MessageEmbed {
	embed.Author.Name = "Command: crab"
	embed.Description = "`crab` lets admins toggle whether any text containing crab / rave (even within words) will send a crab rave gif."
	return embed
}
