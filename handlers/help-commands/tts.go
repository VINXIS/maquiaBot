package helpcommands

import (
	"github.com/bwmarrin/discordgo"
)

// TTS explains the tts functionality
func TTS(embed *discordgo.MessageEmbed) *discordgo.MessageEmbed {
	embed.Author.Name = "Command: tts"
	embed.Description = "`tts [text] [-v <voice>]` will create a tts audio file from https://15.ai/."
	embed.Fields = []*discordgo.MessageEmbedField{
		{
			Name:   "[text]",
			Value:  "The text to turn to speech. No text given will use the latest message that is not from me.\nIf you use ARPAbet, make sure it is valid. More information on text functionality is on the website, and all of them will work through the bot https://15.ai/\n **MAX 300 CHARACTERS**",
			Inline: true,
		},
		{
			Name:  "[-v <voice]",
			Value: "The voice to use. The default used voice is The Stanley Parable's The Narrator.\nThe list of voices allowed can be seen on the website https://15.ai/.",
		},
	}
	return embed
}
