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

	structs "maquiaBot/structs"
	tools "maquiaBot/tools"

	"github.com/bwmarrin/discordgo"
)

// Info gives information about the user
func Info(s *discordgo.Session, m *discordgo.MessageCreate) {
	userRegex, _ := regexp.Compile(`(?i)info\s+(.+)`)

	userTest := ""
	user := m.Author
	nickname := "N/A"
	roles := "N/A"
	var joinDate discordgo.Timestamp
	var err error
	if len(m.Mentions) >= 1 {
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
	member, err := s.GuildMember(m.GuildID, user.ID)
	if userTest != "" {
		members, membersError := s.GuildMembers(m.GuildID, "", 1000)
		if membersError == nil {
			for _, mem := range members {
				if strings.HasPrefix(strings.ToLower(mem.User.Username), strings.ToLower(userTest)) || strings.HasPrefix(strings.ToLower(mem.Nick), strings.ToLower(userTest)) {
					member = mem
					user = mem.User
					err = nil
					break
				}
			}
		}
	}
	if err == nil {
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

	// Reformat joinDate
	joinDateDate, _ := joinDate.Parse()
	serverCreateDate, _ := discordgo.SnowflakeTimestamp(m.GuildID)
	joinDateString := "N/A"
	if joinDateDate.After(serverCreateDate) {
		joinDateString = joinDateDate.Format(time.RFC822Z)
	}

	// Created at date
	createdAt, _ := discordgo.SnowflakeTimestamp(user.ID)

	// Obtain profile cache data
	var cache []structs.PlayerData
	f, err := ioutil.ReadFile("./data/osuData/profileCache.json")
	tools.ErrRead(s, err)
	_ = json.Unmarshal(f, &cache)

	// Obtain player data info
	osuUsername := "N/A"
	points := 0.00
	for _, player := range cache {
		if player.Discord == user.ID && player.Osu.Username != "" {
			osuUsername = player.Osu.Username
			points = player.Currency.Amount
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
			{
				Name:   "ID",
				Value:  user.ID,
				Inline: true,
			},
			{
				Name:   "Nickname",
				Value:  nickname,
				Inline: true,
			},
			{
				Name:   "Account Created",
				Value:  createdAt.UTC().Format(time.RFC822Z),
				Inline: true,
			},
			{
				Name:   "Date Joined",
				Value:  joinDateString,
				Inline: true,
			},
			{
				Name:   "Points",
				Value:  strconv.FormatFloat(points, 'f', 2, 64),
				Inline: true,
			},
			{
				Name:   "Linked on osu! as",
				Value:  osuUsername,
				Inline: true,
			},
			{
				Name:  "Roles",
				Value: roles,
			},
		},
	})
}

// RoleInfo gives information about a role
func RoleInfo(s *discordgo.Session, m *discordgo.MessageCreate) {
	roleRegex, _ := regexp.Compile(`(?i)r(ole)?info\s+(.+)`)

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

	roles, _ := s.GuildRoles(m.GuildID)

	// Get role
	if !roleRegex.MatchString(m.Content) {
		serverData, _ := tools.GetServer(*server, s)
		embed := &discordgo.MessageEmbed{
			Author: &discordgo.MessageEmbedAuthor{
				Name:    server.Name,
				IconURL: serverImg,
			},
		}

		if len(serverData.RoleAutomation) == 0 {
			s.ChannelMessageSend(m.ChannelID, "There is no role automation configured for this server currently! Admins can see `help roleautomation` for details on how to add role automation.")
			return
		}

		for _, roleAuto := range serverData.RoleAutomation {
			var roleNames string
			for _, role := range roleAuto.Roles {
				for _, discordRole := range roles {
					if discordRole.ID == role {
						roleNames += discordRole.Name + ", "
						break
					}
				}
			}
			embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
				Name:  strconv.FormatInt(roleAuto.ID, 10),
				Value: "Trigger: " + roleAuto.Text + "\nRoles: " + strings.TrimSuffix(roleNames, ", "),
			})
		}
		s.ChannelMessageSendEmbed(m.ChannelID, embed)
		return
	}

	roleName := roleRegex.FindStringSubmatch(m.Content)[2]
	role := &discordgo.Role{}
	for _, servrole := range roles {
		if strings.HasPrefix(strings.ToLower(servrole.Name), strings.ToLower(roleName)) {
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

	s.ChannelMessageSendEmbed(m.ChannelID, &discordgo.MessageEmbed{
		Color: role.Color,
		Author: &discordgo.MessageEmbedAuthor{
			Name:    role.Name,
			IconURL: serverImg,
		},
		Thumbnail: &discordgo.MessageEmbedThumbnail{
			URL: serverImg,
		},
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:   "ID",
				Value:  role.ID,
				Inline: true,
			},
			{
				Name:   "Role Created",
				Value:  createdAt.UTC().Format(time.RFC822Z),
				Inline: true,
			},
			{
				Name:   "Position",
				Value:  strconv.Itoa(role.Position),
				Inline: true,
			},
			{
				Name:   "Managed Externally",
				Value:  strconv.FormatBool(role.Managed),
				Inline: true,
			},
			{
				Name:   "Mentionable",
				Value:  strconv.FormatBool(role.Mentionable),
				Inline: true,
			},
			{
				Name:   "Shown Separately",
				Value:  strconv.FormatBool(role.Hoist),
				Inline: true,
			},
			{
				Name:  "Members (" + strconv.Itoa(memberCount) + ")",
				Value: memberList,
			},
		},
	})
}

