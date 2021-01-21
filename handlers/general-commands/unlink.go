package gencommands

import (
	"encoding/json"
	"io/ioutil"
	"maquiaBot/structs"
	tools "maquiaBot/tools"

	"github.com/bwmarrin/discordgo"
)

// Unlink unlinks your account from the player cache.
func Unlink(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Obtain profile cache data
	var cache []structs.PlayerData
	f, err := ioutil.ReadFile("./data/osuData/profileCache.json")
	tools.ErrRead(s, err)
	_ = json.Unmarshal(f, &cache)

	dataExists := false
	for i, cacheUser := range cache {
		if cacheUser.Discord == m.Author.ID {
			dataExists = true
			cache[i] = cache[len(cache)-1]
			cache = cache[:len(cache)-1]
			break
		}
	}

	if !dataExists {
		s.ChannelMessageSend(m.ChannelID, "There is no data corresponding to your discord ID!")
		return
	}

	jsonCache, err := json.Marshal(cache)
	tools.ErrRead(s, err)

	err = ioutil.WriteFile("./data/osuData/profileCache.json", jsonCache, 0644)
	tools.ErrRead(s, err)
	s.ChannelMessageSend(m.ChannelID, "Successfully removed your data.")
}
