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

//Penis gives u a penis size for the day based off of a normal distribution from data obtained at http://penissizes.org/average-penis-size-ethnicity-race-and-country
func Penis(s *discordgo.Session, m *discordgo.MessageCreate) {
	userRegex, _ := regexp.Compile(`penis\s+(.+)`)

	users := m.Mentions
	user := m.Author.ID
	username := ""
	if len(users) == 1 {
		user = users[0].ID
		username = users[0].Username
	} else if userRegex.MatchString(m.Content) {
		userTest := userRegex.FindStringSubmatch(m.Content)[1]
		discordUser, err := s.User(userTest)
		if err == nil {
			user = discordUser.ID
			username = discordUser.Username
		} else {
			server, err := s.Guild(m.GuildID)
			if err != nil {
				s.ChannelMessageSend(m.ChannelID, "This is not a server! Use their ID directly instead.")
				return
			}
			sort.Slice(server.Members, func(i, j int) bool {
				time1, _ := server.Members[i].JoinedAt.Parse()
				time2, _ := server.Members[j].JoinedAt.Parse()
				return time1.Unix() < time2.Unix()
			})
			for _, member := range server.Members {
				if strings.HasPrefix(strings.ToLower(member.User.Username), userTest) {
					discordUser, _ = s.User(member.User.ID)
					user = discordUser.ID
					username = discordUser.Username
				}
			}
			for _, member := range server.Members {
				if strings.HasPrefix(strings.ToLower(member.Nick), userTest) {
					discordUser, _ := s.User(member.User.ID)
					user = discordUser.ID
					username = discordUser.Username
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

	if user == m.Author.ID {
		s.ChannelMessageSend(m.ChannelID, "Your erect size for the day is "+strconv.FormatFloat(penisSize, 'f', 2, 64)+"cm ("+strconv.FormatFloat(penisSize/2.54, 'f', 2, 64)+"in) which is larger than "+strconv.FormatFloat(percentile, 'f', 2, 64)+"% of penises.")
	} else {
		s.ChannelMessageSend(m.ChannelID, username+"'s erect size for the day is "+strconv.FormatFloat(penisSize, 'f', 2, 64)+"cm ("+strconv.FormatFloat(penisSize/2.54, 'f', 2, 64)+"in) which is larger than "+strconv.FormatFloat(percentile, 'f', 2, 64)+"% of penises.")
	}
}
