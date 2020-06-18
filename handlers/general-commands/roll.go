package gencommands

import (
	"math/rand"
	"regexp"
	"strconv"
	"time"

	"github.com/bwmarrin/discordgo"
)

// Roll gives a number between 1 to n
func Roll(s *discordgo.Session, m *discordgo.MessageCreate) {
	rollRegex, _ := regexp.Compile(`(?i)roll\s*(\S+)`)
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
	authorid, _ := strconv.Atoi(m.Author.ID)
	random := rand.New(rand.NewSource(int64(authorid) + time.Now().UnixNano()))
	s.ChannelMessageSend(m.ChannelID, strconv.Itoa(random.Intn(number)+1))
}
