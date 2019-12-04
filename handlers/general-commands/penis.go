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
	userRegex, _ := regexp.Compile(`(penis|cock)\s+(.+)`)

	users := m.Mentions
	user := m.Author.ID
	username := "Your"
	if len(users) == 1 {
		user = users[0].ID
		username = users[0].Username + "'s"
	} else if userRegex.MatchString(m.Content) {
		userTest := userRegex.FindStringSubmatch(m.Content)[2]
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
	emote := ""
	if percentile < 25 {
		emote = ":pinching_hand:"
	} else if percentile > 75 {
		emote = ":eggplant:"
	}
	s.ChannelMessageSend(m.ChannelID, username+" erect size for the day is "+strconv.FormatFloat(penisSize, 'f', 2, 64)+"cm ("+strconv.FormatFloat(penisSize/2.54, 'f', 2, 64)+"in) which is larger than approximately "+strconv.FormatFloat(percentile, 'f', 2, 64)+"% of penises. "+emote)
}

// PenisCompare compares ur penis size to someone else's
func PenisCompare(s *discordgo.Session, m *discordgo.MessageCreate) {
	userRegex, _ := regexp.Compile(`(cc|cp|comparec|comparep|comparecock|comparepenis)\s+(.+)`)
	penisRegex, _ := regexp.Compile(`(penis|cock|cc|cp|comparec|comparep|comparecock|comparepenis)`)

	users := m.Mentions
	user1 := m.Author.ID
	user2 := m.Author.ID
	user2name := ""
	if len(users) == 1 { // Mention
		user2 = users[0].ID
		user2name = users[0].Username
	} else if userRegex.MatchString(m.Content) { // Name / ID written
		userTest := userRegex.FindStringSubmatch(m.Content)[2]
		discordUser, err := s.User(userTest)
		if err == nil { // ID found
			user2 = discordUser.ID
			user2name = discordUser.Username
		} else { // Find member if no ID found
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
					user2 = member.User.ID
					user2name = member.User.Username
				}
			}
		}
	} else { // Compare with latest !penis
		// Get prev messages
		messages, _ := s.ChannelMessages(m.ChannelID, -1, "", "", "")
		for i := 0; i < len(messages)-1; i++ {
			if messages[i].Author.ID == s.State.User.ID && penisRegex.MatchString(messages[i+1].Content) && messages[i+1].Author.ID != m.Author.ID {
				user2 = messages[i+1].Author.ID
				user2name = messages[i+1].Author.Username
				break
			}
		}
	}
	if user2name == "" {
		s.ChannelMessageSend(m.ChannelID, "Could not find anyone to compare to!")
		return
	}

	year, month, day := time.Now().Date()
	id1, _ := strconv.Atoi(user1)
	random1 := rand.New(rand.NewSource(int64(id1 + day + int(month) + year)))
	id2, _ := strconv.Atoi(user2)
	random2 := rand.New(rand.NewSource(int64(id2 + day + int(month) + year)))

	average := 13.91
	stddev := 2.20

	penisSize1 := random1.NormFloat64()*stddev + average
	penisSize2 := random2.NormFloat64()*stddev + average

	percentile1 := 100 * 0.5 * math.Erfc((average-penisSize1)/(math.Sqrt(2.0)*stddev))
	percentile2 := 100 * 0.5 * math.Erfc((average-penisSize2)/(math.Sqrt(2.0)*stddev))

	mainText := "**Your** erect size: " + strconv.FormatFloat(penisSize1, 'f', 2, 64) + "cm (" + strconv.FormatFloat(penisSize1/2.54, 'f', 2, 64) + "in) larger than approximately " + strconv.FormatFloat(percentile1, 'f', 2, 64) + "%\n" +
		"**" + user2name + "'s** erect size: " + strconv.FormatFloat(penisSize2, 'f', 2, 64) + "cm (" + strconv.FormatFloat(penisSize2/2.54, 'f', 2, 64) + "in) larger than approximately " + strconv.FormatFloat(percentile2, 'f', 2, 64) + "%\n"

	if penisSize1 > penisSize2 {
		mainText += "You are fucking MASSIVE compared to **" + user2name + "!** You are " + strconv.FormatFloat(penisSize1-penisSize2, 'f', 2, 64) + "cm (" + strconv.FormatFloat((penisSize1-penisSize2)/2.54, 'f', 2, 64) + "in) larger and " + strconv.FormatFloat(penisSize1/penisSize2*100, 'f', 2, 64) + "% their size! Holy fuck!"
	} else {
		mainText += "You are fucking TINY compared to **" + user2name + "!** They are " + strconv.FormatFloat(penisSize2-penisSize1, 'f', 2, 64) + "cm (" + strconv.FormatFloat((penisSize2-penisSize1)/2.54, 'f', 2, 64) + "in) larger and " + strconv.FormatFloat(penisSize2/penisSize1*100, 'f', 2, 64) + "% your size! LOOOOOL"
	}

	s.ChannelMessageSend(m.ChannelID, mainText)
}
