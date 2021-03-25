package helpcommands

import (
	"github.com/bwmarrin/discordgo"
)

// Conversation explains the conversation functionality
func Conversation(embed *discordgo.MessageEmbed) *discordgo.MessageEmbed {
	embed.Author.Name = "Command: conv / convo / conversation"
	embed.Description = "`(conv|convo|conversation) [num] [-i]` provides a conversation based on quotes stored for the server"
	embed.Fields = []*discordgo.MessageEmbedField{
		{
			Name:   "[num]",
			Value:  "The number of people to add.",
			Inline: true,
		},
		{
			Name:   "[-i]",
			Value:  "Include links. This will change the command from a code block to regular text as well.",
			Inline: true,
		},
		{
			Name:  "Related Commands:",
			Value: "`quoteadd`, `quote`",
		},
	}
	return embed
}
