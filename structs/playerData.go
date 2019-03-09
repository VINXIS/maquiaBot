package structs

import (
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/thehowl/go-osuapi"
)

// PlayerData stores information regarding the discord user, and the osu user
type PlayerData struct {
	Time    time.Time
	Discord discordgo.User
	Osu     osuapi.User
}
