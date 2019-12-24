package helpcommands

import (
	"github.com/bwmarrin/discordgo"
)

// ScorePost explains the score post functionality
func ScorePost(embed *discordgo.MessageEmbed) *discordgo.MessageEmbed {
	embed.Author.Name = "Command: s / sc / scorepost"
	embed.Description = "`[osu] (s|sc|scorepost) [link] <osu! username> [-m <mod>] [-sp]` prints out a scorepost."
	embed.Fields = []*discordgo.MessageEmbedField{
		&discordgo.MessageEmbedField{
			Name:   "[link]",
			Value:  "The map to find the score for. No link will look for a score on the most recently linked map previously.",
			Inline: true,
		},
		&discordgo.MessageEmbedField{
			Name:   "<osu! username>",
			Value:  "The username of the osu! player to find the score for.",
			Inline: true,
		},
		&discordgo.MessageEmbedField{
			Name:   "[-m <mod>]",
			Value:  "The score's mod combination to look for.",
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
