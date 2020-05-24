package helpcommands

import (
	"github.com/bwmarrin/discordgo"
)

// Average explains the average functionality
func Average(embed *discordgo.MessageEmbed) *discordgo.MessageEmbed {
	embed.Author.Name = "Command: ave / average / mean"
	embed.Description = "`[math] (ave|average|mean) <num> <num> [num]...` calculates the harmonic, geometric, and arithmetic averages for a list of numbers given."
	embed.Fields = []*discordgo.MessageEmbedField{
		{
			Name:  "<num> <num> [num]...",
			Value: "The numbers to find the average for (MUST HAVE AT LEAST 2 NUMBERS).",
		},
	}
	return embed
}
