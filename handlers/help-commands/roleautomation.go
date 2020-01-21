package helpcommands

import (
	"github.com/bwmarrin/discordgo"
)

// RoleAutomation explains the role automation functionality
func RoleAutomation(embed *discordgo.MessageEmbed) *discordgo.MessageEmbed {
	embed.Author.Name = "Command: rolea / roleauto / roleautomation"
	embed.Description = "`(rolea|roleauto|roleautomation) <text> <role mentions / IDs>` lets admins create role automations ."
	embed.Fields = []*discordgo.MessageEmbedField{
		&discordgo.MessageEmbedField{
			Name:  "<text>",
			Value: "The text to add roles to the person",
		},
		&discordgo.MessageEmbedField{
			Name:  "<role mentions / IDs>",
			Value: "The roles to add to the person when the text is written.",
		},
		&discordgo.MessageEmbedField{
			Name:  "Related Commands:",
			Value: "`roleinfo`, `serverinfo`",
		},
	}
	return embed
}
