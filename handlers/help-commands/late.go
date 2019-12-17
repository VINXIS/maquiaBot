package helpcommands

import (
	"github.com/bwmarrin/discordgo"
)

// Late explains the late functionality
func Late(embed *discordgo.MessageEmbed) *discordgo.MessageEmbed {
	embed.Author.Name = "Command: late / old / ancient"
	embed.Description = "`(late|old|ancient)` lets you send a late video."
	embed.Fields = []*discordgo.MessageEmbedField{
		&discordgo.MessageEmbedField{
			Name:  "Related Commands:",
			Value: "`latetoggle`",
		},
	}
	return embed
}

// LateToggle explains the late toggle functionality
func LateToggle(embed *discordgo.MessageEmbed) *discordgo.MessageEmbed {
	embed.Author.Name = "Command: lt / latet / ltoggle / latetoggle"
	embed.Description = "`(lt|latet|ltoggle|latetoggle)` lets admins toggle whether any text containing late / old / ancient (even within words) will send a late video."
	embed.Fields = []*discordgo.MessageEmbedField{
		&discordgo.MessageEmbedField{
			Name:  "Related Commands:",
			Value: "`late`",
		},
	}
	return embed
}
