package helpcommands

import (
	"github.com/bwmarrin/discordgo"
)

// Download explains the download
func Download(embed *discordgo.MessageEmbed) *discordgo.MessageEmbed {
	embed.Author.Name = "Command: dl / download"
	embed.Description = "`(dl|download)` lets users download data stored regarding them."
	embed.Fields = []*discordgo.MessageEmbedField{
		{
			Name:   "Data format",
			Value:  "https://github.com/VINXIS/maquiaBot/blob/master/structs/playerData.go#L12",
			Inline: true,
		},
		{
			Name:   "Related Commands:",
			Value:  "`downloadchannel`, `downloadserver`",
			Inline: true,
		},
	}
	return embed
}

// DownloadChannel explains the download channel functionality
func DownloadChannel(embed *discordgo.MessageEmbed) *discordgo.MessageEmbed {
	embed.Author.Name = "Command: dlch / dlchannel / downloadch / downloadchannel"
	embed.Description = "`(dlch|dlchannel|downloadch|downloadchannel)` lets admins download data stored regarding the channel."
	embed.Fields = []*discordgo.MessageEmbedField{
		{
			Name:   "Data format",
			Value:  "https://github.com/VINXIS/maquiaBot/blob/master/structs/channelData.go#L11",
			Inline: true,
		},
		{
			Name:   "Related Commands:",
			Value:  "`download`, `downloadserver`",
			Inline: true,
		},
	}
	return embed
}

// DownloadServer explains the download server functionality
func DownloadServer(embed *discordgo.MessageEmbed) *discordgo.MessageEmbed {
	embed.Author.Name = "Command: dlsv / dlserver / downloadsv / downloadserver"
	embed.Description = "`(dlsv|dlserver|downloadsv|downloadserver)` lets admins download data stored regarding the server."
	embed.Fields = []*discordgo.MessageEmbedField{
		{
			Name:   "Data format",
			Value:  "https://github.com/VINXIS/maquiaBot/blob/master/structs/serverData.go#L11",
			Inline: true,
		},
		{
			Name:   "Related Commands:",
			Value:  "`download`, `downloadchannel`",
			Inline: true,
		},
	}
	return embed
}
