package paginationcommands

import (
	"maquiaBot/structs"
	"regexp"
	"strconv"

	"github.com/bwmarrin/discordgo"
)

// Counters handles the pagination of a list of counters
func Counters(s *discordgo.Session, r *discordgo.MessageReactionAdd, msg *discordgo.Message, serverData structs.ServerData, num, numend int) (*discordgo.MessageEmbed, bool) {
	embed := &discordgo.MessageEmbed{}

	if numend > len(serverData.Counters) {
		numend = len(serverData.Counters)
	}
	counters := serverData.Counters[num:numend]
	for _, counter := range counters {
		counter.Text = `(?i)` + counter.Text
		regex := false
		_, err := regexp.Compile(counter.Text)
		if err == nil {
			regex = true
		}
		valueText := "Counter: " + counter.Text + "\nRegex compatible: " + strconv.FormatBool(regex)
		if len(valueText) > 1024 {
			valueText = valueText[:1021] + "..."
		}
		embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
			Name:  strconv.FormatInt(counter.ID, 10),
			Value: valueText,
		})

		if len(embed.Fields) == 25 {
			break
		}
	}

	return embed, len(embed.Fields) < 25 || numend == len(serverData.Triggers)
}
