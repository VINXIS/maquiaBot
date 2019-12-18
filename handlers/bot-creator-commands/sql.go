package botcreatorcommands

import (
	"fmt"
	"regexp"

	config "../../config"
	tools "../../tools"
	"github.com/bwmarrin/discordgo"
)

// SQL lets the bot creator do SQL stuff
func SQL(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID != config.Conf.BotHoster.UserID {
		s.ChannelMessageSend(m.ChannelID, "YOU ARE NOT "+config.Conf.BotHoster.Username+".........")
		return
	}
	sqlregex, err := regexp.Compile(`sql\s+(.+)`)
	query := sqlregex.FindStringSubmatch(m.Content)[1]

	res, err := tools.DB.Query(query)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, err.Error())
		return
	}
	s.ChannelMessageSend(m.ChannelID, "Success")
	fmt.Println(res)
	res.Close()
	return
}
