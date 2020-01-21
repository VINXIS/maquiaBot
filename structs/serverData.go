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
	Cheers           bool
	Daily            bool
	Late             bool
	NiceIdea         bool
	OsuToggle        bool
	OverIt           bool
	Vibe             bool
	Announce         bool
	Adjectives       []string
	Nouns            []string
	Skills           []string
	AllowAnyoneStats bool
	Quotes           []discordgo.Message
	Genital          GenitalRecordData
	RoleAutomation   []Role
}

// Role holds information for role automation
type Role struct {
	ID    int
	Text  string
	Roles []discordgo.Role
}

// NewServer creates a new ServerData
func NewServer(server discordgo.Guild) ServerData {
	return ServerData{
		Prefix:    "$",
		OsuToggle: true,
		Crab:      true,
		Cheers:    true,
		Daily:     true,
		Late:      true,
		NiceIdea:  true,
		OverIt:    true,
		Announce:  true,
		Server:    server,
		Genital: GenitalRecordData{
			Penis: struct {
				Largest  GenitalData
				Smallest GenitalData
			}{
				Smallest: GenitalData{
					Size: 1e308,
				},
			},
			Vagina: struct {
				Largest  GenitalData
				Smallest GenitalData
			}{
				Smallest: GenitalData{
					Size: 1e308,
				},
			},
		},
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

// AddQuote adds a quote to the server data
func (s *ServerData) AddQuote(message *discordgo.Message) error {
	for _, quote := range s.Quotes {
		if quote.ID == message.ID {
			return errors.New("message already a quote")
		}
	}
	s.Quotes = append(s.Quotes, *message)
	return nil
}

// RemoveQuote removes a quote from the server data
func (s *ServerData) RemoveQuote(ID string) error {
	for i, quote := range s.Quotes {
		if quote.ID == ID {
			s.Quotes[i] = s.Quotes[len(s.Quotes)-1]
			s.Quotes = s.Quotes[:len(s.Quotes)-1]
			return nil
		}
	}
	return errors.New("message is not a quote")
}
