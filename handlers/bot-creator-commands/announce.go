package botcreatorcommands

import (
	"regexp"

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

	announcement := announceRegex.FindStringSubmatch(m.Content)[1]

	for _, guild := range s.State.Guilds {
		if guild.ID == m.GuildID {
			continue
		}

		server := tools.GetServer(*guild)
		if server.Announce {
			sent := false
			for _, channel := range guild.Channels {
				if channel.ID == guild.ID {
					sent = true
					s.ChannelMessageSend(channel.ID, "Admins of the server can always toggle announcements from the bot creator on/off by using `announcetoggle`.\n\n**Announcement below:**\n"+announcement)
					break
				}
			}

			if sent {
				continue
			}

			for _, channel := range guild.Channels {
				_, err := s.ChannelMessageSend(channel.ID, "Admins of the server can always toggle announcements from the bot creator on/off by using `announcetoggle`.\n\n**Announcement below:**\n"+announcement)
				if err == nil {
					break
				}
			}
		}
	}
	s.ChannelMessageSend(m.ChannelID, "Sent announcement to all servers!")
	return
}
