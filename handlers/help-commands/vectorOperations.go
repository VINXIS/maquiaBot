package helpcommands

import (
	"github.com/bwmarrin/discordgo"
)

// DistanceDirection explains the distance / direction functionality
func DistanceDirection(embed *discordgo.MessageEmbed) *discordgo.MessageEmbed {
	embed.Author.Name = "Command: d / dist / distance / dir / direction"
	embed.Description = "`[math] (d|dist|distance|ir|direction) <vector1> <vector2>` will gives the distance between the 2 vectors and the direction of the distance vector."
	embed.Fields = []*discordgo.MessageEmbedField{
		&discordgo.MessageEmbedField{
			Name:   "<vector1> <vector2>",
			Value:  "The vectors to calculate the distance vector for.",
			Inline: true,
		},
		&discordgo.MessageEmbedField{
			Name:   "Vector Notation:",
			Value:  "2 Dimensions: (x, y)\n 3 Dimensions: (x, y, z)",
			Inline: true,
		},
	}
	return embed
}

// VectorAdd explains the vector addition functionality
func VectorAdd(embed *discordgo.MessageEmbed) *discordgo.MessageEmbed {
	embed.Author.Name = "Command: va / vadd / vectora / vectoradd"
	embed.Description = "`[math] (va|vadd|vectora|vectoradd) <vector1> <vector2>` will add the 2 vectors together."
	embed.Fields = []*discordgo.MessageEmbedField{
		&discordgo.MessageEmbedField{
			Name:   "<vector1> <vector2>",
			Value:  "The vectors to add together for.",
			Inline: true,
		},
		&discordgo.MessageEmbedField{
			Name:   "Vector Notation:",
			Value:  "2 Dimensions: (x, y)\n 3 Dimensions: (x, y, z)",
			Inline: true,
		},
		&discordgo.MessageEmbedField{
			Name:  "Related commands:",
			Value: "`vectorcross`, `vectordivide`, `vectordot`, `vectormultiply`, `vectorsubtract`",
		},
	}
	return embed
}

// VectorCross explains the vector cross product functionality
func VectorCross(embed *discordgo.MessageEmbed) *discordgo.MessageEmbed {
	embed.Author.Name = "Command: vc / vcross / vectorc / vectorcross"
	embed.Description = "`[math] (vc|vcross|vectorc|vectorcross) <vector1> <vector2>` will give the cross product of 2 vectors."
	embed.Fields = []*discordgo.MessageEmbedField{
		&discordgo.MessageEmbedField{
			Name:   "<vector1> <vector2>",
			Value:  "The vectors to get the cross product for.",
			Inline: true,
		},
		&discordgo.MessageEmbedField{
			Name:   "Vector Notation:",
			Value:  "2 Dimensions: (x, y)\n 3 Dimensions: (x, y, z)",
			Inline: true,
		},
		&discordgo.MessageEmbedField{
			Name:  "Related commands:",
			Value: "`vectoradd`, `vectordivide`, `vectordot`, `vectormultiply`, `vectorsubtract`",
		},
	}
	return embed
}

// VectorDivide explains the vector division functionality
func VectorDivide(embed *discordgo.MessageEmbed) *discordgo.MessageEmbed {
	embed.Author.Name = "Command: va / vadd / vectora / vectoradd"
	embed.Description = "`[math] (va|vadd|vectora|vectoradd) <vector1> <number>` will divide the vector with the number."
	embed.Fields = []*discordgo.MessageEmbedField{
		&discordgo.MessageEmbedField{
			Name:   "<vector1> <vector2>",
			Value:  "The vector to divide.",
			Inline: true,
		},
		&discordgo.MessageEmbedField{
			Name:   "<number>",
			Value:  "The number to divide the vector with.",
			Inline: true,
		},
		&discordgo.MessageEmbedField{
			Name:   "Vector Notation:",
			Value:  "2 Dimensions: (x, y)\n 3 Dimensions: (x, y, z)",
			Inline: true,
		},
		&discordgo.MessageEmbedField{
			Name:  "Related commands:",
			Value: "`vectoradd`, `vectorcross`, `vectordot`, `vectormultiply`, `vectorsubtract`",
		},
	}
	return embed
}

// VectorDot explains the vector dot product functionality
func VectorDot(embed *discordgo.MessageEmbed) *discordgo.MessageEmbed {
	embed.Author.Name = "Command: vdot / vectordot"
	embed.Description = "`[math] (vdot|vectordot) <vector1> <vector2>` will get the dot product of the 2 vectors."
	embed.Fields = []*discordgo.MessageEmbedField{
		&discordgo.MessageEmbedField{
			Name:   "<vector1> <vector2>",
			Value:  "The vectors to get the dot product for.",
			Inline: true,
		},
		&discordgo.MessageEmbedField{
			Name:   "Vector Notation:",
			Value:  "2 Dimensions: (x, y)\n 3 Dimensions: (x, y, z)",
			Inline: true,
		},
		&discordgo.MessageEmbedField{
			Name:  "Related commands:",
			Value: "`vectoradd`, `vectorcross`, `vectordivide`, `vectormultiply`, `vectorsubtract`",
		},
	}
	return embed
}

// VectorMultiply explains the vector multiplication functionality
func VectorMultiply(embed *discordgo.MessageEmbed) *discordgo.MessageEmbed {
	embed.Author.Name = "Command: va / vadd / vectora / vectoradd"
	embed.Description = "`[math] (va|vadd|vectora|vectoradd) <vector1> <number>` will add the 2 vectors together."
	embed.Fields = []*discordgo.MessageEmbedField{
		&discordgo.MessageEmbedField{
			Name:   "<vector1> <vector2>",
			Value:  "The vector to multiply.",
			Inline: true,
		},
		&discordgo.MessageEmbedField{
			Name:   "<number>",
			Value:  "The number to multiply the vector with.",
			Inline: true,
		},
		&discordgo.MessageEmbedField{
			Name:   "Vector Notation:",
			Value:  "2 Dimensions: (x, y)\n 3 Dimensions: (x, y, z)",
			Inline: true,
		},
		&discordgo.MessageEmbedField{
			Name:  "Related commands:",
			Value: "`vectoradd`, `vectorcross`, `vectordivide`, `vectordot`, `vectorsubtract`",
		},
	}
	return embed
}

// VectorSubtract explains the vector substraction functionality
func VectorSubtract(embed *discordgo.MessageEmbed) *discordgo.MessageEmbed {
	embed.Author.Name = "Command: va / vadd / vectora / vectoradd"
	embed.Description = "`[math] (va|vadd|vectora|vectoradd) <vector1> <vector2>` will subtract the 2 vectors from each other."
	embed.Fields = []*discordgo.MessageEmbedField{
		&discordgo.MessageEmbedField{
			Name:   "<vector1> <vector2>",
			Value:  "The vectors to subtract each other for.",
			Inline: true,
		},
		&discordgo.MessageEmbedField{
			Name:   "Vector Notation:",
			Value:  "2 Dimensions: (x, y)\n 3 Dimensions: (x, y, z)",
			Inline: true,
		},
		&discordgo.MessageEmbedField{
			Name:  "Related commands:",
			Value: "`vectoradd`, `vectorcross`, `vectordivide`, `vectordot`, `vectormultiply`",
		},
	}
	return embed
}
