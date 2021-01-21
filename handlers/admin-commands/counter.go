package admincommands

import (
	"encoding/json"
	"io/ioutil"
	"regexp"
	"strconv"
	"strings"

	structs "maquiaBot/structs"
	tools "maquiaBot/tools"

	"github.com/bwmarrin/discordgo"
)

// Counter adds / removes counters
func Counter(s *discordgo.Session, m *discordgo.MessageCreate) {
	counterRegex, _ := regexp.Compile(`(?i)counter\s+(.+)`)
	deleteRegex, _ := regexp.Compile(`(?i)-d`)

	// Check if params were given
	if !counterRegex.MatchString(m.Content) {
		s.ChannelMessageSend(m.ChannelID, "No params given!")
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
	serverData, _ := tools.GetServer(*server, s)

	text := counterRegex.FindStringSubmatch(m.Content)[1]

	// Delete function
	if deleteRegex.MatchString(m.Content) {
		text = strings.TrimSpace(deleteRegex.ReplaceAllString(text, ""))
		if ID, err := strconv.ParseInt(text, 10, 64); err == nil {
			for i, counter := range serverData.Counters {
				if counter.ID == ID {
					serverData.Counters = append(serverData.Counters[:i], serverData.Counters[i+1:]...)
					break
				}
			}
			jsonCache, err := json.Marshal(serverData)
			tools.ErrRead(s, err)

			err = ioutil.WriteFile("./data/serverData/"+m.GuildID+".json", jsonCache, 0644)
			tools.ErrRead(s, err)
			s.ChannelMessageSend(m.ChannelID, "Removed counter ID: "+text)
		} else {
			s.ChannelMessageSend(m.ChannelID, text+" is an invalid ID!")
		}
		return
	}

	// Create counter data
	counterData := structs.NewCounter(text)

	// Check duplicate and ID
	for _, counter := range serverData.Counters {
		if strings.ToLower(counter.Text) == strings.ToLower(counterData.Text) {
			s.ChannelMessageSend(m.ChannelID, "This counter already exists!")
			return
		}
	}

	serverData.Counters = append(serverData.Counters, counterData)
	jsonCache, err := json.Marshal(serverData)
	tools.ErrRead(s, err)

	err = ioutil.WriteFile("./data/serverData/"+m.GuildID+".json", jsonCache, 0644)
	tools.ErrRead(s, err)

	s.ChannelMessageSend(m.ChannelID, "The counter for `"+counterData.Text+"` has been created!")

}
