package handlers

import (
	"math"
	"strconv"
	"time"

	"github.com/bwmarrin/discordgo"
)

// ServerJoin is to send a message when the bot joins a server
func ServerJoin(s *discordgo.Session, g *discordgo.GuildCreate) {
	s.UpdateStatus(0, "Chillin in "+strconv.Itoa(len(s.State.Guilds))+" servers")

	// Check if bot was already in server or if server is unavailable
	joinTime, _ := g.Guild.JoinedAt.Parse()
	if g.Guild.Unavailable || math.Abs(joinTime.Sub(time.Now()).Seconds()) > 5 {
		return
	}

	for _, channel := range g.Channels {
		if channel.ID == g.Guild.ID {
			s.ChannelMessageSend(channel.ID, "Hello! My default prefix is `$` but you can change it by using `$prefix` or `maquiaprefix`\nPlease remember that this bot is still currently under development so this bot may constantly go on and off as more features are being added!")
			return
		}
	}
	for _, channel := range g.Channels {
		_, err := s.ChannelMessageSend(channel.ID, "Hello! My default prefix is `$` but you can change it by using `$prefix` or `maquiaprefix`\nPlease remember that this bot is still currently under development so this bot may constantly go on and off as more features are being added!")
		if err == nil {
			return
		}
	}
}
