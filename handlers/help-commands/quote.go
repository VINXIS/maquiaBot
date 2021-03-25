package helpcommands

import (
	"github.com/bwmarrin/discordgo"
)

// Quote explains the quote functionality
func Quote(embed *discordgo.MessageEmbed) *discordgo.MessageEmbed {
	embed.Author.Name = "Command: q / quote"
	embed.Description = "`(q|quote) [(@mentions|username) [num]]` gives a quote."
	embed.Fields = []*discordgo.MessageEmbedField{
		{
			Name:   "[username]",
			Value:  "The user to get a quote for. No username will have the bot randomly choose from the users.",
			Inline: true,
		},
		{
			Name:   "[num]",
			Value:  "The user's nth quote to give. No number will result in a random quote to be chosen.",
			Inline: true,
		},
		{
			Name:  "Related Commands:",
			Value: "`quoteadd`, `quoteremove`, `quotes`, `avatarquote`",
		},
	}
	return embed
}

// QuoteAdd explains the quoteadd functionality
func QuoteAdd(embed *discordgo.MessageEmbed) *discordgo.MessageEmbed {
	embed.Author.Name = "Command: qa / qadd / quotea / quoteadd"
	embed.Description = "`(qa|qadd|quotea|quoteadd) ([username] [-r]|[message ID/link] [message ID/link]...|[messageID/link] [num])` adds a quote to the user."
	embed.Fields = []*discordgo.MessageEmbedField{
		{
			Name:   "[username]",
			Value:  "The user / message ID / message link to add a quote for. If you do not give anything, it will add the quote to the latest person who sent a message aside for you.",
			Inline: true,
		},
		{
			Name:   "[-r]",
			Value:  "Randomly choose one of the past 100 messages instead of the latest.",
			Inline: true,
		},
		{
			Name:   "[message ID/link] [message ID/link]...",
			Value:  "Allows you to merge the message text of 2 or more messages into 1 quote.",
			Inline: true,
		},
		{
			Name:   "[messageID/link] [num]",
			Value:  "Allows you to merge the message text of the message ID given alongside the num amount of messages after, or until another person's message appears after. No num will simply quote the given message ID/link",
			Inline: true,
		},
		{
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
		{
			Name:  "[messageID]",
			Value: "The message ID to remove.",
		},
		{
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
		{
			Name:  "[username]",
			Value: "The user to list quotes from. If you do not give anything, it will list quotes from all users.",
		},
		{
			Name:  "Related Commands:",
			Value: "`quote`, `quoteadd`, `quoteremove`",
		},
	}
	return embed
}
