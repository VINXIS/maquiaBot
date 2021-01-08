package helpcommands

import (
	"github.com/bwmarrin/discordgo"
)

// Leaderboard explains the leaderboard functionality
func Leaderboard(embed *discordgo.MessageEmbed) *discordgo.MessageEmbed {
	embed.Author.Name = "Command: l / leader / leaderboard"
	embed.Description = "`(l|leader|leaderboard) [-s] [-n <number>] [-m <mods>]` lets you obtain the leaderboard of the beatmap."
	embed.Fields = []*discordgo.MessageEmbedField{
		{
			Name:   "[-s]",
			Value:  "Add this if you want the leaderboard only containing players in the server.",
			Inline: true,
		},
		{
			Name:   "[-n number]",
			Value:  "Numbers of scores to show (Default: 5).",
			Inline: true,
		},
		{
			Name:   "[-m mods]",
			Value:  "The mods to get pp information for (Default: NM)",
			Inline: true,
		},
	}
	return embed
}
