package helpcommands

import (
	"github.com/bwmarrin/discordgo"
)

// Cheers explains the cheers functionality
func Cheers(embed *discordgo.MessageEmbed) *discordgo.MessageEmbed {
	embed.Author.Name = "Command: cheers"
	embed.Description = "`cheers` lets you send a cheers video."
	embed.Fields = []*discordgo.MessageEmbedField{
		&discordgo.MessageEmbedField{
			Name:  "Related Commands:",
			Value: "`cheerstoggle`",
		},
	}
	return embed
}

// CheersToggle explains the cheers toggle functionality
func CheersToggle(embed *discordgo.MessageEmbed) *discordgo.MessageEmbed {
	embed.Author.Name = "Command: cht / cheerst / chtoggle / cheerstoggle"
	embed.Description = "`(cgt|cheerst|chtoggle|cheerstoggle)` lets admins toggle whether any text containing 🍻 / 🍺 / 🦐 / cheer (even within words) will send a cheers video."
	embed.Fields = []*discordgo.MessageEmbedField{
		&discordgo.MessageEmbedField{
			Name:  "Related Commands:",
			Value: "`cheers`",
		},
	}
	return embed
}
