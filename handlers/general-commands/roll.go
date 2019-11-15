package gencommands

import (
	"crypto/rand"
	"math/big"
	"regexp"
	"strconv"

	"github.com/bwmarrin/discordgo"
)

// Roll gives a number between 1 to n
func Roll(s *discordgo.Session, m *discordgo.MessageCreate) {
	rollRegex, _ := regexp.Compile(`roll\s*(\S+)`)
	number := 100
	var err error
	if rollRegex.MatchString(m.Content) {
		number, err = strconv.Atoi(rollRegex.FindStringSubmatch(m.Content)[1])
		if err != nil {
			number = 100
		}
		if number <= 0 {
			s.ChannelMessageSend(m.ChannelID, "Give a number >=1 mate.")
			return
		}
	}
	roll, _ := rand.Int(rand.Reader, big.NewInt(int64(number)))
	s.ChannelMessageSend(m.ChannelID, strconv.Itoa(int(roll.Int64())+1))
}
