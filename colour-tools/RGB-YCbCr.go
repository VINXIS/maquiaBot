package colourtools

import (
	"image/color"
)

// YCBCRToRGB returns a color.NRGBA from a hex colour string
func YCBCRToRGB(vals []uint8) (c color.NRGBA, err error) {
	c.R, c.G, c.B = color.YCbCrToRGB(vals[0], vals[1], vals[2])
	return c, err
}

// RGBToYCBCR returns a color.YCbCr from a list of values given
func RGBToYCBCR(vals []uint8) (c color.NYCbCrA, err error) {
	c.Y, c.Cb, c.Cr = color.RGBToYCbCr(vals[0], vals[1], vals[2])
	if len(vals) == 4 {
		c.A = vals[3]
	}
	return c, err
}
