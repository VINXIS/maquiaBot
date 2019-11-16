package osucommands

import (
	"encoding/json"
	"io/ioutil"
	"regexp"
	"strings"
	"time"

	osuapi "../../osu-api"
	osutools "../../osu-functions"
	structs "../../structs"
	tools "../../tools"
	"github.com/bwmarrin/discordgo"
)

// Link links an osu! account with the discord user
func Link(s *discordgo.Session, m *discordgo.MessageCreate, args []string, osuAPI *osuapi.Client, cache []structs.PlayerData) {
	usernameRegex, _ := regexp.Compile(`(.+)(link|set)(\s+<@\S+)?(\s+.+)?`)

	discordUser := m.Author
	osuUsername := strings.TrimSpace(usernameRegex.FindStringSubmatch(m.Content)[4])

	// Obtain server and check admin permissions for linking with mentions involved
	server, err := s.Guild(m.GuildID)
	if err == nil {
		if m.Author.ID != server.OwnerID {
			member, _ := s.GuildMember(server.ID, m.Author.ID)
			admin := false
			for _, roleID := range member.Roles {
				role, _ := s.State.Role(m.GuildID, roleID)
				if role.Permissions&discordgo.PermissionAdministrator != 0 || role.Permissions&discordgo.PermissionManageServer != 0 {
					admin = true
					break
				}
			}
			if !admin && len(m.Mentions) >= 1 {
				s.ChannelMessageSend(m.ChannelID, "You must be an admin, server manager, or server owner!")
				return
			}
		}

		if len(m.Mentions) > 0 {
			discordUser = m.Mentions[0]
		}
	}

	// Run through the player cache to find the user using discord ID.
	for i, player := range cache {
		if player.Discord.ID == discordUser.ID {
			if strings.ToLower(player.Osu.Username) == strings.ToLower(osuUsername) {
				if len(m.Mentions) >= 1 {
					s.ChannelMessageSend(m.ChannelID, "osu! account **"+osuUsername+"** already been linked to "+discordUser.Username+"'s account!")
				} else {
					s.ChannelMessageSend(m.ChannelID, "osu! account **"+osuUsername+"** already linked to your discord account!")
				}
				return
			}

			user, err := osuAPI.GetUser(osuapi.GetUserOpts{
				Username: osuUsername,
			})
			if err != nil {
				s.ChannelMessageSend(m.ChannelID, "Player: **"+osuUsername+"** may not exist!")
				return
			}
			player.Time = time.Now()
			player.Osu = *user
			player.Farm = structs.FarmerdogData{}
			cache[i] = player

			jsonCache, err := json.Marshal(cache)
			tools.ErrRead(err)

			err = ioutil.WriteFile("./data/osuData/profileCache.json", jsonCache, 0644)
			tools.ErrRead(err)

			if len(m.Mentions) >= 1 {
				s.ChannelMessageSend(m.ChannelID, "osu! account **"+osuUsername+"** has been linked to "+discordUser.Username+"'s account!")
			} else {
				s.ChannelMessageSend(m.ChannelID, "osu! account **"+osuUsername+"** has been linked to your discord account!")
			}
			return
		}
	}

	// Run through the player cache to find the user using the osu! username if no discord ID exists.
	for i, player := range cache {
		if strings.ToLower(player.Osu.Username) == strings.ToLower(osuUsername) && player.Discord.ID == "" {
			player.Discord = *discordUser
			cache[i] = player

			jsonCache, err := json.Marshal(cache)
			tools.ErrRead(err)

			err = ioutil.WriteFile("./data/osuData/profileCache.json", jsonCache, 0644)
			tools.ErrRead(err)

			if len(m.Mentions) >= 1 {
				s.ChannelMessageSend(m.ChannelID, "osu! account **"+osuUsername+"** has been linked to "+discordUser.Username+"'s account!")
			} else {
				s.ChannelMessageSend(m.ChannelID, "osu! account **"+osuUsername+"** has been linked to your discord account!")
			}
			return
		}
	}

	// Create player
	user, err := osuAPI.GetUser(osuapi.GetUserOpts{
		Username: osuUsername,
	})
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "Player: **"+osuUsername+"** may not exist!")
		return
	}
	player := structs.PlayerData{
		Time:    time.Now(),
		Osu:     *user,
		Discord: *discordUser,
	}

	// Farm stuff
	farmData := structs.FarmData{}
	f, err := ioutil.ReadFile("./data/osuData/mapFarm.json")
	tools.ErrRead(err)
	_ = json.Unmarshal(f, &farmData)
	player = osutools.FarmCalc(player, osuAPI, farmData)

	// Save player
	cache = append(cache, player)
	jsonCache, err := json.Marshal(cache)
	tools.ErrRead(err)

	err = ioutil.WriteFile("./data/osuData/profileCache.json", jsonCache, 0644)
	tools.ErrRead(err)

	if len(m.Mentions) >= 1 {
		s.ChannelMessageSend(m.ChannelID, "osu! account **"+osuUsername+"** has been linked to "+discordUser.Username+"'s account!")
		return
	}

	s.ChannelMessageSend(m.ChannelID, "osu! account **"+osuUsername+"** has been linked to your discord account!")
	return
}
