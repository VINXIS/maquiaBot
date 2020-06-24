package gencommands

import (
	"bytes"
	"image/color"
	"image/png"
	"math/rand"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/fogleman/gg"
	colourtools "maquiaBot/colour-tools"
)

// Colour generates a 256x256 image of that colour
func Colour(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Initiate image generator
	ctx := gg.NewContext(512, 512)
	ctx.DrawRectangle(0, 0, 512, 512)
	var col color.Color
	var text string

	regex, err := regexp.Compile(`(?i)col(ou?r)?\s(.+)`)
	params := ""
	if !regex.MatchString(m.Content) { // No values given, generate random colour
		authorid, _ := strconv.Atoi(m.Author.ID)
		random := rand.New(rand.NewSource(int64(authorid) + time.Now().UnixNano()))
		col = color.NRGBA{
			uint8(random.Intn(256)),
			uint8(random.Intn(256)),
			uint8(random.Intn(256)),
			255,
		}
	} else {
		params = regex.FindStringSubmatch(m.Content)[2]
	}

	// Set colour
	if strings.Contains(params, "-hex") { // HEX
		params = strings.TrimSpace(strings.Replace(params, "-hex", "", -1))
		params = strings.TrimSpace(strings.Replace(params, "#", "", -1))

		col, err = colourtools.HexToRGB(params)
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, "Invalid hex! Make sure the hex value for the colour is either 3, 6, or 8 characters long, and has no illegal characters!")
			return
		}
	} else if strings.Contains(params, "-hsla") || strings.Contains(params, "-hsl") { // HSL(A)
		params = strings.TrimSpace(strings.Replace(params, "-hsla", "", -1))
		params = strings.TrimSpace(strings.Replace(params, "-hsl", "", -1))

		vals := strings.Split(params, " ")
		if len(vals) < 3 || len(vals) > 4 {
			s.ChannelMessageSend(m.ChannelID, "You may only send 3 to 4 values for hsl(a)! Separate the values by spaces.")
			return
		}

		hslavals := []float64{}
		for _, val := range vals {
			valNum, err := strconv.ParseFloat(val, 64)
			if err != nil || valNum < 0 || valNum > 255 {
				s.ChannelMessageSend(m.ChannelID, "Invalid HSL(A) value! Value must be between 0 and 255")
				return
			}
			hslavals = append(hslavals, valNum)
		}

		col, err = colourtools.HSLtoRGB(hslavals)
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, "Invalid HSL(A) values! `"+err.Error()+"`")
			return
		}
	} else if strings.Contains(params, "-hsva") || strings.Contains(params, "-hsv") { // HSV(A)
		params = strings.TrimSpace(strings.Replace(params, "-hsva", "", -1))
		params = strings.TrimSpace(strings.Replace(params, "-hsv", "", -1))

		vals := strings.Split(params, " ")
		if len(vals) < 3 || len(vals) > 4 {
			s.ChannelMessageSend(m.ChannelID, "You may only send 3 to 4 values for hsv(a)! Separate the values by spaces.")
			return
		}

		hsvavals := []float64{}
		for _, val := range vals {
			valNum, err := strconv.ParseFloat(val, 64)
			if err != nil || valNum < 0 || valNum > 255 {
				s.ChannelMessageSend(m.ChannelID, "Invalid HSV(A) value! Value must be between 0 and 255")
				return
			}
			hsvavals = append(hsvavals, valNum)
		}

		col, err = colourtools.HSVtoRGB(hsvavals)
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, "Invalid HSV(A) values! `"+err.Error()+"`")
			return
		}
	} else if strings.Contains(params, "-cmyk") { // CMYK
		params = strings.TrimSpace(strings.Replace(params, "-cmyk", "", -1))

		vals := strings.Split(params, " ")
		if len(vals) != 4 {
			s.ChannelMessageSend(m.ChannelID, "You may only send 3 values for CMYK! Separate the values by spaces.")
			return
		}

		cmykvals := []uint8{}
		for _, val := range vals {
			valNum, err := strconv.ParseUint(val, 10, 8)
			if err != nil || valNum > 100 {
				s.ChannelMessageSend(m.ChannelID, "Invalid CMYK value! Value must be between 0 and 100")
				return
			}
			cmykvals = append(cmykvals, uint8(float64(valNum)*2.5))
		}

		col = color.CMYK{
			cmykvals[0],
			cmykvals[1],
			cmykvals[2],
			cmykvals[3],
		}
	} else if strings.Contains(params, "-ycbcr") { // YCBCR
		params = strings.TrimSpace(strings.Replace(params, "-ycbcr", "", -1))

		vals := strings.Split(params, " ")
		if len(vals) != 3 {
			s.ChannelMessageSend(m.ChannelID, "You may only send 3 values for YCbCr! Separate the values by spaces.")
			return
		}

		ycbcrvals := []uint8{}
		for _, val := range vals {
			valNum, err := strconv.ParseUint(val, 10, 8)
			if err != nil {
				s.ChannelMessageSend(m.ChannelID, "Invalid YCbCr value! Value must be between 0 and 255")
				return
			}
			ycbcrvals = append(ycbcrvals, uint8(valNum))
		}

		col = color.YCbCr{
			ycbcrvals[0],
			ycbcrvals[1],
			ycbcrvals[2],
		}
	} else if params != "" { // RGB(A)
		// In case they tried to use this tag not knowing the default is rgb anyway
		params = strings.TrimSpace(strings.Replace(params, "-rgba", "", -1))
		params = strings.TrimSpace(strings.Replace(params, "-rgb", "", -1))

		vals := strings.Split(params, " ")
		if len(vals) < 3 || len(vals) > 4 {
			s.ChannelMessageSend(m.ChannelID, "You may only send 3 to 4 values for rgb(a)! Separate the values by spaces.")
			return
		}

		rgbavals := []uint8{}
		for _, val := range vals {
			valNum, err := strconv.ParseUint(val, 10, 8)
			if err != nil {
				s.ChannelMessageSend(m.ChannelID, "Invalid RGB value! Value must be between 0 and 255")
				return
			}
			rgbavals = append(rgbavals, uint8(valNum))
		}

		if len(rgbavals) == 4 { // RGBA
			col = color.NRGBA{
				rgbavals[0],
				rgbavals[1],
				rgbavals[2],
				rgbavals[3],
			}
		} else { // RGB
			col = color.NRGBA{
				rgbavals[0],
				rgbavals[1],
				rgbavals[2],
				255,
			}
		}
	}

	// Generate image
	ctx.SetColor(col)
	ctx.Fill()
	img := ctx.Image()

	// Generate text
	r, g, b, a := col.RGBA()
	r >>= 8
	g >>= 8
	b >>= 8
	a >>= 8
	vals := []uint8{uint8(r), uint8(g), uint8(b)}
	if a != 255 {
		vals = append(vals, uint8(a))
	}
	hex, err := colourtools.RGBToHex(vals)
	hsl, err := colourtools.RGBToHSL(vals)
	hsv, err := colourtools.RGBToHSV(vals)
	cmyk, err := colourtools.RGBToCMYK(vals)
	ycbcr, err := colourtools.RGBToYCBCR(vals)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "Error converting values! "+err.Error())
		return
	}
	if a != 255 {
		text += "RGBA: " + strconv.Itoa(int(r)) + " " + strconv.Itoa(int(g)) + " " + strconv.Itoa(int(b)) + " " + strconv.Itoa(int(a)) + "\n" +
			"Hex: " + hex + "\n" +
			"HSLA: " + strconv.Itoa(hsl[0]) + " " + strconv.Itoa(hsl[1]) + " " + strconv.Itoa(hsl[2]) + " " + strconv.Itoa(hsl[3]) + "\n" +
			"HSVA: " + strconv.Itoa(hsv[0]) + " " + strconv.Itoa(hsv[1]) + " " + strconv.Itoa(hsv[2]) + " " + strconv.Itoa(hsv[3]) + "\n" +
			"CMYK: " + strconv.Itoa(int(float64(cmyk.C)/255*100)) + " " + strconv.Itoa(int(float64(cmyk.M)/255*100)) + " " + strconv.Itoa(int(float64(cmyk.Y)/255*100)) + " " + strconv.Itoa(int(float64(cmyk.K)/255*100)) + "\n" +
			"YCbCrA: " + strconv.Itoa(int(ycbcr.Y)) + " " + strconv.Itoa(int(ycbcr.Cb)) + " " + strconv.Itoa(int(ycbcr.Cr)) + " " + strconv.Itoa(int(ycbcr.A))
	} else {
		text += "RGB: " + strconv.Itoa(int(r)) + " " + strconv.Itoa(int(g)) + " " + strconv.Itoa(int(b)) + "\n" +
			"Hex: " + hex + "\n" +
			"HSL: " + strconv.Itoa(hsl[0]) + " " + strconv.Itoa(hsl[1]) + " " + strconv.Itoa(hsl[2]) + "\n" +
			"HSV: " + strconv.Itoa(hsv[0]) + " " + strconv.Itoa(hsv[1]) + " " + strconv.Itoa(hsv[2]) + "\n" +
			"CMYK: " + strconv.Itoa(int(float64(cmyk.C)/255*100)) + " " + strconv.Itoa(int(float64(cmyk.M)/255*100)) + " " + strconv.Itoa(int(float64(cmyk.Y)/255*100)) + " " + strconv.Itoa(int(float64(cmyk.K)/255*100)) + "\n" +
			"YCbCr: " + strconv.Itoa(int(ycbcr.Y)) + " " + strconv.Itoa(int(ycbcr.Cb)) + " " + strconv.Itoa(int(ycbcr.Cr))
	}

	// Send image
	imgBytes := new(bytes.Buffer)
	_ = png.Encode(imgBytes, img)
	s.ChannelMessageSendComplex(m.ChannelID, &discordgo.MessageSend{
		Content: text,
		Files: []*discordgo.File{
			{
				Name:   "image.png",
				Reader: imgBytes,
			},
		},
	})
}
