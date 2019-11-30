package gencommands

import (
	"math"
	"math/rand"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

// Penis gives u a penis size for the day based off of a normal distribution from data obtained at http://penissizes.org/average-penis-size-ethnicity-race-and-country
func Penis(s *discordgo.Session, m *discordgo.MessageCreate) {
	userRegex, _ := regexp.Compile(`penis\s+(.+)`)

	users := m.Mentions
	user := m.Author.ID
	username := "Your"
	if len(users) == 1 {
		user = users[0].ID
		username = users[0].Username + "'s"
	} else if userRegex.MatchString(m.Content) {
		userTest := userRegex.FindStringSubmatch(m.Content)[1]
		discordUser, err := s.User(userTest)
		if err == nil {
			user = discordUser.ID
			username = discordUser.Username + "'s"
		} else {
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
	}

	year, month, day := time.Now().Date()
	authorid, _ := strconv.Atoi(user)
	random := rand.New(rand.NewSource(int64(authorid + day + int(month) + year)))

	average := 13.91
	stddev := 2.20

	penisSize := random.NormFloat64()*stddev + average

	percentile := 100 * 0.5 * math.Erfc((average-penisSize)/(math.Sqrt(2.0)*stddev))

	s.ChannelMessageSend(m.ChannelID, username+" erect size for the day is "+strconv.FormatFloat(penisSize, 'f', 2, 64)+"cm ("+strconv.FormatFloat(penisSize/2.54, 'f', 2, 64)+"in) which is larger than approximately "+strconv.FormatFloat(percentile, 'f', 2, 64)+"% of penises.")
}
