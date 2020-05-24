package helpcommands

import (
	"github.com/bwmarrin/discordgo"
)

// Toggle explains the toggle functionality
func Toggle(embed *discordgo.MessageEmbed) *discordgo.MessageEmbed {
	embed.Author.Name = "Command: toggle"
	embed.Description = "`toggle <[-a] [-d] [-os] [-s] [-v]>` lets admins toggle specific server options."
	embed.Fields = []*discordgo.MessageEmbedField{
		{
			Name:   "[-a]",
			Value:  "Toggles announcements for this bot to show in that channel.",
			Inline: true,
		},
		{
			Name:   "[-d]",
			Value:  "Toggles whether dailies (`penis`, `bpm`, e.t.c) should run in the server.",
			Inline: true,
		},
		{
			Name:   "[-os]",
			Value:  "Toggle whether map links, profile links, and timestamps should be read by the bot.",
			Inline: true,
		},
		{
			Name:   "[-s]",
			Value:  "Toggle whether anyone can add stats, or only admins.",
			Inline: true,
		},
		{
			Name:   "[-v]",
			Value:  "Toggle whether there shuold be a 1/100000 chance for someone to be vibe checked per message.",
			Inline: true,
		},
	}
	return embed
}
