package gencommands

import (
	"bytes"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"io/ioutil"
	"math"
	"math/rand"
	"net/http"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"

	tools "maquiaBot/tools"
	"github.com/bwmarrin/discordgo"
	"github.com/disintegration/imaging"
	pigo "github.com/esimov/pigo/core"
	"github.com/fogleman/gg"
)

// OCR lets people use the tesseract-OCR utility on their images
func OCR(s *discordgo.Session, m *discordgo.MessageCreate) {
	linkRegex, _ := regexp.Compile(`(?i)https?:\/\/\S*`)
	saturationRegex, _ := regexp.Compile(`(?i)-s\s+(-?\d+)`)
	contrastRegex, _ := regexp.Compile(`(?i)-c\s+(-?\d+)`)
	langRegex, _ := regexp.Compile(`(?i)-l\s+(\S+)`)

	var url string
	if len(m.Attachments) > 0 {
		url = m.Attachments[0].URL
	} else if len(m.Embeds) > 0 && m.Embeds[0].Image != nil {
		url = m.Embeds[0].Image.URL
	} else if linkRegex.MatchString(m.Content) {
		url = linkRegex.FindStringSubmatch(m.Content)[0]
	} else {
		// Get prev messages
		messages, err := s.ChannelMessages(m.ChannelID, -1, "", "", "")
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

	saturationValue := float64(-100)
	contrastValue := float64(0)
	lang := "eng"
	if saturationRegex.MatchString(m.Content) {
		saturationValue, _ = strconv.ParseFloat(saturationRegex.FindStringSubmatch(m.Content)[1], 64)
		if saturationValue < -100 {
			saturationValue = -100
		} else if saturationValue > 100 {
			saturationValue = 100
		}
	}
	if contrastRegex.MatchString(m.Content) {
		contrastValue, _ = strconv.ParseFloat(contrastRegex.FindStringSubmatch(m.Content)[1], 64)
		if contrastValue < -100 {
			contrastValue = -100
		} else if contrastValue > 100 {
			contrastValue = 100
		}
	}
	if langRegex.MatchString(m.Content) {
		lang = langRegex.FindStringSubmatch(m.Content)[1]
	} else if strings.Contains(m.Content, "-l") {
		s.ChannelMessageSend(m.ChannelID, "https://en.wikipedia.org/wiki/ISO_639-2")
		return
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

	// Convert image to grayscale and raise contrast
	newImg := imaging.AdjustSaturation(imgSrc, saturationValue)
	newImg = imaging.AdjustContrast(newImg, contrastValue)

	// Check if name already exists, create a new name via integer suffix instead if target name is currently in use
	name := strconv.Itoa(rand.Intn(10000000))
	_, err1 := os.Stat("./" + name + ".png")
	_, err2 := os.Stat("./" + name + ".txt")
	if !os.IsNotExist(err1) || !os.IsNotExist(err2) {
		i := 1
		for {
			_, err1 := os.Stat("./" + name + strconv.Itoa(i) + ".png")
			_, err2 := os.Stat("./" + name + strconv.Itoa(i) + ".txt")
			if os.IsNotExist(err1) && os.IsNotExist(err2) {
				name = name + strconv.Itoa(i)
				break
			} else {
				i++
			}
		}
	}

	// Create the file to write in
	file, err := os.Create("./" + name + ".png")
	tools.ErrRead(s, err)

	// Dump the image data into the file
	png.Encode(file, newImg)

	// Close file and res
	response.Body.Close()
	file.Close()

	// Run tesseract to parse the image
	_, err = exec.Command("tesseract", "./"+name+".png", name, "-l", lang, "--oem", "3", "--psm", "12").Output()
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "Invalid language!")
		tools.DeleteFile("./" + name + ".png")
		tools.DeleteFile("./" + name + ".txt")
		return
	}

	// Read result and parse it
	text, err := ioutil.ReadFile(name + ".txt")
	tools.ErrRead(s, err)

	// Parse result
	str := string(text)

	// Delete files
	if !(m.Author.ID == "92502458588205056" && strings.Contains(m.Content, "-t")) {
		tools.DeleteFile("./" + name + ".png")
		tools.DeleteFile("./" + name + ".txt")
	}

	if len(strings.TrimSpace(str)) <= 1 {
		s.ChannelMessageSend(m.ChannelID, "No text found...")
		return
	}
	_, err = s.ChannelMessageSend(m.ChannelID, "```"+str+"```")
	if err != nil {
		s.ChannelMessageSendComplex(m.ChannelID, &discordgo.MessageSend{
			Files: []*discordgo.File{
				{
					Name:   "ocr.txt",
					Reader: strings.NewReader(str),
				},
			},
		})
	}
	return
}

// Face lets you detect faces
func Face(s *discordgo.Session, m *discordgo.MessageCreate) {
	linkRegex, _ := regexp.Compile(`(?i)https?:\/\/\S*`)

	var url string
	if len(m.Attachments) > 0 {
		url = m.Attachments[0].URL
	} else if len(m.Embeds) > 0 && m.Embeds[0].Image != nil {
		url = m.Embeds[0].Image.URL
	} else if linkRegex.MatchString(m.Content) {
		url = linkRegex.FindStringSubmatch(m.Content)[0]
	} else {
		// Get prev messages
		messages, err := s.ChannelMessages(m.ChannelID, -1, "", "", "")
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
	s.ChannelMessageSendComplex(m.ChannelID, &discordgo.MessageSend{
		Files: []*discordgo.File{
			{
				Name:   "image.png",
				Reader: imgBytes,
			},
		},
	})
}
