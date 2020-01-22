package helpcommands

import (
	"github.com/bwmarrin/discordgo"
)

// Crab explains the crab functionality
func Crab(embed *discordgo.MessageEmbed) *discordgo.MessageEmbed {
	embed.Author.Name = "Command: crab"
	embed.Description = "`crab` lets you send a crab rave gif."
	return embed
}
