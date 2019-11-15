package structs

import (
	"time"

	"github.com/bwmarrin/discordgo"
)

// ServerData stores information regarding the discord server, so that server specific customizations may be used.
type ServerData struct {
	Time       time.Time
	Server     discordgo.Guild
	Prefix     string
	Crab       bool
	OsuToggle  bool
	Possession bool
	Vibe       bool
	Adjectives []string
	Nouns      []string
	Skills     []string
}

// NewServer creates a new ServerData
func NewServer() ServerData {
	return ServerData{
		Prefix:    "$",
		OsuToggle: true,
		Crab:      true,
	}
}
