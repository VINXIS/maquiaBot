package botcreatorcommands

import (
	"regexp"
	"strconv"

	config "maquiaBot/config"

	"github.com/bwmarrin/discordgo"
)

// Update updates the bot discord status
func Update(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID != config.Conf.BotHoster.UserID {
		s.ChannelMessageSend(m.ChannelID, "YOU ARE NOT "+config.Conf.BotHoster.Username+".........")
		return
	}
	updateRegex, _ := regexp.Compile(`(?i)up(date)?\s+(.+)`)
	if !updateRegex.MatchString(m.Content) {
		s.UpdateStatus(0, "maquiahelp | maquiaprefix | "+strconv.Itoa(len(s.State.Guilds))+" servers")
		return
	}

	text := updateRegex.FindStringSubmatch(m.Content)[2]
	s.UpdateStatus(0, text)

	s.ChannelMessageSend(m.ChannelID, "Updated status.")
}
