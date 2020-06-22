package admincommands

import (
	"encoding/json"
	"io/ioutil"
	"strings"
	"time"

	tools "maquiaBot/tools"

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

	if !tools.AdminCheck(s, m, *server) {
		s.ChannelMessageSend(m.ChannelID, "You must be an admin, server manager, or server owner!")
		return
	}

	// Obtain server data
	serverData := tools.GetServer(*server, s)

	// Set new information in server data
	oldPrefix := serverData.Prefix
	args := strings.Split(m.Content, " ")
	if len(args) < 2 {
		s.ChannelMessageSend(m.ChannelID, "No prefix given to change to!")
		return
	}
	prefix := strings.Split(m.Content, " ")[1]
	serverData.Time = time.Now()
	serverData.Prefix = prefix

	jsonCache, err := json.Marshal(serverData)
	tools.ErrRead(s, err)

	err = ioutil.WriteFile("./data/serverData/"+m.GuildID+".json", jsonCache, 0644)
	tools.ErrRead(s, err)

	s.ChannelMessageSend(m.ChannelID, "Prefix changed from "+oldPrefix+" to "+serverData.Prefix)
	return
}
