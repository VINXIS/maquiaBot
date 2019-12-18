package helpcommands

import (
	"github.com/bwmarrin/discordgo"
)

// NiceIdea explains the nice idea functionality
func NiceIdea(embed *discordgo.MessageEmbed) *discordgo.MessageEmbed {
	embed.Author.Name = "Command: idea / niceidea"
	embed.Description = "`(idea|niceidea)` lets you send a nice idea."
	embed.Fields = []*discordgo.MessageEmbedField{
		&discordgo.MessageEmbedField{
			Name:  "Related Commands:",
			Value: "`toggle`",
		},
	}
	return embed
}
