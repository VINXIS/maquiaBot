package helpcommands

import (
	"github.com/bwmarrin/discordgo"
)

// Link explains the link functionality
func Link(embed *discordgo.MessageEmbed) *discordgo.MessageEmbed {
	embed.Author.Name = "Command: link / set"
	embed.Description = "`[osu] (link|set) [@mention] <osu! username>` lets you link an osu! account with the username given to your discord account."
	embed.Fields = []*discordgo.MessageEmbedField{
		{
			Name:   "[@mention]",
			Value:  "The person to link the osu! user to **(REQUIRES ADMIN PERMS)**.",
			Inline: true,
		},
		{
			Name:   "<osu! username>",
			Value:  "The username of the osu! player to link to.",
			Inline: true,
		},
	}
	return embed
}
