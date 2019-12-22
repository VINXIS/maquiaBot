package helpcommands

import (
	"github.com/bwmarrin/discordgo"
)

// Pokemon explains the pokemon functionality
func Pokemon(embed *discordgo.MessageEmbed) *discordgo.MessageEmbed {
	embed.Author.Name = "Command: pokemon"
	embed.Description = "`pokemon <pokemon name / ID>` lets you obtain information about a pokemon."
	embed.Fields = []*discordgo.MessageEmbedField{
		&discordgo.MessageEmbedField{
			Name:  "<pokemon name / ID>",
			Value: "The name / ID of a pokemon to get information for.",
		},
	}
	return embed
}
