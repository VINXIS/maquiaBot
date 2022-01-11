package helpcommands

import (
	"github.com/bwmarrin/discordgo"
)

// ScorePost explains the score post functionality
func ScorePost(embed *discordgo.MessageEmbed) *discordgo.MessageEmbed {
	embed.Author.Name = "Command: s / sc / scorepost"
	embed.Description = "`[osu] (s|sc|scorepost) [link] [osu! username] [-m <mod>] [-mapper] [-sr] [-fc] [-c] [-mc] [-star] [-b] [-add <text>]` prints out a scorepost."
	embed.Fields = []*discordgo.MessageEmbedField{
		{
			Name:   "[link]",
			Value:  "The map to find the score for. No link will look for a score on the most recently linked map previously. If no link / osu! username / mods are given, then the bot will make a scorepost based off of the most recent play sent by the bot.",
			Inline: true,
		},
		{
			Name:   "[osu! username]",
			Value:  "The username of the osu! player to find the score for.",
			Inline: true,
		},
		{
			Name:   "[-m mod]",
			Value:  "The score's mod combination to look for.",
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
		{
			Name:   "[-fc]",
			Value:  "Adds pp for if the score was an FC for the scoprepost generation.",
			Inline: true,
		},
		{
			Name:   "[-c]",
			Value:  "Adds commas instead of a | between mapper and SR.",
			Inline: true,
		},
		{
			Name:   "[-mc]",
			Value:  "Adds commas between each mod the play used.",
			Inline: true,
		},
		{
			Name:   "[-star]",
			Value:  "Adds a unicode star ★ instead of *.",
			Inline: true,
		},
		{
			Name:   "[-b]",
			Value:  "Adds brackets around the mapper and SR.",
			Inline: true,
		},
		{
			Name:   "[-add text]",
			Value:  "Any text to append to the end for the scoprepost generation.",
			Inline: true,
		},
	}
	return embed
}
