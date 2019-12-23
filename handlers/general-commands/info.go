package gencommands

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	structs "../../structs"
	tools "../../tools"
	"github.com/bwmarrin/discordgo"
)

// Info gives information about the user
func Info(s *discordgo.Session, m *discordgo.MessageCreate, cache []structs.PlayerData) {
	userRegex, _ := regexp.Compile(`info\s+(.+)`)

	userTest := ""
	user := m.Author
	nickname := "N/A"
	roles := "N/A"
	var joinDate discordgo.Timestamp
	var err error
	if len(m.Mentions) == 1 {
		user = m.Mentions[0]
	} else {
		if userRegex.MatchString(m.Content) {
			userTest = userRegex.FindStringSubmatch(m.Content)[1]
		}
		user, err = s.User(userTest)
		if err == nil {
			userTest = user.Username
		} else {
			user = m.Author
		}
	}
	members, err := s.GuildMembers(m.GuildID, "", 1000)
	if err == nil {
		if userTest == "" {
			for _, member := range members {
				if member.User.ID == m.Author.ID {
					user, _ = s.User(member.User.ID)
					nickname = member.Nick
					joinDate = member.JoinedAt
					roles = ""
					for _, role := range member.Roles {
						discordRole, err := s.State.Role(m.GuildID, role)
						if err != nil {
							continue
						}
						roles = roles + discordRole.Name + ", "
					}
					if roles != "" {
						roles = roles[:len(roles)-2]
					}
				}
			}
		} else {
			sort.Slice(members, func(i, j int) bool {
				time1, _ := members[i].JoinedAt.Parse()
				time2, _ := members[j].JoinedAt.Parse()
				return time1.Unix() < time2.Unix()
			})
			for _, member := range members {
				if strings.HasPrefix(strings.ToLower(member.User.Username), userTest) || strings.HasPrefix(strings.ToLower(member.Nick), userTest) {
					user, _ = s.User(member.User.ID)
					nickname = member.Nick
					joinDate = member.JoinedAt
					roles = ""
					for _, role := range member.Roles {
						discordRole, err := s.State.Role(m.GuildID, role)
						if err != nil {
							continue
						}
						roles = roles + discordRole.Name + ", "
					}
					if roles != "" {
						roles = roles[:len(roles)-2]
					}
					break
				}
			}
		}
	}

	// Reformat joinDate
	joinDateDate, _ := joinDate.Parse()
	serverCreateDate, _ := discordgo.SnowflakeTimestamp(m.GuildID)
	joinDateString := "N/A"
	if joinDateDate.After(serverCreateDate) {
		joinDateString = joinDateDate.Format(time.RFC822Z)
	}

	// Created at date
	createdAt, _ := discordgo.SnowflakeTimestamp(user.ID)

	// Status
	presence, err := s.State.Presence(m.GuildID, user.ID)
	status := "Offline"
	if err == nil {
		status = strings.Title(string(presence.Status))
	}

	// Obtain osu! info
	osuUsername := "N/A"
	for _, player := range cache {
		if player.Discord.ID == user.ID && player.Osu.Username != "" {
			osuUsername = player.Osu.Username
			break
		}
	}
	// Fix any blanks
	if roles == "" {
		roles = "None"
	}
	if nickname == "" {
		nickname = "None"
	}

	s.ChannelMessageSendEmbed(m.ChannelID, &discordgo.MessageEmbed{
		Author: &discordgo.MessageEmbedAuthor{
			Name:    user.String(),
			IconURL: user.AvatarURL("2048"),
		},
		Thumbnail: &discordgo.MessageEmbedThumbnail{
			URL: user.AvatarURL("2048"),
		},
		Fields: []*discordgo.MessageEmbedField{
			&discordgo.MessageEmbedField{
				Name:   "ID",
				Value:  user.ID,
				Inline: true,
			},
			&discordgo.MessageEmbedField{
				Name:   "Nickname",
				Value:  nickname,
				Inline: true,
			},
			&discordgo.MessageEmbedField{
				Name:   "Account Created",
				Value:  createdAt.UTC().Format(time.RFC822Z),
				Inline: true,
			},
			&discordgo.MessageEmbedField{
				Name:   "Date Joined",
				Value:  joinDateString,
				Inline: true,
			},
			&discordgo.MessageEmbedField{
				Name:   "Status",
				Value:  status,
				Inline: true,
			},
			&discordgo.MessageEmbedField{
				Name:   "Linked on osu! as",
				Value:  osuUsername,
				Inline: true,
			},
			&discordgo.MessageEmbedField{
				Name:  "Roles",
				Value: roles,
			},
		},
	})
}

