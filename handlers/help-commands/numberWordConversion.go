package helpcommands

import (
	"github.com/bwmarrin/discordgo"
)

// Ntow explains the number to word functionality
func Ntow(embed *discordgo.MessageEmbed) *discordgo.MessageEmbed {
	embed.Author.Name = "Command: ntw / numtw / ntow / ntword / numtow / numtword / ntoword / numtoword / numbertoword"
	embed.Description = "`(ntw|numtw|ntow|ntword|numtow|numtword|ntoword|numtoword|numbertoword) <number>` converts a number to word form."
	embed.Fields = []*discordgo.MessageEmbedField{
		{
			Name:  "<number>",
			Value: "The number to convert into words.",
		},
	}
	return embed
}
