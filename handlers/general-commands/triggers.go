package gencommands

import (
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

	if len(serverData.RoleAutomation) == 0 {
		s.ChannelMessageSend(m.ChannelID, "There is no role automation configured for this server currently! Admins can see `help roleautomation` for details on how to add role automation.")
		return
	}

	for _, roleAuto := range serverData.RoleAutomation {
		var roleNames string
		for _, role := range roleAuto.Roles {
			roleNames += role.Name + ", "
		}
		embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
			Name:  strconv.Itoa(roleAuto.ID),
			Value: "Trigger: " + roleAuto.Text + "\nRoles: " + strings.TrimSuffix(roleNames, ", "),
		})
	}
	s.ChannelMessageSendEmbed(m.ChannelID, embed)
	return
}
