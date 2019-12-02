package gencommands

import (
	"regexp"
	"strconv"
	"time"

	"github.com/bwmarrin/discordgo"
)

// Parse parses the discord snowflake ID given
func Parse(s *discordgo.Session, m *discordgo.MessageCreate) {
	parseRegex, _ := regexp.Compile(`parse\s+(\d+)`)

	// Get snowflake value to test
	snowflake := m.Author.ID
	if parseRegex.MatchString(m.Content) {
		snowflake = parseRegex.FindStringSubmatch(m.Content)[1]
	}
	snowflakeInt, err := strconv.ParseInt(snowflake, 10, 64)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "Invalid value!")
		return
	}

	// Parse snowflake
	timeStamp, _ := discordgo.SnowflakeTimestamp(snowflake)
	intWorkerID := (snowflakeInt & 0x3E0000) >> 17
	intProcessID := (snowflakeInt & 0x1F000) >> 12
	Increment := snowflakeInt & 0xFFF

	// Parse twitter snowflake
	timestamp := (snowflakeInt >> 22) + 1288834974657
	timeStamp2 := time.Unix(timestamp/1000, 0)

	s.ChannelMessageSend(m.ChannelID, "Discord Timestamp: "+timeStamp.UTC().Format(time.RFC822Z)+"\nTwitter Timestamp: "+timeStamp2.UTC().Format(time.RFC822Z)+"\nInternal worker ID: "+strconv.FormatInt(intWorkerID, 10)+"\nInternal process ID: "+strconv.FormatInt(intProcessID, 10)+"\nIncrement #: "+strconv.FormatInt(Increment, 10))
}
