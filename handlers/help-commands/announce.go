package helpcommands

import (
	"github.com/bwmarrin/discordgo"
)

// AnnounceToggle explains the announce toggle functionality
func AnnounceToggle(embed *discordgo.MessageEmbed) *discordgo.MessageEmbed {
	embed.Author.Name = "Command: at / announcet / atoggle / announcetoggle"
	embed.Description = "`(at|announcet|atoggle|announcetoggle)` lets admins toggle whether announcements from the bot creator should be sent to their server."
	return embed
}
