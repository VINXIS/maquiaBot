package gencommands

import (
	"bytes"
	"image"
	"image/png"
	"math"
	"net/http"
	"regexp"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/fogleman/gg"
)

// Meme lets you create a meme
func Meme(s *discordgo.Session, m *discordgo.MessageCreate) {
	linkRegex, _ := regexp.Compile(`https?:\/\/\S*`)
	memeRegex, _ := regexp.Compile(`meme\s+(https:\/\/(\S+)\s+)?([^|]+)?(\|)?([^|]+)?`)

	if !memeRegex.MatchString(m.Content) {
		s.ChannelMessageSend(m.ChannelID, "Please give text to add onto the image!")
		return
	}

	url := memeRegex.FindStringSubmatch(m.Content)[1]
	topText := strings.TrimSpace(memeRegex.FindStringSubmatch(m.Content)[3])
	bottomText := strings.TrimSpace(memeRegex.FindStringSubmatch(m.Content)[5])
	if topText == "" && bottomText == "" {
		s.ChannelMessageSend(m.ChannelID, "Please give text to add onto the image!")
		return
	}

	if len(m.Attachments) > 0 {
		url = m.Attachments[0].URL
	} else if len(m.Embeds) > 0 && m.Embeds[0].Image != nil {
		url = m.Embeds[0].Image.URL
	} else if url == "" {
		// Get prev messages
		messages, err := s.ChannelMessages(m.ChannelID, 100, "", "", "")
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, "Error fetching messages.")
			return
		}

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
	img, _, err := image.Decode(response.Body)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "Could not parse any image from data given.")
		return
	}
	r := img.Bounds()
	w := float64(r.Dx())
	h := float64(r.Dy())
	size := math.Min(40, float64(h)/10)

	cxt := gg.NewContext(r.Dx(), r.Dy())
	cxt.DrawImage(img, 0, 0)
	cxt.LoadFontFace("./data/fonts/impact.ttf", size)

	// Apply black stroke
	cxt.SetHexColor("#000")
	strokeSize := math.Min(3, float64(h)*3/200)
	for dy := -strokeSize; dy <= strokeSize; dy++ {
		for dx := -strokeSize; dx <= strokeSize; dx++ {
			// give it rounded corners
			if dx*dx+dy*dy >= strokeSize*strokeSize {
				continue
			}
			x := float64(w/2 + dx)
			y := size + float64(dy)
			cxt.DrawStringAnchored(strings.ToUpper(topText), x, y, 0.5, 0.5)
		}
	}
	for dy := -strokeSize; dy <= strokeSize; dy++ {
		for dx := -strokeSize; dx <= strokeSize; dx++ {
			// give it rounded corners
			if dx*dx+dy*dy >= strokeSize*strokeSize {
				continue
			}
			x := w/2 + dx
			y := float64(h) - size + dy
			cxt.DrawStringAnchored(strings.ToUpper(bottomText), x, y, 0.5, 0.5)
		}
	}

	// Apply white fill
	cxt.SetHexColor("#FFF")
	cxt.DrawStringAnchored(strings.ToUpper(topText), float64(w)/2, size, 0.5, 0.5)
	cxt.DrawStringAnchored(strings.ToUpper(bottomText), float64(w)/2, float64(h)-size, 0.5, 0.5)

	imgBytes := new(bytes.Buffer)
	_ = png.Encode(imgBytes, cxt.Image())
	s.ChannelMessageSendComplex(m.ChannelID, &discordgo.MessageSend{
		Files: []*discordgo.File{
			&discordgo.File{
				Name:   "image.png",
				Reader: imgBytes,
			},
		},
	})
}
