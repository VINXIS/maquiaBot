package gencommands

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"maquiaBot/structs"
	tools "maquiaBot/tools"

	"github.com/bwmarrin/discordgo"
)

// Download lets users download any data stored for them
func Download(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Obtain profile cache data
	var cache []structs.PlayerData
	f, err := ioutil.ReadFile("./data/osuData/profileCache.json")
	tools.ErrRead(s, err)
	_ = json.Unmarshal(f, &cache)

	user := structs.PlayerData{}
	for _, cacheUser := range cache {
		if cacheUser.Discord == m.Author.ID {
			user = cacheUser
			break
		}
	}

	if user.Discord == "" {
		s.ChannelMessageSend(m.ChannelID, "There is currently no data stored for you.")
		return
	}
	dm, err := s.UserChannelCreate(m.Author.ID)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "Error in sending you the data! Please make sure your DMs are open in order to obtain your data!")
		return
	}
	b, err := json.Marshal(user)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "Error in parsing your data. Please contact `@vinxis1` on twitter or `VINXIS#1000` on discord about this!")
		return
	}
	s.ChannelMessageSendComplex(dm.ID, &discordgo.MessageSend{
		Content: "Here is the data stored for you.",
		File: &discordgo.File{
			Name:   m.Author.Username + ".json",
			Reader: bytes.NewReader(b),
		},
	})
}
