package structs

import (
	"strings"
	"time"

	osuapi "maquiaBot/osu-api"
)

// ChannelData stores information regarding the discord channel so that tracking for osu! plays may occur in that channel
type ChannelData struct {
	Time            time.Time
	PPReq           float64
	LeaderboardReq  int
	TopReq          int
	Ranked          bool
	Loved           bool
	Qualified       bool
	Users           []osuapi.User
	Mode            osuapi.Mode
	Pixiv           PixivData
	Daily           bool
	OsuToggle       bool
	TimestampToggle bool
	Vibe            bool
	OsuTracking     bool
	PixivTracking   bool
}

// PixivData stores pixiv tracking
type PixivData struct {
	Rankings []string
	Users    []string
}

// NewChannel creates a new ChannelData
func NewChannel() ChannelData {
	return ChannelData{
		Time:           time.Now(),
		PPReq:          -1,
		LeaderboardReq: 101,
		TopReq:         101,
		Mode:           osuapi.ModeOsu,
		OsuTracking:    true,
	}
}

// ClearList clears the user list and turns off tracking
func (c *ChannelData) ClearList() {
	c.Users = []osuapi.User{}
	c.OsuTracking = false
}

// TrackToggle toggles the tracking on and off
func (c *ChannelData) TrackToggle() {
	if len(c.Users) == 0 {
		c.OsuTracking = false
		return
	}
	c.OsuTracking = !c.OsuTracking
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
func (c *ChannelData) UpdateMapStatus(mapStatus []string) {
	for _, status := range mapStatus {
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
