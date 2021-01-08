package helpcommands

import (
	"github.com/bwmarrin/discordgo"
)

// AllCaps explains the capitalization
func AllCaps(embed *discordgo.MessageEmbed) *discordgo.MessageEmbed {
	embed.Author.Name = "Command: cap / caps / upper"
	embed.Description = "`(cap|caps|upper) <text>` will create the text given into all caps."
	embed.Fields = []*discordgo.MessageEmbedField{
		{
			Name:   "<text>",
			Value:  "The text to completely upper case.",
			Inline: true,
		},
		{
			Name:  "Related commands:",
			Value: "`lower`, `randomcaps`, `swap`, `title`",
		},
	}
	return embed
}

// AllLower explains the lowercasing of all letters
func AllLower(embed *discordgo.MessageEmbed) *discordgo.MessageEmbed {
	embed.Author.Name = "Command: lower"
	embed.Description = "`lower <text>` will create the text given into lower case."
	embed.Fields = []*discordgo.MessageEmbedField{
		{
			Name:   "<text>",
			Value:  "The text to completely lower case.",
			Inline: true,
		},
		{
			Name:  "Related commands:",
			Value: "`caps`, `randomcaps`, `swap`, `title`",
		},
	}
	return embed
}

// RandomCaps explains the random capitalization
func RandomCaps(embed *discordgo.MessageEmbed) *discordgo.MessageEmbed {
	embed.Author.Name = "Command: rcap / rcaps / rupper / rlower / randomcap / randomcaps / randomupper / randomlower"
	embed.Description = "`(rcap|rcaps|rupper|rlower|randomcap|randomcaps|randomupper|randomlower) <text>` will randomly select characters to capitalize / lower case."
	embed.Fields = []*discordgo.MessageEmbedField{
		{
			Name:   "<text>",
			Value:  "The text to randomly upper / lower case.",
			Inline: true,
		},
		{
			Name:  "Related commands:",
			Value: "`caps`, `lower`, `swap`, `title`",
		},
	}
	return embed
}

// Swap explains the swap formatting
func Swap(embed *discordgo.MessageEmbed) *discordgo.MessageEmbed {
	embed.Author.Name = "Command: swap"
	embed.Description = "`swap <text>` will swap letters in the text given."
	embed.Fields = []*discordgo.MessageEmbedField{
		{
			Name:   "<text>",
			Value:  "The text to swap letters in",
			Inline: true,
		},
		{
			Name:  "Related commands:",
			Value: "`caps`, `lower`, `randomcaps`, `title`",
		},
	}
	return embed
}

// Title explains the title formatting
func Title(embed *discordgo.MessageEmbed) *discordgo.MessageEmbed {
	embed.Author.Name = "Command: title"
	embed.Description = "`title <text>` will create the text into title form."
	embed.Fields = []*discordgo.MessageEmbedField{
		{
			Name:   "<text>",
			Value:  "The text to change into title form.",
			Inline: true,
		},
		{
			Name:  "Related commands:",
			Value: "`caps`, `lower`, `randomcaps`, `swap`",
		},
	}
	return embed
}
