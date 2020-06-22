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

	"github.com/bwmarrin/discordgo"
	tools "maquiaBot/tools"
)

type genitals struct {
	member     *discordgo.Member
	size       float64
	percentile float64
}

// Penis gives u a penis size for the day based off of a normal distribution from data obtained at http://penissizes.org/average-penis-size-ethnicity-race-and-country
func Penis(s *discordgo.Session, m *discordgo.MessageCreate) {
	userRegex, _ := regexp.Compile(`(?i)penis\s+(.+)`)

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
				if strings.HasPrefix(strings.ToLower(member.User.Username), strings.ToLower(userTest)) || strings.HasPrefix(strings.ToLower(member.Nick), strings.ToLower(userTest)) {
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
	records(s, m, penisSize, user, "penis")
}

// Vagina gives u a vagina size for teh day based off of a normal distribution from data obtained at https://obgyn.onlinelibrary.wiley.com/doi/full/10.1111/j.1471-0528.2004.00517.x
func Vagina(s *discordgo.Session, m *discordgo.MessageCreate) {
	userRegex, _ := regexp.Compile(`(?i)vagina\s+(.+)`)

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
				if strings.HasPrefix(strings.ToLower(member.User.Username), strings.ToLower(userTest)) || strings.HasPrefix(strings.ToLower(member.Nick), strings.ToLower(userTest)) {
					user = member.User.ID
					username = member.User.Username + "'s"
				}
			}
		}
	}

	year, month, day := time.Now().UTC().Date()
	authorid, _ := strconv.Atoi(user)
	random := rand.New(rand.NewSource(int64(authorid + day + int(month) + year + 1)))

	average := 9.6
	stddev := 1.5

	vaginaSize := random.NormFloat64()*stddev + average

	percentile := 100 * 0.5 * math.Erfc((average-vaginaSize)/(math.Sqrt(2.0)*stddev))
	emote := ""
	if percentile < 25 {
		emote = ":beach:"
	} else if percentile > 75 {
		emote = ":ocean:"
	}
	s.ChannelMessageSend(m.ChannelID, username+" depth for the day is "+strconv.FormatFloat(vaginaSize, 'f', 2, 64)+"cm ("+strconv.FormatFloat(vaginaSize/2.54, 'f', 2, 64)+"in) which is deeper than approximately "+strconv.FormatFloat(percentile, 'f', 2, 64)+"% of vaginas. "+emote)
	records(s, m, vaginaSize, user, "vagina")
}

// PenisCompare compares ur penis size to someone else's
func PenisCompare(s *discordgo.Session, m *discordgo.MessageCreate) {
	userRegex, _ := regexp.Compile(`(?i)(cp|comparep|comparepenis)\s+(.+)`)
	penisRegex, _ := regexp.Compile(`(?i)(penis|cp|comparep|comparepenis)(\s+(.+))?`)

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
				if strings.HasPrefix(strings.ToLower(member.User.Username), strings.ToLower(userTest)) || strings.HasPrefix(strings.ToLower(member.Nick), strings.ToLower(userTest)) {
					user2 = member.User.ID
					user2name = member.User.Username
				}
			}
		}
	} else { // Compare with latest !penis
		// Get prev messages
		messages, _ := s.ChannelMessages(m.ChannelID, -1, "", "", "")
		for i := 0; i < len(messages)-1; i++ {
			if messages[i].Author.ID == s.State.User.ID && penisRegex.MatchString(messages[i+1].Content) {
				if messages[i+1].Author.ID != m.Author.ID {
					user2 = messages[i+1].Author.ID
					user2name = messages[i+1].Author.Username
					break
				} else {
					userTest := penisRegex.FindStringSubmatch(messages[i+1].Content)[3]
					if userTest == "" {
						continue
					}
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
							if strings.HasPrefix(strings.ToLower(member.User.Username), strings.ToLower(userTest)) || strings.HasPrefix(strings.ToLower(member.Nick), strings.ToLower(userTest)) {
								user2 = member.User.ID
								user2name = member.User.Username
							}
						}
					}
					if user2name != "" {
						break
					}
				}
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
	records(s, m, penisSize1, user1, "penis")
}

