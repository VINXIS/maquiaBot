package helpcommands

import (
	"github.com/bwmarrin/discordgo"
)

// Twitter explains the twitter functionality
func Twitter(embed *discordgo.MessageEmbed) *discordgo.MessageEmbed {
	embed.Author.Name = "Command: twitter / twitterdl"
	embed.Description = "`(twitter|twitterdl)` will download the image / video from the latest tweet posted in the channel."
	return embed
}
