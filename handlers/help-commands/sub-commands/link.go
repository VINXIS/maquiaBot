package helpsubcommands

import (
	"github.com/bwmarrin/discordgo"
)

// Link explains the link functionality
func Link(embed *discordgo.MessageEmbed, arg string) *discordgo.MessageEmbed {
	embed.Author.Name = "Command: " + arg
	embed.Description = "`link [mention] <osu! username>` lets you link an osu! account with the username given to your discord account."
	embed.Fields = []*discordgo.MessageEmbedField{
		&discordgo.MessageEmbedField{
			Name:   "[mention]",
			Value:  "The person to link the osu! user to (REQUIRES ADMIN PERMISSIONS)",
			Inline: true,
		},
		&discordgo.MessageEmbedField{
			Name:   "<osu! username>",
			Value:  "The username of the osu! player to link to",
			Inline: true,
		},
	}
	return embed
}
