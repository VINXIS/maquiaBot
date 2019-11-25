package helpcommands

import (
	"github.com/bwmarrin/discordgo"
)

// Kanye explains the kanye functionality
func Kanye(embed *discordgo.MessageEmbed) *discordgo.MessageEmbed {
	embed.Author.Name = "Command: kanye"
	embed.Description = "`kanye` gives a quote from Kanye West."
	return embed
}
