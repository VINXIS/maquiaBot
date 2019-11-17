package tools

import (
	"encoding/json"
	"io/ioutil"
	"os"

	structs "../structs"
	"github.com/bwmarrin/discordgo"
)

// GetChannel obtains a channel using its channel ID
func GetChannel(channel discordgo.Channel) (structs.ChannelData, bool) {
	channelData := structs.NewChannel(channel)
	new := true
	_, err := os.Stat("./data/channelData/" + channel.ID + ".json")
	if err == nil {
		f, err := ioutil.ReadFile("./data/channelData/" + channel.ID + ".json")
		ErrRead(err)
		_ = json.Unmarshal(f, &channelData)
		new = false
	}
	return channelData, new
}
