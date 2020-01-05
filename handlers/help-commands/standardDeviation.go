package helpcommands

import (
	"github.com/bwmarrin/discordgo"
)

// StandardDeviation explains the standard deviation functionality
func StandardDeviation(embed *discordgo.MessageEmbed) *discordgo.MessageEmbed {
	embed.Author.Name = "Command: stddev / standarddev / stddeviation / standarddeviation"
	embed.Description = "`[math] (stddev|standarddev|stddeviation|standarddeviation) <num> <num> [num]...` calculates the population and sample standard deviation for a list of numbers given."
	embed.Fields = []*discordgo.MessageEmbedField{
		&discordgo.MessageEmbedField{
			Name:  "<num> <num> [num]...",
			Value: "The numbers to find the standard deviations for (MUST HAVE AT LEAST 2 NUMBERS).",
		},
	}
	return embed
}
