package tools

import (
	"encoding/json"
	"io/ioutil"
	"os"

	"maquiaBot/structs"

	"github.com/bwmarrin/discordgo"
)

// GetChannel obtains a channel using its channel ID
func GetChannel(channel discordgo.Channel, s *discordgo.Session) (structs.ChannelData, bool) {
	channelData := structs.NewChannel()
	new := true
	_, err := os.Stat("./data/channelData/" + channel.ID + ".json")
	if err == nil {
		f, err := ioutil.ReadFile("./data/channelData/" + channel.ID + ".json")
		ErrRead(s, err)
		_ = json.Unmarshal(f, &channelData)
		new = false
	}
	return channelData, new
}
