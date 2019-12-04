package helpcommands

import (
	"github.com/bwmarrin/discordgo"
)

// NiceIdea explains the nice idea functionality
func NiceIdea(embed *discordgo.MessageEmbed) *discordgo.MessageEmbed {
	embed.Author.Name = "Command: idea / niceidea"
	embed.Description = "`(idea|niceidea)` lets you send a nice idea if an admin disabled automatic nice idea posting."
	embed.Fields = []*discordgo.MessageEmbedField{
		&discordgo.MessageEmbedField{
			Name:  "Related Commands:",
			Value: "`ideatoggle`",
		},
	}
	return embed
}

// NiceIdeaToggle explains the nice idea toggle functionality
func NiceIdeaToggle(embed *discordgo.MessageEmbed) *discordgo.MessageEmbed {
	embed.Author.Name = "Command: it / ideat / itoggle / ideatoggle"
	embed.Description = "`(it|ideat|itoggle|ideatoggle)` lets admins toggle whether any text containing nice idea, good idea, or gud idea will send a nice idea."
	embed.Fields = []*discordgo.MessageEmbedField{
		&discordgo.MessageEmbedField{
			Name:  "Related Commands:",
			Value: "`idea`",
		},
	}
	return embed
}
