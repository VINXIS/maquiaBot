package helpcommands

import (
	"github.com/bwmarrin/discordgo"
)

// TTS explains the tts functionality
func TTS(embed *discordgo.MessageEmbed) *discordgo.MessageEmbed {
	embed.Author.Name = "Command: tts"
	embed.Description = "`tts <text> [-v <voice>]` will create a tts audio file from https://15.ai/."
	embed.Fields = []*discordgo.MessageEmbedField{
		{
			Name:   "<text>",
			Value:  "The text to turn to speech.",
			Inline: true,
		},
		{
			Name:  "[-v <voice]",
			Value: "The voice to use. The list of voices allowed can be seen on the website https://15.ai/. This is case-sensitive so make sure you use the same exact case form used on the website.",
		},
	}
	return embed
}
