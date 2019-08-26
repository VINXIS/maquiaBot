package structs

import (
	"errors"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/thehowl/go-osuapi"
)

// ChannelData stores information regarding the discord channel so that tracking for osu! plays may occur in that channel
type ChannelData struct {
	Time     time.Time
	Channel  discordgo.Channel
	PPLimit  int
	TopPlay  int
	Users    []osuapi.User
	Tracking bool
}

func (c ChannelData) trackToggle() {
	c.Tracking = !c.Tracking
}

func (c ChannelData) addUser(u osuapi.User) {
	c.Users = append(c.Users, u)
}

func (c ChannelData) removeUser(users []string) error {
	if len(users) == 0 {
		return errors.New("No users given")
	}
	for _, user := range users {
		for i, osuUser := range c.Users {
			if strings.ToLower(osuUser.Username) == strings.ToLower(user) {
				c.Users = append(c.Users[:i], c.Users[i+1:]...)
			}
		}
	}
	return nil
}

func (c ChannelData) clearList() {
	c.Users = []osuapi.User{}
}
