package helpcommands

import (
	"github.com/bwmarrin/discordgo"
)

// Similar explains the score post functionality
func Similar(embed *discordgo.MessageEmbed) *discordgo.MessageEmbed {
	embed.Author.Name = "Command: similar"
	embed.Description = "`[osu] similar [link] [-list]` gives a random beatmap (or a list of beatmaps) that is/are considered similar according to osuMapMatcher https://github.com/Xarib/OsuMapMatcher."
	embed.Fields = []*discordgo.MessageEmbedField{
		{
			Name:   "[link]",
			Value:  "The map to find similar maps of. No link will look for the most recently linked map previously.",
			Inline: true,
		},
		{
			Name:   "[-list]",
			Value:  "Will provide a list of 10 beatmaps instead of a random one.",
			Inline: true,
		},
		{
			Name:  "**WARNING**",
			Value: "Please note that this feature only works for ranked/approved/loved beatmaps that are not **recently** ranked/approved/loved as well. If you wish to see the list of beatmap IDs that are contained within this API, see https://omm.xarib.ch/api/knn/maps",
		},
	}
	return embed
}
