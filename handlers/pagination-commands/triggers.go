package paginationcommands

import (
	"maquiaBot/structs"
	"regexp"
	"strconv"

	"github.com/bwmarrin/discordgo"
)

// Triggers handles the pagination of a list of triggers
func Triggers(s *discordgo.Session, r *discordgo.MessageReactionAdd, msg *discordgo.Message, serverData structs.ServerData, num, numend int) (*discordgo.MessageEmbed, bool) {
	embed := &discordgo.MessageEmbed{}

	triggers := serverData.Triggers[num:numend]
	for _, trigger := range triggers {
		trigger.Cause = `(?i)` + trigger.Cause
		regex := false
		_, err := regexp.Compile(trigger.Cause)
		if err == nil {
			regex = true
		}
		valueText := "Trigger: " + trigger.Cause + "\nResult: " + trigger.Result + "\nRegex compatible: " + strconv.FormatBool(regex)
		if len(valueText) > 1024 {
			valueText = valueText[:1021] + "..."
		}
		embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
			Name:  strconv.FormatInt(trigger.ID, 10),
			Value: valueText,
		})

		if len(embed.Fields) == 25 {
			break
		}
	}

	return embed, len(embed.Fields) < 25 || numend == len(serverData.Triggers)
}
