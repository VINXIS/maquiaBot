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
		{
			Name:   "Related Commands:",
			Value:  "`unlink`",
			Inline: true,
		},
	}
	return embed
}

// Unlink explains the unlink functionality
func Unlink(embed *discordgo.MessageEmbed) *discordgo.MessageEmbed {
	embed.Author.Name = "Command: u / unlink"
	embed.Description = "`(u|unlink)` lets you delete all data related to your discord ID."
	embed.Fields = []*discordgo.MessageEmbedField{
		{
			Name:   "Related Commands:",
			Value:  "`link`",
			Inline: true,
		},
	}
	return embed
}
