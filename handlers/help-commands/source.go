package helpcommands

import (
	"github.com/bwmarrin/discordgo"
)

// Source explains the source functionality
func Source(embed *discordgo.MessageEmbed) *discordgo.MessageEmbed {
	embed.Author.Name = "Command: src / source"
	embed.Description = "`(src|source)` sends a link to the github repository for this bot, which is https://github.com/VINXIS/maquiaBot"
	return embed
}
