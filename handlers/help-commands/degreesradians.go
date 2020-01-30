package helpcommands

import (
	"github.com/bwmarrin/discordgo"
)

// DegreesRadians explains the degrees to radians functionality
func DegreesRadians(embed *discordgo.MessageEmbed) *discordgo.MessageEmbed {
	embed.Author.Name = "Command: dr / degrad / degreesradians"
	embed.Description = "`[math] (dr|degrad|degreesradians) <number>` will conver the number from degrees to radians."
	embed.Fields = []*discordgo.MessageEmbedField{
		&discordgo.MessageEmbedField{
			Name:   "<number>",
			Value:  "The number to convert.",
			Inline: true,
		},
		&discordgo.MessageEmbedField{
			Name:  "Related commands:",
			Value: "`radiansdegrees`",
		},
	}
	return embed
}

// RadiansDegrees explains the radians to degrees functionality
func RadiansDegrees(embed *discordgo.MessageEmbed) *discordgo.MessageEmbed {
	embed.Author.Name = "Command: rd / raddeg / radiansdegrees"
	embed.Description = "`[math] (rd|raddeg|radiansdegrees) <number>` will conver the number from radians to degrees."
	embed.Fields = []*discordgo.MessageEmbedField{
		&discordgo.MessageEmbedField{
			Name:   "<number>",
			Value:  "The number to convert.",
			Inline: true,
		},
		&discordgo.MessageEmbedField{
			Name:  "Related commands:",
			Value: "`degreesradians`",
		},
	}
	return embed
}