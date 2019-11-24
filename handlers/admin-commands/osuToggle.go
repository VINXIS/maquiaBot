package admincommands

import (
	"encoding/json"
	"io/ioutil"
	"time"

	tools "../../tools"
	"github.com/bwmarrin/discordgo"
)

// OsuToggle toggles beatmap/user messages on/off
func OsuToggle(s *discordgo.Session, m *discordgo.MessageCreate) {
	server, err := s.Guild(m.GuildID)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "This is not a server!")
		return
	}

	if !tools.AdminCheck(s, m, *server) {
		s.ChannelMessageSend(m.ChannelID, "You must be an admin, server manager, or server owner!")
		return
	}

	// Obtain server data
	serverData := tools.GetServer(*server)

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
