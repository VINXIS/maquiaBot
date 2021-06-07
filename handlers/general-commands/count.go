package gencommands

import (
	"io/ioutil"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
)

// CharCount counts the number of characters in text
func CharCount(s *discordgo.Session, m *discordgo.MessageCreate) {
	charRegex, _ := regexp.Compile(`(?is)(c(har)?)?(count)?\s+(.+)`)
	text := ""

	if !charRegex.MatchString(m.Content) {
		if len(m.Attachments) == 0 {
			s.ChannelMessageSend(m.ChannelID, "No text provided to char count for!")
			return
		}
		res, err := http.Get(m.Attachments[0].URL)
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, "Unable to access URL provided")
			return
		}
		defer res.Body.Close()
		b, err := ioutil.ReadAll(res.Body)
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, "Unable to parse URL provided")
			return
		}
		if !strings.Contains(http.DetectContentType(b), "text/plain") {
			s.ChannelMessageSend(m.ChannelID, "Attachment provided is not a valid txt file!")
			return
		}
		text = string(b)
	} else {
		text = charRegex.FindStringSubmatch(m.Content)[4]
	}
	s.ChannelMessageSend(m.ChannelID, strconv.Itoa(len(text)))
}

// WordCount counts the number of words in text
func WordCount(s *discordgo.Session, m *discordgo.MessageCreate) {
	wordRegex, _ := regexp.Compile(`(?is)w(ord)?(count)?\s+(.+)`)
	text := ""

	if !wordRegex.MatchString(m.Content) {
		if len(m.Attachments) == 0 {
			s.ChannelMessageSend(m.ChannelID, "No text provided to word count for!")
			return
		}
		res, err := http.Get(m.Attachments[0].URL)
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, "Unable to access URL provided")
			return
		}
		defer res.Body.Close()
		b, err := ioutil.ReadAll(res.Body)
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, "Unable to parse URL provided")
			return
		}
		if !strings.Contains(http.DetectContentType(b), "text/plain") {
			s.ChannelMessageSend(m.ChannelID, "Attachment provided is not a valid txt file!")
			return
		}
		text = string(b)
	} else {
		text = wordRegex.FindStringSubmatch(m.Content)[3]
	}

	s.ChannelMessageSend(m.ChannelID, strconv.Itoa(len(strings.Fields(text))))
}
