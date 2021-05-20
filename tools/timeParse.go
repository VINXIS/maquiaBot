package tools

import (
	"strings"
	"time"
)

// TimeParse tries to parses a date / time
func TimeParse(datetime string) (timestamp time.Time, err error) {
	datetime = strings.Replace(datetime, "am", "AM", -1)
	datetime = strings.Replace(datetime, "pm", "PM", -1)
	timeFormats := []string{
		time.ANSIC,
		time.UnixDate,
		time.RubyDate,
		time.RFC822,
		time.RFC822Z,
		time.RFC850,
		time.RFC1123,
		time.RFC1123Z,
		time.RFC3339,
		time.RFC3339Nano,
		time.Kitchen,
		time.Stamp,
		time.StampMilli,
		time.StampMicro,
		time.StampNano,
	}

	customDate := []string{
		"January 2 2006",
		"Jan 2 2006",
		"1 2 2006",
		"01 2 2006",
		"1 02 2006",
		"01 02 2006",

		"2 January 2006",
		"2 Jan 2006",
		"2 1 2006",
		"2 01 2006",
		"02 1 2006",
		"02 01 2006",

		"2006 January 2",
		"2006 Jan 2",
		"2006 1 2",
		"2006 01 2",
		"2006 1 02",
		"2006 01 02",

		"January 2",
		"Jan 2",
		"1 2",
		"01 2",
		"1 02",
		"01 02",

		"2 January",
		"2 Jan",
		"2 1",
		"2 01",
		"02 1",
		"02 01",
	}

	customTime := []string{
		"15:04:05",
		"15:04",

		"3:04 PM",
		"03:04 PM",
		"3 PM",
		"03 PM",

		"3:04PM",
		"03:04PM",
		"3PM",
		"03PM",
	}

	customZone := []string{
		"MST",

		"GMT-0700",
		"GMT-7",
		"GMT-07",
		"GMT-07:00",
		"GMT-7:00",

		"UTC-0700",
		"UTC-7",
		"UTC-07",
		"UTC-07:00",
		"UTC-7:00",
	}

	for _, timeFormat := range timeFormats {
		timestamp, err = time.Parse(timeFormat, datetime)
		if err == nil {
			return timestamp, nil
		}
	}

	// Run custom formats only if none of the default formats work
	for _, date := range customDate {
		timestamp, err = time.Parse(date, datetime)
		if err == nil {
			return timestamp, nil
		}

		for _, timer := range customTime {
			timestamp, err = time.Parse(timer, datetime)
			if err == nil {
				timestamp = timestamp.AddDate(time.Now().Year(), int(time.Now().Month())-1, time.Now().Day())
				return timestamp, nil
			}

			timestamp, err = time.Parse(date+" "+timer, datetime)
			if err == nil {
				return timestamp, nil
			}

			timestamp, err = time.Parse(timer+" "+date, datetime)
			if err == nil {
				return timestamp, nil
			}

			for _, zone := range customZone {
				timestamp, err = time.Parse(date+" "+zone, datetime)
				if err == nil {
					return timestamp, nil
				}

				timestamp, err = time.Parse(timer+" "+zone, datetime)
				if err == nil {
					timestamp = timestamp.AddDate(time.Now().Year(), int(time.Now().Month())-1, time.Now().Day())
					return timestamp, nil
				}

				timestamp, err = time.Parse(zone+" "+date, datetime)
				if err == nil {
					return timestamp, nil
				}

				timestamp, err = time.Parse(zone+" "+timer, datetime)
				if err == nil {
					timestamp = timestamp.AddDate(time.Now().Year(), int(time.Now().Month())-1, time.Now().Day())
					return timestamp, nil
				}

				timestamp, err = time.Parse(date+" "+timer+" "+zone, datetime)
				if err == nil {
					return timestamp, nil
				}

				timestamp, err = time.Parse(date+" "+zone+" "+timer, datetime)
				if err == nil {
					return timestamp, nil
				}

				timestamp, err = time.Parse(timer+" "+date+" "+zone, datetime)
				if err == nil {
					return timestamp, nil
				}

				timestamp, err = time.Parse(timer+" "+zone+" "+date, datetime)
				if err == nil {
					return timestamp, nil
				}

				timestamp, err = time.Parse(zone+" "+timer+" "+date, datetime)
				if err == nil {
					return timestamp, nil
				}

				timestamp, err = time.Parse(zone+" "+date+" "+timer, datetime)
				if err == nil {
					return timestamp, nil
				}
			}
		}
	}

	// Run with dashed date now if none of the non-dashed work
	for _, date := range customDate {
		date = dashed(date)

		timestamp, err = time.Parse(date, datetime)
		if err == nil {
			return timestamp, nil
		}

		for _, timer := range customTime {
			timestamp, err = time.Parse(timer, datetime)
			if err == nil {
				return timestamp, nil
			}

			timestamp, err = time.Parse(date+" "+timer, datetime)
			if err == nil {
				return timestamp, nil
			}

			timestamp, err = time.Parse(timer+" "+date, datetime)
			if err == nil {
				return timestamp, nil
			}

			for _, zone := range customZone {
				timestamp, err = time.Parse(date+" "+zone, datetime)
				if err == nil {
					return timestamp, nil
				}

				timestamp, err = time.Parse(timer+" "+zone, datetime)
				if err == nil {
					return timestamp, nil
				}

				timestamp, err = time.Parse(zone+" "+date, datetime)
				if err == nil {
					return timestamp, nil
				}

				timestamp, err = time.Parse(zone+" "+timer, datetime)
				if err == nil {
					return timestamp, nil
				}

				timestamp, err = time.Parse(date+" "+timer+" "+zone, datetime)
				if err == nil {
					return timestamp, nil
				}

				timestamp, err = time.Parse(date+" "+zone+" "+timer, datetime)
				if err == nil {
					return timestamp, nil
				}

				timestamp, err = time.Parse(timer+" "+date+" "+zone, datetime)
				if err == nil {
					return timestamp, nil
				}

				timestamp, err = time.Parse(timer+" "+zone+" "+date, datetime)
				if err == nil {
					return timestamp, nil
				}

				timestamp, err = time.Parse(zone+" "+timer+" "+date, datetime)
				if err == nil {
					return timestamp, nil
				}

				timestamp, err = time.Parse(zone+" "+date+" "+timer, datetime)
				if err == nil {
					return timestamp, nil
				}
			}
		}
	}

	return timestamp, err
}

func dashed(date string) string {
	return strings.Replace(date, " ", "-", -1)
}
