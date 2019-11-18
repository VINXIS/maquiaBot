package gencommands

import (
	"regexp"
	"sort"
	"strings"
	"time"

	structs "../../structs"
	"github.com/bwmarrin/discordgo"
)

// Info gives information about the user
func Info(s *discordgo.Session, m *discordgo.MessageCreate, cache []structs.PlayerData) {
	userRegex, _ := regexp.Compile(`info\s+(.+)`)

	user := m.Author
	nickname := "N/A"
	roles := "N/A"
	var joinDate discordgo.Timestamp
	var err error
	if len(m.Mentions) == 1 {
		user = m.Mentions[0]
	} else {
		userTest := ""
		if userRegex.MatchString(m.Content) {
			userTest = userRegex.FindStringSubmatch(m.Content)[1]
		}
		user, err = s.User(userTest)
		if err != nil {
			server, err := s.Guild(m.GuildID)
			if err != nil {
				s.ChannelMessageSend(m.ChannelID, "This is not a server! Use their ID directly instead.")
				return
			}
			if userTest == "" {
				for _, member := range server.Members {
					if member.User.ID == m.Author.ID {
						user, _ = s.User(member.User.ID)
						nickname = member.Nick
						joinDate = member.JoinedAt
						roles = ""
						for _, role := range member.Roles {
							discordRole, _ := s.State.Role(server.ID, role)
							roles = roles + discordRole.Name + ", "
						}
						if roles != "" {
							roles = roles[:len(roles)-2]
						}
					}
				}
			} else {
				sort.Slice(server.Members, func(i, j int) bool {
					time1, _ := server.Members[i].JoinedAt.Parse()
					time2, _ := server.Members[j].JoinedAt.Parse()
					return time1.Unix() < time2.Unix()
				})
				for _, member := range server.Members {
					if strings.HasPrefix(strings.ToLower(member.User.Username), userTest) || strings.HasPrefix(strings.ToLower(member.Nick), userTest) {
						user, _ = s.User(member.User.ID)
						nickname = member.Nick
						joinDate = member.JoinedAt
						roles = ""
						for _, role := range member.Roles {
							discordRole, _ := s.State.Role(server.ID, role)
							roles = roles + discordRole.Name + ", "
						}
						if roles != "" {
							roles = roles[:len(roles)-2]
						}
						break
					}
				}
			}
		} else {
			server, err := s.Guild(m.GuildID)
			if err == nil {
				for _, member := range server.Members {
					if member.User.ID == user.ID {
						nickname = member.Nick
						joinDate = member.JoinedAt
						roles = ""
						for _, role := range member.Roles {
							discordRole, _ := s.State.Role(server.ID, role)
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
	}

	// Reformat joinDate
	joinDateDate, _ := joinDate.Parse()

	// Created at date
	createdAt, err := discordgo.SnowflakeTimestamp(user.ID)

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
			Name:    user.Username + "#" + user.Discriminator,
			IconURL: user.AvatarURL(""),
		},
		Thumbnail: &discordgo.MessageEmbedThumbnail{
			URL: user.AvatarURL(""),
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
				Value:  joinDateDate.Format(time.RFC822Z),
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
