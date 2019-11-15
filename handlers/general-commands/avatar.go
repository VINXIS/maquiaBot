package gencommands

import (
	"bytes"
	"image/png"
	"regexp"
	"sort"
	"strings"

	"github.com/bwmarrin/discordgo"
)

// Avatar gets the avatar of the user/referenced user
func Avatar(s *discordgo.Session, m *discordgo.MessageCreate) {
	negateRegex, _ := regexp.Compile(`-(np|noprev(iew)?)`)
	userRegex, _ := regexp.Compile(`(a|ava|avatar)\s+(.+)`)
	serverRegex, _ := regexp.Compile(`(-s\s|-s$)`)

	users := m.Mentions
	if serverRegex.MatchString(m.Content) {
		ava, err := s.GuildIcon(m.GuildID)
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, "This is not a server!")
			return
		}

		imgBytes := new(bytes.Buffer)
		_ = png.Encode(imgBytes, ava)
		s.ChannelMessageSendComplex(m.ChannelID, &discordgo.MessageSend{
			Content: "Here is the server avatar:",
			Files: []*discordgo.File{
				&discordgo.File{
					Name:   "image.png",
					Reader: imgBytes,
				},
			},
		})
	} else if len(users) > 0 {
		var avatarURLs strings.Builder
		if negateRegex.MatchString(m.Content) {
			for _, mention := range users {
				avatarURLs.WriteString(mention.Username + "'s avatar is: <" + mention.AvatarURL("") + ">\n")
			}
		} else {
			for _, mention := range users {
				avatarURLs.WriteString(mention.Username + "'s avatar is: " + mention.AvatarURL("") + "\n")
			}
		}
		s.ChannelMessageSend(m.ChannelID, avatarURLs.String())
		return
	} else if userRegex.MatchString(m.Content) {
		username := userRegex.FindStringSubmatch(m.Content)[2]
		discordUser, err := s.User(username)
		if err == nil {
			if negateRegex.MatchString(m.Content) {
				s.ChannelMessageSend(m.ChannelID, discordUser.Username+"'s avatar is: <"+discordUser.AvatarURL("")+">")
				return
			}
			s.ChannelMessageSend(m.ChannelID, discordUser.Username+"'s avatar is: "+discordUser.AvatarURL(""))
			return
		}

		server, err := s.Guild(m.GuildID)
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, "This is not a server!")
			return
		}
		sort.Slice(server.Members, func(i, j int) bool {
			time1, _ := server.Members[i].JoinedAt.Parse()
			time2, _ := server.Members[j].JoinedAt.Parse()
			return time1.Unix() < time2.Unix()
		})
		for _, member := range server.Members {
			if strings.HasPrefix(strings.ToLower(member.User.Username), username) {
				discordUser, _ = s.User(member.User.ID)
				if negateRegex.MatchString(m.Content) {
					s.ChannelMessageSend(m.ChannelID, member.Nick+"'s avatar is: <"+discordUser.AvatarURL("")+">")
					return
				}
				s.ChannelMessageSend(m.ChannelID, member.Nick+"'s avatar is: "+discordUser.AvatarURL(""))
				return
			}
		}
		for _, member := range server.Members {
			if strings.HasPrefix(strings.ToLower(member.Nick), username) {
				discordUser, _ := s.User(member.User.ID)
				if negateRegex.MatchString(m.Content) {
					s.ChannelMessageSend(m.ChannelID, member.Nick+"'s avatar is: <"+discordUser.AvatarURL("")+">")
					return
				}
				s.ChannelMessageSend(m.ChannelID, member.Nick+"'s avatar is: "+discordUser.AvatarURL(""))
				return
			}
		}
		if negateRegex.MatchString(m.Content) {
			s.ChannelMessageSend(m.ChannelID, "No person named "+username+", Your avatar is: <"+m.Author.AvatarURL("")+">")
			return
		}
		s.ChannelMessageSend(m.ChannelID, "No person named "+username+", Your avatar is: "+m.Author.AvatarURL(""))
		return
	} else {
		if negateRegex.MatchString(m.Content) {
			s.ChannelMessageSend(m.ChannelID, "Your avatar is: <"+m.Author.AvatarURL("")+">")
			return
		}
		s.ChannelMessageSend(m.ChannelID, "Your avatar is: "+m.Author.AvatarURL(""))
		return
	}
}