// RoleInfo gives information about a role
func RoleInfo(s *discordgo.Session, m *discordgo.MessageCreate) {
	roleRegex, _ := regexp.Compile(`r(ole)?info\s+(.+)`)

	// Get server
	server, err := s.Guild(m.GuildID)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "This is not a server!")
		return
	}
	serverImg := "https://cdn.discordapp.com/icons/" + server.ID + "/" + server.Icon
	if strings.Contains(server.Icon, "a_") {
		serverImg += ".gif"
	} else {
		serverImg += ".png"
	}

	// Get role
	if !roleRegex.MatchString(m.Content) {
		s.ChannelMessageSend(m.ChannelID, "No role name given!")
		return
	}

	roles, _ := s.GuildRoles(m.GuildID)
	roleName := roleRegex.FindStringSubmatch(m.Content)[2]
	role := &discordgo.Role{}
	for _, servrole := range roles {
		if strings.HasPrefix(strings.ToLower(servrole.Name), roleName) {
			role = servrole
		}
	}
	if role.ID == "" {
		s.ChannelMessageSend(m.ChannelID, "Could not find a role called **"+roleName+"**!")
		return
	}

	// Get list of people who have the role
	members, err := s.GuildMembers(m.GuildID, "", 1000)
	memberList := ""
	memberCount := 0
	for _, member := range members {
		for _, memberRole := range member.Roles {
			if memberRole == role.ID {
				memberCount++
				memberList += member.User.Username + ", "
				break
			}
		}
	}
	memberList = strings.TrimSpace(strings.TrimSuffix(strings.TrimSpace(memberList), ","))
	if memberList == "" {
		memberList = "N/A"
	}

	// Created at date
	createdAt, _ := discordgo.SnowflakeTimestamp(role.ID)

	_, err = s.ChannelMessageSendEmbed(m.ChannelID, &discordgo.MessageEmbed{
		Color: role.Color,
		Author: &discordgo.MessageEmbedAuthor{
			Name:    role.Name,
			IconURL: serverImg,
		},
		Thumbnail: &discordgo.MessageEmbedThumbnail{
			URL: serverImg,
		},
		Fields: []*discordgo.MessageEmbedField{
			&discordgo.MessageEmbedField{
				Name:   "ID",
				Value:  role.ID,
				Inline: true,
			},
			&discordgo.MessageEmbedField{
				Name:   "Role Created",
				Value:  createdAt.UTC().Format(time.RFC822Z),
				Inline: true,
			},
			&discordgo.MessageEmbedField{
				Name:   "Position",
				Value:  strconv.Itoa(role.Position),
				Inline: true,
			},
			&discordgo.MessageEmbedField{
				Name:   "Managed Externally",
				Value:  strconv.FormatBool(role.Managed),
				Inline: true,
			},
			&discordgo.MessageEmbedField{
				Name:   "Mentionable",
				Value:  strconv.FormatBool(role.Mentionable),
				Inline: true,
			},
			&discordgo.MessageEmbedField{
				Name:   "Shown Separately",
				Value:  strconv.FormatBool(role.Hoist),
				Inline: true,
			},
			&discordgo.MessageEmbedField{
				Name:  "Members (" + strconv.Itoa(memberCount) + ")",
				Value: memberList,
			},
		},
	})
}

