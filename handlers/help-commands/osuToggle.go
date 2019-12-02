package helpcommands

import (
	"github.com/bwmarrin/discordgo"
)

// OsuToggle explains the osu! toggle functionality
func OsuToggle(embed *discordgo.MessageEmbed) *discordgo.MessageEmbed {
	embed.Author.Name = "Command: ot, osut, otoggle, osutoggle"
	embed.Description = "`(ot|osut|otoggle|osutoggle)` lets admins toggle if beatmap / profile info and timestamps should be sent by this bot or not."
	return embed
}
