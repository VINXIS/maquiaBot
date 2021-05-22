package helpcommands

import (
	"github.com/bwmarrin/discordgo"
)

// Task explains the task functionality
func Task(embed *discordgo.MessageEmbed) *discordgo.MessageEmbed {
	embed.Author.Name = "Command: task"
	embed.Description = "`task [text] [in] [every] <time duration> [<starting> [at] <datetime>])` sends messages every <time duration> amount of time about [text] starting at <datetime> if specified."
	embed.Fields = []*discordgo.MessageEmbedField{
		{
			Name:   "[text]",
			Value:  "The text about the task. Not required.",
			Inline: true,
		},
		{
			Name:   "[in] [every] <time duration>",
			Value:  "The duration until you want to be messaged about the task.",
			Inline: true,
		},
		{
			Name:   "[<starting> [at] <datetime>]",
			Value:  "The specific date and/or time the recurring messages should start.",
			Inline: true,
		},
		{
			Name:   "Example format (without a starting point):",
			Value:  "`$task osu! in 1 week` Will message you about `osu!` every week starting from now.",
			Inline: true,
		},
		{
			Name:   "Example format (with a starting point):",
			Value:  "`$task osu! in 1 week starting at January 1 2021 10:00 PM` Will message you about `osu!` every week starting from January 1 2021 10:00 PM. That means the next message would be sent at January 8 2021 10:00 PM.",
			Inline: true,
		},
		{
			Name:  "Related Commands:",
			Value: "`tasks`, `taskremove`",
		},
	}
	return embed
}

// Tasks explains the tasks functionality
func Tasks(embed *discordgo.MessageEmbed) *discordgo.MessageEmbed {
	embed.Author.Name = "Command: tasks"
	embed.Description = "`tasks` shows you all of your currently running tasks."
	embed.Fields = []*discordgo.MessageEmbedField{
		{
			Name:  "Related Commands:",
			Value: "`task`, `taskremove`",
		},
	}
	return embed
}

// TaskRemove explains the taskremove functionality
func TaskRemove(embed *discordgo.MessageEmbed) *discordgo.MessageEmbed {
	embed.Author.Name = "Command: taskremove / tremove"
	embed.Description = "`(taskremove|tremove) <task id|all>` removes a task / all of your tasks."
	embed.Fields = []*discordgo.MessageEmbedField{
		{
			Name:  "<task id|all>",
			Value: "The ID of the task you want to remove which is obtainable from `tasks`. You can also state `all` instead of an ID to remove all of your currently running tasks.",
		},
		{
			Name:  "Related Commands:",
			Value: "`task`, `tasks`",
		},
	}
	return embed
}
