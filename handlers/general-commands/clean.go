package gencommands

import (
	"encoding/json"
	"io/ioutil"

	structs "../../structs"
	tools "../../tools"
	"github.com/bwmarrin/discordgo"
)

// Clean cleans the caches
func Clean(s *discordgo.Session, m *discordgo.MessageCreate, cache []structs.PlayerData) {
	if m.Author.ID != "92502458588205056" {
		s.ChannelMessageSend(m.ChannelID, "YOU ARE NOT VINXIS.........")
		return
	}

	newPlayerCache := []structs.PlayerData{}
	for _, player := range cache {
		if player.Discord.ID != "" {
			newPlayerCache = append(newPlayerCache, player)
		}
	}

	jsonCache, err := json.Marshal(newPlayerCache)
	tools.ErrRead(err)

	err = ioutil.WriteFile("./data/osuData/profileCache.json", jsonCache, 0644)
	tools.ErrRead(err)
	s.ChannelMessageSend(m.ChannelID, "Cleaned player cache!")
}
