package tools

import (
	"image"
	"image/color"
	"image/draw"
	"io/ioutil"

	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
)

// AddLabel adds a label onto an image
func AddLabel(img image.Image, x, y int, label string) {
	col := color.Black
	point := fixed.Point26_6{
		X: fixed.Int26_6(x * 64),
		Y: fixed.Int26_6(y * 64),
	}
	b, _ := ioutil.ReadFile("./data/fonts/Roboto-Regular.ttf")
	writingFont, _ := truetype.Parse(b)

	dimg, _ := img.(draw.Image)
	d := &font.Drawer{
		Dst: dimg,
		Src: image.NewUniform(col),
		Face: truetype.NewFace(writingFont, &truetype.Options{
			Size: 14,
		}),
		Dot: point,
	}
	d.DrawString(label)
}
