package helpcommands

import (
	"github.com/bwmarrin/discordgo"
)

// Info explains the info functionality
func Info(embed *discordgo.MessageEmbed) *discordgo.MessageEmbed {
	embed.Author.Name = "Command: info"
	embed.Description = "`info [username]` gets the information for a user."
	embed.Fields = []*discordgo.MessageEmbedField{
		&discordgo.MessageEmbedField{
			Name:  "[username]",
			Value: "Gets the info for the given username / nickname / ID. Gives your info if no username / nickname / ID is given",
		},
		&discordgo.MessageEmbedField{
			Name:  "Related Commands:",
			Value: "`serverinfo`",
		},
	}
	return embed
}

// ServerInfo explains the server info functionality
func ServerInfo(embed *discordgo.MessageEmbed) *discordgo.MessageEmbed {
	embed.Author.Name = "Command: sinfo / serverinfo"
	embed.Description = "`(sinfo|serverinfo)` gets the information for this server."
	embed.Fields = []*discordgo.MessageEmbedField{
		&discordgo.MessageEmbedField{
			Name:  "Related Commands:",
			Value: "`info`",
		},
	}
	return embed
}
