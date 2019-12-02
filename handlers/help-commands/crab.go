package helpcommands

import (
	"github.com/bwmarrin/discordgo"
)

// Crab explains the crab functionality
func Crab(embed *discordgo.MessageEmbed) *discordgo.MessageEmbed {
	embed.Author.Name = "Command: crab"
	embed.Description = "`crab` lets you send a crab rave gif if an admin disabled automatic crab rave posting."
	embed.Fields = []*discordgo.MessageEmbedField{
		&discordgo.MessageEmbedField{
			Name:  "Related Commands:",
			Value: "`crabtoggle`",
		},
	}
	return embed
}

// CrabToggle explains the crab toggle functionality
func CrabToggle(embed *discordgo.MessageEmbed) *discordgo.MessageEmbed {
	embed.Author.Name = "Command: ct / crabt / ctoggle / crabtoggle"
	embed.Description = "`(ct|crabt|ctoggle|crabtoggle)` lets admins toggle whether any text containing crab / rave (even within words) will send a crab rave gif."
	embed.Fields = []*discordgo.MessageEmbedField{
		&discordgo.MessageEmbedField{
			Name:  "Related Commands:",
			Value: "`crab`",
		},
	}
	return embed
}
