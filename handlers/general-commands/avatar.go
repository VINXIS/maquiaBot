package gencommands

import (
	"regexp"
	"strings"

	"github.com/bwmarrin/discordgo"
)

// Avatar gets the avatar of the user/referenced user
func Avatar(s *discordgo.Session, m *discordgo.MessageCreate) {
	negateRegex, _ := regexp.Compile(`-(np|noprev(iew)?)`)

	users := m.Mentions
	if len(users) == 0 {
		if negateRegex.MatchString(m.Content) {
			s.ChannelMessageSend(m.ChannelID, "Your avatar is: <"+m.Author.AvatarURL("")+">")
			return
		}
		s.ChannelMessageSend(m.ChannelID, "Your avatar is: "+m.Author.AvatarURL(""))
		return
	}
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
}
