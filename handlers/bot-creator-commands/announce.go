package botcreatorcommands

import (
	"regexp"
	"strings"

	config "../../config"
	tools "../../tools"
	"github.com/bwmarrin/discordgo"
)

// Announce announces new stuff
func Announce(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID != config.Conf.BotHoster.UserID {
		s.ChannelMessageSend(m.ChannelID, "YOU ARE NOT "+config.Conf.BotHoster.Username+".........")
		return
	}

	announceRegex, _ := regexp.Compile(`announce\s+(.+)`)
	if !announceRegex.MatchString(m.Content) {
		s.ChannelMessageSend(m.ChannelID, "ur dumb as hell")
		return
	}

	announcement := announceRegex.FindStringSubmatch(m.Content)[1] + strings.Join(strings.Split(m.Content, "\n")[1:], "\n")

	for _, guild := range s.State.Guilds {
		if guild.ID == m.GuildID {
			continue
		}

		server := tools.GetServer(*guild, s)
		if server.AnnounceChannel != "" {
			s.ChannelMessageSend(server.AnnounceChannel, "Admins of the server can always toggle announcements from the bot creator on/off by using `toggle -a`.\n\n**Announcement below:**\n"+announcement)
		}
	}
	s.ChannelMessageSend(m.ChannelID, "Sent announcement to all servers!")
	return
}
