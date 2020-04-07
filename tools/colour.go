package tools

import (
	"errors"
	"fmt"
	"image/color"
	"math"
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

// HSLtoRGB changes HSL(A) values to RGB(A) values
func HSLtoRGB(vals []float64) (c color.NRGBA, err error) {
	h := vals[0]
	s := vals[1] / 100
	l := vals[2] / 100

	// val checks
	if h < 0 || h >= 360 {
		return c, errors.New("invalid hue value, must be between 0 and 360")
	}
	if s < 0 || s > 1 {
		return c, errors.New("invalid saturation value, must be between 0 and 100")
	}
	if l < 0 || l > 1 {
		return c, errors.New("invalid lighting value, must be between 0 and 100")
	}

	// Conversion time
	ch := s * (1 - math.Abs(2*l-1))
	x := ch * (1 - math.Abs(math.Mod(h/60.0, 2)-1))
	m := l - ch/2.0

	var r1, g1, b1 float64
	if h >= 0 && h < 60 {
		r1 = ch
		g1 = x
	} else if h >= 60 && h < 120 {
		r1 = x
		g1 = ch
	} else if h >= 120 && h < 180 {
		g1 = ch
		b1 = x
	} else if h >= 180 && h < 240 {
		g1 = x
		b1 = ch
	} else if h >= 240 && h < 300 {
		r1 = x
		b1 = ch
	} else if h >= 300 && h < 360 {
		r1 = ch
		b1 = x
	}

	c.R = uint8(255 * (r1 + m))
	c.G = uint8(255 * (g1 + m))
	c.B = uint8(255 * (b1 + m))

	// Alpha check
	c.A = 0xff
	if len(vals) == 4 {
		c.A = uint8(vals[4])
	}

	fmt.Println(c)

	return c, err
}
