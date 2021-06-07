package helpcommands

import (
	"github.com/bwmarrin/discordgo"
)

// CharCount explains the char count functionality
func CharCount(embed *discordgo.MessageEmbed) *discordgo.MessageEmbed {
	embed.Author.Name = "Command: charcount / ccount / count"
	embed.Description = "`(charcount|ccount|count) <text>` counts the number of characters in the text provided"
	embed.Fields = []*discordgo.MessageEmbedField{
		{
			Name:  "<text>",
			Value: "The text to count on.",
		},
		{
			Name:  "Related commands:",
			Value: "`wordcount`",
		},
	}
	return embed
}

// WordCount explains the word count functionality
func WordCount(embed *discordgo.MessageEmbed) *discordgo.MessageEmbed {
	embed.Author.Name = "Command: wordcount / wcount"
	embed.Description = "`(wordcount|wcount) <text>` counts the number of words in the text provided"
	embed.Fields = []*discordgo.MessageEmbedField{
		{
			Name:  "<text>",
			Value: "The text to count on.",
		},
		{
			Name:  "Related commands:",
			Value: "`charcount`",
		},
	}
	return embed
}
