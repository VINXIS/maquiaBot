package helpcommands

import (
	"github.com/bwmarrin/discordgo"
)

// Toggle explains the toggle functionality
func Toggle(embed *discordgo.MessageEmbed) *discordgo.MessageEmbed {
	embed.Author.Name = "Command: toggle"
	embed.Description = "`toggle <[-a] [-ch] [-cr] [-d] [-i] [-l] [-o] [-s] [-v]>` lets admins toggle specific server options."
	embed.Fields = []*discordgo.MessageEmbedField{
		&discordgo.MessageEmbedField{
			Name:   "[-a]",
			Value:  "Toggle announces from the bot creator on and off",
			Inline: true,
		},
		&discordgo.MessageEmbedField{
			Name:   "[-ch]",
			Value:  "Toggle whether any message containing 🍻, 🍺, 🦐, and / or cheer should trigger a cheers video message.",
			Inline: true,
		},
		&discordgo.MessageEmbedField{
			Name:   "[-cr]",
			Value:  "Toggle whether any message containing crab or rave should trigger a crab rave message.",
			Inline: true,
		},
		&discordgo.MessageEmbedField{
			Name:   "[-d]",
			Value:  "Toggles whether dailies (`penis`, `bpm`, e.t.c) should run in the server.",
			Inline: true,
		},
		&discordgo.MessageEmbedField{
			Name:   "[-i]",
			Value:  "Toggle whether any message containing nice idea, good idea, or gud idea should send a nice idea video message.",
			Inline: true,
		},
		&discordgo.MessageEmbedField{
			Name:   "[-l]",
			Value:  "Toggle whether any message containing late or ancient should send a late video message.",
			Inline: true,
		},
		&discordgo.MessageEmbedField{
			Name:   "[-o]",
			Value:  "Toggle whether map links, profile links, and timestamps should be read by the bot.",
			Inline: true,
		},
		&discordgo.MessageEmbedField{
			Name:   "[-s]",
			Value:  "Toggle whether anyone can add stats, or only admins.",
			Inline: true,
		},
		&discordgo.MessageEmbedField{
			Name:   "[-v]",
			Value:  "Toggle whether there shuold be a 1/100000 chance for someone to be vibe checked per message.",
			Inline: true,
		},
	}
	return embed
}
