package gencommands

import (
	"bytes"
	"image"
	"image/draw"
	"image/png"
	"math"
	"net/http"
	"regexp"
	"strings"

	"github.com/bwmarrin/discordgo"
)

// Merge merges 2+ images into 1
func Merge(s *discordgo.Session, m *discordgo.MessageCreate) {
	linkRegex, _ := regexp.Compile(`(?i)https?:\/\/\S*`)

	links := strings.Split(m.Content, " ")

	// Fetch images
	var images []image.Image
	length := 0
	height := 0
	for _, link := range links {
		if !linkRegex.MatchString(link) {
			continue
		}

		response, err := http.Get(link)
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, "Could not reach the URL "+link)
			return
		}
		imgSrc, _, err := image.Decode(response.Body)
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, "Could not parse the image from "+link)
			return
		}
		length += imgSrc.Bounds().Dx()
		height = int(math.Max(float64(height), float64(imgSrc.Bounds().Dy())))

		images = append(images, imgSrc)
		response.Body.Close()
	}

	if len(images) < 2 {
		s.ChannelMessageSend(m.ChannelID, "Please provide at least 2 image links!")
		return
	}

	size := image.Rectangle{image.Point{0, 0}, image.Point{length, height}}
	finalImage := image.NewRGBA(size)

	area := images[0].Bounds()
	for i, img := range images {

		draw.Draw(finalImage, area, img, image.Point{0, 0}, draw.Src)
		sp := image.Point{img.Bounds().Dx(), 0}
		if i+1 < len(images) {
			area = image.Rectangle{sp, sp.Add(images[i+1].Bounds().Size())}
		}
	}

	imgBytes := new(bytes.Buffer)
	_ = png.Encode(imgBytes, finalImage)
	s.ChannelMessageSendComplex(m.ChannelID, &discordgo.MessageSend{
		Files: []*discordgo.File{
			{
				Name:   "image.png",
				Reader: imgBytes,
			},
		},
	})
}
