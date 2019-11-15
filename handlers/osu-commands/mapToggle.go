package osucommands

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"time"

	structs "../../structs"
	tools "../../tools"
	"github.com/bwmarrin/discordgo"
)

// MapToggle toggles crab messages on/off
func MapToggle(s *discordgo.Session, m *discordgo.MessageCreate) {
	server, err := s.Guild(m.GuildID)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "This is not a server!")
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

	// Set new information in server data
	serverData.Time = time.Now()
	serverData.OsuToggle = !serverData.OsuToggle

	jsonCache, err := json.Marshal(serverData)
	tools.ErrRead(err)

	err = ioutil.WriteFile("./data/serverData/"+m.GuildID+".json", jsonCache, 0644)
	tools.ErrRead(err)

	if serverData.OsuToggle {
		s.ChannelMessageSend(m.ChannelID, "Enabled map/user links O_o")
	} else {
		s.ChannelMessageSend(m.ChannelID, "Disabled map/user links O_o")
	}
	return
}
