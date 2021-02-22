package botcreatorcommands

import (
	"strings"
	"time"

	config "maquiaBot/config"

	"github.com/bwmarrin/discordgo"
)

// Servers gives the list of servers
func Servers(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID != config.Conf.BotHoster.UserID {
		s.ChannelMessageSend(m.ChannelID, "YOU ARE NOT "+config.Conf.BotHoster.Username+".........")
		return
	}

	dm, err := s.UserChannelCreate(config.Conf.BotHoster.UserID)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "Could not create DM channel!")
		return
	}
	msg, err := s.ChannelMessageSend(dm.ID, "Fetching servers...")
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "Could not send message in the DM channel!")
		return
	}
	message := "ID,Name,Owner,Date Joined\n"
	for _, guild := range s.State.Guilds {
		ownerName := s.State.User.String()
		owner, err := s.User(guild.OwnerID)
		if err == nil {
			ownerName = owner.String()
		}

		var joinDate discordgo.Timestamp
		member, err := s.GuildMember(guild.ID, s.State.User.ID)
		if err == nil {
			joinDate = member.JoinedAt
		}
		joinDateDate, _ := joinDate.Parse()
		joinDateString := strings.Replace(joinDateDate.Format(time.RFC822Z), " +0000", "", -1)

		message += guild.ID + "," + guild.Name + "," + ownerName + "," + joinDateString + "\n"
	}
	go s.ChannelMessageDelete(msg.ChannelID, msg.ID)
	s.ChannelMessageSendComplex(dm.ID, &discordgo.MessageSend{
		Files: []*discordgo.File{
			{
				Name:   "servers.csv",
				Reader: strings.NewReader(message),
			},
		},
	})
}
