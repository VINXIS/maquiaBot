package gencommands

import (
	"net/http"
	"regexp"

	"github.com/bwmarrin/discordgo"
)

// Emoji gives the image to an emoji
func Emoji(s *discordgo.Session, m *discordgo.MessageCreate) {
	emojiRegex, _ := regexp.Compile(`(?i)<(a)?:(.+):(\d+)>`)
	if !emojiRegex.MatchString(m.Content) {
		s.ChannelMessageSend(m.ChannelID, "No emoji given!")
		return
	}

	emojiID := emojiRegex.FindStringSubmatch(m.Content)[3]
	animated := emojiRegex.FindStringSubmatch(m.Content)[1] == "a"
	message := &discordgo.MessageSend{}
	if animated {
		res, err := http.Get("https://cdn.discordapp.com/emojis/" + emojiID + ".gif")
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, "Error obtaining emote!")
			return
		}
		message.Files = append(message.Files, &discordgo.File{
			Name:   emojiID + ".gif",
			Reader: res.Body,
		})
	} else {
		res, err := http.Get("https://cdn.discordapp.com/emojis/" + emojiID + ".png")
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, "Error obtaining emote!")
			return
		}
		message.Files = append(message.Files, &discordgo.File{
			Name:   emojiID + ".png",
			Reader: res.Body,
		})
	}
	s.ChannelMessageSendComplex(m.ChannelID, message)
}
