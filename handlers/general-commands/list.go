package gencommands

import (
	"fmt"
	"math/rand"
	"strings"

	"github.com/bwmarrin/discordgo"
)

// List randomizes a list of objects
func List(s *discordgo.Session, m *discordgo.MessageCreate) {
	fmt.Println(m.Content)
	list := strings.Split(m.Content, "\n")[1:]

	if len(list) <= 1 {
		s.ChannelMessageSend(m.ChannelID, "Please give a list of lines to randomize!")
		return
	}

	rand.Shuffle(len(list), func(i, j int) { list[i], list[j] = list[j], list[i] })

	s.ChannelMessageSend(m.ChannelID, strings.Join(list, "\n"))
}
