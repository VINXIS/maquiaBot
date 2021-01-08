package helpcommands

import (
	"github.com/bwmarrin/discordgo"
)

// Map explains the map functionality
func Map(embed *discordgo.MessageEmbed) *discordgo.MessageEmbed {
	embed.Author.Name = "Command: m / map"
	embed.Description = "`(m|map|<beatmap link>) [-m <mods>] [-acc <accuracy>] [-c <combo>] [-x <misses>] [-s <score>]` lets you obtain beatmap information."
	embed.Fields = []*discordgo.MessageEmbedField{
		{
			Name:   "<beatmap link>",
			Value:  "You may link a map instead of using `m` or `map` to get beatmap information.",
			Inline: true,
		},
		{
			Name:   "[-m mods]",
			Value:  "The mods to get pp information for (Default: NM)",
			Inline: true,
		},
		{
			Name:   "[-acc accuracy]",
			Value:  "The accuracy to get pp information for. No `-acc` will give pp information for 95, 97, 98, 99, SS acc.",
			Inline: true,
		},
		{
			Name:   "[-c combo]",
			Value:  "The combo to get pp information for (Default: Max Combo).",
			Inline: true,
		},
		{
			Name:   "[-x misses]",
			Value:  "The miss count to get pp information for (Default: 0).",
			Inline: true,
		},
		{
			Name:   "[-s score]",
			Value:  "The score to get the pp value for. Not for use alongside `-acc`. **osu!mania ONLY!**",
			Inline: true,
		},
	}
	return embed
}
