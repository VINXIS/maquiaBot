package helpcommands

import (
	"github.com/bwmarrin/discordgo"
)

// Twitter explains the twitter functionality
func Twitter(embed *discordgo.MessageEmbed) *discordgo.MessageEmbed {
	embed.Author.Name = "Command: twitter / twitterdl"
	embed.Description = "`(twitter|twitterdl) [link]` will download the image / video from the latest tweet posted in the channel."
	embed.Fields = []*discordgo.MessageEmbedField{
		&discordgo.MessageEmbedField{
			Name:  "[message]",
			Value: "A video / gif tweet. If you want to download an image, then link the tweet first, and then do the `twitter` command.",
		},
	}
	return embed
}
