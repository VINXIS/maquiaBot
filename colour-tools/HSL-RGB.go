package colourtools

import (
	"errors"
	"fmt"
	"image/color"
	"math"
)

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
		if vals[3] < 0 || vals[3] > 255 {
			return c, errors.New("invalid alpha value, must be between 0 and 255")
		}
		c.A = uint8(vals[3])
	}

	fmt.Println(c)

	return c, err
}

// RGBToHSL converts RGB(A) values to HSL(A)
func RGBToHSL(vals []uint8) (hsl []int, err error) {
	r1 := float64(vals[0]) / 255.0
	g1 := float64(vals[1]) / 255.0
	b1 := float64(vals[2]) / 255.0

	// val checks
	if r1 < 0 || r1 > 1 {
		return hsl, errors.New("invalid red value, must be between 0 and 255")
	}
	if g1 < 0 || g1 > 1 {
		return hsl, errors.New("invalid green value, must be between 0 and 255")
	}
	if b1 < 0 || b1 > 1 {
		return hsl, errors.New("invalid blue value, must be between 0 and 255")
	}

	// Obtain values
	cMax := math.Max(r1, math.Max(g1, b1))
	cMin := math.Min(r1, math.Min(g1, b1))
	delta := cMax - cMin

	var h, s, l int

	// Hue calc
	if delta != 0 {
		switch cMax {
		case r1:
			h = int(60 * math.Mod((g1-b1)/delta, 6.0))
		case g1:
			h = int(60 * ((b1-r1)/delta + 2.0))
		case b1:
			h = int(60 * ((r1-g1)/delta + 4.0))
		}
	}

	// Lightness calc
	lFloat := (cMax + cMin) / 2.0
	l = int(100 * lFloat)

	// Saturation calc
	if delta != 0 {
		s = int(100 * delta / (1 - math.Abs(2*lFloat-1)))
	}

	hsl = []int{h, s, l}

	// Alpha check
	if len(vals) == 4 {
		if vals[3] < 0 || vals[3] > 255 {
			return hsl, errors.New("invalid alpha value, must be between 0 and 255")
		}
		hsl = append(hsl, int(vals[3]))
	}
	return hsl, err
}
