package gencommands

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/hex"
	"image"
	"image/png"
	"net/http"
	"regexp"
	"strings"

	"github.com/bwmarrin/discordgo"
)

// Decrypt gives a string of a hash
func Decrypt(s *discordgo.Session, m *discordgo.MessageCreate) {
	decryptRegex, _ := regexp.Compile(`decrypt\s+(.+)`)
	linkRegex, _ := regexp.Compile(`https?:\/\/\S+`)
	keyRegex, _ := regexp.Compile(`-k\s+(.+)`)

	if !decryptRegex.MatchString(m.Content) {
		s.ChannelMessageSend(m.ChannelID, "Nothing to decrypt!")
		return
	}
	key := []byte("Nxb]^NSc;L*qn3K(/tN{6N7%4n32fF#@")
	text, err := hex.DecodeString(decryptRegex.FindStringSubmatch(m.Content)[1])
	if err != nil {
		if keyRegex.MatchString(m.Content) {
			match := keyRegex.FindStringSubmatch(m.Content)
			key = []byte(match[1])
			text, err = hex.DecodeString(strings.TrimSpace(strings.Replace(decryptRegex.FindStringSubmatch(m.Content)[1], match[0], "", 1)))
			if err != nil {
				s.ChannelMessageSend(m.ChannelID, "Invalid value!")
				return
			}
		} else {
			s.ChannelMessageSend(m.ChannelID, "Invalid value!")
			return
		}
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "Key must be a length of 16, 24, or 32 characters!")
		return
	}
	gcm, _ := cipher.NewGCM(block)
	nonceSize := gcm.NonceSize()
	if len(text) < nonceSize {
		s.ChannelMessageSend(m.ChannelID, "Invalid value!")
		return
	}
	nonce, ciphertext := text[:nonceSize], text[nonceSize:]
	result, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "Invalid value!")
		return
	}
	resultText := string(result)
	if linkRegex.MatchString(resultText) {
		response, err := http.Get(linkRegex.FindStringSubmatch(resultText)[0])
		defer response.Body.Close()
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, "```"+resultText+"```")
			return
		}
		img, _, err := image.Decode(response.Body)
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, "```"+resultText+"```")
			return
		}
		imgBytes := new(bytes.Buffer)
		err = png.Encode(imgBytes, img)
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, "```"+resultText+"```")
			return
		}
		_, err = s.ChannelMessageSendComplex(m.ChannelID, &discordgo.MessageSend{
			Files: []*discordgo.File{
				&discordgo.File{
					Name:   "image.png",
					Reader: imgBytes,
				},
			},
		})
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, "```"+resultText+"```")
		}
		return
	}
	s.ChannelMessageSend(m.ChannelID, "```"+resultText+"```")
	return
}
