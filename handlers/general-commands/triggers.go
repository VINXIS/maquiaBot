package gencommands

import (
	"regexp"
	"strconv"
	"strings"

	tools "../../tools"
	"github.com/bwmarrin/discordgo"
)

// Triggers list out the triggers enabled in the server
func Triggers(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Get server
	server, err := s.Guild(m.GuildID)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "This is not a server!")
		return
	}
	serverImg := "https://cdn.discordapp.com/icons/" + server.ID + "/" + server.Icon
	if strings.Contains(server.Icon, "a_") {
		serverImg += ".gif"
	} else {
		serverImg += ".png"
	}

	serverData := tools.GetServer(*server)
	embed := &discordgo.MessageEmbed{
		Author: &discordgo.MessageEmbedAuthor{
			Name:    server.Name,
			IconURL: serverImg,
		},
	}

	if len(serverData.Triggers) == 0 {
		s.ChannelMessageSend(m.ChannelID, "There are no triggers configuered for this server currently! Admins can see `help trigger` for details on how to add triggers.")
		return
	}

	for _, trigger := range serverData.Triggers {
		regex := false
		_, err := regexp.Compile(trigger.Cause)
		if err == nil {
			regex = true
		}
		embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
			Name:  strconv.FormatInt(trigger.ID, 10),
			Value: "Trigger: " + trigger.Cause + "\nResult: " + trigger.Result + "\nRegex compatible: " + strconv.FormatBool(regex),
		})
	}
	s.ChannelMessageSendEmbed(m.ChannelID, embed)
	return
}
