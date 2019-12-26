package helpcommands

import (
	"github.com/bwmarrin/discordgo"
)

// ScorePost explains the score post functionality
func ScorePost(embed *discordgo.MessageEmbed) *discordgo.MessageEmbed {
	embed.Author.Name = "Command: s / sc / scorepost"
	embed.Description = "`[osu] (s|sc|scorepost) [link] [osu! username] [-m <mod>]` prints out a scorepost."
	embed.Fields = []*discordgo.MessageEmbedField{
		&discordgo.MessageEmbedField{
			Name:   "[link]",
			Value:  "The map to find the score for. No link will look for a score on the most recently linked map previously. If no link / osu! username / mods are given, then the bot will make a scorepost based off of the most recent play sent by the bot.",
			Inline: true,
		},
		&discordgo.MessageEmbedField{
			Name:   "[osu! username]",
			Value:  "The username of the osu! player to find the score for.",
			Inline: true,
		},
		&discordgo.MessageEmbedField{
			Name:   "[-m <mod>]",
			Value:  "The score's mod combination to look for.",
			Inline: true,
		},
	}
	return embed
}
