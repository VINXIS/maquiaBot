package gencommands

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"strings"
	"time"

	structs "../../structs"
	tools "../../tools"

	"github.com/bwmarrin/discordgo"
)

// Prefix sets a new prefix for the bot
func Prefix(s *discordgo.Session, m *discordgo.MessageCreate) {
	if len(m.Mentions) > 0 {
		s.ChannelMessageSend(m.ChannelID, "Please don't try mentioning people with the bot! >:/")
		return
	}

	server, err := s.Guild(m.GuildID)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "This is not a server so custom prefixes are unavailable! Please use `$` instead for commands!")
		return
	}

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
		if role.Permissions&discordgo.PermissionAdministrator != 0 || role.Permissions&discordgo.PermissionManageServer != 0 {
			admin = true
			break
		}
	}

	if !admin && m.Author.ID != server.OwnerID {
		s.ChannelMessageSend(m.ChannelID, "You must be an admin, server manager, or server owner!")
		return
	}

	// Obtain server data
	serverData := structs.NewServer()
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
	oldPrefix := serverData.Prefix

	// Set new information in server data
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
