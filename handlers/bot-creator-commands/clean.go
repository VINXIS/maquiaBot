package botcreatorcommands

import (
	"encoding/json"
	"io/ioutil"

	config "maquiaBot/config"
	osuapi "maquiaBot/osu-api"
	structs "maquiaBot/structs"
	tools "maquiaBot/tools"
	osucommands "maquiaBot/handlers/osu-commands"
	"github.com/bwmarrin/discordgo"
)

// Clean cleans the caches
func Clean(s *discordgo.Session, m *discordgo.MessageCreate, cache []structs.PlayerData) {
	if m.Author.ID != config.Conf.BotHoster.UserID {
		s.ChannelMessageSend(m.ChannelID, "YOU ARE NOT "+config.Conf.BotHoster.Username+".........")
		return
	}

	keys := make(map[string]bool)
	newPlayerCache := []structs.PlayerData{}
	for _, player := range cache {
		if player.Discord.ID != "" {
			if _, value := keys[player.Discord.ID]; !value {
				keys[player.Discord.ID] = true
				newPlayerCache = append(newPlayerCache, player)
			}
		}
	}

	jsonCache, err := json.Marshal(newPlayerCache)
	tools.ErrRead(s, err)

	err = ioutil.WriteFile("./data/osuData/profileCache.json", jsonCache, 0644)
	tools.ErrRead(s, err)
	s.ChannelMessageSend(m.ChannelID, "Cleaned player cache!")
}

// CleanFarm cleans the all farmerdog ratings
func CleanFarm(s *discordgo.Session, m *discordgo.MessageCreate, cache []structs.PlayerData) {
	if m.Author.ID != config.Conf.BotHoster.UserID {
		s.ChannelMessageSend(m.ChannelID, "YOU ARE NOT "+config.Conf.BotHoster.Username+".........")
		return
	}

	// Farm Data
	farmData := structs.FarmData{}
	f, err := ioutil.ReadFile("./data/osuData/mapFarm.json")
	tools.ErrRead(s, err)
	_ = json.Unmarshal(f, &farmData)

	// Update
	for i := range cache {
		if cache[i].Osu.Username != "" && cache[i].Farm.Rating == 0.00 {
			cache[i].Osu = osuapi.User{}
		}
		cache[i].FarmCalc(osucommands.OsuAPI, farmData)
	}

	// Save
	jsonCache, err := json.Marshal(cache)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "Error with wiping data!")
		tools.ErrRead(s, err)
		return
	}
	err = ioutil.WriteFile("./data/osuData/profileCache.json", jsonCache, 0644)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "Error with wiping data!")
		tools.ErrRead(s, err)
		return
	}

	s.ChannelMessageSend(m.ChannelID, "Updated farmerdog ratings!")
}

// CleanEmpty removes any users with no discord or osu! account
func CleanEmpty(s *discordgo.Session, m *discordgo.MessageCreate, cache []structs.PlayerData) {
	if m.Author.ID != config.Conf.BotHoster.UserID {
		s.ChannelMessageSend(m.ChannelID, "YOU ARE NOT "+config.Conf.BotHoster.Username+".........")
		return
	}

	for i := 0; i < len(cache); i++ {
		if cache[i].Discord.ID == "" && cache[i].Osu.Username == "" {
			cache = append(cache[:i], cache[i+1:]...)
			i--
		}
	}

	jsonCache, err := json.Marshal(cache)
	tools.ErrRead(s, err)

	err = ioutil.WriteFile("./data/osuData/profileCache.json", jsonCache, 0644)
	tools.ErrRead(s, err)
	s.ChannelMessageSend(m.ChannelID, "Removed empty users!")
}
