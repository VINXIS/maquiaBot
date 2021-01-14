package tools

import "math"

var base = []string{
	"zero", "one", "two", "three", "four", "five", "six", "seven", "eight", "nine", "ten", "eleven", "twelve", "thirteen", "fourteen", "fifteen", "sixteen", "seventeen", "eighteen", "nineteen",
}

var tens = []string{
	"", "", "twenty", "thirty", "forty", "fifty", "sixty", "seventy", "eighty", "ninety",
}

var scales = []string{
	"", "thousand", "million", "billion", "trillion", "quadrillion", "quintillion",
}

// Ntow changes numbers to words
func Ntow(n int64) (w string) {
	// Zero
	if n == 0 {
		return base[n]
	}

	// Negative
	if n < 0 {
		w = "negative "
	}
	pos := math.Abs(float64(n))

	// groups of 3
	var groups [7]int64
	for i := 0; i < 7; i++ {
		groups[i] = int64(math.Mod(pos, 1000))
		pos /= 1000
	}

	// text version of each group
	var groupsText [7]string
	for i := 0; i < 7; i++ {
		group := groups[i]

		hundred := group / 100
		tenUnits := int64(math.Mod(float64(group), 100))

		if hundred != 0 {
			groupsText[i] += base[hundred] + " hundred"

			if tenUnits != 0 {
				groupsText[i] += " "
			}
		}

		ten := tenUnits / 10
		units := int64(math.Mod(float64(tenUnits), 10))

		if ten >= 2 {
			groupsText[i] += tens[ten]

			if units != 0 {
				groupsText[i] += "-" + base[units]
			}
		} else if tenUnits != 0 {
			groupsText[i] += base[tenUnits]
		}
	}

	// Combining the groups
	w += groupsText[0]
	for i := 1; i < 7; i++ {
		if groups[i] != 0 {
			p := groupsText[i] + " " + scales[i]

			if len(w) != 0 {
				p += " "
			}

			w = p + w
		}
	}

	return w
}

// Wton changes words to numbers
func Wton(w string) (n int64, err error) {
	return -1, nil
}
