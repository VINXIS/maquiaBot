package gencommands

import (
	"math/rand"
	"strings"
	"unicode"

	"github.com/bwmarrin/discordgo"
)

// Capitalization capitalizes / lowercases text
func Capitalization(s *discordgo.Session, m *discordgo.MessageCreate, capsType string) {
	message := strings.Join(strings.Split(m.Content, " ")[1:], " ")

	if len(message) == 0 {
		msgs, err := s.ChannelMessages(m.ChannelID, -1, m.ID, "", "")
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, "Ur really funny mate")
			return
		}
		for _, msg := range msgs {
			if msg.Author.ID != m.Author.ID {
				message = msgs[0].Content
				break
			}
		}
	}

	switch capsType {
	case "allLower":
		message = strings.ToLower(message)
	case "allCaps":
		message = strings.ToUpper(message)
	case "title":
		message = strings.Title(strings.ToLower(message))
	case "random":
		message = strings.Map(func(r rune) rune {
			val := rand.Intn(2)
			switch val {
			case 0:
				r = unicode.ToUpper(r)
			case 1:
				r = unicode.ToLower(r)
			}
			return r
		}, message)
	}

	s.ChannelMessageSend(m.ChannelID, message)
}
