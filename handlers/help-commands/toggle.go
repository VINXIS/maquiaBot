package helpcommands

import (
	"github.com/bwmarrin/discordgo"
)

// Toggle explains the toggle functionality
func Toggle(embed *discordgo.MessageEmbed) *discordgo.MessageEmbed {
	embed.Author.Name = "Command: toggle"
	embed.Description = "`toggle [-ch] <[-a] [-d] [-os] [-s] [-t] [-v]>` lets admins toggle specific server / channel options."
	embed.Fields = []*discordgo.MessageEmbedField{
		{
			Name:   "[-ch]",
			Value:  "Changes toggling to the specific channel instead of the server. **Please note that server options override channel options.**",
			Inline: true,
		},
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
			Value:  "Toggle whether map links, and profile links should be read by the bot.",
			Inline: true,
		},
		{
			Name:   "[-s]",
			Value:  "Toggle whether anyone can add stats, or only admins.",
			Inline: true,
		},
		{
			Name:   "[-t]",
			Value:  "Toggle whether osu! timestamp links should be generated or not.",
			Inline: true,
		},
		{
			Name:   "[-v]",
			Value:  "Toggle whether there should be a 1/100000 chance for someone to be vibe checked per message.",
			Inline: true,
		},
	}
	return embed
}
