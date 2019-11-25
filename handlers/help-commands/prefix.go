package helpcommands

import (
	"github.com/bwmarrin/discordgo"
)

// Prefix explains the prefix functionality
func Prefix(embed *discordgo.MessageEmbed) *discordgo.MessageEmbed {
	embed.Author.Name = "Command: maquiaprefix / prefix / newprefix"
	embed.Description = "`(maquiaprefix|prefix|newprefix) <prefix>` lets admins change the prefix for the bot in the server."
	embed.Fields = []*discordgo.MessageEmbedField{
		&discordgo.MessageEmbedField{
			Name:   "<prefix>",
			Value:  "The prefix to change to.",
		},
	}
	return embed
}
