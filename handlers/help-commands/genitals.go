package helpcommands

import (
	"github.com/bwmarrin/discordgo"
)

// Penis explains the penis functionality
func Penis(embed *discordgo.MessageEmbed) *discordgo.MessageEmbed {
	embed.Author.Name = "Command: penis"
	embed.Description = "`penis [username]` calculates your erect length for today."
	embed.Fields = []*discordgo.MessageEmbedField{
		{
			Name:  "[username]",
			Value: "Gets the erect length for the given username / nickname / ID. Gives your erect length if no username / nickname / ID is given.",
		},
		{
			Name:  "Related Commands:",
			Value: "`vagina`, `comparepenis`, `comparevagina`, `rankpenis`, `rankvagina`, `history`",
		},
	}
	return embed
}

// Vagina explains the vagina functionality
func Vagina(embed *discordgo.MessageEmbed) *discordgo.MessageEmbed {
	embed.Author.Name = "Command: vagina"
	embed.Description = "`vagina [username]` calculates your length for today."
	embed.Fields = []*discordgo.MessageEmbedField{
		{
			Name:  "[username]",
			Value: "Gets the erect length for the given username / nickname / ID. Gives your length if no username / nickname / ID is given.",
		},
		{
			Name:  "Related Commands:",
			Value: "`penis`, `comparepenis`, `comparevagina`, `rankpenis`, `rankvagina`, `history`",
		},
	}
	return embed
}

// PenisCompare explains the penis compare functionality
func PenisCompare(embed *discordgo.MessageEmbed) *discordgo.MessageEmbed {
	embed.Author.Name = "Command: cp / comparep / comparepenis"
	embed.Description = "`(cp|comparep|comparepenis) [username]` compares your erect length with someone else's."
	embed.Fields = []*discordgo.MessageEmbedField{
		{
			Name:  "[username]",
			Value: "Compares your erect length with the given username / nickname / ID. Compares your erect length with the last user to do a penis / penis compare command if no username / nickname / ID is given.",
		},
		{
			Name:  "Related Commands:",
			Value: "`penis`, `vagina`, `comparevagina`, `rankpenis`, `rankvagina`, `history`",
		},
	}
	return embed
}

// VaginaCompare explains the vagina compare functionality
func VaginaCompare(embed *discordgo.MessageEmbed) *discordgo.MessageEmbed {
	embed.Author.Name = "Command: cv / comparev / comparevagina"
	embed.Description = "`(cv|comparev|comparevagina) [username]` compares your length with someone else's."
	embed.Fields = []*discordgo.MessageEmbedField{
		{
			Name:  "[username]",
			Value: "Compares your length with the given username / nickname / ID. Compares your length with the last user to do a vagina / vagina compare command if no username / nickname / ID is given.",
		},
		{
			Name:  "Related Commands:",
			Value: "`penis`, `vagina`, `comparepenis`, `rankpenis`, `rankvagina`, `history`",
		},
	}
	return embed
}

// PenisRank explains the penis rank functionality
func PenisRank(embed *discordgo.MessageEmbed) *discordgo.MessageEmbed {
	embed.Author.Name = "Command: rp / rankp / rankpenis"
	embed.Description = "`(rp|rankp|rankpenis) [number] [-s]` ranks the largest / smallest penis sizes in the server."
	embed.Fields = []*discordgo.MessageEmbedField{
		{
			Name:   "[number]",
			Value:  "Display a certain number of largest / smallest penises (Default: 1).",
			Inline: true,
		},
		{
			Name:   "[-s]",
			Value:  "Add this to show the smallest sizes. If this is not added, it will show the largest sizes instead.",
			Inline: true,
		},
		{
			Name:  "Related Commands:",
			Value: "`penis`, `vagina`, `comparepenis`, `comparevagina`, `rankvagina`, `history`",
		},
	}
	return embed
}

// VaginaRank explains the vagina rank functionality
func VaginaRank(embed *discordgo.MessageEmbed) *discordgo.MessageEmbed {
	embed.Author.Name = "Command: rv / rankv / rankvagina"
	embed.Description = "`(rv|rankv|rankvagina) [number] [-s]` ranks the largest / smallest vagina sizes in the server."
	embed.Fields = []*discordgo.MessageEmbedField{
		{
			Name:   "[number]",
			Value:  "Display a certain number of largest / smallest vagina sizes (Default: 1).",
			Inline: true,
		},
		{
			Name:   "[-s]",
			Value:  "Add this to show the smallest sizes. If this is not added, it will show the largest sizes instead.",
			Inline: true,
		},
		{
			Name:  "Related Commands:",
			Value: "`penis`, `vagina`, `comparepenis`, `comparevagina`, `rankpenis`, `history`",
		},
	}
	return embed
}

// History explains the penis history functionality
func History(embed *discordgo.MessageEmbed) *discordgo.MessageEmbed {
	embed.Author.Name = "Command: history"
	embed.Description = "`history [-s]` shows the largest and smallest sizes ever recorded."
	embed.Fields = []*discordgo.MessageEmbedField{
		{
			Name:   "[-s]",
			Value:  "Show the largest and smallest in the server instead of all servers.",
			Inline: true,
		},
		{
			Name:  "Related Commands:",
			Value: "`penis`, `vagina`, `comparepenis`, `comparevagina`, `rankpenis`, `rankvagina`",
		},
	}
	return embed
}
