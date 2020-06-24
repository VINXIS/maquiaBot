package osucommands

import (
	"encoding/json"
	"io/ioutil"
	"regexp"
	"strings"
	"time"

	osuapi "maquiaBot/osu-api"
	structs "maquiaBot/structs"
	tools "maquiaBot/tools"

	"github.com/bwmarrin/discordgo"
)

// Link links an osu! account with the discord user
func Link(s *discordgo.Session, m *discordgo.MessageCreate, args []string, cache []structs.PlayerData) {
	usernameRegex, _ := regexp.Compile(`(?i)(link|set)(\s+<@\S+)?(\s+.+)?`)

	discordUser := m.Author
	osuUsername := strings.TrimSpace(usernameRegex.FindStringSubmatch(m.Content)[3])

	farmData := structs.FarmData{}
	f, err := ioutil.ReadFile("./data/osuData/mapFarm.json")
	tools.ErrRead(s, err)
	_ = json.Unmarshal(f, &farmData)

	// Obtain server and check admin permissions for linking with mentions involved
	server, err := s.Guild(m.GuildID)
	if err == nil {
		if !tools.AdminCheck(s, m, *server) && len(m.Mentions) > 0 {
			s.ChannelMessageSend(m.ChannelID, "You must be an admin, server manager, or server owner!")
			return
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

			user, err := OsuAPI.GetUser(osuapi.GetUserOpts{
				Username: osuUsername,
			})
			if err != nil {
				s.ChannelMessageSend(m.ChannelID, "Player: **"+osuUsername+"** may not exist!")
				return
			}
			player.Time = time.Now()
			player.Osu = *user
			player.FarmCalc(OsuAPI, farmData)
			cache[i] = player

			// Remove any accounts of the same user or empty osu! user and with no discord linked
			for j := 0; j < len(cache); j++ {
				if player.Discord.ID == "" && (cache[j].Osu.Username == "" || strings.ToLower(cache[j].Osu.Username) == strings.ToLower(osuUsername)) {
					cache = append(cache[:j], cache[j+1:]...)
					j--
				}
			}

			jsonCache, err := json.Marshal(cache)
			tools.ErrRead(s, err)

			err = ioutil.WriteFile("./data/osuData/profileCache.json", jsonCache, 0644)
			tools.ErrRead(s, err)

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
			player.FarmCalc(OsuAPI, farmData)
			cache[i] = player

			jsonCache, err := json.Marshal(cache)
			tools.ErrRead(s, err)

			err = ioutil.WriteFile("./data/osuData/profileCache.json", jsonCache, 0644)
			tools.ErrRead(s, err)

			if len(m.Mentions) >= 1 {
				s.ChannelMessageSend(m.ChannelID, "osu! account **"+osuUsername+"** has been linked to "+discordUser.Username+"'s account!")
			} else {
				s.ChannelMessageSend(m.ChannelID, "osu! account **"+osuUsername+"** has been linked to your discord account!")
			}
			return
		}
	}

	// Create player
	user, err := OsuAPI.GetUser(osuapi.GetUserOpts{
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

	// Farm calc
	player.FarmCalc(OsuAPI, farmData)

	// Save player
	cache = append(cache, player)
	jsonCache, err := json.Marshal(cache)
	tools.ErrRead(s, err)

	err = ioutil.WriteFile("./data/osuData/profileCache.json", jsonCache, 0644)
	tools.ErrRead(s, err)

	if len(m.Mentions) >= 1 {
		s.ChannelMessageSend(m.ChannelID, "osu! account **"+osuUsername+"** has been linked to "+discordUser.Username+"'s account!")
		return
	}

	s.ChannelMessageSend(m.ChannelID, "osu! account **"+osuUsername+"** has been linked to your discord account!")
	return
}