// VaginaCompare compares ur vagina depth to someone else's
func VaginaCompare(s *discordgo.Session, m *discordgo.MessageCreate) {
	userRegex, _ := regexp.Compile(`(?i)(cv|comparev|comparevagina)\s+(.+)`)
	vaginaRegex, _ := regexp.Compile(`(?i)(vagina|cv|comparev|comparevagina)(\s+(.+))?`)

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
				if strings.HasPrefix(strings.ToLower(member.User.Username), strings.ToLower(userTest)) || strings.HasPrefix(strings.ToLower(member.Nick), strings.ToLower(userTest)) {
					user2 = member.User.ID
					user2name = member.User.Username
				}
			}
		}
	} else { // Compare with latest !vagina
		// Get prev messages
		messages, _ := s.ChannelMessages(m.ChannelID, -1, "", "", "")
		for i := 0; i < len(messages)-1; i++ {
			if messages[i].Author.ID == s.State.User.ID && vaginaRegex.MatchString(messages[i+1].Content) {
				if messages[i+1].Author.ID != m.Author.ID {
					user2 = messages[i+1].Author.ID
					user2name = messages[i+1].Author.Username
					break
				} else {
					userTest := vaginaRegex.FindStringSubmatch(messages[i+1].Content)[3]
					if userTest == "" {
						continue
					}
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
							if strings.HasPrefix(strings.ToLower(member.User.Username), strings.ToLower(userTest)) || strings.HasPrefix(strings.ToLower(member.Nick), strings.ToLower(userTest)) {
								user2 = member.User.ID
								user2name = member.User.Username
							}
						}
					}
					if user2name != "" {
						break
					}
				}
			}
		}
	}
	if user2name == "" {
		s.ChannelMessageSend(m.ChannelID, "Could not find anyone to compare to!")
		return
	}

	year, month, day := time.Now().UTC().Date()
	id1, _ := strconv.Atoi(user1)
	random1 := rand.New(rand.NewSource(int64(id1 + day + int(month) + year + 1)))
	id2, _ := strconv.Atoi(user2)
	random2 := rand.New(rand.NewSource(int64(id2 + day + int(month) + year + 1)))

	average := 9.6
	stddev := 1.5

	vaginaSize1 := random1.NormFloat64()*stddev + average
	vaginaSize2 := random2.NormFloat64()*stddev + average

	percentile1 := 100 * 0.5 * math.Erfc((average-vaginaSize1)/(math.Sqrt(2.0)*stddev))
	percentile2 := 100 * 0.5 * math.Erfc((average-vaginaSize2)/(math.Sqrt(2.0)*stddev))

	mainText := "**Your** depth: " + strconv.FormatFloat(vaginaSize1, 'f', 2, 64) + "cm (" + strconv.FormatFloat(vaginaSize1/2.54, 'f', 2, 64) + "in) deeper than approximately " + strconv.FormatFloat(percentile1, 'f', 2, 64) + "%\n" +
		"**" + user2name + "'s** depth: " + strconv.FormatFloat(vaginaSize2, 'f', 2, 64) + "cm (" + strconv.FormatFloat(vaginaSize2/2.54, 'f', 2, 64) + "in) deeper than approximately " + strconv.FormatFloat(percentile2, 'f', 2, 64) + "%\n"

	if vaginaSize1 > vaginaSize2 {
		mainText += "You are fucking DEEP compared to **" + user2name + "!** You are " + strconv.FormatFloat(vaginaSize1-vaginaSize2, 'f', 2, 64) + "cm (" + strconv.FormatFloat((vaginaSize1-vaginaSize2)/2.54, 'f', 2, 64) + "in) deeper and " + strconv.FormatFloat(vaginaSize1/vaginaSize2*100, 'f', 2, 64) + "% their depth! Holy fuck!"
	} else {
		mainText += "You are fucking SHALLOW compared to **" + user2name + "!** They are " + strconv.FormatFloat(vaginaSize2-vaginaSize1, 'f', 2, 64) + "cm (" + strconv.FormatFloat((vaginaSize2-vaginaSize1)/2.54, 'f', 2, 64) + "in) deeper and " + strconv.FormatFloat(vaginaSize2/vaginaSize1*100, 'f', 2, 64) + "% your depth! TINY BITCH!"
	}

	s.ChannelMessageSend(m.ChannelID, mainText)
	records(s, m, vaginaSize1, user1, "vagina")
}

