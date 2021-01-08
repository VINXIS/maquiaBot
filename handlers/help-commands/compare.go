package helpcommands

import (
	"github.com/bwmarrin/discordgo"
)

// Compare explains the compare functionality
func Compare(embed *discordgo.MessageEmbed) *discordgo.MessageEmbed {
	embed.Author.Name = "Command: c / compare"
	embed.Description = "`[osu] (c|compare) [link] <osu! username> [-m <mod> [-nostrict]|-all] [-sp [-mapper] [-sr] [-fc] [-add <text>]]` lets you show your score(s) on a map."
	embed.Fields = []*discordgo.MessageEmbedField{
		{
			Name:   "[link]",
			Value:  "The map to find the score for. No link will look for a score on the most recently linked map previously.",
			Inline: true,
		},
		{
			Name:   "<osu! username>",
			Value:  "The username of the osu! player to find the score for.",
			Inline: true,
		},
		{
			Name:   "[-m mod]",
			Value:  "The score's mod combination to look for.",
			Inline: true,
		},
		{
			Name:   "[-nostrict]",
			Value:  "If the score should have that mod combination exactly, or if it can have other mods included.",
			Inline: true,
		},
		{
			Name:   "[-all]",
			Value:  "Show all scores the user has made on the map.",
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
		{
			Name:   "[-fc]",
			Value:  "Adds pp for if the score was an FC for the scoprepost generation.",
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
