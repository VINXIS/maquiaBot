package gencommands

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/hex"
	"regexp"
	"strings"

	"github.com/bwmarrin/discordgo"
)

// Encrypt gives a hash of a string
func Encrypt(s *discordgo.Session, m *discordgo.MessageCreate) {
	encryptRegex, _ := regexp.Compile(`encrypt\s+(.+)`)
	keyRegex, _ := regexp.Compile(`-k\s+(.+)`)

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
