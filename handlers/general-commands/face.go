package gencommands

import (
	"bytes"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"io/ioutil"
	"math"
	"net/http"
	"regexp"
	"sort"

	"github.com/bwmarrin/discordgo"
	pigo "github.com/esimov/pigo/core"
	"github.com/fogleman/gg"
)

// Face lets you detect faces
func Face(s *discordgo.Session, m *discordgo.MessageCreate) {
	linkRegex, _ := regexp.Compile(`https?:\/\/\S*`)
	// saturationRegex, _ := regexp.Compile(`-s\s+(-?\d+)`)
	// contrastRegex, _ := regexp.Compile(`-c\s+(-?\d+)`)
	// langRegex, _ := regexp.Compile(`-l\s+(\S+)`)

	var url string
	if len(m.Attachments) > 0 {
		url = m.Attachments[0].URL
	} else if len(m.Embeds) > 0 && m.Embeds[0].Image != nil {
		url = m.Embeds[0].Image.URL
	} else if linkRegex.MatchString(m.Content) {
		url = linkRegex.FindStringSubmatch(m.Content)[0]
	} else {
		// Get prev messages
		messages, err := s.ChannelMessages(m.ChannelID, 100, "", "", "")
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, "Error fetching messages.")
			return
		}

		// Sort by date
		sort.Slice(messages, func(i, j int) bool {
			time1, _ := messages[i].Timestamp.Parse()
			time2, _ := messages[j].Timestamp.Parse()
			return time1.After(time2)
		})
		for _, msg := range messages {
			if len(msg.Attachments) > 0 {
				url = msg.Attachments[0].URL
				break
			} else if len(msg.Embeds) > 0 && msg.Embeds[0].Image != nil {
				url = msg.Embeds[0].Image.URL
				break
			} else if linkRegex.MatchString(msg.Content) {
				url = linkRegex.FindStringSubmatch(msg.Content)[0]
				break
			}
		}
		if url == "" {
			s.ChannelMessageSend(m.ChannelID, "No link/image given.")
			return
		}
	}

	// Fetch the image data
	response, err := http.Get(url)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "Could not reach URL.")
		return
	}
	imgSrc, _, err := image.Decode(response.Body)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "Could not parse any image from data given.")
		return
	}

	// Fetch the facefinder
	cascade, err := ioutil.ReadFile("./facefinder")
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "Could not access the classifier's data.")
		return
	}

	// Get params
	colourPxls := pigo.ImgToNRGBA(imgSrc)
	draw.Draw(colourPxls, colourPxls.Bounds(), imgSrc, imgSrc.Bounds().Min, draw.Src)
	pxls := pigo.RgbToGrayscale(colourPxls)
	cols, rows := colourPxls.Bounds().Max.X, colourPxls.Bounds().Max.Y
	ctx := gg.NewContext(cols, rows)
	ctx.DrawImage(colourPxls, 0, 0)
	cParams := pigo.CascadeParams{
		MinSize:     20,
		MaxSize:     2000,
		ShiftFactor: 0.1,
		ScaleFactor: 1.1,
		ImageParams: pigo.ImageParams{
			Pixels: pxls,
			Rows:   rows,
			Cols:   cols,
			Dim:    cols,
		},
	}

	// Get faces
	newPigo := pigo.NewPigo()
	classifier, err := newPigo.Unpack(cascade)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "Could not access the classifier's data.")
		return
	}

	faces := classifier.RunCascade(cParams, 0)
	faces = classifier.ClusterDetections(faces, 0)

	// Write faces to image
	var (
		rects []image.Rectangle
		found bool
	)
	for _, face := range faces {
		if face.Q > 5 {
			found = true
			ctx.DrawArc(
				float64(face.Col),
				float64(face.Row),
				float64(face.Scale/2),
				0,
				2*math.Pi,
			)
			rects = append(rects, image.Rect(
				face.Col-face.Scale/2,
				face.Row-face.Scale/2,
				face.Scale,
				face.Scale,
			))
			ctx.SetLineWidth(2.0)
			ctx.SetStrokeStyle(gg.NewSolidPattern(color.RGBA{R: 255, G: 0, B: 0, A: 255}))
			ctx.Stroke()
		}
	}

	// See if any faces were found, send image otherwise
	if !found {
		s.ChannelMessageSend(m.ChannelID, "No faces detected!")
		return
	}
	img := ctx.Image()
	imgBytes := new(bytes.Buffer)
	err = png.Encode(imgBytes, img)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "Error in encoding image!")
		return
	}
	_, err = s.ChannelMessageSendComplex(m.ChannelID, &discordgo.MessageSend{
		Files: []*discordgo.File{
			&discordgo.File{
				Name:   "image.png",
				Reader: imgBytes,
			},
		},
	})
}
