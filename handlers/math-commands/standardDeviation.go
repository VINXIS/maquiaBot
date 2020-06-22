package mathcommands

import (
	"math"
	"strconv"
	"strings"

	mathtools "maquiaBot/math-tools"
	"github.com/bwmarrin/discordgo"
)

// StandardDeviation gives the apopulation and sample standard deviation of a list of values
func StandardDeviation(s *discordgo.Session, m *discordgo.MessageCreate) {
	nums := strings.Split(m.Content, " ")[1:]
	var actualNums []float64

	// Get numbers
	for _, num := range nums {
		actualNum, err := strconv.ParseFloat(num, 64)
		if err == nil {
			actualNums = append(actualNums, actualNum)
		}
	}

	// Check for more than 1 number
	if len(actualNums) <= 1 {
		s.ChannelMessageSend(m.ChannelID, "Please provide 2 or more numbers to get the average for!")
	}

	stddev := mathtools.StandardDeviation(actualNums, false)
	stddevSample := mathtools.StandardDeviation(actualNums, true)

	text := "Standard deviation for: "
	for _, num := range actualNums {
		text += strconv.FormatFloat(num, 'f', 2, 64) + ", "
	}
	text = strings.TrimSuffix(text, ", ") + "\n\n"

	text += "Population: **" + strconv.FormatFloat(stddev, 'f', 2, 64) + "** σ | **" + strconv.FormatFloat(math.Pow(stddev, 2.0), 'f', 2, 64) + "** σ²\n" +
		"Sample: **" + strconv.FormatFloat(stddevSample, 'f', 2, 64) + "** σ | **" + strconv.FormatFloat(math.Pow(stddevSample, 2.0), 'f', 2, 64) + "** σ²\n"

	s.ChannelMessageSend(m.ChannelID, text)
}
