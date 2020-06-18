package botcreatorcommands

import (
	"regexp"
	"strconv"

	config "../../config"
	"github.com/bwmarrin/discordgo"
)

// Update updates the bot discord status
func Update(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID != config.Conf.BotHoster.UserID {
		s.ChannelMessageSend(m.ChannelID, "YOU ARE NOT "+config.Conf.BotHoster.Username+".........")
		return
	}
	updateRegex, _ := regexp.Compile(`(?i)updatestatus\s+(.+)`)
	if !updateRegex.MatchString(m.Content) {
		s.UpdateStatus(0, strconv.Itoa(len(s.State.Guilds))+" servers")
		return
	}

	text := updateRegex.FindStringSubmatch(m.Content)[1]
	s.UpdateStatus(0, text)
}
