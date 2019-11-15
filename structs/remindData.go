package structs

import (
	"time"

	"github.com/bwmarrin/discordgo"
)

// ReminderTimer stores the reminder alongside the timer
type ReminderTimer struct {
	Reminder Reminder
	Timer    time.Timer
}

// Reminder stores information about a remind a user asked for
type Reminder struct {
	ID     int64
	Target time.Time
	Info   string
	User   discordgo.User
	Active bool
}

// NewReminder creates a new Reminder with a snowflake ID similar to Discord's
func NewReminder(target time.Time, user discordgo.User, info string) Reminder {
	return Reminder{
		ID:     generateSnowflake(time.Now()),
		Target: target,
		User:   user,
		Info:   info,
		Active: true,
	}
}

func generateSnowflake(t time.Time) int64 {
	ms := t.Unix()*1000 - 1420070400000
	ms <<= 22
	return ms
}
