package gencommands

import (
	"encoding/json"
	"io/ioutil"
	"math"
	"math/rand"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	tools "../../tools"
	"github.com/bwmarrin/discordgo"
)

type penis struct {
	member     *discordgo.Member
	size       float64
	percentile float64
}

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

	year, month, day := time.Now().UTC().Date()
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
	penisRecords(s, m, penisSize, user)
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

	year, month, day := time.Now().UTC().Date()
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
	penisRecords(s, m, penisSize1, user1)
}

// PenisRank displays the largest / smallest penises in the server
func PenisRank(s *discordgo.Session, m *discordgo.MessageCreate) {
	rankRegex, _ := regexp.Compile(`(rc|rp|rankc|rankp|rankcock|rankpenis)\s+(\d+)`)
	smallRegex, _ := regexp.Compile(`-s`)
	// Get members
	members, err := s.GuildMembers(m.GuildID, "", 1000)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "This is not a server!")
		return
	}

	// Get number of people to show
	num := 1
	if rankRegex.MatchString(m.Content) {
		num, err = strconv.Atoi(rankRegex.FindStringSubmatch(m.Content)[2])
		if err != nil {
			num = 1
		}
		if num > len(members) {
			num = len(members)
		}
	}

	penisSizes := []penis{}
	average := 13.91
	stddev := 2.20
	year, month, day := time.Now().UTC().Date()

	// Get average server size, member sizes, and percentiles
	var avgSize float64
	for _, member := range members {
		authorid, _ := strconv.Atoi(member.User.ID)
		random := rand.New(rand.NewSource(int64(authorid + day + int(month) + year)))

		penisSize := random.NormFloat64()*stddev + average
		percentile := 100 * 0.5 * math.Erfc((average-penisSize)/(math.Sqrt(2.0)*stddev))

		penisSizes = append(penisSizes, penis{
			member:     member,
			size:       penisSize,
			percentile: percentile,
		})
		avgSize += penisSize
	}
	avgSize /= float64(len(members))

	// Sort and obtain wanted amount
	var text string
	if smallRegex.MatchString(m.Content) {
		sort.Slice(penisSizes, func(i, j int) bool { return penisSizes[i].size < penisSizes[j].size })
		if num <= 1 {
			emote := ""
			if penisSizes[0].percentile < 25 {
				emote = ":pinching_hand:"
			} else if penisSizes[0].percentile > 75 {
				emote = ":eggplant: WTF"
			}
			s.ChannelMessageSend(m.ChannelID, "**"+penisSizes[0].member.User.Username+"'s** erect size for the day is "+strconv.FormatFloat(penisSizes[0].size, 'f', 2, 64)+"cm ("+strconv.FormatFloat(penisSizes[0].size/2.54, 'f', 2, 64)+"in) which is larger than approximately "+strconv.FormatFloat(penisSizes[0].percentile, 'f', 2, 64)+"% of penises. Their size is the smallest in this server today! "+emote+"\n**Average size in the server:** "+strconv.FormatFloat(avgSize, 'f', 2, 64)+"cm ("+strconv.FormatFloat(avgSize/2.54, 'f', 2, 64)+"in)")
			penisRecords(s, m, penisSizes[0].size, penisSizes[0].member.User.ID)
			return
		}

		text = "Smallest **" + strconv.Itoa(num) + "** sizes in this server: \n"
	} else {
		sort.Slice(penisSizes, func(i, j int) bool { return penisSizes[i].size > penisSizes[j].size })
		if num <= 1 {
			emote := ""
			if penisSizes[0].percentile < 25 {
				emote = ":pinching_hand: LOLL"
			} else if penisSizes[0].percentile > 75 {
				emote = ":eggplant:"
			}
			s.ChannelMessageSend(m.ChannelID, "**"+penisSizes[0].member.User.Username+"'s** erect size for the day is "+strconv.FormatFloat(penisSizes[0].size, 'f', 2, 64)+"cm ("+strconv.FormatFloat(penisSizes[0].size/2.54, 'f', 2, 64)+"in) which is larger than approximately "+strconv.FormatFloat(penisSizes[0].percentile, 'f', 2, 64)+"% of penises. Their size is the largest in this server today! "+emote+"\n**Average size in the server:** "+strconv.FormatFloat(avgSize, 'f', 2, 64)+"cm ("+strconv.FormatFloat(avgSize/2.54, 'f', 2, 64)+"in)")
			penisRecords(s, m, penisSizes[0].size, penisSizes[0].member.User.ID)
			return
		}

		text = "Largest **" + strconv.Itoa(num) + "** sizes in this server: \n"
	}
	for i := 0; i < num; i++ {
		emote := ""
		if penisSizes[i].percentile < 25 {
			emote = ":pinching_hand:"
		} else if penisSizes[i].percentile > 75 {
			emote = ":eggplant:"
		}
		text += "**" + penisSizes[i].member.User.Username + ":** " + strconv.FormatFloat(penisSizes[i].size, 'f', 2, 64) + "cm (" + strconv.FormatFloat(penisSizes[i].size/2.54, 'f', 2, 64) + "in) " + emote + "\n"
	}
	s.ChannelMessageSend(m.ChannelID, text+"**Average size in the server:** "+strconv.FormatFloat(avgSize, 'f', 2, 64)+"cm ("+strconv.FormatFloat(avgSize/2.54, 'f', 2, 64)+"in)")
	penisRecords(s, m, penisSizes[0].size, penisSizes[0].member.User.ID)
}

