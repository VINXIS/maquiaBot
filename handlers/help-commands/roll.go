package helpcommands

import (
	"github.com/bwmarrin/discordgo"
)

// Roll explains the roll functionality
func Roll(embed *discordgo.MessageEmbed) *discordgo.MessageEmbed {
	embed.Author.Name = "Command: roll"
	embed.Description = "`roll [num]` gives a number between 1 to num."
	embed.Fields = []*discordgo.MessageEmbedField{
		{
			Name:  "[num]",
			Value: "The max number to roll for (Default: 100).",
		},
	}
	return embed
}
