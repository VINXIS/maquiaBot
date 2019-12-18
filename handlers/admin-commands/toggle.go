package admincommands

import (
	"encoding/json"
	"io/ioutil"
	"strings"
	"time"

	tools "../../tools"
	gencommands "../general-commands"
	"github.com/bwmarrin/discordgo"
)

// Toggle toggles server options on/off
func Toggle(s *discordgo.Session, m *discordgo.MessageCreate) {
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
	if strings.Contains(m.Content, "-a") || strings.Contains(m.Content, "-announce") {
		serverData.Announce = !serverData.Announce
	}
	if strings.Contains(m.Content, "-ch") || strings.Contains(m.Content, "-cheers") {
		serverData.Cheers = !serverData.Cheers
	}
	if strings.Contains(m.Content, "-c") || strings.Contains(m.Content, "-crab") {
		serverData.Crab = !serverData.Crab
	}
	if strings.Contains(m.Content, "-d") || strings.Contains(m.Content, "-daily") {
		serverData.Daily = !serverData.Daily
	}
	if strings.Contains(m.Content, "-i") || strings.Contains(m.Content, "-idea") {
		serverData.NiceIdea = !serverData.NiceIdea
	}
	if strings.Contains(m.Content, "-l") || strings.Contains(m.Content, "-late") {
		serverData.Late = !serverData.Late
	}
	if strings.Contains(m.Content, "-o") || strings.Contains(m.Content, "-osu") {
		serverData.OsuToggle = !serverData.OsuToggle
	}
	if strings.Contains(m.Content, "-s") || strings.Contains(m.Content, "-stats") {
		serverData.AllowAnyoneStats = !serverData.AllowAnyoneStats
	}
	if strings.Contains(m.Content, "-v") || strings.Contains(m.Content, "-vibe") {
		serverData.Vibe = !serverData.Vibe
	}

	jsonCache, err := json.Marshal(serverData)
	tools.ErrRead(err)

	err = ioutil.WriteFile("./data/serverData/"+m.GuildID+".json", jsonCache, 0644)
	tools.ErrRead(err)

	go gencommands.ServerInfo(s, m)
}
