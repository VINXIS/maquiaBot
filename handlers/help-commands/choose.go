package helpcommands

import (
	"github.com/bwmarrin/discordgo"
)

// Choose explains the choose functionality
func Choose(embed *discordgo.MessageEmbed) *discordgo.MessageEmbed {
	embed.Author.Name = "Command: ch / choose"
	embed.Description = "`(ch|choose) <option1> | <option2> | [option3]...` lets the bot choose from one of the options given."
	embed.Fields = []*discordgo.MessageEmbedField{
		{
			Name:  "<option1> | <option2> | [option3]...",
			Value: "List out the options separated by ` | `.",
		},
	}
	return embed
}
