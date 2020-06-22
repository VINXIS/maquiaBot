package mathcommands

import (
	"strconv"
	"strings"

	mathtools "maquiaBot/math-tools"
	"github.com/bwmarrin/discordgo"
)

// Average gives the average of a list of numbers
func Average(s *discordgo.Session, m *discordgo.MessageCreate) {
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

	harmonicMean := mathtools.HarmonicMean(actualNums)
	geometricMean := mathtools.GeometricMean(actualNums)
	arithmeticMean := mathtools.ArithmeticMean(actualNums)

	text := "Averages for: "
	for _, num := range actualNums {
		text += strconv.FormatFloat(num, 'f', 2, 64) + ", "
	}
	text = strings.TrimSuffix(text, ", ") + "\n\n"

	text += "Arithmetic Mean: **" + strconv.FormatFloat(arithmeticMean, 'f', 2, 64) + "**\n" +
		"Geometric Mean: **" + strconv.FormatFloat(geometricMean, 'f', 2, 64) + "**\n" +
		"Harmonic Mean: **" + strconv.FormatFloat(harmonicMean, 'f', 2, 64) + "**\n"

	s.ChannelMessageSend(m.ChannelID, text)
}
