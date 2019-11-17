package structs

import (
	"strings"
	"time"

	osuapi "../osu-api"
	"github.com/bwmarrin/discordgo"
)

// ChannelData stores information regarding the discord channel so that tracking for osu! plays may occur in that channel
type ChannelData struct {
	Time           time.Time
	Channel        discordgo.Channel
	PPReq          float64
	LeaderboardReq int
	TopReq         int
	Ranked         bool
	Loved          bool
	Qualified      bool
	Users          []osuapi.User
	Mode           osuapi.Mode
	Tracking       bool
}

// NewChannel creates a new ChannelData
func NewChannel(Channel discordgo.Channel) ChannelData {
	return ChannelData{
		Time:           time.Now(),
		Channel:        Channel,
		PPReq:          -1,
		LeaderboardReq: 101,
		TopReq:         101,
		Mode:           osuapi.ModeOsu,
		Tracking:       true,
	}
}

// ClearList clears the user list and turns off tracking
func (c *ChannelData) ClearList() {
	c.Users = []osuapi.User{}
	c.Tracking = false
}

// TrackToggle toggles the tracking on and off
func (c *ChannelData) TrackToggle() {
	if len(c.Users) == 0 {
		c.Tracking = false
		return
	}
	c.Tracking = !c.Tracking
}

// AddUser adds a user to the list of people to track
func (c *ChannelData) AddUser(u osuapi.User) {
	for _, user := range c.Users {
		if user.Username == u.Username {
			return
		}
	}
	c.Users = append(c.Users, u)
}

// RemoveUser removes users from the list of people to track
func (c *ChannelData) RemoveUser(users []string) {
	if len(users) == 0 {
		return
	}
	for _, user := range users {
		for i := 0; i < len(c.Users); i++ {
			if strings.ToLower(c.Users[i].Username) == strings.ToLower(user) {
				c.Users = append(c.Users[:i], c.Users[i+1:]...)
				i--
			}
		}
	}
}

// UpdateMapStatus updates the map statuses allowed
func (c *ChannelData) UpdateMapStatus(mapTypes []string) {
	for _, status := range mapTypes {
		switch status {
		case "r", "rank", "ranked":
			c.Ranked = !c.Ranked
		case "q", "qual", "qualified":
			c.Qualified = !c.Qualified
		case "l", "love", "loved":
			c.Loved = !c.Loved
		}
	}
}
