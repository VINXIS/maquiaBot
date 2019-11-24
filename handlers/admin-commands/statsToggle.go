package admincommands

import (
	"encoding/json"
	"io/ioutil"
	"time"

	tools "../../tools"
	"github.com/bwmarrin/discordgo"
)

// StatsToggle allows users to add/remove nouns, adjectives, and skills
func StatsToggle(s *discordgo.Session, m *discordgo.MessageCreate) {
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
	serverData.AllowAnyoneStats = !serverData.AllowAnyoneStats

	jsonCache, err := json.Marshal(serverData)
	tools.ErrRead(err)

	err = ioutil.WriteFile("./data/serverData/"+m.GuildID+".json", jsonCache, 0644)
	tools.ErrRead(err)

	if serverData.AllowAnyoneStats {
		s.ChannelMessageSend(m.ChannelID, "Anyone may add nouns/adjectives/skills for stats now.")
	} else {
		s.ChannelMessageSend(m.ChannelID, "Only admins, server managers, or the server owner may add nouns/adjectives/skills for stats now.")
	}
	return
}
