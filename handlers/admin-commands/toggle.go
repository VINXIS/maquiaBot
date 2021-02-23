package admincommands

import (
	"encoding/json"
	"io/ioutil"
	"strconv"
	"strings"
	"time"

	gencommands "maquiaBot/handlers/general-commands"
	tools "maquiaBot/tools"

	"github.com/bwmarrin/discordgo"
)

// Toggle toggles server options on/off
func Toggle(s *discordgo.Session, m *discordgo.MessageCreate) {
	server, err := s.Guild(m.GuildID)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "This is not a server!")
		return
	}

	channel, err := s.Channel(m.ChannelID)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "This is not an allowed channel!")
		return
	}

	if !tools.AdminCheck(s, m, *server) {
		s.ChannelMessageSend(m.ChannelID, "You must be an admin, server manager, or server owner!")
		return
	}

	// Obtain server data
	serverData, _ := tools.GetServer(*server, s)
	channelData, _ := tools.GetChannel(*channel, s)

	// Set new information in server data
	serverData.Time = time.Now()
	flagged := false
	target := ""
	status := "false"
	if strings.Contains(m.Content, "-a") || strings.Contains(m.Content, "-announce") {
		if serverData.AnnounceChannel == m.ChannelID {
			serverData.AnnounceChannel = ""
			status = "N/A"
		} else {
			serverData.AnnounceChannel = m.ChannelID
			status = "this channel"
		}
		flagged = true
		target = "The announcement channel"
	}
	if strings.Contains(m.Content, "-s") || strings.Contains(m.Content, "-stats") {
		serverData.AllowAnyoneStats = !serverData.AllowAnyoneStats
		flagged = true
		target = "`AllowAnyoneStats`"
		status = strconv.FormatBool(serverData.AllowAnyoneStats)
	}
	if strings.Contains(m.Content, "-d") || strings.Contains(m.Content, "-daily") {
		if strings.Contains(m.Content, "-ch") || strings.Contains(m.Content, "-channel") {
			channelData.Daily = !channelData.Daily
		} else {
			serverData.Daily = !serverData.Daily
		}
		flagged = true
		target = "`daily`"
		status = strconv.FormatBool(serverData.Daily)
	}
	if strings.Contains(m.Content, "-os") || strings.Contains(m.Content, "-osu") {
		if strings.Contains(m.Content, "-ch") || strings.Contains(m.Content, "-channel") {
			channelData.OsuToggle = !channelData.OsuToggle
		} else {
			serverData.OsuToggle = !serverData.OsuToggle
		}
		flagged = true
		target = "`osuToggle`"
		status = strconv.FormatBool(serverData.OsuToggle)
	}
	if strings.Contains(m.Content, "-t") || strings.Contains(m.Content, "-time") || strings.Contains(m.Content, "-timestamp") {
		if strings.Contains(m.Content, "-ch") || strings.Contains(m.Content, "-channel") {
			channelData.TimestampToggle = !channelData.TimestampToggle
		} else {
			serverData.TimestampToggle = !serverData.TimestampToggle
		}
		flagged = true
		target = "`timestampToggle`"
		status = strconv.FormatBool(serverData.TimestampToggle)
	}
	if strings.Contains(m.Content, "-v") || strings.Contains(m.Content, "-vibe") {
		if strings.Contains(m.Content, "-ch") || strings.Contains(m.Content, "-channel") {
			channelData.Vibe = !channelData.Vibe
		} else {
			serverData.Vibe = !serverData.Vibe
		}
		flagged = true
		target = "`Vibe`"
		status = strconv.FormatBool(serverData.Vibe)
	}
	if !flagged {
		s.ChannelMessageSend(m.ChannelID, "No flags given! Please use one of the flags listed in `help toggle`!")
		return
	}

	jsonCache, err := json.Marshal(serverData)
	tools.ErrRead(s, err)

	err = ioutil.WriteFile("./data/serverData/"+m.GuildID+".json", jsonCache, 0644)
	tools.ErrRead(s, err)

	s.ChannelMessageSend(m.ChannelID, target+" has now been set to "+status+". Obtaining server info...")
	go gencommands.ServerInfo(s, m)
}
