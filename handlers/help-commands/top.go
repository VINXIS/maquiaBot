package helpcommands

import (
	"github.com/bwmarrin/discordgo"
)

// Top explains the top functionality
func Top(embed *discordgo.MessageEmbed) *discordgo.MessageEmbed {
	embed.Author.Name = "Command: t / top"
	embed.Description = "`[osu] (t|top) [osu! username] [num] [-m mod] [-sp]` shows the player's top 100 pp score."
	embed.Fields = []*discordgo.MessageEmbedField{
		&discordgo.MessageEmbedField{
			Name:   "[osu! username]",
			Value:  "The osu! user to check. No user given will use the account linked to your discord account.",
			Inline: true,
		},
		&discordgo.MessageEmbedField{
			Name:   "[num]",
			Value:  "The nth most pp score to find (Default: Top play).",
			Inline: true,
		},
		&discordgo.MessageEmbedField{
			Name:   "[-m mod]",
			Value:  "The mods to check for.",
			Inline: true,
		},
		&discordgo.MessageEmbedField{
			Name:   "[-sp]",
			Value:  "Print out the score in a scorepost format after.",
			Inline: true,
		},
	}
	return embed
}
