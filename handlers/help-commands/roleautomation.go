package helpcommands

import (
	"github.com/bwmarrin/discordgo"
)

// RoleAutomation explains the role automation functionality
func RoleAutomation(embed *discordgo.MessageEmbed) *discordgo.MessageEmbed {
	embed.Author.Name = "Command: rolea / roleauto / roleautomation"
	embed.Description = "`(rolea|roleauto|roleautomation) (<-d <number>>|<<text> <role mentions / IDs>>` lets admins create role automations ."
	embed.Fields = []*discordgo.MessageEmbedField{
		{
			Name:  "<text> <role mentions / IDs>",
			Value: "The text / regex to add roles to the person, and then the list of role mentions / IDs to add to the person when they send the text / regex",
		},
		{
			Name:   "<-d <number>>",
			Value:  "`-d` followed be the role automation ID found in `roleinfo`",
			Inline: true,
		},
		{
			Name:  "Related Commands:",
			Value: "`roleinfo`",
		},
	}
	return embed
}
