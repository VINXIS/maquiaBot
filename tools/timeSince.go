package tools

import (
	"math"
	"strconv"
	"time"
)

// TimeSince gives back a parsed string of the time elpased since
func TimeSince(timeParse time.Time) (timeString string) {
	timeSince := time.Since(timeParse)
	if timeSince.Hours() > 24 {
		days := strconv.FormatFloat(math.Floor(timeSince.Hours()/24.0), 'f', 0, 64)
		hours := strconv.FormatFloat(math.Mod(timeSince.Hours(), 24), 'f', 0, 64)

		if days == "1" {
			days = days + " day"
		} else {
			days = days + " days"
		}

		if hours == "1" {
			hours = hours + " hour"
		} else {
			hours = hours + " hours"
		}

		timeString = days + " and " + hours + " ago."
	} else if timeSince.Hours() > 1 {
		hours := strconv.FormatFloat(timeSince.Hours(), 'f', 0, 64)
		minutes := strconv.FormatFloat(math.Mod(timeSince.Minutes(), 60), 'f', 0, 64)

		if hours == "1" {
			hours = hours + " hour"
		} else {
			hours = hours + " hours"
		}

		if minutes == "1" {
			minutes = minutes + " minute"
		} else {
			minutes = minutes + " minutes"
		}

		timeString = hours + " and " + minutes + " ago."
	} else if timeSince.Minutes() > 1 {
		minutes := strconv.FormatFloat(timeSince.Minutes(), 'f', 0, 64)
		seconds := strconv.FormatFloat(math.Mod(timeSince.Seconds(), 60), 'f', 0, 64)

		if minutes == "1" {
			minutes = minutes + " minute"
		} else {
			minutes = minutes + " minutes"
		}

		if seconds == "1" {
			seconds = seconds + " second"
		} else {
			seconds = seconds + " seconds"
		}

		timeString = minutes + " and " + seconds + " ago."
	} else {
		seconds := strconv.FormatFloat(math.Abs(timeSince.Seconds()), 'f', 0, 64)

		if seconds == "1" {
			seconds = seconds + " second"
		} else {
			seconds = seconds + " seconds"
		}

		timeString = seconds + " ago."
	}
	return timeString
}
