package helpcommands

import (
	"github.com/bwmarrin/discordgo"
)

// Weather explains the weather functionality
func Weather(embed *discordgo.MessageEmbed) *discordgo.MessageEmbed {
	embed.Author.Name = "Command: w / weather"
	embed.Description = "`(w|weather) <location> [-l lang]` provides current weather information for a location."
	embed.Fields = []*discordgo.MessageEmbedField{
		{
			Name: "<location>",
			Value: "A location. This can be the following:\n" +
				"Latitude and Longitude (e.g 48.8567,2.3508)\n" +
				"City name (e.g New York City)\n" +
				"US Zip Code (e.g 10001)\n" +
				"Postal Code (e.g G2J 2X1)\n" +
				"Metar Code (e.g EGLL)\n" +
				"Airport Code (e.g DXB)\n" +
				"IP Address (e.g 127.0.0.1)",
			Inline: true,
		},
		{
			Name:   "[-l lang]",
			Value:  "Give information in a different language if you want. Although it only changes the \"sunny, partly cloudy, e.t.c\" text.",
			Inline: true,
		},
	}
	return embed
}
