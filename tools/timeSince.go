package tools

import (
	"math"
	"strconv"
	"time"
)

// TimeSince gives back a parsed string of the time elpased since
func TimeSince(timeParse time.Time) (timeString string) {
	timeSince := time.Since(timeParse)
	if timeSince.Hours() > 8760 {
		years := strconv.FormatFloat(math.Floor(timeSince.Hours()/8760.0), 'f', 0, 64)
		months := strconv.FormatFloat(math.Floor(math.Mod(timeSince.Hours(), 8760)/730.0), 'f', 0, 64)

		if years == "1" {
			years = years + " year"
		} else {
			years = years + " years"
		}

		if months == "1" {
			months = months + " month"
		} else {
			months = months + " months"
		}

		if months == "0 months" {
			timeString = years + " ago."
		} else {
			timeString = years + " and " + months + " ago."
		}
	} else if timeSince.Hours() > 730 {
		months := strconv.FormatFloat(math.Floor(timeSince.Hours()/730.0), 'f', 0, 64)
		days := strconv.FormatFloat(math.Floor(math.Mod(timeSince.Hours(), 730)/24.0), 'f', 0, 64)

		if months == "1" {
			months = months + " month"
		} else {
			months = months + " months"
		}

		if days == "1" {
			days = days + " day"
		} else {
			days = days + " days"
		}

		if days == "0 days" {
			timeString = months + " ago."
		} else {
			timeString = months + " and " + days + " ago."
		}
	} else if timeSince.Hours() > 24 {
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

		if hours == "0 hours" {
			timeString = days + " ago."
		} else {
			timeString = days + " and " + hours + " ago."
		}
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

		if minutes == "0 minutes" {
			timeString = hours + " ago."
		} else {
			timeString = hours + " and " + minutes + " ago."
		}
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

		if seconds == "0 seconds" {
			timeString = minutes + " ago."
		} else {
			timeString = minutes + " and " + seconds + " ago."
		}
	} else {
		seconds := strconv.FormatFloat(math.Min(0, timeSince.Seconds()), 'f', 0, 64)

		if seconds == "1" {
			seconds = seconds + " second"
		} else {
			seconds = seconds + " seconds"
		}

		if seconds == "0 seconds" {
			timeString = "just now."
		} else {
			timeString = seconds + " ago."
		}
	}
	return timeString
}
