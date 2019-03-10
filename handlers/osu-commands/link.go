package osucommands

import (
	"encoding/json"
	"io/ioutil"
	"time"

	structs "../../structs"
	tools "../../tools"
	"github.com/bwmarrin/discordgo"
	"github.com/thehowl/go-osuapi"
)

// Link links an osu! account with the discord user
func Link(s *discordgo.Session, m *discordgo.MessageCreate, args []string, osuAPI *osuapi.Client, cache []structs.PlayerData) {
	discordUser := m.Author
	osuUsername := args[2]

	for _, player := range cache {
		if player.Discord == *discordUser {
			if player.Osu.Username == osuUsername {
				s.ChannelMessageSend(m.ChannelID, "Player: **"+osuUsername+"** already assigned to this account!")
				return
			}

			user, err := osuAPI.GetUser(osuapi.GetUserOpts{
				Username: osuUsername,
			})
			if err != nil {
				s.ChannelMessageSend(m.ChannelID, "Player: **"+osuUsername+"** does not exist!")
				return
			}
			player.Time = time.Now()
			player.Osu = *user

			jsonCache, err := json.Marshal(cache)
			tools.ErrRead(err)

			err = ioutil.WriteFile("./data/osuData/profileCache.json", jsonCache, 0644)
			tools.ErrRead(err)

			s.ChannelMessageSend(m.ChannelID, "Player: **"+osuUsername+"** has now been assigned to this account!")
			return
		}
		emptyDiscordUser := discordgo.User{}
		if player.Osu.Username == osuUsername && player.Discord == emptyDiscordUser {
			player.Discord = *discordUser

			jsonCache, err := json.Marshal(cache)
			tools.ErrRead(err)

			err = ioutil.WriteFile("./data/osuData/profileCache.json", jsonCache, 0644)
			tools.ErrRead(err)

			s.ChannelMessageSend(m.ChannelID, "Player: **"+osuUsername+"** has now been assigned to this account!")
			return
		}
	}

	user, err := osuAPI.GetUser(osuapi.GetUserOpts{
		Username: osuUsername,
	})
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "Player: **"+osuUsername+"** does not exist!")
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

	s.ChannelMessageSend(m.ChannelID, "Player: **"+osuUsername+"** has now been assigned to this account!")
	return
}
