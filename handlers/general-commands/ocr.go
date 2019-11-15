package gencommands

import (
	"image"
	"image/png"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"os/exec"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	tools "../../tools"
	"github.com/bwmarrin/discordgo"
	"github.com/disintegration/imaging"
)

// OCR lets people use the tesseract-OCR utility on their images
func OCR(s *discordgo.Session, m *discordgo.MessageCreate) {
	linkRegex, _ := regexp.Compile(`https?:\/\/\S*`)
	saturationRegex, _ := regexp.Compile(`-s\s+(-?\d+)`)
	contrastRegex, _ := regexp.Compile(`-c\s+(-?\d+)`)
	langRegex, _ := regexp.Compile(`-l\s+(\S+)`)

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
			time1, err := time.Parse(time.RFC3339, string(messages[i].Timestamp))
			tools.ErrRead(err)
			time2, err := time.Parse(time.RFC3339, string(messages[j].Timestamp))
			tools.ErrRead(err)
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
	tools.ErrRead(err)

	// Dump the image data into the file
	png.Encode(file, newImg)

	// Close file and res
	response.Body.Close()
	file.Close()

	// Run tesseract to parse the image
	_, err = exec.Command("tesseract", "./"+name+".png", name, "--dpi", "96", "-l", lang).Output()
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "Invalid language!")
		tools.DeleteFile("./" + name + ".png")
		tools.DeleteFile("./" + name + ".txt")
		return
	}

	// Read result and parse it
	text, err := ioutil.ReadFile(name + ".txt")
	tools.ErrRead(err)

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
	s.ChannelMessageSend(m.ChannelID, "```"+str+"```")
	return
}
