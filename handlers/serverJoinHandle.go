package handlers

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"math"
	"strconv"
	"time"

	"maquiaBot/config"
	tools "maquiaBot/tools"

	"github.com/bwmarrin/discordgo"
)

// ServerJoin is to send a message when the bot joins a server
func ServerJoin(s *discordgo.Session, g *discordgo.GuildCreate) {
	s.UpdateStatus(0, "maquiahelp | maquiaprefix | "+strconv.Itoa(len(s.State.Guilds))+" servers")

	// Check for a guild ID
	if g.ID == "" {
		dm, err := s.UserChannelCreate(config.Conf.BotHoster.UserID)
		if err == nil {
			s.ChannelMessageSend(dm.ID, "An error occured in obtaining information from a server join.")
		}
		log.Println("An error occured in obtaining information from a server join Guild info below:")
		log.Println(g)
		return
	}

	// Obtain server data
	server, err := s.Guild(g.ID)
	if err != nil {
		dm, err := s.UserChannelCreate(config.Conf.BotHoster.UserID)
		if err == nil {
			s.ChannelMessageSend(dm.ID, "An error occured in obtaining information from a server join.")
			return
		}
		log.Println("An error occured in obtaining information from a server join Guild info below:")
		log.Println(g)
		return
	}
	serverData, _ := tools.GetServer(*server, s)

	// Check if bot was already in server or if server is unavailable
	joinTime, _ := g.JoinedAt.Parse()
	if g.Guild.Unavailable || math.Abs(joinTime.Sub(time.Now()).Seconds()) > 5 {
		return
	}

	_, err = s.ChannelMessageSend(g.ID, "Hello! My default prefix is `$` but you can change it by using `$prefix` or `maquiaprefix`.\nPlease note that if you remove this bot from the server, all custom configurations for the server such as custom prefix, quotes, stats, e.t.c will be removed. The `$downloadserver` command will allow you to obtain the server's data related to the bot if you ever kick the bot and re-invite it back to the server.\n\nFor information about this bot's commands, check out `$help` or `maquiahelp` to see the variety of commands created.\nFor any questions or concerns about this bot, please contact `@vinxis1` on twitter, or `VINXIS#1000` on discord.")
	if err != nil {
		for _, channel := range g.Channels {
			_, err := s.ChannelMessageSend(channel.ID, "Hello! My default prefix is `$` but you can change it by using `$prefix` or `maquiaprefix`.\nPlease note that if you remove this bot from the server, all custom configurations for the server such as custom prefix, quotes, stats, e.t.c will be removed. The `$downloadserver` command will allow you to obtain the server's data related to the bot if you ever kick the bot and re-invite it back to the server.\n\nFor information about this bot's commands, check out `$help` or `maquiahelp` to see the variety of commands created.\nFor any questions or concerns about this bot, please contact `@vinxis1` on twitter, or `VINXIS#1000` on discord.")
			if err == nil {
				serverData.AnnounceChannel = channel.ID
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
