package helpcommands

import (
	"github.com/bwmarrin/discordgo"
)

// RemoveChannel explains the remove channel functionality
func RemoveChannel(embed *discordgo.MessageEmbed) *discordgo.MessageEmbed {
	embed.Author.Name = "Command: rmch / rmchannel / removech / removechannel"
	embed.Description = "`(rmch|rmchannel|removech|removechannel)` lets admins remove data stored regarding the channel."
	embed.Fields = []*discordgo.MessageEmbedField{
		{
			Name:  "Related Commands:",
			Value: "`removeserver`",
		},
	}
	return embed
}

// RemoveServer explains the remove server functionality
func RemoveServer(embed *discordgo.MessageEmbed) *discordgo.MessageEmbed {
	embed.Author.Name = "Command: rmsv / rmserver / removesv / removeserver"
	embed.Description = "`(rmsv|rmserver|removesv|removeserver)` lets admins remove data stored regarding the server."
	embed.Fields = []*discordgo.MessageEmbedField{
		{
			Name:  "Related Commands:",
			Value: "`removechannel`",
		},
	}
	return embed
}
