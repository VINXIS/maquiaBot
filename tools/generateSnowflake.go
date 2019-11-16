package tools

import "time"

// GenerateSnowflake generates a snowflake identical to Discord's
func GenerateSnowflake(t time.Time) int64 {
	ms := t.Unix()*1000 - 1420070400000
	ms <<= 22
	return ms
}
