package osucommands

import (
	"encoding/json"
	"io/ioutil"
	"strings"
	"time"

	structs "../../structs"
	tools "../../tools"
	"github.com/bwmarrin/discordgo"
	"github.com/thehowl/go-osuapi"
)

// Link links an osu! account with the discord user
func Link(s *discordgo.Session, m *discordgo.MessageCreate, args []string, osuAPI *osuapi.Client, cache []structs.PlayerData) {
	discordUser := m.Author
	osuUsername := ""
	if len(args) > 2 {
		if len(m.Mentions) >= 1 && len(args) > 3 {
			osuUsername = args[3]
		} else if len(m.Mentions) >= 1 {
			osuUsername = args[2]
		}
	} else {
		osuUsername = args[1]
	}

	server, err := s.Guild(m.GuildID)
	member := &discordgo.Member{}
	admin := false

	if err == nil {
		for _, guildMember := range server.Members {
			if guildMember.User.ID == m.Author.ID {
				member = guildMember
			}
		}

		for _, roleID := range member.Roles {
			role, err := s.State.Role(m.GuildID, roleID)
			tools.ErrRead(err)
			if role.Permissions&discordgo.PermissionAdministrator == discordgo.PermissionAdministrator {
				admin = true
				break
			}
		}
		if m.Author.ID == server.OwnerID {
			admin = true
		}
	}

	if len(m.Mentions) >= 1 && admin {
		discordUser = m.Mentions[0]
	} else if len(m.Mentions) >= 1 && !admin {
		s.ChannelMessageSend(m.ChannelID, "You are not an admin!")
		return
	}

	for i, player := range cache {
		if player.Discord.ID == discordUser.ID {
			if strings.ToLower(player.Osu.Username) == strings.ToLower(osuUsername) {
				if len(m.Mentions) >= 1 {
					s.ChannelMessageSend(m.ChannelID, "osu! account **"+osuUsername+"** already been linked to "+discordUser.Username+"'s account!")
					return
				}
				s.ChannelMessageSend(m.ChannelID, "osu! account **"+osuUsername+"** already linked to your discord account!")
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
			cache[i] = player

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
		if strings.ToLower(player.Osu.Username) == strings.ToLower(osuUsername) && player.Discord.ID == "" {
			player.Discord = *discordUser
			cache[i] = player

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
	}

	user, err := osuAPI.GetUser(osuapi.GetUserOpts{
		Username: osuUsername,
	})
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "Player: **"+osuUsername+"** may not exist!")
		return
	}

	cache = append(cache, structs.PlayerData{
		Time:    time.Now(),
		Osu:     *user,
		Discord: *discordUser,
	})
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
