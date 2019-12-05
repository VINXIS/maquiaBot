package helpcommands

import (
	"github.com/bwmarrin/discordgo"
)

// BPM explains the BPM functionality
func BPM(embed *discordgo.MessageEmbed) *discordgo.MessageEmbed {
	embed.Author.Name = "Command: bpm"
	embed.Description = "`[osu] bpm [username]` calculates your BPM for today."
	embed.Fields = []*discordgo.MessageEmbedField{
		&discordgo.MessageEmbedField{
			Name:  "[username]",
			Value: "Gets the BPM for the given osu! username. Gives your BPM if no osu! username is given.",
		},
	}
	return embed
}
