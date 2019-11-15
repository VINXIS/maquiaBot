package tools

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"math/big"
	"regexp"
	"strconv"

	"github.com/bwmarrin/discordgo"
)

// FunnyData gets teh data for funny values
func FunnyData(s *discordgo.Session) {
	funnyData := []UserFunny{}
	funnyValues := []float64{}
	guilds := s.State.Guilds
	emojiRegex, _ := regexp.Compile(`<(:.+:)\d+>`)

	fmt.Println("Collecting average funny data...")

	for i, guild := range guilds {
		fmt.Println("Running guild #" + strconv.Itoa(i+1))
		guild, err := s.Guild(guild.ID)
		if err == nil {
			for j, channel := range guild.Channels {
				fmt.Println("Running channel #" + strconv.Itoa(j+1))
				if channel.Type == discordgo.ChannelTypeGuildText {
					messages, err := s.ChannelMessages(channel.ID, 100, "", "", "")
					if err == nil {
						for _, message := range messages {
							exists := false
							for k, user := range funnyData {
								if message.Author.ID == user.User.ID {
									funnyData[k].UserMessages = append(funnyData[k].UserMessages, *message)
									exists = true
								}
							}
							if !exists {
								funnyData = append(funnyData, UserFunny{
									User:         *message.Author,
									UserMessages: []discordgo.Message{*message},
								})
							}
						}
					}
				}
				fmt.Println("Finished channel #" + strconv.Itoa(j+1))
			}
		}
		fmt.Println("Finished guild #" + strconv.Itoa(i+1))
	}
	fmt.Println("Finished all guilds! Obtaining funny values")

	for i, user := range funnyData {
		fmt.Println("Running user #" + strconv.Itoa(i+1))
		size := len(user.UserMessages)
		lengthValue := big.NewInt(1)
		levenshteinValue := 0.0
		for i, message := range user.UserMessages {
			messageLevenVal := 0.0
			messageLevenSize := 0.0
			for j, message1 := range user.UserMessages {
				if i == j {
					break
				}
				if message.ID != message1.ID {
					if emojiRegex.MatchString(message.Content) {
						message.Content = emojiRegex.ReplaceAllString(message.Content, emojiRegex.FindStringSubmatch(message.Content)[1])
					}
					if emojiRegex.MatchString(message1.Content) {
						message1.Content = emojiRegex.ReplaceAllString(message1.Content, emojiRegex.FindStringSubmatch(message1.Content)[1])
					}
					messageLevenVal += Levenshtein(message.Content, message1.Content)
					messageLevenSize++
				}
			}
			levenshteinValue += messageLevenVal / math.Max(1.0, float64(messageLevenSize))
			lengthValue.Mul(lengthValue, big.NewInt(int64(len(message.Content))))
		}
		lengthValue.Exp(lengthValue, big.NewInt(int64(1.0/float64(size))), nil)
		lengthVal := float64(lengthValue.Int64())
		funnyValues = append(funnyValues, math.Sqrt(lengthVal*levenshteinValue/float64(size)))
		fmt.Println("Finished user #" + strconv.Itoa(i+1))
	}

	fmt.Println("Finished funny values! Obtaining mean and stddev")

	total := 0.00
	size := len(funnyValues)
	for _, val := range funnyValues {
		total += val
	}
	mean := total / float64(size)
	stddev := 0.00
	total = 0.00

	for _, val := range funnyValues {
		total += math.Pow(val-mean, 2.0)
	}
	stddev = math.Sqrt(total / float64(size-1))
	fmt.Println("Mean: " + strconv.FormatFloat(mean, 'f', 16, 64))
	fmt.Println("Standard Deviation: " + strconv.FormatFloat(stddev, 'f', 16, 64))

	jsonCache, err := json.Marshal(funnyData)
	ErrRead(err)

	err = ioutil.WriteFile("./data/funny.json", jsonCache, 0644)
	ErrRead(err)
}

// UserFunny contains data about the user's messages
type UserFunny struct {
	User         discordgo.User
	UserMessages []discordgo.Message
}