// PenisRank displays the largest / smallest penises in the server
func PenisRank(s *discordgo.Session, m *discordgo.MessageCreate) {
	rankRegex, _ := regexp.Compile(`(?i)(rp|rankp|rankpenis)\s+(\d+)`)
	smallRegex, _ := regexp.Compile(`(?i)-s`)
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

	penisSizes := []genitals{}
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

		penisSizes = append(penisSizes, genitals{
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
			records(s, m, penisSizes[0].size, penisSizes[0].member.User.ID, "penis")
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
			records(s, m, penisSizes[0].size, penisSizes[0].member.User.ID, "penis")
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
	records(s, m, penisSizes[0].size, penisSizes[0].member.User.ID, "penis")
}

// VaginaRank displays the largest / smallest vaginas in the server
func VaginaRank(s *discordgo.Session, m *discordgo.MessageCreate) {
	rankRegex, _ := regexp.Compile(`(?i)(rv|rankv|rankvagina)\s+(\d+)`)
	smallRegex, _ := regexp.Compile(`(?i)-s`)
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

	vaginaSizes := []genitals{}
	average := 9.6
	stddev := 1.5
	year, month, day := time.Now().UTC().Date()

	// Get average server size, member sizes, and percentiles
	var avgSize float64
	for _, member := range members {
		authorid, _ := strconv.Atoi(member.User.ID)
		random := rand.New(rand.NewSource(int64(authorid + day + int(month) + year + 1)))

		vaginaSize := random.NormFloat64()*stddev + average
		percentile := 100 * 0.5 * math.Erfc((average-vaginaSize)/(math.Sqrt(2.0)*stddev))

		vaginaSizes = append(vaginaSizes, genitals{
			member:     member,
			size:       vaginaSize,
			percentile: percentile,
		})
		avgSize += vaginaSize
	}
	avgSize /= float64(len(members))

	// Sort and obtain wanted amount
	var text string
	if smallRegex.MatchString(m.Content) {
		sort.Slice(vaginaSizes, func(i, j int) bool { return vaginaSizes[i].size < vaginaSizes[j].size })
		if num <= 1 {
			emote := ""
			if vaginaSizes[0].percentile < 25 {
				emote = ":beach:"
			} else if vaginaSizes[0].percentile > 75 {
				emote = ":ocean: WTF"
			}
			s.ChannelMessageSend(m.ChannelID, "**"+vaginaSizes[0].member.User.Username+"'s** depth for the day is "+strconv.FormatFloat(vaginaSizes[0].size, 'f', 2, 64)+"cm ("+strconv.FormatFloat(vaginaSizes[0].size/2.54, 'f', 2, 64)+"in) which is deeper than approximately "+strconv.FormatFloat(vaginaSizes[0].percentile, 'f', 2, 64)+"% of vaginas. Their depth is the smallest in this server today! "+emote+"\n**Average depth in the server:** "+strconv.FormatFloat(avgSize, 'f', 2, 64)+"cm ("+strconv.FormatFloat(avgSize/2.54, 'f', 2, 64)+"in)")
			records(s, m, vaginaSizes[0].size, vaginaSizes[0].member.User.ID, "vagina")
			return
		}

		text = "Smallest **" + strconv.Itoa(num) + "** depths in this server: \n"
	} else {
		sort.Slice(vaginaSizes, func(i, j int) bool { return vaginaSizes[i].size > vaginaSizes[j].size })
		if num <= 1 {
			emote := ""
			if vaginaSizes[0].percentile < 25 {
				emote = ":beach: LOLL"
			} else if vaginaSizes[0].percentile > 75 {
				emote = ":ocean:"
			}
			s.ChannelMessageSend(m.ChannelID, "**"+vaginaSizes[0].member.User.Username+"'s** depth for the day is "+strconv.FormatFloat(vaginaSizes[0].size, 'f', 2, 64)+"cm ("+strconv.FormatFloat(vaginaSizes[0].size/2.54, 'f', 2, 64)+"in) which is deeper than approximately "+strconv.FormatFloat(vaginaSizes[0].percentile, 'f', 2, 64)+"% of vaginas. Their depth is the largest in this server today! "+emote+"\n**Average depth in the server:** "+strconv.FormatFloat(avgSize, 'f', 2, 64)+"cm ("+strconv.FormatFloat(avgSize/2.54, 'f', 2, 64)+"in)")
			records(s, m, vaginaSizes[0].size, vaginaSizes[0].member.User.ID, "vagina")
			return
		}

		text = "Largest **" + strconv.Itoa(num) + "** depths in this server: \n"
	}
	for i := 0; i < num; i++ {
		emote := ""
		if vaginaSizes[i].percentile < 25 {
			emote = ":beach:"
		} else if vaginaSizes[i].percentile > 75 {
			emote = ":ocean:"
		}
		text += "**" + vaginaSizes[i].member.User.Username + ":** " + strconv.FormatFloat(vaginaSizes[i].size, 'f', 2, 64) + "cm (" + strconv.FormatFloat(vaginaSizes[i].size/2.54, 'f', 2, 64) + "in) " + emote + "\n"
	}
	s.ChannelMessageSend(m.ChannelID, text+"**Average depth in the server:** "+strconv.FormatFloat(avgSize, 'f', 2, 64)+"cm ("+strconv.FormatFloat(avgSize/2.54, 'f', 2, 64)+"in)")
	records(s, m, vaginaSizes[0].size, vaginaSizes[0].member.User.ID, "vagina")
}

// History shows the largest and smallest penis and vagina sizes ever recorded
func History(s *discordgo.Session, m *discordgo.MessageCreate) {
	serverRegex, _ := regexp.Compile(`(?i)-s`)
	averagePenis := 13.91
	stddevPenis := 2.20
	averageVagina := 9.6
	stddevVagina := 1.5

	genitalRecord := tools.GetGenitalRecord(s)
	text := "Records for all servers:\n"

	// Check server flag
	if serverRegex.MatchString(m.Content) {
		text = "Records for this server:\n"
		server, err := s.Guild(m.GuildID)
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, "This is not a server!")
			return
		}
		serverData := tools.GetServer(*server, s)
		genitalRecord.Penis = serverData.Genital.Penis
		genitalRecord.Vagina = serverData.Genital.Vagina
	}
	penisText := ""
	vaginaText := ""

	if genitalRecord.Penis.Largest.Size != 0 {
		user, err := s.User(genitalRecord.Penis.Largest.UserID)
		if err != nil {
			penisText += "Unknown user (" + genitalRecord.Penis.Largest.UserID + ") "
		} else {
			penisText += "**" + user.Username + "** "
		}
		percentile := 100 * 0.5 * math.Erfc((averagePenis-genitalRecord.Penis.Largest.Size)/(math.Sqrt(2.0)*stddevPenis))
		penisText += "on " + strings.Replace(genitalRecord.Penis.Largest.Date.Format(time.RFC822Z), "+0000", "UTC", -1) + ": " + strconv.FormatFloat(genitalRecord.Penis.Largest.Size, 'f', 2, 64) + "cm (" + strconv.FormatFloat(genitalRecord.Penis.Largest.Size/2.54, 'f', 2, 64) + "in) larger than approximately " + strconv.FormatFloat(percentile, 'f', 2, 64) + "% of penises.\n"
	}
	if genitalRecord.Penis.Smallest.Size != 1e308 {
		user, err := s.User(genitalRecord.Penis.Smallest.UserID)
		if err != nil {
			penisText += "Unknown user (" + genitalRecord.Penis.Smallest.UserID + ") "
		} else {
			penisText += "**" + user.Username + "** "
		}
		percentile := 100 * 0.5 * math.Erfc((averagePenis-genitalRecord.Penis.Smallest.Size)/(math.Sqrt(2.0)*stddevPenis))
		penisText += "on " + strings.Replace(genitalRecord.Penis.Smallest.Date.Format(time.RFC822Z), "+0000", "UTC", -1) + ": " + strconv.FormatFloat(genitalRecord.Penis.Smallest.Size, 'f', 2, 64) + "cm (" + strconv.FormatFloat(genitalRecord.Penis.Smallest.Size/2.54, 'f', 2, 64) + "in) larger than approximately " + strconv.FormatFloat(percentile, 'f', 2, 64) + "% of penises.\n"
	}

	if genitalRecord.Vagina.Largest.Size != 0 {
		user, err := s.User(genitalRecord.Vagina.Largest.UserID)
		if err != nil {
			vaginaText += "Unknown user (" + genitalRecord.Vagina.Largest.UserID + ") "
		} else {
			vaginaText += "**" + user.Username + "** "
		}
		percentile := 100 * 0.5 * math.Erfc((averageVagina-genitalRecord.Vagina.Largest.Size)/(math.Sqrt(2.0)*stddevVagina))
		vaginaText += "on " + strings.Replace(genitalRecord.Vagina.Largest.Date.Format(time.RFC822Z), "+0000", "UTC", -1) + ": " + strconv.FormatFloat(genitalRecord.Vagina.Largest.Size, 'f', 2, 64) + "cm (" + strconv.FormatFloat(genitalRecord.Vagina.Largest.Size/2.54, 'f', 2, 64) + "in) deeper than approximately " + strconv.FormatFloat(percentile, 'f', 2, 64) + "% of vaginas.\n"
	}
	if genitalRecord.Vagina.Smallest.Size != 1e308 {
		user, err := s.User(genitalRecord.Vagina.Smallest.UserID)
		if err != nil {
			vaginaText += "Unknown user (" + genitalRecord.Vagina.Smallest.UserID + ") "
		} else {
			vaginaText += "**" + user.Username + "** "
		}
		percentile := 100 * 0.5 * math.Erfc((averageVagina-genitalRecord.Vagina.Smallest.Size)/(math.Sqrt(2.0)*stddevVagina))
		vaginaText += "on " + strings.Replace(genitalRecord.Vagina.Smallest.Date.Format(time.RFC822Z), "+0000", "UTC", -1) + ": " + strconv.FormatFloat(genitalRecord.Vagina.Smallest.Size, 'f', 2, 64) + "cm (" + strconv.FormatFloat(genitalRecord.Vagina.Smallest.Size/2.54, 'f', 2, 64) + "in) deeper than approximately " + strconv.FormatFloat(percentile, 'f', 2, 64) + "% of vaginas.\n"
	}
	text += "__Penis:__\n" + penisText + "\n\n__Vaginas:__\n" + vaginaText

	s.ChannelMessageSend(m.ChannelID, text)
}

