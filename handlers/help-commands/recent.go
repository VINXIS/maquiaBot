package helpcommands

import (
	"github.com/bwmarrin/discordgo"
)

// Recent explains the recent functionality
func Recent(embed *discordgo.MessageEmbed) *discordgo.MessageEmbed {
	embed.Author.Name = "Command: r / rs / recent"
	embed.Description = "`[osu] (r|rs|recent) [osu! username] [num] [-m mod] [-sp]` shows the player's recent score."
	embed.Fields = []*discordgo.MessageEmbedField{
		&discordgo.MessageEmbedField{
			Name:   "[osu! username]",
			Value:  "The osu! user to check. No user given will use the account linked to your discord account.",
			Inline: true,
		},
		&discordgo.MessageEmbedField{
			Name:   "[num]",
			Value:  "The nth recent score to find (Default: Latest).",
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
		&discordgo.MessageEmbedField{
			Name:  "Related Commands:",
			Value: "`recentbest`",
		},
	}
	return embed
}

// RecentBest explains the recentbest functionality
func RecentBest(embed *discordgo.MessageEmbed) *discordgo.MessageEmbed {
	embed.Author.Name = "Command: rb / recentb / recentbest"
	embed.Description = "`[osu] (rb|recentb|recentbest) [osu! username] [num] [-m mod] [-sp]` shows the player's recent top 100 pp score."
	embed.Fields = []*discordgo.MessageEmbedField{
		&discordgo.MessageEmbedField{
			Name:   "[osu! username]",
			Value:  "The osu! user to check. No user given will use the account linked to your discord account.",
			Inline: true,
		},
		&discordgo.MessageEmbedField{
			Name:   "[num]",
			Value:  "The nth recent top score to find (Default: Latest).",
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
		&discordgo.MessageEmbedField{
			Name:  "Related Commands:",
			Value: "`recentbest`",
		},
	}
	return embed
}
