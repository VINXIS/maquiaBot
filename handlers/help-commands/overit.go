package helpcommands

import (
	"github.com/bwmarrin/discordgo"
)

// OverIt explains the over it functionality
func OverIt(embed *discordgo.MessageEmbed) *discordgo.MessageEmbed {
	embed.Author.Name = "Command: over"
	embed.Description = "`over` lets you send an over it pic."
	return embed
}
