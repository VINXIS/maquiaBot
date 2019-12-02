package helpcommands

import (
	"github.com/bwmarrin/discordgo"
)

// Penis explains the penis functionality
func Penis(embed *discordgo.MessageEmbed) *discordgo.MessageEmbed {
	embed.Author.Name = "Command: penis / cock"
	embed.Description = "`(penis|cock) [username]` calculates your erect length for today."
	embed.Fields = []*discordgo.MessageEmbedField{
		&discordgo.MessageEmbedField{
			Name:  "[username]",
			Value: "Gets the erect length for the given username / nickname / ID. Gives your erect length if no username / nickname / ID is given.",
		},
		&discordgo.MessageEmbedField{
			Name:  "Related Commands:",
			Value: "`comparepenis`",
		},
	}
	return embed
}

// PenisCompare explains the penis compare functionality
func PenisCompare(embed *discordgo.MessageEmbed) *discordgo.MessageEmbed {
	embed.Author.Name = "Command: cc / cp / comparec / comparep / comparecock / comparepenis"
	embed.Description = "`(cc|cp|comparec|comparep|comparecock|comparepenis) [username]` compares your erect length with someone else's."
	embed.Fields = []*discordgo.MessageEmbedField{
		&discordgo.MessageEmbedField{
			Name:  "[username]",
			Value: "Compares your erect length with the given username / nickname / ID. Compares your erect length with the last user to do a penis / penis compare command if no username / nickname / ID is given.",
		},
		&discordgo.MessageEmbedField{
			Name:  "Related Commands:",
			Value: "`penis`",
		},
	}
	return embed
}
