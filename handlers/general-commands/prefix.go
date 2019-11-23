package gencommands

import (
	"encoding/json"
	"io/ioutil"
	"strings"
	"time"

	tools "../../tools"

	"github.com/bwmarrin/discordgo"
)

// Prefix sets a new prefix for the bot
func Prefix(s *discordgo.Session, m *discordgo.MessageCreate) {
	if len(m.Mentions) > 0 {
		s.ChannelMessageSend(m.ChannelID, "Please don't try making mentions a prefix with the bot! >:/")
		return
	}

	// Check if server exists
	server, err := s.Guild(m.GuildID)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "This is not a server so custom prefixes are unavailable! Please use `$` instead for commands!")
		return
	}

	// Check if member has proper server privileges
	if m.Author.ID != server.OwnerID {
		member, _ := s.GuildMember(server.ID, m.Author.ID)
		admin := false
		for _, roleID := range member.Roles {
			role, _ := s.State.Role(m.GuildID, roleID)
			if role.Permissions&discordgo.PermissionAdministrator != 0 || role.Permissions&discordgo.PermissionManageServer != 0 {
				admin = true
				break
			}
		}
		if !admin {
			s.ChannelMessageSend(m.ChannelID, "You must be an admin, server manager, or server owner!")
			return
		}
	}

	// Obtain server data
	serverData := tools.GetServer(*server)

	// Set new information in server data
	oldPrefix := serverData.Prefix
	prefix := strings.Split(m.Content, " ")[1]
	serverData.Time = time.Now()
	serverData.Prefix = prefix

	jsonCache, err := json.Marshal(serverData)
	tools.ErrRead(err)

	err = ioutil.WriteFile("./data/serverData/"+m.GuildID+".json", jsonCache, 0644)
	tools.ErrRead(err)

	s.ChannelMessageSend(m.ChannelID, "Prefix changed from "+oldPrefix+" to "+serverData.Prefix)
	return
}
