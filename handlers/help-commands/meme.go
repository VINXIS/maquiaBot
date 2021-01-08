package helpcommands

import (
	"github.com/bwmarrin/discordgo"
)

// Meme explains the meme functionality
func Meme(embed *discordgo.MessageEmbed) *discordgo.MessageEmbed {
	embed.Author.Name = "Command: meme"
	embed.Description = "`meme [link] <[top text] | [bottom text]>` let's you generate a bad meme."
	embed.Fields = []*discordgo.MessageEmbedField{
		{
			Name:   "[link]",
			Value:  "The image link, no link given will use the attachment given if there is one, or the most recent image.",
			Inline: true,
		},
		{
			Name:   "<top text | bottom text>",
			Value:  "The text you want to add, separated by a `|`. You do not need both. If you only want the bottom text for example, simply only do `| [bottom text]`.",
			Inline: true,
		},
	}
	return embed
}
