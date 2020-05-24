package helpcommands

import (
	"github.com/bwmarrin/discordgo"
)

// Berry explains the berry functionality
func Berry(embed *discordgo.MessageEmbed) *discordgo.MessageEmbed {
	embed.Author.Name = "Command: b / berry"
	embed.Description = "`[pokemon] (b|berry) <berry name>` lets you obtain information about a pokemon berry."
	embed.Fields = []*discordgo.MessageEmbedField{
		{
			Name:  "<berry name>",
			Value: "The berry to get information for.",
		},
	}
	return embed
}
