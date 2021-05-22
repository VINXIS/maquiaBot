package structs

import (
	"time"
)

// TaskCache stores the task alongside if they are active or not
type TaskCache struct {
	Task   Task
	Active bool
}

// Task stores information about a task a user asked for
type Task struct {
	ID      int64
	Info    string
	User    string
	Seconds float64
	LastRun time.Time
}

// NewTask creates a new Task with a snowflake ID similar to Discord's
func NewTask(id string, info string, duration time.Duration, lastRun time.Time) Task {
	ID := time.Now().Unix()*1000 - 1420070400000
	ID <<= 22
	return Task{
		ID:      ID,
		User:    id,
		Info:    info,
		Seconds: duration.Seconds(),
		LastRun: lastRun,
	}
}
