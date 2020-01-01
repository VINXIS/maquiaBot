package helpcommands

import (
	"github.com/bwmarrin/discordgo"
)

// Compare explains the compare functionality
func Compare(embed *discordgo.MessageEmbed) *discordgo.MessageEmbed {
	embed.Author.Name = "Command: c / compare"
	embed.Description = "`[osu] (c|compare) [link] <osu! username> [-m <mod> [-nostrict]|-all] [-sp [-mapper] [-sr]]` lets you show your score(s) on a map."
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
			Name:   "[-nostrict]",
			Value:  "If the score should have that mod combination exactly, or if it can have other mods included.",
			Inline: true,
		},
		&discordgo.MessageEmbedField{
			Name:   "[-all]",
			Value:  "Show all scores the user has made on the map.",
			Inline: true,
		},
		&discordgo.MessageEmbedField{
			Name:   "[-sp]",
			Value:  "Print out the score in a scorepost format after.",
			Inline: true,
		},
		&discordgo.MessageEmbedField{
			Name:   "[-mapper]",
			Value:  "Remove the mapset host from the scorepost generation.",
			Inline: true,
		},
		&discordgo.MessageEmbedField{
			Name:   "[-sr]",
			Value:  "Remove the star rating from the scorepost generation.",
			Inline: true,
		},
	}
	return embed
}
