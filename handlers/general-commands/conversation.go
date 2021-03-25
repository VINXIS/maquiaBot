package gencommands

import (
	"crypto/rand"
	tools "maquiaBot/tools"
	"math/big"
	"regexp"
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
)

// Conversation creates a conversation based off of quotes
func Conversation(s *discordgo.Session, m *discordgo.MessageCreate) {
	convoRegex, _ := regexp.Compile(`(?i)(convo?|conversation)?\s+(\d+)`)
	linkRegex, _ := regexp.Compile(`(?i)(https://www|https://|www)\S+`)

	// Get server
	server, err := s.Guild(m.GuildID)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "This is not a server!")
		return
	}
	serverData, _ := tools.GetServer(*server, s)
	if len(serverData.Quotes) == 0 {
		s.ChannelMessageSend(m.ChannelID, "No quotes saved for this server! Please see `help quoteadd` to see how to add quotes!")
		return
	}

	// Get number of quotes and if links should be removed
	num := 2
	excludeLinks := true
	if convoRegex.MatchString(m.Content) {
		num, err = strconv.Atoi(convoRegex.FindStringSubmatch(m.Content)[2])
		if err != nil {
			num = 2
		}
	}

	if strings.Contains(m.Content, "-i") {
		excludeLinks = false
	}

	if num > len(serverData.Quotes) {
		s.ChannelMessageSend(m.ChannelID, "You don't have enough quotes in the server! Please see `help quoteadd` to see how to add quotes!")
		return
	}

	// Create the Convo .
	convo := []string{}
	for i := 0; i < num; i++ {
		roll, _ := rand.Int(rand.Reader, big.NewInt(int64(len(serverData.Quotes))))
		j := roll.Int64()
		if len(serverData.Quotes[j].Attachments) <= 0 || excludeLinks {
			text := serverData.Quotes[j].ContentWithMentionsReplaced()
			if linkRegex.MatchString(text) {
				text = strings.TrimSpace(linkRegex.ReplaceAllString(text, ""))
			}
			if text == "" {
				i--
				continue
			}
			convo = append(convo, "**"+serverData.Quotes[j].Author.Username+"**: "+serverData.Quotes[j].ContentWithMentionsReplaced())
		} else if len(serverData.Quotes[j].Attachments) > 0 {
			convo = append(convo, "**"+serverData.Quotes[j].Author.Username+"**: "+serverData.Quotes[j].ContentWithMentionsReplaced()+" "+serverData.Quotes[j].Attachments[0].URL)
		} else {
			i--
			continue
		}

		if len(strings.Join(convo, "\n")) > 2000 {
			convo = convo[:len(convo)-1]
			break
		}
	}
	if excludeLinks {
		s.ChannelMessageSend(m.ChannelID, "```md\n"+strings.Join(convo, "\n")+"```")
	} else {
		s.ChannelMessageSend(m.ChannelID, strings.Join(convo, "\n"))
	}
}
