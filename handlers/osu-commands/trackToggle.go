package osucommands

import (
	"encoding/json"
	"io/ioutil"
	"os"

	osutools "../../osu-functions"
	structs "../../structs"
	tools "../../tools"
	"github.com/bwmarrin/discordgo"
)

// TrackToggle stops tracking for the channel
func TrackToggle(s *discordgo.Session, m *discordgo.MessageCreate, mapCache []structs.MapData) {
	// Check perms
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
		if role.Permissions&discordgo.PermissionAdministrator == discordgo.PermissionAdministrator || role.Permissions&discordgo.PermissionManageServer == discordgo.PermissionManageServer {
			admin = true
			break
		}
	}

	if !admin && m.Author.ID != server.OwnerID {
		s.ChannelMessageSend(m.ChannelID, "You are not an admin or server manager!")
		return
	}

	// Obtain channel data
	channelData := structs.ChannelData{}
	_, err = os.Stat("./data/channelData/" + m.ChannelID + ".json")
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "No tracking information currently for this channel!")
		return
	}
	f, err := ioutil.ReadFile("./data/channelData/" + m.ChannelID + ".json")
	tools.ErrRead(err)
	_ = json.Unmarshal(f, &channelData)

	channelData.Tracking = !channelData.Tracking

	// Write data to JSON
	jsonCache, err := json.Marshal(channelData)
	tools.ErrRead(err)

	err = ioutil.WriteFile("./data/channelData/"+m.ChannelID+".json", jsonCache, 0644)
	tools.ErrRead(err)

	if channelData.Tracking {
		go osutools.TrackPost("data/channelData/"+m.ChannelID+".json", s, mapCache)
		s.ChannelMessageSend(m.ChannelID, "Successfully started tracking for this channel!")
	} else {
		s.ChannelMessageSend(m.ChannelID, "Successfully stopped tracking for this channel!")
	}
}
