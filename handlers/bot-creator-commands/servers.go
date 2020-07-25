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

	dm, _ := s.UserChannelCreate(config.Conf.BotHoster.UserID)
	message := "ID,Name,Owner,Date Joined\n"
	for _, guild := range s.State.Guilds {
		owner, _ := s.User(guild.OwnerID)

		var joinDate discordgo.Timestamp
		members, err := s.GuildMembers(guild.ID, "", 1000)
		if err == nil {
			for _, member := range members {
				if member.User.ID == s.State.User.ID {
					joinDate = member.JoinedAt
					break
				}
			}
		}
		joinDateDate, _ := joinDate.Parse()
		joinDateString := strings.Replace(joinDateDate.Format(time.RFC822Z), " +0000", "", -1)

		message += guild.ID + "," + guild.Name + "," + owner.String() + "," + joinDateString + "\n"
	}
	s.ChannelMessageSendComplex(dm.ID, &discordgo.MessageSend{
		Files: []*discordgo.File{
			{
				Name:   "servers.csv",
				Reader: strings.NewReader(message),
			},
		},
	})
}
