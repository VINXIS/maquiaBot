package helpcommands

import (
	"github.com/bwmarrin/discordgo"
)

// Counter explains the counter functionality
func Counter(embed *discordgo.MessageEmbed) *discordgo.MessageEmbed {
	embed.Author.Name = "Command: counter"
	embed.Description = "`counter (<text>|<-d <number>>)` let's you create counters for phrases to see how often people say it.\nYou may also use regex for counters! Test your regex here to create a valid regex: https://regex101.com/"
	embed.Fields = []*discordgo.MessageEmbedField{
		{
			Name:   "<text>",
			Value:  "The text / regex to track.",
			Inline: true,
		},
		{
			Name:   "<-d number>",
			Value:  "`-d` followed by the counter's ID found in `counters`",
			Inline: true,
		},
		{
			Name:  "Related Commands:",
			Value: "`countrank`, `counters`",
		},
	}
	return embed
}

// CountRank explains the countrank functionality
func CountRank(embed *discordgo.MessageEmbed) *discordgo.MessageEmbed {
	embed.Author.Name = "Command: countrank"
	embed.Description = "`(cr|countr|crank|countrank) [num]` showcases the top number of people who have sent the tracked word/regex."
	embed.Fields = []*discordgo.MessageEmbedField{
		{
			Name:   "[num]",
			Value:  "The number of people to show",
			Inline: true,
		},
		{
			Name:  "Related Commands:",
			Value: "`counter`, `counters`",
		},
	}
	return embed
}

// Counters explains the counters functionality
func Counters(embed *discordgo.MessageEmbed) *discordgo.MessageEmbed {
	embed.Author.Name = "Command: countrank"
	embed.Description = "`(cs|counters)` shows all the currently tracked words/regex in the server."
	embed.Fields = []*discordgo.MessageEmbedField{
		{
			Name:  "Related Commands:",
			Value: "`counter`, `countrank`",
		},
	}
	return embed
}
