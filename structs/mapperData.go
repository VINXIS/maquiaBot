package structs

import (
	"github.com/bwmarrin/discordgo"
	osuapi "maquiaBot/osu-api"
)

// MapperData holds information about a mapper for tracking
type MapperData struct {
	Mapper   osuapi.User
	Beatmaps []osuapi.Beatmap
	Channels []string
}

// AddChannel adds a channel to that user
func (m *MapperData) AddChannel(c discordgo.Channel) {
	for _, ch := range m.Channels {
		if c.ID == ch {
			return
		}
	}
	m.Channels = append(m.Channels, c.ID)
}

// RemoveChannel removes a channel from that user
func (m *MapperData) RemoveChannel(c discordgo.Channel) {
	for i := 0; i < len(m.Channels); i++ {
		if c.ID == m.Channels[i] {
			m.Channels = append(m.Channels[:i], m.Channels[i+1:]...)
			i--
		}
	}
}
