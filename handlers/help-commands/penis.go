package helpcommands

import (
	"github.com/bwmarrin/discordgo"
)

// Penis explains the penis functionality
func Penis(embed *discordgo.MessageEmbed) *discordgo.MessageEmbed {
	embed.Author.Name = "Command: penis"
	embed.Description = "`penis [username]` calculates your erect length for today."
	embed.Fields = []*discordgo.MessageEmbedField{
		&discordgo.MessageEmbedField{
			Name:  "[username]",
			Value: "Gets the erect length for the given username / nickname / ID. Gives your erect length if no username / nickname / ID is given.",
		},
	}
	return embed
}
