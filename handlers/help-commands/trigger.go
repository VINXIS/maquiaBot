package helpcommands

import (
	"github.com/bwmarrin/discordgo"
)

// Trigger explains the trigger functionality
func Trigger(embed *discordgo.MessageEmbed) *discordgo.MessageEmbed {
	embed.Author.Name = "Command: trigger"
	embed.Description = "`trigger (<[text] | [result]>|<<word> <result>>|<-d <number>>)` let's you create custom word / line triggers (*technically* custom functions).\nYou may also use regex for triggers! Test your regex here to create a valid regex: https://regex101.com/"
	embed.Fields = []*discordgo.MessageEmbedField{
		{
			Name:   "<[text] | [result]>",
			Value:  "The text / regex, then a |, then the result to send.",
			Inline: true,
		},
		{
			Name:   "<<word> <result>>",
			Value:  "The word to trigger, followed by the result to send",
			Inline: true,
		},
		{
			Name:   "<-d <number>>",
			Value:  "`-d` followed be the trigger ID found in `triggers`",
			Inline: true,
		},
		{
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
		{
			Name:  "Related Commands:",
			Value: "`trigger`",
		},
	}
	return embed
}
