package helpcommands

import (
	"github.com/bwmarrin/discordgo"
)

// AvatarQuote explains the avatarquote functionality
func AvatarQuote(embed *discordgo.MessageEmbed) *discordgo.MessageEmbed {
	embed.Author.Name = "Command: aq / avaquote / quoteava / avatarquote / quoteavatar"
	embed.Description = "`(aq|avaquote|quoteava|avatarquote|quoteavatar) (@mentions|username)` runs the avatar command and then the quote command one after another."
	embed.Fields = []*discordgo.MessageEmbedField{
		{
			Name:  "(@mentions|username)",
			Value: "Mention / provide their username / ID to have the avatar command and quote command run simultaneously.",
		},
		{
			Name:  "Related commands:",
			Value: "`avatar`, `quote`",
		},
	}
	return embed
}
