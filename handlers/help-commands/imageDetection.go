package helpcommands

import (
	"github.com/bwmarrin/discordgo"
)

// OCR explains the OCR functionality
func OCR(embed *discordgo.MessageEmbed) *discordgo.MessageEmbed {
	embed.Author.Name = "Command: ocr"
	embed.Description = "`ocr [link] [-l lang]` detects for text in an image."
	embed.Fields = []*discordgo.MessageEmbedField{
		{
			Name:   "[link]",
			Value:  "Gets the image from this link, otherwise gets the latest posted image if a link isn't given.",
			Inline: true,
		},
		{
			Name:   "[-l lang]",
			Value:  "Give the language to detect using the `-l` flag. Language codes are 3 letters, and most can be found here: https://en.wikipedia.org/wiki/ISO_639-2",
			Inline: true,
		},
		{
			Name:   "What is OCR?",
			Value:  "Short for Optical Character Recognition, it converts text from an image into text that can be manipulated by people.\nhttps://en.wikipedia.org/wiki/Optical_character_recognition",
			Inline: true,
		},
		{
			Name:  "Related Commands:",
			Value: "`face`",
		},
	}
	return embed
}

// Face explains the face functionality
func Face(embed *discordgo.MessageEmbed) *discordgo.MessageEmbed {
	embed.Author.Name = "Command: face"
	embed.Description = "`face [link]` detects for faces in an image."
	embed.Fields = []*discordgo.MessageEmbedField{
		{
			Name:  "[link]",
			Value: "Gets the image from this link, otherwise gets the latest posted image if a link isn't given.",
		},
		{
			Name:  "Related Commands:",
			Value: "`ocr`",
		},
	}
	return embed
}
