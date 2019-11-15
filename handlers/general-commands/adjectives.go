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

// Adjectives allows users to add/change/see their custom adjectives
func Adjectives(s *discordgo.Session, m *discordgo.MessageCreate) {
	adjectivesRegex, _ := regexp.Compile(`(adj|adjectives)\s*(add|remove)?\s+(.+)`)

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
	if adjectivesRegex.MatchString(m.Content) {
		matches := adjectivesRegex.FindStringSubmatch(m.Content)
		if matches[2] != "" {
			mode = matches[2]
		}
		word = matches[3]
	} else if len(serverData.Adjectives) != 0 {
		text := "There are " + strconv.Itoa(len(serverData.Adjectives)) + " adjectives listed for this server! The adjectives are:\n"
		for i, adjective := range serverData.Adjectives {
			if i != len(serverData.Adjectives)-1 {
				text = text + adjective + ", "
			} else {
				text = text + adjective
			}
		}
		s.ChannelMessageSend(m.ChannelID, text)
		return
	}

	if word == "" {
		s.ChannelMessageSend(m.ChannelID, "No word to add/remove from the server's adjective list found!")
		return
	}

	// Commence operation
	if mode == "add" {
		for _, adjective := range serverData.Adjectives {
			if adjective == word {
				s.ChannelMessageSend(m.ChannelID, "`"+word+"` is already in the server's adjective list!")
				return
			}
		}
		serverData.Adjectives = append(serverData.Adjectives, word)
	} else if mode == "remove" {
		contained := false
		for i, adjective := range serverData.Adjectives {
			if adjective == word {
				serverData.Adjectives[i] = serverData.Adjectives[len(serverData.Adjectives)-1]
				serverData.Adjectives = serverData.Adjectives[:len(serverData.Adjectives)-1]
				contained = true
				break
			}
		}
		if !contained {
			s.ChannelMessageSend(m.ChannelID, "`"+word+"` does not exist in the server's adjective list!")
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
		s.ChannelMessageSend(m.ChannelID, "`"+word+"` is now added to the server's adjective list!")
	} else if mode == "remove" {
		s.ChannelMessageSend(m.ChannelID, "`"+word+"` is now removed from the server's adjective list!")
	}
	return
}
