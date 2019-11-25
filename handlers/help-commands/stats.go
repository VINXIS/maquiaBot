package helpcommands

import (
	"github.com/bwmarrin/discordgo"
)

// Stats explains the stats functionality
func Stats(embed *discordgo.MessageEmbed) *discordgo.MessageEmbed {
	embed.Author.Name = "Command: stats"
	embed.Description = "`stats [text] [num]` lets admins change the prefix for the bot in the server."
	embed.Fields = []*discordgo.MessageEmbedField{
		&discordgo.MessageEmbedField{
			Name:   "[text]",
			Value:  "The text to get the stats for. No text will give your own stats instead.",
			Inline: true,
		},
		&discordgo.MessageEmbedField{
			Name:   "[num]",
			Value:  "The number of skills to print (Default is 4)",
			Inline: true,
		},
		&discordgo.MessageEmbedField{
			Name:  "Related Commands",
			Value: "`adjective`, `nouns`, `skills`, `statstoggle`",
		},
	}
	return embed
}

// Adjectives explains the adjective functionality
func Adjectives(embed *discordgo.MessageEmbed) *discordgo.MessageEmbed {
	embed.Author.Name = "Command: adj / adjective / adjectives"
	embed.Description = "`(adj|adjective|adjectives) [remove] <adj>` lets you add/remove adjectives from the stats feature."
	embed.Fields = []*discordgo.MessageEmbedField{
		&discordgo.MessageEmbedField{
			Name:   "[remove]",
			Value:  "State remove to remove the adj.",
			Inline: true,
		},
		&discordgo.MessageEmbedField{
			Name:   "<adj>",
			Value:  "The word to add/remove",
			Inline: true,
		},
		&discordgo.MessageEmbedField{
			Name:  "Related Commands",
			Value: "`stats`, `nouns`, `skills`, `statstoggle`",
		},
	}
	return embed
}

// Nouns explains the noun functionality
func Nouns(embed *discordgo.MessageEmbed) *discordgo.MessageEmbed {
	embed.Author.Name = "Command: noun / nouns"
	embed.Description = "`(noun|nouns) [remove] <noun>` lets you add/remove nouns from the stats feature."
	embed.Fields = []*discordgo.MessageEmbedField{
		&discordgo.MessageEmbedField{
			Name:   "[remove]",
			Value:  "State remove to remove the noun.",
			Inline: true,
		},
		&discordgo.MessageEmbedField{
			Name:   "<noun>",
			Value:  "The word to add/remove",
			Inline: true,
		},
		&discordgo.MessageEmbedField{
			Name:  "Related Commands",
			Value: "`stats`, `adjective`, `skills`, `statstoggle`",
		},
	}
	return embed
}

// Skills explains the skills functionality
func Skills(embed *discordgo.MessageEmbed) *discordgo.MessageEmbed {
	embed.Author.Name = "Command: skill / skills"
	embed.Description = "`(skill|skills) [remove] <skill>` lets you add/remove skills from the stats feature."
	embed.Fields = []*discordgo.MessageEmbedField{
		&discordgo.MessageEmbedField{
			Name:   "[remove]",
			Value:  "State remove to remove the skill.",
			Inline: true,
		},
		&discordgo.MessageEmbedField{
			Name:   "<skill>",
			Value:  "The word to add/remove",
			Inline: true,
		},
		&discordgo.MessageEmbedField{
			Name:  "Related Commands",
			Value: "`stats`, `adjective`, `nouns`, `statstoggle`",
		},
	}
	return embed
}

// StatsToggle explains the statstoggle functionality
func StatsToggle(embed *discordgo.MessageEmbed) *discordgo.MessageEmbed {
	embed.Author.Name = "Command: statst / statstoggle"
	embed.Description = "`(statst|statstoggle)` lets admins toggle if anyone can add adj/nouns/skills, or if only admins are allowed."
	embed.Fields = []*discordgo.MessageEmbedField{
		&discordgo.MessageEmbedField{
			Name:  "Related Commands",
			Value: "`stats`, `adjective`, `nouns`, `skills`",
		},
	}
	return embed
}
