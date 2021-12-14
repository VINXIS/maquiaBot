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
	userRegex, _ := regexp.Compile(`(?i)(quote)?(a|ava|avatar)(q|quote)?\s+(.+)`)
	serverRegex, _ := regexp.Compile(`(?i)(-s\s|-s$)`)

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
			Content: "Here is the server's avatar:",
			Files: []*discordgo.File{
				{
					Name:   "image.png",
					Reader: imgBytes,
				},
			},
		})
		return
	} else if len(users) > 0 {
		var names []string
		var avatars []string
		for _, user := range users {
			names = append(names, user.Username)
			avatars = append(avatars, user.AvatarURL("2048"))
		}
		postAva(s, m, names, avatars, true)
		return
	} else if userRegex.MatchString(m.Content) {
		username := userRegex.FindStringSubmatch(m.Content)[4]
		discordUser, err := s.User(username)
		if err == nil {
			postAva(s, m, []string{discordUser.Username}, []string{discordUser.AvatarURL("2048")}, true)
			return
		}

		members, err := s.GuildMembers(m.GuildID, "", 1000)
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, "This is not a server!")
			return
		}

		// Run through usernames, if no match is found, run through member names, if no match is found, send the message author's avatar
		sort.Slice(members, func(i, j int) bool {
			time1, _ := members[i].JoinedAt.Parse()
			time2, _ := members[j].JoinedAt.Parse()
			return time1.Unix() < time2.Unix()
		})
		for _, member := range members {
			if strings.HasPrefix(strings.ToLower(member.User.Username), strings.ToLower(username)) || strings.HasPrefix(strings.ToLower(member.Nick), strings.ToLower(username)) {
				discordUser, _ = s.User(member.User.ID)
				postAva(s, m, []string{member.User.Username}, []string{discordUser.AvatarURL("2048")}, true)
				return
			}
		}
		postAva(s, m, []string{username}, []string{m.Author.AvatarURL("2048")}, false)
		return
	} else {
		postAva(s, m, []string{}, []string{m.Author.AvatarURL("2048")}, true)
	}
}

func postAva(s *discordgo.Session, m *discordgo.MessageCreate, name, avatarURL []string, found bool) {
	negateRegex, _ := regexp.Compile(`-(np|noprev(iew)?)`)
	if len(name) == 0 {
		if negateRegex.MatchString(m.Content) {
			s.ChannelMessageSend(m.ChannelID, "Your general avatar is: <"+avatarURL[0]+">")
		} else {
			s.ChannelMessageSend(m.ChannelID, "Your general avatar is: "+avatarURL[0])
		}
	} else if len(name) == 1 {
		if found {
			if negateRegex.MatchString(m.Content) {
				s.ChannelMessageSend(m.ChannelID, name[0]+"'s general avatar is: <"+avatarURL[0]+">")
			} else {
				s.ChannelMessageSend(m.ChannelID, name[0]+"'s general avatar is: "+avatarURL[0])
			}
		} else {
			if negateRegex.MatchString(m.Content) {
				s.ChannelMessageSend(m.ChannelID, "No person named "+name[0]+", Your general avatar is: <"+avatarURL[0]+">")
			} else {
				s.ChannelMessageSend(m.ChannelID, "No person named "+name[0]+", Your general avatar is: "+avatarURL[0])
			}
		}
	} else {
		message := ""
		if negateRegex.MatchString(m.Content) {
			for i := range name {
				message += name[i] + "'s general avatar is: <" + avatarURL[i] + ">\n"
			}
		} else {
			for i := range name {
				message += name[i] + "'s general avatar is: " + avatarURL[i] + "\n"
			}
		}
		s.ChannelMessageSend(m.ChannelID, message)
	}
}
