package gencommands

import (
	"encoding/json"
	"io/ioutil"

	structs "../../structs"
	tools "../../tools"
	"github.com/bwmarrin/discordgo"
)

// CleanFarm cleans the all farmerdog ratings
func CleanFarm(s *discordgo.Session, m *discordgo.MessageCreate, cache []structs.PlayerData) {
	if m.Author.ID != "92502458588205056" {
		s.ChannelMessageSend(m.ChannelID, "YOU ARE NOT VINXIS.........")
		return
	}

	for i := range cache {
		cache[i].Farm = structs.FarmerdogData{}
	}

	jsonCache, err := json.Marshal(cache)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "Error with wiping data!")
		tools.ErrRead(err)
		return
	}

	err = ioutil.WriteFile("./data/osuData/profileCache.json", jsonCache, 0644)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "Error with wiping data!")
		tools.ErrRead(err)
		return
	}
	s.ChannelMessageSend(m.ChannelID, "Cleaned farmerdog ratings!")
	return
}
