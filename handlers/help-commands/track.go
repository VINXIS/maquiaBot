package helpcommands

import "github.com/bwmarrin/discordgo"

// Track explains the track functionality
func Track(embed *discordgo.MessageEmbed) *discordgo.MessageEmbed {
	embed.Author.Name = "Command: track / tr"
	embed.Description = "`[osu] (track|tr) [remove] <-u users> [-pp ppreq] [-l leaderboardreq] [-t topreq] [-s mapstatuses] [-m mode]` lets admins track plays for osu! users. If you want to change parts of the tracker, simply just use those flags. You do not need to restate all previous changes again."
	embed.Fields = []*discordgo.MessageEmbedField{
		{
			Name:  "[remove]",
			Value: "Write remove before any flags to remove users instead of adding them.",
		},
		{
			Name:   "<-u users>",
			Value:  "Users to add / remove.",
			Inline: true,
		},
		{
			Name:   "[-pp ppreq]",
			Value:  "The lowest pp required to post (Default: N/A).",
			Inline: true,
		},
		{
			Name:   "[-l leaderboardreq]",
			Value:  "The largest # required on a leaderboard to post with the largest # possible being 100 (Default: Top 100).",
			Inline: true,
		},
		{
			Name:   "[-t topreq]",
			Value:  "The largest # required on a best performance to post with the largest # possible being 100 (Default: Top 100).",
			Inline: true,
		},
		{
			Name:   "[-s mapstatuses]",
			Value:  "The types of map statuses to allow for posting. The options are `qualified`, `loved`, and `ranked` (Default: All 3).",
			Inline: true,
		},
		{
			Name:   "[-m mode]",
			Value:  "The type of mode to show. The options are `standard`, `taiko`, `catch`, `mania` (Default: standard).",
			Inline: true,
		},
		{
			Name:  "Example Usage:",
			Value: "`track -u VINXIS -pp 400 -l 50 -t 50 -s loved, qualified, ranked -m 0` Will track any play of VINXIS's which are either 400pp, top 50 on leaderboard, or top 50 in his tops, allowing all 3 types of maps to post, and for the mode osu!standard.",
		},
		{
			Name:  "Related Commands:",
			Value: "`trackinfo`, `tracktoggle`",
		},
	}
	return embed
}

// TrackToggle explains the track toggle functionality
func TrackToggle(embed *discordgo.MessageEmbed) *discordgo.MessageEmbed {
	embed.Author.Name = "Command: tt / trackt / ttoggle / tracktoggle"
	embed.Description = "`(tt|trackt|ttoggle|tracktoggle)` lets admins toggle whether tracking should run or not without removing the tracking information."
	embed.Fields = []*discordgo.MessageEmbedField{
		{
			Name:  "Related Commands:",
			Value: "`track`, `trackinfo`",
		},
	}
	return embed
}

// TrackInfo explains the track info functionality
func TrackInfo(embed *discordgo.MessageEmbed) *discordgo.MessageEmbed {
	embed.Author.Name = "Command: ti / tinfo / tracking / trackinfo"
	embed.Description = "`(ti|tinfo|tracking|trackinfo)` gives you the information regarding what types of scores are being tracked for the channel."
	embed.Fields = []*discordgo.MessageEmbedField{
		{
			Name:  "Related Commands:",
			Value: "`track`, `tracktoggle`",
		},
	}
	return embed
}
