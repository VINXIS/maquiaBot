package helpcommands

import (
	"github.com/bwmarrin/discordgo"
)

// Top explains the top functionality
func Top(embed *discordgo.MessageEmbed) *discordgo.MessageEmbed {
	embed.Author.Name = "Command: t / top"
	embed.Description = "`[osu] (t|top) [osu! username] [num] [-m mod] [-sp [-mapper] [-sr]]` shows the player's top 100 pp score."
	embed.Fields = []*discordgo.MessageEmbedField{
		{
			Name:   "[osu! username]",
			Value:  "The osu! user to check. No user given will use the account linked to your discord account.",
			Inline: true,
		},
		{
			Name:   "[num]",
			Value:  "The nth most pp score to find (Default: Top play).",
			Inline: true,
		},
		{
			Name:   "[-m mod]",
			Value:  "The mods to check for.",
			Inline: true,
		},
		{
			Name:   "[-sp]",
			Value:  "Print out the score in a scorepost format after.",
			Inline: true,
		},
		{
			Name:   "[-mapper]",
			Value:  "Remove the mapset host from the scorepost generation.",
			Inline: true,
		},
		{
			Name:   "[-sr]",
			Value:  "Remove the star rating from the scorepost generation.",
			Inline: true,
		},
	}
	return embed
}
