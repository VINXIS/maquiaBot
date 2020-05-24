package helpcommands

import (
	"github.com/bwmarrin/discordgo"
)

// Profile explains the profile functionality
func Profile(embed *discordgo.MessageEmbed) *discordgo.MessageEmbed {
	embed.Author.Name = "Command: osu / profile"
	embed.Description = "`(osu|[osu] profile|<profile link>) [osu! username] [-m <mode>]` lets you obtain user information."
	embed.Fields = []*discordgo.MessageEmbedField{
		{
			Name:   "<profile link>",
			Value:  "You may link a map instead of using `osu` or `profile` to get user information.",
			Inline: true,
		},
		{
			Name:   "[osu! username]",
			Value:  "The username to look for. Using a link will use the user linked instead. No user linked for `osu` or `profile` messages will use the user linked to your discord account.",
			Inline: true,
		},
		{
			Name:   "[-m <mode>]",
			Value:  "The mode to show user information for (Default: osu!standard).",
			Inline: true,
		},
		{
			Name:  "Related Commands:",
			Value: "`osudetail`, `osutop`",
		},
	}
	return embed
}

// ProfileDetail explains the profiledetail functionality
func ProfileDetail(embed *discordgo.MessageEmbed) *discordgo.MessageEmbed {
	embed.Author.Name = "Command: osudetail"
	embed.Description = "`osudetail [osu! username] [-m <mode>]` lets you obtain detailed user information."
	embed.Fields = []*discordgo.MessageEmbedField{
		{
			Name:   "[osu! username]",
			Value:  "The username to look for. Using a link will use the user linked instead. No user linked for `osu` or `profile` messages will use the user linked to your discord account.",
			Inline: true,
		},
		{
			Name:   "[-m <mode>]",
			Value:  "The mode to show user information for (Default: osu!standard).",
			Inline: true,
		},
		{
			Name:  "Related Commands:",
			Value: "`osu`, `osutop`",
		},
	}
	return embed
}

// ProfileTop explains the profiletop functionality
func ProfileTop(embed *discordgo.MessageEmbed) *discordgo.MessageEmbed {
	embed.Author.Name = "Command: osutop"
	embed.Description = "`osutop [osu! username] [-m <mode>] [-r]` lets you obtain user top information."
	embed.Fields = []*discordgo.MessageEmbedField{
		{
			Name:   "[osu! username]",
			Value:  "The username to look for. Using a link will use the user linked instead. No user linked for `osu` or `profile` messages will use the user linked to your discord account.",
			Inline: true,
		},
		{
			Name:   "[-m <mode>]",
			Value:  "The mode to show user information for (Default: osu!standard).",
			Inline: true,
		},
		{
			Name:   "[-r]",
			Value:  "Show in chronological order instead of by PP (includes the graph).",
			Inline: true,
		},
		{
			Name:  "Related Commands:",
			Value: "`osu`, `osudetail`",
		},
	}
	return embed
}
