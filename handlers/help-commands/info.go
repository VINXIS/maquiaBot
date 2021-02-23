package helpcommands

import (
	"github.com/bwmarrin/discordgo"
)

// Info explains the info functionality
func Info(embed *discordgo.MessageEmbed) *discordgo.MessageEmbed {
	embed.Author.Name = "Command: info"
	embed.Description = "`info [username]` gets the information for a user."
	embed.Fields = []*discordgo.MessageEmbedField{
		{
			Name:  "[username]",
			Value: "Gets the info for the given username / nickname / ID. Gives your info if no username / nickname / ID is given",
		},
		{
			Name:  "Related Commands:",
			Value: "`channelinfo`, `roleinfo`, `serverinfo`",
		},
	}
	return embed
}

// RoleInfo explains the role info functionality
func RoleInfo(embed *discordgo.MessageEmbed) *discordgo.MessageEmbed {
	embed.Author.Name = "Command: rinfo / roleinfo"
	embed.Description = "`(rinfo|roleinfo) [rolename]` gets the information for a role / the server's role automations"
	embed.Fields = []*discordgo.MessageEmbedField{
		{
			Name:  "[rolename]",
			Value: "The role to get the information for. No role will list the server's role automations instead.",
		},
		{
			Name:  "Related Commands:",
			Value: "`info`, `channelinfo`, `serverinfo`",
		},
	}
	return embed
}

// ChannelInfo explains the channel info functionality
func ChannelInfo(embed *discordgo.MessageEmbed) *discordgo.MessageEmbed {
	embed.Author.Name = "Command: cinfo / chinfo / channelinfo"
	embed.Description = "`(cinfo|chinfo|channelinfo)` gets the information for this channel."
	embed.Fields = []*discordgo.MessageEmbedField{
		{
			Name:  "Related Commands:",
			Value: "`info`, `roleinfo`, `serverinfo`",
		},
	}
	return embed
}

// ServerInfo explains the server info functionality
func ServerInfo(embed *discordgo.MessageEmbed) *discordgo.MessageEmbed {
	embed.Author.Name = "Command: sinfo / serverinfo"
	embed.Description = "`(sinfo|serverinfo)` gets the information for this server."
	embed.Fields = []*discordgo.MessageEmbedField{
		{
			Name:  "Related Commands:",
			Value: "`info`, `channelinfo`, `roleinfo`",
		},
	}
	return embed
}