// ServerInfo gives information about the server
func ServerInfo(s *discordgo.Session, m *discordgo.MessageCreate) {
	var server *discordgo.Guild
	for _, stateServer := range s.State.Guilds {
		if m.GuildID == stateServer.ID {
			server = stateServer
			break
		}
	}
	if server.ID == "" {
		s.ChannelMessageSend(m.ChannelID, "This is not a server!")
		return
	}

	// Stats info
	serverData, _ := tools.GetServer(*server, s)
	statsInfo := strconv.Itoa(len(serverData.Nouns)) + " nouns\n" + strconv.Itoa(len(serverData.Adjectives)) + " adjectives\n" + strconv.Itoa(len(serverData.Skills)) + " skills\nAllowAnyoneAdd: " + strconv.FormatBool(serverData.AllowAnyoneStats)

	// Server Options Information
	serverOptions := "Daily refreshing commands: " + strconv.FormatBool(serverData.Daily) + "\n" +
		"osu!: " + strconv.FormatBool(serverData.OsuToggle) + "\n" +
		"osu! timestamps: " + strconv.FormatBool(serverData.TimestampToggle) + "\n" +
		"Vibe Check: " + strconv.FormatBool(serverData.Vibe) + "\n"
	if serverData.AnnounceChannel == "" {
		serverOptions += "Announcements: N/A\n"
	} else {
		announceChannel, err := s.Channel(serverData.AnnounceChannel)
		if err == nil {
			serverOptions += "Announcements: #" + announceChannel.Name + "\n"
		} else {
			serverOptions += "Announcements: N/A\n"
		}
	}

	// Created at date
	createdAt, err := discordgo.SnowflakeTimestamp(server.ID)

	// Server owner
	owner, _ := s.User(server.OwnerID)

	// Channel info
	channels, _ := s.GuildChannels(m.GuildID)
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
	if channelInfo == "" {
		channelInfo = "Currently unavailable."
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
			quoteData = append(quoteData, struct {
				Name  string
				ID    string
				Count int
			}{
				Name:  quote.Author.Username,
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
			{
				Name:   "ID",
				Value:  server.ID,
				Inline: true,
			},
			{
				Name:   "Region",
				Value:  server.Region,
				Inline: true,
			},
			{
				Name:   "Server Created",
				Value:  createdAt.UTC().Format(time.RFC822Z),
				Inline: true,
			},
			{
				Name:   "Server Owner",
				Value:  owner.String(),
				Inline: true,
			},
			{
				Name:   "AFK Channel",
				Value:  afkChName + " **After " + timeoutString + "**",
				Inline: true,
			},
			{
				Name:   "Trigger Count",
				Value:  strconv.Itoa(len(serverData.Triggers)),
				Inline: true,
			},
			{
				Name:   "Emoji Count",
				Value:  strconv.Itoa(len(server.Emojis)),
				Inline: true,
			},
			{
				Name:   "Member Count",
				Value:  strconv.Itoa(server.MemberCount),
				Inline: true,
			},
			{
				Name:   "Role Count",
				Value:  strconv.Itoa(len(server.Roles)),
				Inline: true,
			},
			{
				Name:   "Channels (" + strconv.Itoa(len(channels)) + ")",
				Value:  channelInfo,
				Inline: true,
			},
			{
				Name:   "`" + serverData.Prefix + "stats` Information:",
				Value:  statsInfo,
				Inline: true,
			},
			{
				Name:   "Server Options",
				Value:  serverOptions,
				Inline: true,
			},
			{
				Name:   "Quotes (" + quoteCount + ")",
				Value:  quoteInfo,
				Inline: true,
			},
		},
	})
}

// ChannelInfo gives information about the channel
func ChannelInfo(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Get channel
	channel, err := s.Channel(m.ChannelID)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "This is not an allowed channel!")
		return
	}

	channelData, _ := tools.GetChannel(*channel, s)

	// Created at date
	createdAt, _ := discordgo.SnowflakeTimestamp(channel.ID)

	// Last pin timestamp
	lastPin, _ := channel.LastPinTimestamp.Parse()

	// Channel Options Information
	channelOptions := "Daily refreshing commands: " + strconv.FormatBool(channelData.Daily) + "\n" +
		"osu!: " + strconv.FormatBool(channelData.OsuToggle) + "\n" +
		"osu! timestamps: " + strconv.FormatBool(channelData.TimestampToggle) + "\n" +
		"Vibe Check: " + strconv.FormatBool(channelData.Vibe) + "\n" +
		"Options are default false since server options overwrite."

	s.ChannelMessageSendEmbed(m.ChannelID, &discordgo.MessageEmbed{
		Author: &discordgo.MessageEmbedAuthor{
			Name: "#" + channel.Name,
		},
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:   "ID",
				Value:  channel.ID,
				Inline: true,
			},
			{
				Name:   "Topic",
				Value:  channel.Topic,
				Inline: true,
			},
			{
				Name:   "Channel Created",
				Value:  createdAt.UTC().Format(time.RFC822Z),
				Inline: true,
			},
			{
				Name:   "Last Pin Timestamp",
				Value:  lastPin.UTC().Format(time.RFC822Z),
				Inline: true,
			},
			{
				Name:   "Rate Limit Per User",
				Value:  strconv.Itoa(channel.RateLimitPerUser),
				Inline: true,
			},
			{
				Name:   "NSFW",
				Value:  strconv.FormatBool(channel.NSFW),
				Inline: true,
			},
			{
				Name:   "Channel Position",
				Value:  strconv.Itoa(channel.Position),
				Inline: true,
			},
			{
				Name:   "Channel Options",
				Value:  channelOptions,
				Inline: true,
			},
		},
	})
}
