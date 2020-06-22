package botcreatorcommands

import (
	config "maquiaBot/config"
	"github.com/bwmarrin/discordgo"
)

// Servers gives the list of servers
func Servers(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID != config.Conf.BotHoster.UserID {
		s.ChannelMessageSend(m.ChannelID, "YOU ARE NOT "+config.Conf.BotHoster.Username+".........")
		return
	}

	dm, _ := s.UserChannelCreate(config.Conf.BotHoster.UserID)
	message := ""
	for _, guild := range s.State.Guilds {
		owner, _ := s.User(guild.OwnerID)
		message += "`" + guild.ID + "` **" + guild.Name + "** - " + owner.String() + "\n"
	}
	s.ChannelMessageSend(dm.ID, message)
}
