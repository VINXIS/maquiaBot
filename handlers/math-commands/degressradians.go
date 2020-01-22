package mathcommands

import (
	"math"
	"regexp"
	"strconv"

	"github.com/bwmarrin/discordgo"
)

// DegreesRadians converts from deg to rad
func DegreesRadians(s *discordgo.Session, m *discordgo.MessageCreate) {
	reg, _ := regexp.Compile(`(dr|degrad|degreesradians)\s+((\d|\.)+)`)
	if !reg.MatchString(m.Content) {
		s.ChannelMessageSend(m.ChannelID, "Please send a valid number!")
		return
	}

	val, err := strconv.ParseFloat(reg.FindStringSubmatch(m.Content)[1], 10)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "Please send a valid number!")
		return
	}

	val *= math.Pi / 180.0

	s.ChannelMessageSend(m.ChannelID, strconv.FormatFloat(val, 'f', 6, 64))
}

// RadiansDegrees converts from rad to deg
func RadiansDegrees(s *discordgo.Session, m *discordgo.MessageCreate) {
	reg, _ := regexp.Compile(`(rd|raddeg|radiansdegrees)\s+((\d|\.)+)`)
	if !reg.MatchString(m.Content) {
		s.ChannelMessageSend(m.ChannelID, "Please send a valid number!")
		return
	}

	val, err := strconv.ParseFloat(reg.FindStringSubmatch(m.Content)[1], 10)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "Please send a valid number!")
		return
	}

	val *= 180.0 / math.Pi

	s.ChannelMessageSend(m.ChannelID, strconv.FormatFloat(val, 'f', 6, 64))
}
