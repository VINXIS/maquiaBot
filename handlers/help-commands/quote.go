package helpcommands

import (
	"github.com/bwmarrin/discordgo"
)

// Quote explains the quote functionality
func Quote(embed *discordgo.MessageEmbed) *discordgo.MessageEmbed {
	embed.Author.Name = "Command: q / quote"
	embed.Description = "`(q|quote) [<username> [num]]` gives a quote."
	embed.Fields = []*discordgo.MessageEmbedField{
		&discordgo.MessageEmbedField{
			Name:   "[username]",
			Value:  "The user to get a quote for. No username will have the bot randomly choose from the users.",
			Inline: true,
		},
		&discordgo.MessageEmbedField{
			Name:   "[num]",
			Value:  "The user's nth quote to give. No number will result in a random quote to be chosen.",
			Inline: true,
		},
		&discordgo.MessageEmbedField{
			Name:  "Related Commands:",
			Value: "`quoteadd`, `quoteremove`, `quotes`",
		},
	}
	return embed
}

// QuoteAdd explains the quoteadd functionality
func QuoteAdd(embed *discordgo.MessageEmbed) *discordgo.MessageEmbed {
	embed.Author.Name = "Command: qa / qadd / quotea / quoteadd"
	embed.Description = "`(qa|qadd|quotea|quoteadd) [username] [-r]` adds a quote to the user."
	embed.Fields = []*discordgo.MessageEmbedField{
		&discordgo.MessageEmbedField{
			Name:   "[username]",
			Value:  "The user / message ID / message link to add a quote for. If you do not give anything, it will add the quote to the latest person who sent a message aside for you.",
			Inline: true,
		},
		&discordgo.MessageEmbedField{
			Name:   "[-r]",
			Value:  "Randomly choose one of the messages instead of the latest.",
			Inline: true,
		},
		&discordgo.MessageEmbedField{
			Name:  "Related Commands:",
			Value: "`quote`, `quoteremove`, `quotes`",
		},
	}
	return embed
}

// QuoteRemove explains the quoteremove functionality
func QuoteRemove(embed *discordgo.MessageEmbed) *discordgo.MessageEmbed {
	embed.Author.Name = "Command: qd / qr / qremove / qdelete / quoteremove / quotedelete"
	embed.Description = "`(qd|qr|qremove|qdelete|quoteremove|quotedelete) <messageID>` removes a quote."
	embed.Fields = []*discordgo.MessageEmbedField{
		&discordgo.MessageEmbedField{
			Name:  "[messageID]",
			Value: "The message ID to remove.",
		},
		&discordgo.MessageEmbedField{
			Name:  "Related Commands:",
			Value: "`quote`, `quoteadd`, `quotes`",
		},
	}
	return embed
}

// Quotes explains the quotes functionality
func Quotes(embed *discordgo.MessageEmbed) *discordgo.MessageEmbed {
	embed.Author.Name = "Command: qs / quotes"
	embed.Description = "`(qs|quotes) [username]` adds a quote to the user."
	embed.Fields = []*discordgo.MessageEmbedField{
		&discordgo.MessageEmbedField{
			Name:  "[username]",
			Value: "The user to list quotes from. If you do not give anything, it will list quotes from all users.",
		},
		&discordgo.MessageEmbedField{
			Name:  "Related Commands:",
			Value: "`quote`, `quoteadd`, `quoteremove`",
		},
	}
	return embed
}
