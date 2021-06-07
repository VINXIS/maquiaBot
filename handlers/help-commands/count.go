package helpcommands

import (
	"github.com/bwmarrin/discordgo"
)

// CharCount explains the char count functionality
func CharCount(embed *discordgo.MessageEmbed) *discordgo.MessageEmbed {
	embed.Author.Name = "Command: charcount / ccount / count / char"
	embed.Description = "`(charcount|ccount|count|char) <text|attachment>` counts the number of characters in the text provided"
	embed.Fields = []*discordgo.MessageEmbedField{
		{
			Name:  "<text|attachment>",
			Value: "The text to count on (you may also provide an attachment instead, those take higher prio).",
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
	embed.Author.Name = "Command: wordcount / wcount / word"
	embed.Description = "`(wordcount|wcount|word) <text|attachment>` counts the number of words in the text provided"
	embed.Fields = []*discordgo.MessageEmbedField{
		{
			Name:  "<text|attachment>",
			Value: "The text to count on (you may also provide an attachment instead, those take higher prio).",
		},
		{
			Name:  "Related commands:",
			Value: "`charcount`",
		},
	}
	return embed
}
