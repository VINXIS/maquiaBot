package tools

import "math"

// BarCreation creates a bar that is filled by percentage
func BarCreation(percent float64) (bar string) {
	barCount := int(math.Round(percent * 50.0))
	bar = "["
	for i := 0; i < 50; i++ {
		if i < barCount {
			bar += "|"
		} else {
			bar += " "
		}
	}
	bar += "]"
	return bar
}
