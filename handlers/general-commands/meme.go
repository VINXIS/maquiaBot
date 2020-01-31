package gencommands

import (
	"net/http"
	"regexp"
	"strings"

	"github.com/bwmarrin/discordgo"
)

// Meme lets you create a meme
func Meme(s *discordgo.Session, m *discordgo.MessageCreate) {
	emojiRegex, _ := regexp.Compile(`<(a)?:(.+):(\d+)>`)
	linkRegex, _ := regexp.Compile(`https?:\/\/\S*`)
	memeRegex, _ := regexp.Compile(`meme\s+(https:\/\/(\S+)\s+)?([^|]+)?(\|)?([^|]+)?`)

	if !memeRegex.MatchString(m.Content) {
		s.ChannelMessageSend(m.ChannelID, "Please give text to add and maybe an image as well!")
		return
	}

	url := memeRegex.FindStringSubmatch(m.Content)[1]
	topText := strings.TrimSpace(memeRegex.FindStringSubmatch(m.Content)[3])
	bottomText := strings.TrimSpace(memeRegex.FindStringSubmatch(m.Content)[5])
	if topText == "" && bottomText == "" {
		s.ChannelMessageSend(m.ChannelID, "Please give text to add onto the image!")
		return
	} else if topText != "" && bottomText == "" {
		words := strings.Split(topText, " ")
		if len(words) > 1 {
			topText = strings.Join(words[:len(words)/2], " ")
			bottomText = strings.Join(words[len(words)/2:], " ")
		}
	}

	// Find URL
	if linkRegex.MatchString(m.Content) {
		url = linkRegex.FindStringSubmatch(m.Content)[0]
		topText = strings.Replace(topText, url, "", -1)
		bottomText = strings.Replace(bottomText, url, "", -1)
	} else if emojiRegex.MatchString(m.Content) {
		emojiID := emojiRegex.FindStringSubmatch(m.Content)[3]
		animated := emojiRegex.FindStringSubmatch(m.Content)[1] == "a"
		if animated {
			url = "https://cdn.discordapp.com/emojis/" + emojiID + ".gif"
		} else {
			url = "https://cdn.discordapp.com/emojis/" + emojiID + ".png"
		}
		topText = strings.Replace(topText, emojiRegex.FindStringSubmatch(m.Content)[0], "", -1)
		bottomText = strings.Replace(bottomText, emojiRegex.FindStringSubmatch(m.Content)[0], "", -1)
	} else if len(m.Attachments) > 0 {
		url = m.Attachments[0].URL
	} else if len(m.Embeds) > 0 && m.Embeds[0].Image != nil {
		url = m.Embeds[0].Image.URL
	} else if len(m.Mentions) > 0 {
		url = m.Mentions[0].AvatarURL("")
		topText = strings.Replace(topText, m.Mentions[0].Mention(), "", -1)
		bottomText = strings.Replace(bottomText, m.Mentions[0].Mention(), "", -1)
	}

	// Look at prev messages if no url is given
	if url == "" {
		messages, err := s.ChannelMessages(m.ChannelID, -1, "", "", "")
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, "Error fetching messages.")
			return
		}

		for _, msg := range messages {
			if len(msg.Attachments) > 0 {
				url = msg.Attachments[0].URL
				break
			} else if len(msg.Embeds) > 0 && msg.Embeds[0].Image != nil {
				url = msg.Embeds[0].Image.URL
				break
			} else if linkRegex.MatchString(msg.Content) {
				url = linkRegex.FindStringSubmatch(msg.Content)[0]
				break
			}
		}
		if url == "" {
			s.ChannelMessageSend(m.ChannelID, "No link/image given.")
			return
		}
	}

	// Replace special characters for bottom text
	bottomText = strings.Replace(bottomText, "?", "~q", -1)
	bottomText = strings.Replace(bottomText, "%", "~p", -1)
	bottomText = strings.Replace(bottomText, "#", "~h", -1)
	bottomText = strings.Replace(bottomText, "/", "~s", -1)
	bottomText = strings.Replace(bottomText, "\"", "''", -1)
	bottomText = strings.Replace(bottomText, "_", "__", -1)
	bottomText = strings.Replace(bottomText, "-", "--", -1)
	bottomText = strings.TrimSpace(bottomText)

	// Replace special characters for top text
	topText = strings.Replace(topText, "?", "~q", -1)
	topText = strings.Replace(topText, "%", "~p", -1)
	topText = strings.Replace(topText, "#", "~h", -1)
	topText = strings.Replace(topText, "/", "~s", -1)
	topText = strings.Replace(topText, "\"", "''", -1)
	topText = strings.Replace(topText, "_", "__", -1)
	topText = strings.Replace(topText, "-", "--", -1)
	topText = strings.TrimSpace(topText)

	if topText == "" {
		topText = "_"
	}

	// Fetch the image data
	msg, err := s.ChannelMessageSend(m.ChannelID, "Generating meme...")
	if err != nil {
		return
	}
	response, err := http.Get("http://memegen.link/custom/" + topText + "/" + bottomText + ".jpg?alt=" + url + "&watermark=none")
	s.ChannelMessageDelete(msg.ChannelID, msg.ID)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "Could not reach URL.")
		return
	}
	s.ChannelMessageSendComplex(m.ChannelID, &discordgo.MessageSend{
		Files: []*discordgo.File{
			&discordgo.File{
				Name:   "image.png",
				Reader: response.Body,
			},
		},
	})
	response.Body.Close()
	return
}
