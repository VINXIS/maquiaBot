package gencommands

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"regexp"
	"strconv"
	"time"

	structs "../../structs"
	tools "../../tools"
	"github.com/bwmarrin/discordgo"
)

// Nouns allows users to add/change/see their custom nouns
func Nouns(s *discordgo.Session, m *discordgo.MessageCreate) {
	nounsRegex, _ := regexp.Compile(`(nouns)\s*(add|remove)?\s+(.+)`)

	server, err := s.Guild(m.GuildID)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "This is not a server!")
		return
	}

	// Obtain server data
	serverData := structs.ServerData{
		Prefix:    "$",
		Crab:      true,
		OsuToggle: true,
	}
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

	// Obtain word and if they want to add/remove it
	mode := "add"
	word := ""
	if nounsRegex.MatchString(m.Content) {
		matches := nounsRegex.FindStringSubmatch(m.Content)
		if matches[2] != "" {
			mode = matches[2]
		}
		word = matches[3]
	} else if len(serverData.Nouns) != 0 {
		text := "There are " + strconv.Itoa(len(serverData.Nouns)) + " nouns listed for this server! The nouns are:\n"
		for i, noun := range serverData.Nouns {
			if i != len(serverData.Nouns)-1 {
				text = text + noun + ", "
			} else {
				text = text + noun
			}
		}
		s.ChannelMessageSend(m.ChannelID, text)
		return
	}

	if word == "" {
		s.ChannelMessageSend(m.ChannelID, "No word to add/remove from the server's noun list found!")
		return
	}

	// Commence operation
	if mode == "add" {
		for _, noun := range serverData.Nouns {
			if noun == word {
				s.ChannelMessageSend(m.ChannelID, "`"+word+"` is already in the server's noun list!")
				return
			}
		}
		serverData.Nouns = append(serverData.Nouns, word)
	} else if mode == "remove" {
		contained := false
		for i, noun := range serverData.Nouns {
			if noun == word {
				serverData.Nouns[i] = serverData.Nouns[len(serverData.Nouns)-1]
				serverData.Nouns = serverData.Nouns[:len(serverData.Nouns)-1]
				contained = true
				break
			}
		}
		if !contained {
			s.ChannelMessageSend(m.ChannelID, "`"+word+"` does not exist in the server's noun list!")
			return
		}
	} else {
		s.ChannelMessageSend(m.ChannelID, "Something broke!!!!!!")
		return
	}

	// Set new information in server data
	serverData.Time = time.Now()

	jsonCache, err := json.Marshal(serverData)
	tools.ErrRead(err)

	err = ioutil.WriteFile("./data/serverData/"+m.GuildID+".json", jsonCache, 0644)
	tools.ErrRead(err)
	if mode == "add" {
		s.ChannelMessageSend(m.ChannelID, "`"+word+"` is now added to the server's noun list!")
	} else if mode == "remove" {
		s.ChannelMessageSend(m.ChannelID, "`"+word+"` is now removed from the server's noun list!")
	}
	return
}
