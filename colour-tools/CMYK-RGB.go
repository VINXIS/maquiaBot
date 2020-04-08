package colourtools

import (
	"image/color"
)

// CMYKToRGB returns a color.NRGBA from a hex colour string
func CMYKToRGB(vals []uint8) (c color.NRGBA, err error) {
	c.R, c.G, c.B = color.CMYKToRGB(vals[0], vals[1], vals[2], vals[3])
	return c, err
}

// RGBToCMYK returns a color.CMYK from a list of values given
func RGBToCMYK(vals []uint8) (c color.CMYK, err error) {
	c.C, c.M, c.Y, c.K = color.RGBToCMYK(vals[0], vals[1], vals[2])
	return c, err
}
