package helpcommands

import (
	"github.com/bwmarrin/discordgo"
)

// Emoji explains the emoji functionality
func Emoji(embed *discordgo.MessageEmbed) *discordgo.MessageEmbed {
	embed.Author.Name = "Command: e / emoji / emote"
	embed.Description = "`(e|emoji|emote) <emote>` gives you the emoji as an image."
	embed.Fields = []*discordgo.MessageEmbedField{
		&discordgo.MessageEmbedField{
			Name:  "<emote>",
			Value: "The emoji to give as an image.",
		},
	}
	return embed
}
