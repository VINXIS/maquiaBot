package gencommands

import (
	"bytes"
	"io/ioutil"
	"math/rand"
	"net/http"
	"strings"

	"github.com/bwmarrin/discordgo"
)

// List randomizes a list of objects
func List(s *discordgo.Session, m *discordgo.MessageCreate) {
	list := strings.Split(m.Content, "\n")[1:]

	// Use txt file if given
	if len(m.Attachments) > 0 {
		res, err := http.Get(m.Attachments[0].URL)
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, "Unable to get file information!")
			return
		}
		defer res.Body.Close()

		b, err := ioutil.ReadAll(res.Body)
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, "Unable to parse file information!")
			return
		}

		list = strings.Split(string(b), "\n")
		if strings.Contains(list[0], "list") {
			list = list[1:]
		}
	}

	if len(list) <= 1 {
		s.ChannelMessageSend(m.ChannelID, "Please give a list of lines to randomize!")
		return
	}

	rand.Shuffle(len(list), func(i, j int) { list[i], list[j] = list[j], list[i] })

	_, err := s.ChannelMessageSend(m.ChannelID, strings.Join(list, "\n"))
	if err != nil {
		buf := new(bytes.Buffer)
		buf.Write([]byte(strings.Join(list, "\n")))
		s.ChannelMessageSendComplex(m.ChannelID, &discordgo.MessageSend{
			Files: []*discordgo.File{
				&discordgo.File{
					Name:   "shuffle.txt",
					Reader: buf,
				},
			},
		})
	}
}
