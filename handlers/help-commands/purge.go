package helpcommands

import (
	"github.com/bwmarrin/discordgo"
)

// Purge explains the purge functionality
func Purge(embed *discordgo.MessageEmbed) *discordgo.MessageEmbed {
	embed.Author.Name = "Command: purge"
	embed.Description = "`purge [users] ([num] | <since <datetime>> | <since <time duration> [ago]>)` lets admins delete previous messages."
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
			Name:   "<since <datetime>>",
			Value:  "The time to delete messages to.",
			Inline: true,
		},
		{
			Name:   "<since <time duration> [ago]>",
			Value:  "The time since to delete messages from.",
			Inline: true,
		},
		{
			Name:   "Example format (since datetime):",
			Value:  "`$purge since march 10` will purge messages up to march 10 2021 00:00 UTC.",
			Inline: true,
		},
		{
			Name:   "Example format (since time duration ago):",
			Value:  "`$purge since 5 sec` will purge messages sent within the last 5 seconds.",
			Inline: true,
		},
	}
	return embed
}
