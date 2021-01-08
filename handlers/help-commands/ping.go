package helpcommands

import (
	"github.com/bwmarrin/discordgo"
)

// Ping explains the ping functionality
func Ping(embed *discordgo.MessageEmbed) *discordgo.MessageEmbed {
	embed.Author.Name = "Command: ping"
	embed.Description = "`ping` checks how long it takes for the bot to send and edit a message."
	return embed
}
