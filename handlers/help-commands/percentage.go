package helpcommands

import (
	"github.com/bwmarrin/discordgo"
)

// Percentage explains the percent functionality
func Percentage(embed *discordgo.MessageEmbed) *discordgo.MessageEmbed {
	embed.Author.Name = "Command: p / per / percent / percentage"
	embed.Description = "`(p|per|percent|percentage) [text]` gives a percentage."
	embed.Fields = []*discordgo.MessageEmbedField{
		&discordgo.MessageEmbedField{
			Name:  "[text]",
			Value: "The text to give a percentage for. No text will not really change anything.",
		},
	}
	return embed
}
