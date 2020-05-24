package helpcommands

import (
	"github.com/bwmarrin/discordgo"
)

// Twitch explains the twitch functionality
func Twitch(embed *discordgo.MessageEmbed) *discordgo.MessageEmbed {
	embed.Author.Name = "Command: twitch / twitchdl"
	embed.Description = "`(twitch|twitchdl) [link]` will download a twitch clip from the latest posted clip linked in the channel."
	embed.Fields = []*discordgo.MessageEmbedField{
		{
			Name:  "[link]",
			Value: "A twitch clip link. No link will look for the latest posted clip instead.",
		},
	}
	return embed
}