func records(s *discordgo.Session, m *discordgo.MessageCreate, size float64, userID string, genital string) {
	// Check full records
	genitalRecords := tools.GetGenitalRecord(s)
	recordBroken := false

	if genital == "penis" {
		if size > genitalRecords.Penis.Largest.Size {
			genitalRecords.Penis.Largest.Size = size
			genitalRecords.Penis.Largest.UserID = userID
			genitalRecords.Penis.Largest.Date = time.Now().UTC()
			recordBroken = true
		} else if size < genitalRecords.Penis.Smallest.Size {
			genitalRecords.Penis.Smallest.Size = size
			genitalRecords.Penis.Smallest.UserID = userID
			genitalRecords.Penis.Smallest.Date = time.Now().UTC()
			recordBroken = true
		}
	} else if genital == "vagina" {
		if size > genitalRecords.Vagina.Largest.Size {
			genitalRecords.Vagina.Largest.Size = size
			genitalRecords.Vagina.Largest.UserID = userID
			genitalRecords.Vagina.Largest.Date = time.Now().UTC()
			recordBroken = true
		} else if size < genitalRecords.Vagina.Smallest.Size {
			genitalRecords.Vagina.Smallest.Size = size
			genitalRecords.Vagina.Smallest.UserID = userID
			genitalRecords.Vagina.Smallest.Date = time.Now().UTC()
			recordBroken = true
		}
	}

	if recordBroken {
		jsonCache, err := json.Marshal(genitalRecords)
		tools.ErrRead(s, err)

		err = ioutil.WriteFile("./data/genitalRecords.json", jsonCache, 0644)
		tools.ErrRead(s, err)
	}

	// Check server records
	server, err := s.Guild(m.GuildID)
	if err != nil {
		return
	}
	serverData := tools.GetServer(*server, s)
	recordBroken = false

	if genital == "penis" {
		if size > serverData.Genital.Penis.Largest.Size {
			serverData.Genital.Penis.Largest.Size = size
			serverData.Genital.Penis.Largest.UserID = userID
			serverData.Genital.Penis.Largest.Date = time.Now().UTC()
			recordBroken = true
		} else if size < serverData.Genital.Penis.Smallest.Size {
			serverData.Genital.Penis.Smallest.Size = size
			serverData.Genital.Penis.Smallest.UserID = userID
			serverData.Genital.Penis.Smallest.Date = time.Now().UTC()
			recordBroken = true
		}
	} else if genital == "vagina" {
		if size > serverData.Genital.Vagina.Largest.Size {
			serverData.Genital.Vagina.Largest.Size = size
			serverData.Genital.Vagina.Largest.UserID = userID
			serverData.Genital.Vagina.Largest.Date = time.Now().UTC()
			recordBroken = true
		} else if size < serverData.Genital.Vagina.Smallest.Size {
			serverData.Genital.Vagina.Smallest.Size = size
			serverData.Genital.Vagina.Smallest.UserID = userID
			serverData.Genital.Vagina.Smallest.Date = time.Now().UTC()
			recordBroken = true
		}
	}

	if recordBroken {
		jsonCache, err := json.Marshal(serverData)
		tools.ErrRead(s, err)

		err = ioutil.WriteFile("./data/serverData/"+m.GuildID+".json", jsonCache, 0644)
		tools.ErrRead(s, err)
	}
}
