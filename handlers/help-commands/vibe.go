package helpcommands

import (
	"github.com/bwmarrin/discordgo"
)

// Vibe explains the vibe functionality
func Vibe(embed *discordgo.MessageEmbed) *discordgo.MessageEmbed {
	embed.Author.Name = "Command: vibe / vibec / vibecheck"
	embed.Description = "`(vibe|vibec|vibecheck) [@mention]` runs a vibe check on a user."
	embed.Fields = []*discordgo.MessageEmbedField{
		&discordgo.MessageEmbedField{
			Name:  "[@mention]",
			Value: "The user to run the vibe check on. No mention will run the vibe check on the person who messaged previously.",
		},
		&discordgo.MessageEmbedField{
			Name:  "Related Commands:",
			Value: "`vibetoggle`",
		},
	}
	return embed
}

// VibeToggle explains the vibe toggle functionality
func VibeToggle(embed *discordgo.MessageEmbed) *discordgo.MessageEmbed {
	embed.Author.Name = "Command: vibet / vibetoggle"
	embed.Description = "`(vibet|vibetoggle)` lets admins toggle as to whether a vibe check should randomly run or not (1/100000 chance)."
	embed.Fields = []*discordgo.MessageEmbedField{
		&discordgo.MessageEmbedField{
			Name:  "Related Commands:",
			Value: "`vibe`",
		},
	}
	return embed
}
