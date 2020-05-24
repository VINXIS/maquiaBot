package helpcommands

import (
	"github.com/bwmarrin/discordgo"
)

// Farm explains the farm functionality
func Farm(embed *discordgo.MessageEmbed) *discordgo.MessageEmbed {
	embed.Author.Name = "Command: farm"
	embed.Description = "`[osu] farm [username] [-n num]` shows your farm rating alongside your most farmy scores."
	embed.Fields = []*discordgo.MessageEmbedField{
		{
			Name:   "[username]",
			Value:  "The osu! user to check their farm rating for.",
			Inline: true,
		},
		{
			Name:   "[-a num]",
			Value:  "Number of plays to show (Default: 5).",
			Inline: true,
		},
		{
			Name:  "Related Commands:",
			Value: "`bottomfarm`, `topfarm`",
		},
	}
	return embed
}

// BottomFarm explains the bottom farm functionality
func BottomFarm(embed *discordgo.MessageEmbed) *discordgo.MessageEmbed {
	embed.Author.Name = "Command: bfarm / bottomfarm"
	embed.Description = "`[osu] (bfarm|bottomfarm) [-s] [num]` shows the best farmerdogs."
	embed.Fields = []*discordgo.MessageEmbedField{
		{
			Name:   "[-s]",
			Value:  "Add this flag if you only want the farm ranking for people in the server.",
			Inline: true,
		},
		{
			Name:   "[num]",
			Value:  "Number of players to show (Default: 1).",
			Inline: true,
		},
		{
			Name:  "Related Commands:",
			Value: "`farm`, `topfarm`",
		},
	}
	return embed
}

// TopFarm explains the top farm functionality
func TopFarm(embed *discordgo.MessageEmbed) *discordgo.MessageEmbed {
	embed.Author.Name = "Command: tfarm / topfarm"
	embed.Description = "`[osu] (tfarm|topfarm) [-s] [num]` shows the worst farmerdogs."
	embed.Fields = []*discordgo.MessageEmbedField{
		{
			Name:   "[-s]",
			Value:  "Add this flag if you only want the farm ranking for people in the server.",
			Inline: true,
		},
		{
			Name:   "[num]",
			Value:  "Number of players to show (Default: 1).",
			Inline: true,
		},
		{
			Name:  "Related Commands:",
			Value: "`farm`, `bottomfarm`",
		},
	}
	return embed
}
