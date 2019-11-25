package helpcommands

import (
	"github.com/bwmarrin/discordgo"
)

// Funny explains the face functionality
func Funny(embed *discordgo.MessageEmbed) *discordgo.MessageEmbed {
	embed.Author.Name = "Command: funny"
	embed.Description = "`funny [username]` calculates how funny you are."
	embed.Fields = []*discordgo.MessageEmbedField{
		&discordgo.MessageEmbedField{
			Name:  "[username]",
			Value: "Gets the funny value for the given username / nickname. Gives your funny value if no username / nickname is given.",
		},
	}
	return embed
}
