package helpcommands

import (
	"github.com/bwmarrin/discordgo"
)

// Parse explains the parse functionality
func Parse(embed *discordgo.MessageEmbed) *discordgo.MessageEmbed {
	embed.Author.Name = "Command: parse"
	embed.Description = "`parse <snowflake ID>` parses a discord / twitter ID."
	embed.Fields = []*discordgo.MessageEmbedField{
		{
			Name:  "<snowflake ID>",
			Value: "The discord / twitter ID to parse.",
		},
	}
	return embed
}
