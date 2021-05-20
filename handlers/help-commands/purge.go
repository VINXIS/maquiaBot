package helpcommands

import (
	"github.com/bwmarrin/discordgo"
)

// Purge explains the purge functionality
func Purge(embed *discordgo.MessageEmbed) *discordgo.MessageEmbed {
	embed.Author.Name = "Command: purge"
	embed.Description = "`purge [users] ([num]| <since <time>>)` lets admins delete previous messages."
	embed.Fields = []*discordgo.MessageEmbedField{
		{
			Name:   "[users]",
			Value:  "A list of mentions / users / nicknames to remove separated by spaces. In the case where users have spaces in both their nick and actual username, use mentions instead. No user given will remove the latest num for all users.",
			Inline: true,
		},
		{
			Name:   "[num]",
			Value:  "The number of previous messages to delete. (Default: 3)",
			Inline: true,
		},
		{
			Name:   "<since <time>>",
			Value:  "The time to delete messages to.",
			Inline: true,
		},
	}
	return embed
}
