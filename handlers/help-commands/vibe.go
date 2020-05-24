package helpcommands

import (
	"github.com/bwmarrin/discordgo"
)

// Vibe explains the vibe functionality
func Vibe(embed *discordgo.MessageEmbed) *discordgo.MessageEmbed {
	embed.Author.Name = "Command: vibe / vibec / vibecheck"
	embed.Description = "`(vibe|vibec|vibecheck) [@mention]` runs a vibe check on a user."
	embed.Fields = []*discordgo.MessageEmbedField{
		{
			Name:  "[@mention]",
			Value: "The user to run the vibe check on. No mention will run the vibe check on the person who messaged previously.",
		},
		{
			Name:  "Related Commands:",
			Value: "`toggle`",
		},
	}
	return embed
}
