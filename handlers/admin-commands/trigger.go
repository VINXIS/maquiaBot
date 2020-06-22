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

// Trigger adds / removes triggers
func Trigger(s *discordgo.Session, m *discordgo.MessageCreate) {
	triggerRegex, _ := regexp.Compile(`(?i)trigger\s+(.+)`)
	deleteRegex, _ := regexp.Compile(`(?i)-d`)

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

	// Check if params were given
	if !triggerRegex.MatchString(m.Content) {
		s.ChannelMessageSend(m.ChannelID, "No params given!")
		return
	}

	text := triggerRegex.FindStringSubmatch(m.Content)[1]

	// Delete function
	if deleteRegex.MatchString(m.Content) {
		text = strings.TrimSpace(deleteRegex.ReplaceAllString(text, ""))
		if ID, err := strconv.ParseInt(text, 10, 64); err == nil {
			for i, trigger := range serverData.Triggers {
				if trigger.ID == ID {
					serverData.Triggers = append(serverData.Triggers[:i], serverData.Triggers[i+1:]...)
					break
				}
			}
			jsonCache, err := json.Marshal(serverData)
			tools.ErrRead(s, err)

			err = ioutil.WriteFile("./data/serverData/"+m.GuildID+".json", jsonCache, 0644)
			tools.ErrRead(s, err)
			s.ChannelMessageSend(m.ChannelID, "Removed trigger ID: "+text)
		} else {
			s.ChannelMessageSend(m.ChannelID, text+" is an invalid ID!")
		}
		return
	}

	// Make sure there's enough params in the first place
	words := strings.Split(text, " ")
	if len(words) == 0 {
		s.ChannelMessageSend(m.ChannelID, "No trigger word found!")
		return
	} else if len(words) == 1 {
		s.ChannelMessageSend(m.ChannelID, "No trigger result found!")
		return
	}

	var (
		cause, result []string
	)

	// Split for multi-word triggers
	if strings.Contains(text, " | ") {
		shift := false
		for _, word := range words {
			if word == "|" {
				shift = true
				continue
			}

			if !shift {
				cause = append(cause, word)
			} else {
				result = append(result, word)
			}
		}
	} else {
		cause = words[:1]
		result = words[1:]
	}

	// Check if a cause and result were found
	if len(cause) == 0 {
		s.ChannelMessageSend(m.ChannelID, "No trigger word found!")
		return
	} else if len(result) == 0 {
		s.ChannelMessageSend(m.ChannelID, "No trigger result found!")
		return
	}

	// Create trigger data
	triggerData := structs.NewTrigger(strings.Join(cause, " "), strings.Join(result, " "))
	triggerData.Cause = strings.TrimSpace(strings.ToLower(triggerData.Cause))
	triggerData.Result = strings.TrimSpace(triggerData.Result)

	// Check duplicate and ID
	for _, trigger := range serverData.Triggers {
		if trigger.Cause == triggerData.Cause && trigger.Result == triggerData.Result {
			s.ChannelMessageSend(m.ChannelID, "This trigger already exists!")
			return
		}

		if trigger.ID == triggerData.ID {
			triggerData.ID++
		}
	}

	serverData.Triggers = append(serverData.Triggers, triggerData)
	jsonCache, err := json.Marshal(serverData)
	tools.ErrRead(s, err)

	err = ioutil.WriteFile("./data/serverData/"+m.GuildID+".json", jsonCache, 0644)
	tools.ErrRead(s, err)

	s.ChannelMessageSend(m.ChannelID, "The trigger `"+triggerData.Cause+"` for `"+triggerData.Result+"` has been created!")

}
