package gencommands

import (
	"math"
	"math/big"
	"regexp"
	"sort"
	"strconv"
	"strings"

	tools "../../tools"
	"github.com/bwmarrin/discordgo"
)

// Funny determines how funny u are
func Funny(s *discordgo.Session, m *discordgo.MessageCreate) {
	userRegex, _ := regexp.Compile(`funny\s+(.+)`)
	emojiRegex, _ := regexp.Compile(`<(:.+:)\d+>`)
	userMessages := []*discordgo.Message{}
	users := m.Mentions
	user := m.Author.ID
	username := "Your"
	if len(users) == 1 {
		user = users[0].ID
		username = users[0].Username + "'s"
	} else if userRegex.MatchString(m.Content) {
		userTest := userRegex.FindStringSubmatch(m.Content)[1]
		members, err := s.GuildMembers(m.GuildID, "", 1000)
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, "This is not a server! Use their ID directly instead.")
			return
		}
		sort.Slice(members, func(i, j int) bool {
			time1, _ := members[i].JoinedAt.Parse()
			time2, _ := members[j].JoinedAt.Parse()
			return time1.Unix() < time2.Unix()
		})
		for _, member := range members {
			if strings.HasPrefix(strings.ToLower(member.User.Username), userTest) || strings.HasPrefix(strings.ToLower(member.Nick), userTest) {
				user = member.User.ID
				username = member.User.Username + "'s"
			}
		}
	}

	guild, err := s.Guild(m.GuildID)
	if err != nil {
		messages, err := s.ChannelMessages(m.ChannelID, 100, "", "", "")
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, "Error fetching messages.")
			return
		}
		for _, msg := range messages {
			if msg.Author.ID == user {
				userMessages = append(userMessages, msg)
			}
		}
	} else {
		for _, textchannel := range guild.Channels {
			if textchannel.Type == discordgo.ChannelTypeGuildText {
				messages, err := s.ChannelMessages(textchannel.ID, 100, "", "", "")
				if err == nil {
					for _, msg := range messages {
						if msg.Author.ID == user {
							userMessages = append(userMessages, msg)
						}
					}
				}
			}
		}
	}

	size := len(userMessages)
	totalLength := big.NewInt(1)
	totalLeven := 0.0
	totalFunny := 0.0
	for i, msg := range userMessages {
		messageLevenVal := 0.0
		messageLevenSize := 0.0
		for j, msg2 := range userMessages {
			if i == j {
				break
			}
			if msg.ID != msg2.ID {
				if emojiRegex.MatchString(msg.Content) {
					msg.Content = emojiRegex.ReplaceAllString(msg.Content, emojiRegex.FindStringSubmatch(msg.Content)[1])
				}
				if emojiRegex.MatchString(msg2.Content) {
					msg2.Content = emojiRegex.ReplaceAllString(msg2.Content, emojiRegex.FindStringSubmatch(msg2.Content)[1])
				}
				messageLevenVal += tools.Levenshtein(msg.Content, msg2.Content) - math.Abs(float64(len(msg.Content)-len(msg2.Content)))
				messageLevenSize++
			}
		}
		totalLeven += messageLevenVal / math.Max(1.0, float64(messageLevenSize))
		totalLength.Mul(totalLength, big.NewInt(int64(len(msg.Content))))
	}
	if size > 0 {
		totalLength.Exp(totalLength, big.NewInt(int64(1.0/float64(size))), nil)
		lengthVal := float64(totalLength.Int64())
		totalFunny = math.Sqrt(lengthVal * totalLeven / float64(size))
	}

	average := 2.3457696791197051
	stddev := 1.5711786950407161

	percentile := 100 * 0.5 * math.Erfc((average-float64(totalFunny))/(math.Sqrt(2.0)*stddev))

	s.ChannelMessageSend(m.ChannelID, username+" funny value is "+strconv.FormatFloat(totalFunny, 'f', 2, 64)+" which is approximately funnier than "+strconv.FormatFloat(percentile, 'f', 2, 64)+"% of people on discord.")
}
