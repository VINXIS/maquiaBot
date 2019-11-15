package osucommands

import (
	"encoding/json"
	"fmt"
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
func Track(s *discordgo.Session, m *discordgo.MessageCreate, args []string, osuAPI *osuapi.Client, mapCache []structs.MapData) {
	// Check perms
	channel, err := s.Channel(m.ChannelID)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "This is not a server so custom prefixes are unavailable! Please use `$` instead for commands!")
		return
	}

	server, err := s.Guild(m.GuildID)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "This is not a server so custom prefixes are unavailable! Please use `$` instead for commands!")
		return
	}

	member := &discordgo.Member{}
	for _, guildMember := range server.Members {
		if guildMember.User.ID == m.Author.ID {
			member = guildMember
		}
	}

	if member.User.ID == "" {
		return
	}

	admin := false
	for _, roleID := range member.Roles {
		role, err := s.State.Role(m.GuildID, roleID)
		tools.ErrRead(err)
		if role.Permissions&discordgo.PermissionAdministrator != 0 || role.Permissions&discordgo.PermissionManageServer != 0 {
			admin = true
			break
		}
	}

	if !admin && m.Author.ID != server.OwnerID {
		s.ChannelMessageSend(m.ChannelID, "You are not an admin or server manager!")
		return
	}

	// Obtain channel data
	channelData := structs.ChannelData{
		Channel: *channel,
	}
	new := true
	_, err = os.Stat("./data/channelData/" + m.ChannelID + ".json")
	if err == nil {
		f, err := ioutil.ReadFile("./data/channelData/" + m.ChannelID + ".json")
		tools.ErrRead(err)
		_ = json.Unmarshal(f, &channelData)
		new = false
	} else if os.IsNotExist(err) {
		channelData.Channel = *channel
	} else {
		fmt.Println(err)
		s.ChannelMessageSend(m.ChannelID, "An error occurred! VINXIS has obtained error info.")
		return
	}

	// Get params
	pp := ""
	top := ""
	users := []string{}
	addition := true
	removal := false
	multiUser := false
	for i, arg := range args {
		if arg == "replace" {
			addition = false
		} else if arg == "remove" {
			addition = false
			removal = true
		}
		if i != len(args)-1 {
			switch strings.ToLower(arg) {
			case "pp":
				pp = args[i+1]
			case "top":
				top = args[i+1]
			case "user":
				users = append(users, args[i+1])
			case "users":
				multiUser = true
			}
		}
		if multiUser {
			users = append(users, arg)
		}
	}

	// Check if any params were NOT given or if add and remove were both stated in the message
	if pp == "" && top == "" && len(users) == 0 && !removal {
		s.ChannelMessageSend(m.ChannelID, "Not enough params given! Need at least one of `pp` or `top` or `user(s)`")
		return
	}
	if addition && removal {
		s.ChannelMessageSend(m.ChannelID, "You cannot add and remove at the same time!")
		return
	}

	// Add/Remove as needed
	if !addition && !removal {
		channelData.Users = []osuapi.User{}
	}
	if removal {
		if len(users) == 0 {
			tools.DeleteFile("./data/channelData/" + m.ChannelID + ".json")
			s.ChannelMessageSend(m.ChannelID, "Completely removed tracking for this channel!")
			return
		}
		text := "Removed: "
		tracked := false
		for _, user := range users {
			for i, osuUser := range channelData.Users {
				if strings.ToLower(osuUser.Username) == strings.ToLower(user) {
					tracked = true
					text = text + user + " "
					channelData.Users = append(channelData.Users[:i], channelData.Users[i+1:]...)
				}
			}
		}
		if !tracked {
			s.ChannelMessageSend(m.ChannelID, "User is not being tracked currently!")
			return
		}
		// Write data to JSON
		jsonCache, err := json.Marshal(channelData)
		tools.ErrRead(err)

		err = ioutil.WriteFile("./data/channelData/"+m.ChannelID+".json", jsonCache, 0644)
		tools.ErrRead(err)
		s.ChannelMessageSend(m.ChannelID, text)
		return
	}

	// Assign params to data
	text := "Now tracking "
	if len(users) != 0 {
		for _, user := range users {
			user = strings.ReplaceAll(user, ",", "")
			userRes, err := osuAPI.GetUser(osuapi.GetUserOpts{
				Username: user,
			})
			if err != nil {
				s.ChannelMessageSend(m.ChannelID, "User "+user+" may not exist! Check to make sure they aren't banned/ that you typed their name properly!")
				return
			}
			if addition {
				for _, osuUser := range channelData.Users {
					if osuUser.Username == userRes.Username {
						s.ChannelMessageSend(m.ChannelID, "User "+user+" is already being tracked!")
						return
					}
				}
			}
			channelData.Users = append(channelData.Users, *userRes)
			text = text + user + " "
		}
	}
	if pp != "" {
		ppint, err := strconv.Atoi(pp)
		if ppint < 0 {
			s.ChannelMessageSend(m.ChannelID, "Invalid paramater for `pp`")
			return
		}
		if err == nil {
			channelData.PPLimit = ppint
		}
		text = text + "with a pp limit of at least " + pp + "pp "
	}
	if top != "" {
		topint, err := strconv.Atoi(top)
		if topint < 0 {
			s.ChannelMessageSend(m.ChannelID, "Invalid paramater for `top`")
			return
		}
		if err == nil && topint <= 100 {
			channelData.TopPlay = topint
		} else {
			topint = 100
			channelData.TopPlay = topint
		}
		if pp != "" {
			text = text + "or if the score is a top " + strconv.Itoa(topint) + " score"
		}
		text = text + "if their score is a top " + strconv.Itoa(topint) + " score"
	} else {
		channelData.TopPlay = 100
	}
	channelData.Tracking = true

	// Write data to JSON
	jsonCache, err := json.Marshal(channelData)
	tools.ErrRead(err)

	err = ioutil.WriteFile("./data/channelData/"+m.ChannelID+".json", jsonCache, 0644)
	tools.ErrRead(err)

	if new {
		go osutools.TrackPost("data/channelData/"+m.ChannelID+".json", s, mapCache)
	}

	s.ChannelMessageSend(m.ChannelID, text)
	return
}
