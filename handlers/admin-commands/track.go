package admincommands

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"strconv"
	"strings"

	osuapi "../../osu-api"
	osutools "../../osu-functions"
	structs "../../structs"
	tools "../../tools"
	"github.com/bwmarrin/discordgo"
)

// Track executes the track command, used for when people want to track/untrack users/pp/etc
func Track(s *discordgo.Session, m *discordgo.MessageCreate, osuAPI *osuapi.Client, mapCache []structs.MapData) {
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
	channelData, new := tools.GetChannel(*channel)

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
			osuUser, err := osuAPI.GetUser(osuapi.GetUserOpts{
				Username: user,
				Mode:     channelData.Mode,
			})
			if err != nil {
				continue
			}
			channelData.AddUser(*osuUser)
		}
	}

	// Map Status
	channelData.UpdateMapStatus(mapStatus)

	// Write data to JSON
	jsonCache, err := json.Marshal(channelData)
	tools.ErrRead(err)

	err = ioutil.WriteFile("./data/channelData/"+m.ChannelID+".json", jsonCache, 0644)
	tools.ErrRead(err)

	// Call trackpost if new, otherwise just post track information
	if new {
		go osutools.TrackPost(*channel, s, mapCache)
	}
	TrackInfo(s, m)
}

// TrackInfo gives info about what's being tracked in the channel currently
func TrackInfo(s *discordgo.Session, m *discordgo.MessageCreate) {
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

	// Obtain channel data
	channelData, new := tools.GetChannel(*channel)
	if new {
		s.ChannelMessageSend(m.ChannelID, "There is no tracking info for this channel currently!")
		return
	}

	serverImg := "https://cdn.discordapp.com/icons/" + server.ID + "/" + server.Icon
	if strings.Contains(server.Icon, "a_") {
		serverImg += ".gif"
	} else {
		serverImg += ".png"
	}

	// Create embed
	embed := discordgo.MessageEmbed{
		Author: &discordgo.MessageEmbedAuthor{
			Name:    channelData.Channel.Name,
			IconURL: serverImg,
		},
		Description: "Any admin, server moderator, or server owner can update this using `track` again! \n**You do not need to readd everything again to update parts of the tracker.** Toggle tracking on or off using `tt`, `trackt`, `ttoggle`, or `tracktoggle`",
		Color:       osutools.ModeColour(channelData.Mode),
	}

	// Warnings
	warning := ""
	if !channelData.Ranked && !channelData.Loved && !channelData.Qualified {
		warning += "\n**WARNING:** You do not have any map rank statuses with leaderboards enabled! Please enable at least one in order for tracking to work!"
	}
	if channelData.LeaderboardReq == 101 && channelData.TopReq == 101 && channelData.PPReq == -1 {
		warning += "\n**WARNING:** You do not have any leader/top/pp requirement for scores! Any score submitted by the users listed on eligible maps will be posted as a result!"
	}
	if !channelData.Tracking {
		warning += "\n**WARNING:** Tracking is currently turned off!"
	}

	// Add users
	userList := ""
	for i, user := range channelData.Users {
		if i == len(channelData.Users)-1 {
			userList += user.Username
		} else {
			userList += user.Username + ", "
		}
	}
	embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
		Name:  "Users Tracked:",
		Value: userList,
	})

	// Add command info

	statuses := ""
	if channelData.Ranked {
		statuses += "r"
	}
	if channelData.Loved {
		statuses += "l"
	}
	if channelData.Qualified {
		statuses += "q"
	}
	statuses = strings.Join(strings.Split(statuses, ""), ", ")

	command :=
		"`-u " + userList +
			" -pp " + strconv.FormatFloat(channelData.PPReq, 'f', 0, 64) +
			" -l " + strconv.Itoa(channelData.LeaderboardReq) +
			" -t " + strconv.Itoa(channelData.TopReq) +
			" -m " + channelData.Mode.String() +
			" -s " + statuses + "`"

	embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
		Name:  "Command for This Config:",
		Value: command,
	})

	// Add PP req
	if channelData.PPReq != -1 {
		embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{"PP Req:", strconv.FormatFloat(channelData.PPReq, 'f', 0, 64), true})
	} else {
		embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{"PP Req:", "N/A", true})
	}

	// Add Leaderboard req
	if channelData.LeaderboardReq != 101 {
		embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{"Leaderboard Req:", strconv.Itoa(channelData.LeaderboardReq), true})
	} else {
		embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{"Leaderboard Req:", "N/A", true})
	}

	// Add Top req
	if channelData.TopReq != 101 {
		embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{"Top Perf. Req:", strconv.Itoa(channelData.TopReq), true})
	} else {
		embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{"Top Perf. Req:", "N/A", true})
	}

	// Map types
	embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{"Ranked: ", strconv.FormatBool(channelData.Ranked), true})
	embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{"Loved: ", strconv.FormatBool(channelData.Loved), true})
	embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{"Qualified: ", strconv.FormatBool(channelData.Qualified), true})

	// Misc.
	embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{"Mode: ", channelData.Mode.String(), true})
	embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{"Tracking: ", strconv.FormatBool(channelData.Tracking), true})

	s.ChannelMessageSendComplex(m.ChannelID, &discordgo.MessageSend{
		Content: warning,
		Embed:   &embed,
	})
	return
}

// TrackToggle stops tracking for the channel
func TrackToggle(s *discordgo.Session, m *discordgo.MessageCreate, mapCache []structs.MapData) {
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
	channelData, new := tools.GetChannel(*channel)
	if new {
		s.ChannelMessageSend(m.ChannelID, "There is no tracking info for this channel currently!")
		return
	}

	// The Main Event
	channelData.Tracking = !channelData.Tracking

	// Write data to JSON
	jsonCache, err := json.Marshal(channelData)
	tools.ErrRead(err)

	err = ioutil.WriteFile("./data/channelData/"+m.ChannelID+".json", jsonCache, 0644)
	tools.ErrRead(err)

	if channelData.Tracking {
		go osutools.TrackPost(*channel, s, mapCache)
		s.ChannelMessageSend(m.ChannelID, "Started tracking for this channel!")
	} else {
		s.ChannelMessageSend(m.ChannelID, "Stopped tracking for this channel!")
	}
}
