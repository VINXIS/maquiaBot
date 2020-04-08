package helpcommands

import (
	"github.com/bwmarrin/discordgo"
)

// Colour explains the colour functionality
func Colour(embed *discordgo.MessageEmbed) *discordgo.MessageEmbed {
	embed.Author.Name = "Command: col / color / colour"
	embed.Description = "`(col|color|colour) <vals> -<hex|hsla|cmyk|ycbcr>` has the bot generate an image of the colour given. See below for different options. No colour being given will generate an image of a random colour."
	embed.Fields = []*discordgo.MessageEmbedField{
		&discordgo.MessageEmbedField{
			Name:   "<r> <g> <b> [a]",
			Value:  "To use rgb values, provide the rgb values, the alpha value is optional. No hyphen tag is required, but you may add `-rgb` or `-rgba` to the end if you want.",
			Inline: true,
		},
		&discordgo.MessageEmbedField{
			Name:   "[#]<hex value> -hex",
			Value:  "To use hex, provide the proper hex code (3, 6, or 8 values long), and add `-hex` to the end.",
			Inline: true,
		},
		&discordgo.MessageEmbedField{
			Name:   "<h> <s> <l> [a] -hsla",
			Value:  "To use hsl values, provide the hsl values, the alpha value is optional. Add `-hsla` to the end.",
			Inline: true,
		},
		&discordgo.MessageEmbedField{
			Name:   "<c> <m> <y> <k> -cmyk",
			Value:  "To use cmyk values, provide the cmyk values, and add `-cmyk` to the end.",
			Inline: true,
		},
		&discordgo.MessageEmbedField{
			Name:   "<y> <cb> <cr> -ycbcr",
			Value:  "To use ycbcr, provide the ycbcr values, and add `-ycbcr` to the end.",
			Inline: true,
		},
	}
	return embed
}
