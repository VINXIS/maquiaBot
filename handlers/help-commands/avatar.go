package helpcommands

import (
	"github.com/bwmarrin/discordgo"
)

// Avatar explains the avatar functionality
func Avatar(embed *discordgo.MessageEmbed) *discordgo.MessageEmbed {
	embed.Author.Name = "Command: avatar / ava / a"
	embed.Description = "`(avatar|ava|a) [@mentions|username|-s]` lets you obtain the avatar of yours or someone else's as an image link."
	embed.Fields = []*discordgo.MessageEmbedField{
		{
			Name:  "[@mentions|username|-s]",
			Value: "@mentions: Mention the person / people to get their username.\nusername: One user's username / nickname / ID to get the avatar for.\n-s: Use this flag for the server's icon.",
		},
	}
	return embed
}
