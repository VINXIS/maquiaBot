package helpcommands

import (
	"github.com/bwmarrin/discordgo"
)

// List explains the list functionality
func List(embed *discordgo.MessageEmbed) *discordgo.MessageEmbed {
	embed.Author.Name = "Command: list"
	embed.Description = "`list <NEW LINE> <option> [<NEW LINE> [option] ...]` randomizes a list separated by lines."
	embed.Fields = []*discordgo.MessageEmbedField{
		{
			Name:   "<NEW LINE> <option>",
			Value:  "A new line (shift-enter) followed by the text for 1 of the options.",
			Inline: true,
		},
		{
			Name:   "Example usage:",
			Value:  "$list\na\nb\nc\nd",
			Inline: true,
		},
	}
	return embed
}
