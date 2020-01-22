package helpcommands

import (
	"github.com/bwmarrin/discordgo"
)

// Cheers explains the cheers functionality
func Cheers(embed *discordgo.MessageEmbed) *discordgo.MessageEmbed {
	embed.Author.Name = "Command: cheers"
	embed.Description = "`cheers` lets you send a cheers video."
	return embed
}
