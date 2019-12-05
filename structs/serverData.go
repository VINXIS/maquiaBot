package structs

import (
	"errors"
	"time"

	"github.com/bwmarrin/discordgo"
)

// ServerData stores information regarding the discord server, so that server specific customizations may be used.
type ServerData struct {
	Time             time.Time
	Server           discordgo.Guild
	Prefix           string
	Crab             bool
	NiceIdea         bool
	OsuToggle        bool
	Possession       bool
	Vibe             bool
	Adjectives       []string
	Nouns            []string
	Skills           []string
	AllowAnyoneStats bool
}

// NewServer creates a new ServerData
func NewServer(server discordgo.Guild) ServerData {
	return ServerData{
		Prefix:    "$",
		OsuToggle: true,
		Crab:      true,
		NiceIdea:  true,
		Server:    server,
	}
}

// Word adds the word to the specified list
func (s *ServerData) Word(word, mode, list string) error {
	var targetList []string
	switch list {
	case "adjective":
		targetList = s.Adjectives
	case "noun":
		targetList = s.Nouns
	case "skill":
		targetList = s.Skills
	}
	if mode == "add" {
		for _, existingWord := range targetList {
			if existingWord == word {
				return errors.New("`" + word + "` is already in the server's " + list + " list!")
			}
		}
		targetList = append(targetList, word)
	} else if mode == "remove" {
		contained := false
		for i, existingWord := range targetList {
			if existingWord == word {
				targetList[i] = targetList[len(targetList)-1]
				targetList = targetList[:len(targetList)-1]
				contained = true
				break
			}
		}
		if !contained {
			return errors.New("`" + word + "` does not exist in the server's " + list + " list!")
		}
	}
	switch list {
	case "adjective":
		s.Adjectives = targetList
	case "noun":
		s.Nouns = targetList
	case "skill":
		s.Skills = targetList
	}
	return nil
}
