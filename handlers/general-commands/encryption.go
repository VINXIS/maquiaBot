package gencommands

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/hex"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

// Encrypt gives a hash of a string
func Encrypt(s *discordgo.Session, m *discordgo.MessageCreate) {
	encryptRegex, _ := regexp.Compile(`(?i)encrypt\s+(.+)`)
	keyRegex, _ := regexp.Compile(`(?i)-k\s+(.+)`)

	text := []byte{}
	if !encryptRegex.MatchString(m.Content) {
		if len(m.Attachments) == 0 {
			s.ChannelMessageSend(m.ChannelID, "Nothing to encrypt!")
			return
		}
		text = []byte(m.Attachments[0].URL)
	} else {
		text = []byte(encryptRegex.FindStringSubmatch(m.Content)[1])
	}

	key := []byte("Nxb]^NSc;L*qn3K(/tN{6N7%4n32fF#@")
	if keyRegex.MatchString(m.Content) {
		match := keyRegex.FindStringSubmatch(m.Content)
		key = []byte(match[1])
		text = []byte(strings.Replace(string(text), match[0], "", 1))
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "Key must be a length of 16, 24, or 32 characters!")
		return
	}
	gcm, _ := cipher.NewGCM(block)

	nonce := make([]byte, gcm.NonceSize())

	result := gcm.Seal(nonce, nonce, text, nil)
	s.ChannelMessageSend(m.ChannelID, "```"+hex.EncodeToString(result)+"```")
}

// Decrypt gives a string of a hash
func Decrypt(s *discordgo.Session, m *discordgo.MessageCreate) {
	decryptRegex, _ := regexp.Compile(`(?i)decrypt\s+(.+)`)
	linkRegex, _ := regexp.Compile(`(?i)https?:\/\/\S+`)
	keyRegex, _ := regexp.Compile(`(?i)-k\s+(.+)`)

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
		client := http.Client{
			Timeout: 10 * time.Second,
		}
		response, err := client.Get(linkRegex.FindStringSubmatch(resultText)[0])
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, "```"+resultText+"```")
			response.Body.Close()
			return
		}
		defer response.Body.Close()
		_, err = s.ChannelMessageSendComplex(m.ChannelID, &discordgo.MessageSend{
			Files: []*discordgo.File{
				{
					Name:   "image.png",
					Reader: response.Body,
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
