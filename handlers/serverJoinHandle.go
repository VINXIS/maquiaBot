package handlers

import (
	"encoding/json"
	"io/ioutil"
	"math"
	"strconv"
	"time"

	tools "maquiaBot/tools"

	"github.com/bwmarrin/discordgo"
)

// ServerJoin is to send a message when the bot joins a server
func ServerJoin(s *discordgo.Session, g *discordgo.GuildCreate) {
	s.UpdateStatus(0, strconv.Itoa(len(s.State.Guilds))+" servers")

	// Obtain server data
	server, _ := s.Guild(g.ID)
	serverData := tools.GetServer(*server, s)

	// Check if bot was already in server or if server is unavailable
	joinTime, _ := g.JoinedAt.Parse()
	if g.Guild.Unavailable || math.Abs(joinTime.Sub(time.Now()).Seconds()) > 5 {
		return
	}

	_, err := s.ChannelMessageSend(g.ID, "Hello! My default prefix is `$` but you can change it by using `$prefix` or `maquiaprefix`\nPlease remember that this bot is still currently under development so this bot may constantly go on and off as more features are being added!")
	if err != nil {
		for _, channel := range g.Channels {
			serverData.AnnounceChannel = channel.ID
			_, err := s.ChannelMessageSend(channel.ID, "Hello! My default prefix is `$` but you can change it by using `$prefix` or `maquiaprefix`\nPlease remember that this bot is still currently under development so this bot may constantly go on and off as more features are being added!")
			if err == nil {
				break
			}
		}
	} else {
		serverData.AnnounceChannel = g.ID
	}

	jsonCache, err := json.Marshal(serverData)
	tools.ErrRead(s, err)

	err = ioutil.WriteFile("./data/serverData/"+g.ID+".json", jsonCache, 0644)
	tools.ErrRead(s, err)
}
