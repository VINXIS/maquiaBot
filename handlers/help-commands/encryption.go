package helpcommands

import (
	"github.com/bwmarrin/discordgo"
)

// Encrypt explains the encrypt functionality
func Encrypt(embed *discordgo.MessageEmbed) *discordgo.MessageEmbed {
	embed.Author.Name = "Command: encrypt"
	embed.Description = "`encrypt <text> [-k <key>]` lets you encrypt some text with AES-GCM encryption."
	embed.Fields = []*discordgo.MessageEmbedField{
		&discordgo.MessageEmbedField{
			Name:   "<text>",
			Value:  "The text you want to encrypt",
			Inline: true,
		},
		&discordgo.MessageEmbedField{
			Name:   "[-k <key>]",
			Value:  "The key to use to encrypt (Default: use `key` to see the default key).",
			Inline: true,
		},
		&discordgo.MessageEmbedField{
			Name:  "Related Commands:",
			Value: "`decrypt`, `key` (no help command for key)",
		},
	}
	return embed
}

// Decrypt explains the decrypt functionality
func Decrypt(embed *discordgo.MessageEmbed) *discordgo.MessageEmbed {
	embed.Author.Name = "Command: decrypt"
	embed.Description = "`decrypt <text> [-k <key>]` lets you decrypt some text with AES-GCM encryption."
	embed.Fields = []*discordgo.MessageEmbedField{
		&discordgo.MessageEmbedField{
			Name:   "<text>",
			Value:  "The text you want to decrypt",
			Inline: true,
		},
		&discordgo.MessageEmbedField{
			Name:   "[-k <key>]",
			Value:  "The key to use to decrypt (Default: use `key` to see the default key).",
			Inline: true,
		},
		&discordgo.MessageEmbedField{
			Name:  "Related Commands:",
			Value: "`encrypt`, `key` (no help command for key)",
		},
	}
	return embed
}
