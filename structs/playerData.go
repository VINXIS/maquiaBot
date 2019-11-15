package structs

import (
	"time"

	osuapi "../osu-api"
	"github.com/bwmarrin/discordgo"
)

// PlayerData stores information regarding the discord user, and the osu user
type PlayerData struct {
	Time    time.Time
	Discord discordgo.User
	Osu     osuapi.User
	Farm    FarmerdogData
}

// FarmerdogData is how much of a farmerdog the player is
type FarmerdogData struct {
	Rating float64
	List   []PlayerScore
}

// PlayerScore is the score by the player, it tells you how farmy the score is as well
type PlayerScore struct {
	BeatmapSet int
	PP         float64
	FarmScore  float64
	Name       string
}