// PenisHistory shows the largest and smallest penis sizes ever recorded
func PenisHistory(s *discordgo.Session, m *discordgo.MessageCreate) {
	serverRegex, _ := regexp.Compile(`-s`)
	average := 13.91
	stddev := 2.20

	penisRecord := tools.GetPenisRecord()
	text := "Records for all servers:\n"
	largestUsername := ""
	smallestUsername := ""

	// Check server flag
	if serverRegex.MatchString(m.Content) {
		text = "Records for this server:\n"
		server, err := s.Guild(m.GuildID)
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, "This is not a server!")
			return
		}
		serverData := tools.GetServer(*server)
		penisRecord.Largest = serverData.Largest
		penisRecord.Smallest = serverData.Smallest
	}

	// Get usernames
	largestUser, err := s.User(penisRecord.Largest.UserID)
	if err != nil {
		largestUsername = "Unknown user (" + penisRecord.Largest.UserID + ")"
	} else {
		largestUsername = largestUser.Username
	}

	smallestUser, err := s.User(penisRecord.Smallest.UserID)
	if err != nil {
		smallestUsername = "Unknown user (" + penisRecord.Smallest.UserID + ")"
	} else {
		smallestUsername = smallestUser.Username
	}

	largestPercentile := 100 * 0.5 * math.Erfc((average-penisRecord.Largest.Size)/(math.Sqrt(2.0)*stddev))
	smallestPercentile := 100 * 0.5 * math.Erfc((average-penisRecord.Smallest.Size)/(math.Sqrt(2.0)*stddev))
	s.ChannelMessageSend(m.ChannelID, text+
		"**" + largestUsername+"** on "+strings.Replace(penisRecord.Largest.Date.Format(time.RFC822Z), " +0000", "UTC", -1)+": "+strconv.FormatFloat(penisRecord.Largest.Size, 'f', 2, 64)+"cm ("+strconv.FormatFloat(penisRecord.Largest.Size/2.54, 'f', 2, 64)+"in) larger than approximately "+strconv.FormatFloat(largestPercentile, 'f', 2, 64)+"% of penises.\n"+
		"**" + smallestUsername+"** on "+strings.Replace(penisRecord.Smallest.Date.Format(time.RFC822Z), " +0000", "UTC", -1)+": "+strconv.FormatFloat(penisRecord.Smallest.Size, 'f', 2, 64)+"cm ("+strconv.FormatFloat(penisRecord.Smallest.Size/2.54, 'f', 2, 64)+"in) larger than approximately "+strconv.FormatFloat(smallestPercentile, 'f', 2, 64)+"% of penises.")
}

func penisRecords(s *discordgo.Session, m *discordgo.MessageCreate, penisSize float64, userID string) {
	// Check full records
	penisRecords := tools.GetPenisRecord()
	recordBroken := false

	if penisSize > penisRecords.Largest.Size {
		penisRecords.Largest.Size = penisSize
		penisRecords.Largest.UserID = userID
		penisRecords.Largest.Date = time.Now().UTC()
		recordBroken = true
	} else if penisSize < penisRecords.Smallest.Size {
		penisRecords.Smallest.Size = penisSize
		penisRecords.Smallest.UserID = userID
		penisRecords.Smallest.Date = time.Now().UTC()
		recordBroken = true
	}

	if recordBroken {
		jsonCache, err := json.Marshal(penisRecords)
		tools.ErrRead(err)

		err = ioutil.WriteFile("./data/penisRecords.json", jsonCache, 0644)
		tools.ErrRead(err)
	}

	// Check server records
	server, err := s.Guild(m.GuildID)
	if err != nil {
		return
	}
	serverData := tools.GetServer(*server)
	recordBroken = false

	if penisSize > serverData.Largest.Size {
		serverData.Largest.Size = penisSize
		serverData.Largest.UserID = userID
		serverData.Largest.Date = time.Now().UTC()
		recordBroken = true
	} else if penisSize < serverData.Smallest.Size {
		serverData.Smallest.Size = penisSize
		serverData.Smallest.UserID = userID
		serverData.Smallest.Date = time.Now().UTC()
		recordBroken = true
	}

	if recordBroken {
		jsonCache, err := json.Marshal(serverData)
		tools.ErrRead(err)

		err = ioutil.WriteFile("./data/serverData/"+m.GuildID+".json", jsonCache, 0644)
		tools.ErrRead(err)
	}
}
