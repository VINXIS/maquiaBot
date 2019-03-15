package structs

import (
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/thehowl/go-osuapi"
)

// ChannelData stores information regarding the discord channel so that specific tracking may occur assigned for that channel
type ChannelData struct {
	Time     time.Time
	Channel  discordgo.Channel
	PPLimit  int
	TopPlay  int
	Users    []osuapi.User
	Tracking bool
}
