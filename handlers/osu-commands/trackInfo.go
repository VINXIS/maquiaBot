package osucommands

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"strconv"

	structs "../../structs"
	tools "../../tools"
	"github.com/bwmarrin/discordgo"
)

// TrackInfo gives info about what's being tracked in the channel currently
func TrackInfo(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Obtain channel data
	channelData := structs.ChannelData{}
	_, err := os.Stat("./data/channelData/" + m.ChannelID + ".json")
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "Nothing being tracked in this channel currently!")
		return
	}
	f, err := ioutil.ReadFile("./data/channelData/" + m.ChannelID + ".json")
	tools.ErrRead(err)
	_ = json.Unmarshal(f, &channelData)

	// Create message
	trackText := "Currently tracking "
	for _, user := range channelData.Users {
		trackText = trackText + user.Username + ", "
	}

	if channelData.PPLimit != 0 {
		trackText = trackText + "with a pp limit of " + strconv.Itoa(channelData.PPLimit)
	} else {
		trackText = trackText + "with no pp limit"
	}

	if channelData.TopPlay != 0 {
		trackText = trackText + ", and the top " + strconv.Itoa(channelData.TopPlay) + " scores"
	}
	s.ChannelMessageSend(m.ChannelID, trackText)
	return
}
