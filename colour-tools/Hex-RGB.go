package colourtools

import (
	"errors"
	"fmt"
	"image/color"
	"strings"
)

// HexToRGB returns a color.RGBA from a hex colour string
func HexToRGB(s string) (c color.NRGBA, err error) {
	c.A = 0xff

	hexToByte := func(b byte) byte {
		switch {
		case b >= '0' && b <= '9':
			return b - '0'
		case b >= 'a' && b <= 'f':
			return b - 'a' + 10
		case b >= 'A' && b <= 'F':
			return b - 'A' + 10
		}
		err = errors.New("invalid format")
		return 0
	}

	switch len(s) {
	case 6, 8:
		c.R = hexToByte(s[0])<<4 + hexToByte(s[1])
		c.G = hexToByte(s[2])<<4 + hexToByte(s[3])
		c.B = hexToByte(s[4])<<4 + hexToByte(s[5])
		if len(s) == 8 {
			c.A = hexToByte(s[6])<<4 + hexToByte(s[7])
		}
	case 3:
		c.R = hexToByte(s[0]) * 17
		c.G = hexToByte(s[1]) * 17
		c.B = hexToByte(s[2]) * 17
	default:
		err = errors.New("invalid format")
	}
	return c, err
}

// RGBToHex returns a hex string from values given
func RGBToHex(vals []uint8) (s string, err error) {
	for i, val := range vals {
		if i == 3 && val == 255 {
			break
		}
		s += fmt.Sprintf("%02x", val)
	}
	s = "#" + strings.ToUpper(s)
	return s, nil
}
