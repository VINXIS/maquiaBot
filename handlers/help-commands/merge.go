package helpcommands

import (
	"github.com/bwmarrin/discordgo"
)

// Merge explains the merge functionality
func Merge(embed *discordgo.MessageEmbed) *discordgo.MessageEmbed {
	embed.Author.Name = "Command: merge"
	embed.Description = "`merge [link] [link 2] <[link 3] ...>` let's you merge images into 1 side by side."
	embed.Fields = []*discordgo.MessageEmbedField{
		{
			Name:   "[link] [link 2] <[link 3] ...>",
			Value:  "The links to merge. They will be merged in order from left to right.",
			Inline: true,
		},
	}
	return embed
}
