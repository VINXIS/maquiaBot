package botcreatorcommands

import (
	"regexp"
	"strings"

	config "maquiaBot/config"
	tools "maquiaBot/tools"

	"github.com/bwmarrin/discordgo"
)

// Announce announces new stuff
func Announce(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID != config.Conf.BotHoster.UserID {
		s.ChannelMessageSend(m.ChannelID, "YOU ARE NOT "+config.Conf.BotHoster.Username+".........")
		return
	}

	announceRegex, _ := regexp.Compile(`(?i)announce\s+(.+)`)
	if !announceRegex.MatchString(m.Content) {
		s.ChannelMessageSend(m.ChannelID, "ur dumb as hell")
		return
	}

	msg, err := s.ChannelMessageSend(m.ChannelID, "Sending announcement to servers...")
	if err != nil {
		return
	}

	announcement := announceRegex.FindStringSubmatch(m.Content)[1] + strings.Join(strings.Split(m.Content, "\n")[1:], "\n")

	for _, guild := range s.State.Guilds {
		if guild.ID == m.GuildID {
			continue
		}

		server, _ := tools.GetServer(*guild, s)
		if server.AnnounceChannel != "" {
			s.ChannelMessageSend(server.AnnounceChannel, "Admins of the server can always toggle announcements from the bot creator on/off by using `toggle -a`.\n\n**Announcement below:**\n"+announcement)
		}
	}

	go s.ChannelMessageDelete(msg.ChannelID, msg.ID)
	s.ChannelMessageSend(m.ChannelID, "Sent announcement to all servers!")
	return
}
