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

// Skills allows users to add/change/see their custom skills
func Skills(s *discordgo.Session, m *discordgo.MessageCreate) {
	skillsRegex, _ := regexp.Compile(`(skills)\s*(add|remove)?\s+(.+)`)

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
	if skillsRegex.MatchString(m.Content) {
		matches := skillsRegex.FindStringSubmatch(m.Content)
		if matches[2] != "" {
			mode = matches[2]
		}
		word = matches[3]
	} else if len(serverData.Skills) != 0 {
		text := "There are " + strconv.Itoa(len(serverData.Skills)) + " skills listed for this server! The skills are:\n"
		for i, skill := range serverData.Skills {
			if i != len(serverData.Skills)-1 {
				text = text + skill + ", "
			} else {
				text = text + skill
			}
		}
		s.ChannelMessageSend(m.ChannelID, text)
		return
	}

	if word == "" {
		s.ChannelMessageSend(m.ChannelID, "No word to add/remove from the server's skill list found!")
		return
	}

	// Commence operation
	if mode == "add" {
		for _, skill := range serverData.Skills {
			if skill == word {
				s.ChannelMessageSend(m.ChannelID, "`"+word+"` is already in the server's skill list!")
				return
			}
		}
		serverData.Skills = append(serverData.Skills, word)
	} else if mode == "remove" {
		contained := false
		for i, skill := range serverData.Skills {
			if skill == word {
				serverData.Skills[i] = serverData.Skills[len(serverData.Skills)-1]
				serverData.Skills = serverData.Skills[:len(serverData.Skills)-1]
				contained = true
				break
			}
		}
		if !contained {
			s.ChannelMessageSend(m.ChannelID, "`"+word+"` does not exist in the server's skill list!")
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
		s.ChannelMessageSend(m.ChannelID, "`"+word+"` is now added to the server's skill list!")
	} else if mode == "remove" {
		s.ChannelMessageSend(m.ChannelID, "`"+word+"` is now removed from the server's skill list!")
	}
	return
}
