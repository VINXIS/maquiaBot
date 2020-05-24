package helpcommands

import (
	"github.com/bwmarrin/discordgo"
)

// Remind explains the remind functionality
func Remind(embed *discordgo.MessageEmbed) *discordgo.MessageEmbed {
	embed.Author.Name = "Command: remind / reminder"
	embed.Description = "`(remind|reminder) [text] [in time]` reminds you in some amount of time."
	embed.Fields = []*discordgo.MessageEmbedField{
		{
			Name:   "[text]",
			Value:  "The text to remind you about. Not required.",
			Inline: true,
		},
		{
			Name:   "[in time]",
			Value:  "The time until you want to be reminded.",
			Inline: true,
		},
		{
			Name:   "Example format:",
			Value:  "`$remind play osu! in 5 hours` Will remind you about `play osu!` in 5 hours.",
			Inline: true,
		},
		{
			Name:  "Related Commands:",
			Value: "`reminders`, `remindremove`",
		},
	}
	return embed
}

// Reminders explains the reminders functionality
func Reminders(embed *discordgo.MessageEmbed) *discordgo.MessageEmbed {
	embed.Author.Name = "Command: reminders"
	embed.Description = "`reminders` shows you all of your currently running reminders."
	embed.Fields = []*discordgo.MessageEmbedField{
		{
			Name:  "Related Commands:",
			Value: "`remind`, `remindremove`",
		},
	}
	return embed
}

// RemindRemove explains the remindremove functionality
func RemindRemove(embed *discordgo.MessageEmbed) *discordgo.MessageEmbed {
	embed.Author.Name = "Command: remindremove / rremove"
	embed.Description = "`(remindremove|rremove) <reminder id|all>` removes a reminder / all of your reminders."
	embed.Fields = []*discordgo.MessageEmbedField{
		{
			Name:  "<reminder id|all>",
			Value: "The ID of the reminder you want to remove which is obtainable from `reminders`. You can also state `all` instead of an ID to remove all of your currently running reminders.",
		},
		{
			Name:  "Related Commands:",
			Value: "`remind`, `reminders`",
		},
	}
	return embed
}