// ServerInfo gives information about the server
func ServerInfo(s *discordgo.Session, m *discordgo.MessageCreate) {

	server, err := s.Guild(m.GuildID)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "This is not a server!")
		return
	}

	// Stats info
	serverData := tools.GetServer(*server)
	statsInfo := strconv.Itoa(len(serverData.Nouns)) + " nouns\n" + strconv.Itoa(len(serverData.Adjectives)) + " adjectives\n" + strconv.Itoa(len(serverData.Skills)) + " skills\nAllowAnyoneAdd: " + strconv.FormatBool(serverData.AllowAnyoneStats)

	// Toggle Information
	toggleInfo := "Cheers: " + strconv.FormatBool(serverData.Cheers) + "\n" +
		"Crab: " + strconv.FormatBool(serverData.Crab) + "\n" +
		"Daily: " + strconv.FormatBool(serverData.Daily) + "\n" +
		"Late: " + strconv.FormatBool(serverData.Late) + "\n" +
		"Nice Idea: " + strconv.FormatBool(serverData.NiceIdea) + "\n" +
		"osu!: " + strconv.FormatBool(serverData.OsuToggle) + "\n" +
		"Over it: " + strconv.FormatBool(serverData.OverIt) + "\n" +
		"Vibe Check: " + strconv.FormatBool(serverData.Vibe) + "\n" +
		"Announcements: " + strconv.FormatBool(serverData.Announce) + "\n"

	// Created at date
	createdAt, err := discordgo.SnowflakeTimestamp(server.ID)

	// Server owner
	owner, _ := s.User(server.OwnerID)

	// Channel info
	channels := server.Channels
	text := 0
	voice := 0
	category := 0
	other := 0
	for _, ch := range channels {
		switch ch.Type {
		case discordgo.ChannelTypeGuildText:
			text++
		case discordgo.ChannelTypeGuildVoice:
			voice++
		case discordgo.ChannelTypeGuildCategory:
			category++
		default:
			other++
		}
	}
	channelInfo := ""
	if text != 0 {
		channelInfo += strconv.Itoa(text) + " text\n"
	}
	if voice != 0 {
		channelInfo += strconv.Itoa(voice) + " voice\n"
	}
	if category != 0 {
		if category == 1 {
			channelInfo += strconv.Itoa(category) + " category\n"
		} else {
			channelInfo += strconv.Itoa(category) + " categories\n"
		}
	}
	if other != 0 {
		channelInfo += strconv.Itoa(other) + " other\n"
	}

	// Quote info
	quoteCount := strconv.Itoa(len(serverData.Quotes))
	quoteInfo := "N/A"
	quoteData := []struct {
		Name  string
		ID    string
		Count int
	}{}
	for _, quote := range serverData.Quotes {
		included := false
		for i, user := range quoteData {
			if quote.Author.ID == user.ID {
				included = true
				quoteData[i].Count++
				break
			}
		}
		if !included {
			user, _ := s.User(quote.Author.ID)
			quoteData = append(quoteData, struct {
				Name  string
				ID    string
				Count int
			}{
				Name:  user.Username,
				ID:    quote.Author.ID,
				Count: 1,
			})
		}
	}
	sort.Slice(quoteData, func(i, j int) bool { return quoteData[i].Count > quoteData[j].Count })

	if len(quoteData) != 0 {
		quoteInfo = ""
	}
	for i, user := range quoteData {
		if (i+1)%4 == 0 {
			quoteInfo += fmt.Sprintf("**%s:** %d,\n", user.Name, user.Count)
		} else {
			quoteInfo += fmt.Sprintf("**%s:** %d, ", user.Name, user.Count)
		}
	}
	quoteInfo = strings.TrimSuffix(quoteInfo, ", ")
	quoteInfo = strings.TrimSuffix(quoteInfo, ",\n")

	// AFK Timeout
	timeout := server.AfkTimeout
	timeoutString := "N/A"
	if timeout == int(time.Hour.Seconds()) {
		timeoutString = "1 hour"
	} else {
		timeoutString = strconv.Itoa(timeout/60) + " min"
	}

	// AFK Channel
	afkCh, err := s.Channel(server.AfkChannelID)
	afkChName := "N/A"
	if err == nil {
		afkChName = afkCh.Name
	}

	serverImg := "https://cdn.discordapp.com/icons/" + server.ID + "/" + server.Icon
	if strings.Contains(server.Icon, "a_") {
		serverImg += ".gif"
	} else {
		serverImg += ".png"
	}

	s.ChannelMessageSendEmbed(m.ChannelID, &discordgo.MessageEmbed{
		Author: &discordgo.MessageEmbedAuthor{
			Name:    server.Name,
			IconURL: serverImg,
		},
		Thumbnail: &discordgo.MessageEmbedThumbnail{
			URL: serverImg,
		},
		Fields: []*discordgo.MessageEmbedField{
			&discordgo.MessageEmbedField{
				Name:   "ID",
				Value:  server.ID,
				Inline: true,
			},
			&discordgo.MessageEmbedField{
				Name:   "Region",
				Value:  server.Region,
				Inline: true,
			},
			&discordgo.MessageEmbedField{
				Name:   "Server Created",
				Value:  createdAt.UTC().Format(time.RFC822Z),
				Inline: true,
			},
			&discordgo.MessageEmbedField{
				Name:   "Server Owner",
				Value:  owner.String(),
				Inline: true,
			},
			&discordgo.MessageEmbedField{
				Name:   "AFK Timeout",
				Value:  timeoutString,
				Inline: true,
			},
			&discordgo.MessageEmbedField{
				Name:   "AFK Channel",
				Value:  afkChName,
				Inline: true,
			},
			&discordgo.MessageEmbedField{
				Name:   "Emoji Count",
				Value:  strconv.Itoa(len(server.Emojis)),
				Inline: true,
			},
			&discordgo.MessageEmbedField{
				Name:   "Member Count",
				Value:  strconv.Itoa(len(server.Members)),
				Inline: true,
			},
			&discordgo.MessageEmbedField{
				Name:   "Role Count",
				Value:  strconv.Itoa(len(server.Roles)),
				Inline: true,
			},
			&discordgo.MessageEmbedField{
				Name:   "Channels (" + strconv.Itoa(len(server.Channels)) + ")",
				Value:  channelInfo,
				Inline: true,
			},
			&discordgo.MessageEmbedField{
				Name:   "`" + serverData.Prefix + "stats` Information:",
				Value:  statsInfo,
				Inline: true,
			},
			&discordgo.MessageEmbedField{
				Name:   "Server Options",
				Value:  toggleInfo,
				Inline: true,
			},
			&discordgo.MessageEmbedField{
				Name:   "Quotes (" + quoteCount + ")",
				Value:  quoteInfo,
				Inline: true,
			},
		},
	})

	// Save new server data
	serverData.Time = time.Now()
	jsonCache, err := json.Marshal(serverData)
	tools.ErrRead(err)

	err = ioutil.WriteFile("./data/serverData/"+m.GuildID+".json", jsonCache, 0644)
	tools.ErrRead(err)
}
