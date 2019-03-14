package gencommands

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"time"

	structs "../../structs"
	tools "../../tools"

	"github.com/bwmarrin/discordgo"
)

// NewPrefix sets a new prefix for the bot
func NewPrefix(s *discordgo.Session, m *discordgo.MessageCreate, args []string, serverPrefix string) {
	server, err := s.Guild(m.GuildID)
	tools.ErrRead(err)

	member := &discordgo.Member{}
	for _, guildMember := range server.Members {
		if guildMember.User.ID == m.Author.ID {
			member = guildMember
		}
	}

	if member.User.ID == "" {
		return
	}

	admin := false
	for _, roleID := range member.Roles {
		role, err := s.State.Role(m.GuildID, roleID)
		tools.ErrRead(err)
		if role.Permissions&discordgo.PermissionAdministrator == discordgo.PermissionAdministrator {
			admin = true
			break
		}
	}

	if !admin {
		s.ChannelMessageSend(m.ChannelID, "You are not an admin!")
		return
	}

	// Obtain server data
	serverData := structs.ServerData{}
	_, err = os.Stat("./data/serverData/" + m.GuildID + ".json")
	if err == nil {
		f, err := ioutil.ReadFile("./data/serverData/" + m.GuildID + ".json")
		tools.ErrRead(err)
		_ = json.Unmarshal(f, &serverData)
	} else if os.IsNotExist(err) {
		serverData.Server = *server
	} else {
		tools.ErrRead(err)
		return
	}

	// Set new information in server data
	serverData.Time = time.Now()
	serverData.Prefix = args[1]

	jsonCache, err := json.Marshal(serverData)
	tools.ErrRead(err)

	err = ioutil.WriteFile("./data/serverData/"+m.GuildID+".json", jsonCache, 0644)
	tools.ErrRead(err)

	s.ChannelMessageSend(m.ChannelID, "Prefix changed from "+serverPrefix+" to "+args[1])
	return
}
