package admincommands

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
	osucommands "maquiaBot/handlers/osu-commands"
	osuapi "maquiaBot/osu-api"
	osutools "maquiaBot/osu-tools"
	tools "maquiaBot/tools"
)

// Track executes the track command, used for when people want to track/untrack users/pp/etc
func Track(s *discordgo.Session, m *discordgo.MessageCreate) {
	channel, err := s.Channel(m.ChannelID)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "This is not an allowed channel!")
		return
	}

	server, err := s.Guild(m.GuildID)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "This is not a server!")
		return
	}

	if !tools.AdminCheck(s, m, *server) {
		s.ChannelMessageSend(m.ChannelID, "You must be an admin, server manager, or server owner!")
		return
	}

	// Obtain channel data
	channelData, new := tools.GetChannel(*channel, s)

	// Get params
	args := strings.Split(m.Content, " ")[1:]
	if len(args) == 0 {
		s.ChannelMessageSend(m.ChannelID, "No options given!")
		return
	}

	// Check if everything should just be removed
	if len(args) == 1 && (args[0] == "r" || args[0] == "rem" || args[0] == "remove") {
		_, err := os.Stat("./data/channelData/" + channel.ID + ".json")
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, "No tracking exists for this channel currently!")
			return
		}
		err = os.Remove("./data/channelData/" + channel.ID + ".json")
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, "An error has occurred in removing your channel's tracking.")
			return
		}
		s.ChannelMessageSend(m.ChannelID, "Successfully removed channel!")
		return
	}

	mode := ""
	var (
		users, mapStatus []string
		pp               float64
		leaderboard, top int
	)

	for _, arg := range args {
		arg = strings.TrimSpace(arg)

		// If flag change mode
		switch arg {
		case "-u":
			mode = "user"
			continue
		case "-pp":
			mode = "pp"
			continue
		case "-l", "-leader", "-leaderboard":
			mode = "leaderboard"
			continue
		case "-t", "-top":
			mode = "top"
			continue
		case "-s", "-status":
			mode = "status"
			continue
		case "-m", "-mode":
			mode = "mode"
			continue
		}

		// Obtain value if not flag
		switch mode {
		case "user":
			users = append(users, arg)
		case "pp":
			pp, err = strconv.ParseFloat(arg, 64)
			if err != nil || pp <= 0 {
				pp = -1
			}
			channelData.PPReq = pp
			mode = ""
		case "leaderboard":
			leaderboard, err = strconv.Atoi(arg)
			if err != nil || leaderboard <= 0 || leaderboard > 100 {
				leaderboard = 101
			}
			channelData.LeaderboardReq = leaderboard
			mode = ""
		case "top":
			top, err = strconv.Atoi(arg)
			if err != nil || top <= 0 || top > 100 {
				top = 101
			}
			channelData.TopReq = top
			mode = ""
		case "status":
			mapStatus = append(mapStatus, arg)
		case "mode":
			switch arg {
			case "0", "s", "std", "standard", "osu!s", "osu!std", "osu!standard":
				channelData.Mode = osuapi.ModeOsu
			case "1", "t", "tko", "taiko", "osu!t", "osu!tko", "osu!taiko":
				channelData.Mode = osuapi.ModeTaiko
			case "2", "c", "ctb", "catch", "osu!c", "catchthebeat", "osu!ctb", "osu!catch", "osu!catchthebeat":
				channelData.Mode = osuapi.ModeCatchTheBeat
			case "3", "m", "man", "mania", "osu!m", "osu!man", "osu!mania":
				channelData.Mode = osuapi.ModeOsuMania
			}
			mode = ""
		}
	}

	// Obtain users and map types given
	users = strings.Split(strings.Join(users, " "), ", ")
	mapStatus = strings.Split(strings.Join(mapStatus, " "), ", ")

	// Users
	if args[0] == "r" || args[0] == "rem" || args[0] == "remove" {
		channelData.RemoveUser(users)
	} else {
		for _, user := range users {
			osuUser, err := osucommands.OsuAPI.GetUser(osuapi.GetUserOpts{
				Username: user,
				Mode:     channelData.Mode,
			})
			if err != nil {
				continue
			}
			channelData.AddUser(*osuUser)
		}
	}

	if len(channelData.Users) == 0 && new {
		s.ChannelMessageSend(m.ChannelID, "No users given! Please use the `-u` flag to add users, and separate users with commas!")
		return
	}

	// Map Status
	channelData.UpdateMapStatus(mapStatus)

	// Write data to JSON
	jsonCache, err := json.Marshal(channelData)
	tools.ErrRead(s, err)

	err = ioutil.WriteFile("./data/channelData/"+m.ChannelID+".json", jsonCache, 0644)
	tools.ErrRead(s, err)

	// Call trackpost if new, otherwise just post track information
	if new {
		go osutools.TrackPost(*channel, s)
	}
	osucommands.TrackInfo(s, m)
}

// TrackToggle stops tracking for the channel
func TrackToggle(s *discordgo.Session, m *discordgo.MessageCreate) {
	channel, err := s.Channel(m.ChannelID)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "This is not an allowed channel!")
		return
	}

	// Check perms
	server, err := s.Guild(m.GuildID)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "This is not a server so custom prefixes are unavailable! Please use `$` instead for commands!")
		return
	}

	if !tools.AdminCheck(s, m, *server) {
		s.ChannelMessageSend(m.ChannelID, "You must be an admin, server manager, or server owner!")
		return
	}

	// Obtain channel data
	channelData, new := tools.GetChannel(*channel, s)
	if new {
		s.ChannelMessageSend(m.ChannelID, "There is no tracking info for this channel currently!")
		return
	}

	// The Main Event
	channelData.Tracking = !channelData.Tracking

	// Write data to JSON
	jsonCache, err := json.Marshal(channelData)
	tools.ErrRead(s, err)

	err = ioutil.WriteFile("./data/channelData/"+m.ChannelID+".json", jsonCache, 0644)
	tools.ErrRead(s, err)

	if channelData.Tracking {
		go osutools.TrackPost(*channel, s)
		s.ChannelMessageSend(m.ChannelID, "Started tracking for this channel!")
	} else {
		s.ChannelMessageSend(m.ChannelID, "Stopped tracking for this channel!")
	}
}
