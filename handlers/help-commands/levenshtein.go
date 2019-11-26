package helpcommands

import (
	"github.com/bwmarrin/discordgo"
)

// Levenshtein explains the levenshtein functionality
func Levenshtein(embed *discordgo.MessageEmbed) *discordgo.MessageEmbed {
	embed.Author.Name = "Command: l / leven / levenshtein"
	embed.Description = "`(l|leven|levenshtein) <word1> <word2>` gives the levenshtein value between 2 words."
	embed.Fields = []*discordgo.MessageEmbedField{
		&discordgo.MessageEmbedField{
			Name:  "<word1> <word2>",
			Value: "The 2 words to complete the value of.",
		},
		&discordgo.MessageEmbedField{
			Name:  "What is the levenshtein distance?",
			Value: "It is a way to calculate how different two words / phrases are from each other.\nhttps://en.wikipedia.org/wiki/Levenshtein_distance",
		},
	}
	return embed
}
