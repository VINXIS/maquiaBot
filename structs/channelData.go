package structs

import (
	"github.com/bwmarrin/discordgo"
	"github.com/thehowl/go-osuapi"
)

// ChannelData stores information regarding the discord channel so that specific tracking may occur assigned for that channel
type ChannelData struct {
	Channel discordgo.Channel
	Type    string
	Users   []osuapi.User
}
