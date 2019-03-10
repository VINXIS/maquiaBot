package structs

import (
	"time"

	"github.com/bwmarrin/discordgo"
)

// ServerData stores information regarding the discord server, so thta server specific customizations may be used.
type ServerData struct {
	Time   time.Time
	Server discordgo.Guild
	Prefix string
}
