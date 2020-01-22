package helpcommands

import (
	"github.com/bwmarrin/discordgo"
)

// Trigger explains the trigger functionality
func Trigger(embed *discordgo.MessageEmbed) *discordgo.MessageEmbed {
	embed.Author.Name = "Command: trigger"
	embed.Description = "`trigger (<[trigger] | [result]>|<<word> <result>>)` let's you custom word / line triggers (*technically* custom functions).\nYou may also use regex for triggers! Test your regex here to create a valid regex: https://regex101.com/"
	embed.Fields = []*discordgo.MessageEmbedField{
		&discordgo.MessageEmbedField{
			Name:   "<[trigger] | [result]>",
			Value:  "The text to trigger, then a |, then the result to send.",
			Inline: true,
		},
		&discordgo.MessageEmbedField{
			Name:   "<<word> <result>>",
			Value:  "The word to trigger, followed by the result to send",
			Inline: true,
		},
		&discordgo.MessageEmbedField{
			Name:  "Related Commands:",
			Value: "`triggers`",
		},
	}
	return embed
}

// Triggers explains the triggers functionality
func Triggers(embed *discordgo.MessageEmbed) *discordgo.MessageEmbed {
	embed.Author.Name = "Command: triggers"
	embed.Description = "`trigger` lists out the currently enabled triggers for the server."
	embed.Fields = []*discordgo.MessageEmbedField{
		&discordgo.MessageEmbedField{
			Name:  "Related Commands:",
			Value: "`trigger`",
		},
	}
	return embed
}
