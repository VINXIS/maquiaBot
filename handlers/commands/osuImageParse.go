package commands

import (
	"fmt"
	"image"
	"image/png"
	"io/ioutil"
	"log"
	"math"
	"math/rand"
	"net/http"
	"os"
	"os/exec"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/disintegration/imaging"

	tools "../../tools"
	"github.com/bwmarrin/discordgo"
	"github.com/thehowl/go-osuapi"
)

// OsuImageParse detects for an osu image
func OsuImageParse(s *discordgo.Session, m *discordgo.MessageCreate, osu *osuapi.Client, hash map[int][5]int) {

	// Create regexps for checks
	mapperRegex, _ := regexp.Compile(`(?i)b?e?a?t?mapp?e?d? by (\S*)`)
	titleRegex, _ := regexp.Compile(`\- (.*) \[`)
	diagnosisRegex, _ := regexp.Compile(` -v`)

	var (
		name string
		url  string
	)

	if len(m.Attachments) > 0 {
		log.Println("Someone sent an image! The image URL is: " + m.Attachments[0].URL)

		name = strconv.Itoa(rand.Intn(10000000))
		url = m.Attachments[0].URL
	} else {
		regex, err := regexp.Compile(`https?:\/\/\S*`)
		tools.ErrRead(err)

		link := regex.FindStringSubmatch(m.Content)[0]
		log.Println("Someone sent a link! The URL is: " + link)

		name = strconv.Itoa(rand.Intn(10000000))
		url = link
	}

	// Fetch the image data
	response, err := http.Get(url)
	tools.ErrRead(err)
	imgSrc, _, err := image.Decode(response.Body)
	if err != nil {
		return
	}

	// Convert image to grayscale and raise contrast
	newImg := imaging.AdjustSaturation(imgSrc, -100)
	newImg = imaging.AdjustContrast(newImg, 100)
	b := newImg.Bounds()
	newImg = imaging.Crop(newImg, image.Rect(0, 0, int(math.Max(2.0*float64(b.Dx())/3.0, 1280.0)), int(math.Max(120.0*float64(b.Dy())/969.0, 120.0))))

	// Check if name already exists, create a new name via integer suffix instead if target name is currently in use
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
				i = i + 1
			}
		}
	}

	// Create the file to write in
	file, err := os.Create("./" + name + ".png")
	tools.ErrRead(err)

	// Dump the image data into the file
	png.Encode(file, newImg)
	tools.ErrRead(err)

	// Close file and res
	response.Body.Close()
	file.Close()

	// Run tesseract to parse the image
	_, err = exec.Command("tesseract", "./"+name+".png", name).Output()
	tools.ErrRead(err)

	// Read result and parse it
	text, err := ioutil.ReadFile(name + ".txt")
	tools.ErrRead(err)

	// Parse result
	raw := string(text)
	str := strings.Split(raw, "\n")
	if len(str) < 2 {
		deleteFile("./" + name + ".png")
		deleteFile("./" + name + ".txt")
		return
	}
	var (
		title      string
		mapperName string
	)

	for _, line := range str {
		if mapperRegex.MatchString(line) {
			mapperName = mapperRegex.FindStringSubmatch(line)[1]
		} else if titleRegex.MatchString(line) {
			title = titleRegex.FindStringSubmatch(line)[1]
		}
	}

	// See if the result was clean with a few checks
	if mapperName != "" && title != "" {
		message, _ := s.ChannelMessageSend(m.ChannelID, "Processing image...")
		var beatmap osuapi.Beatmap
		beatmaps, err := osu.GetBeatmaps(osuapi.GetBeatmapsOpts{
			Username: mapperName,
		})
		if err != nil {
			if diagnosisRegex.MatchString(m.Message.Content) {
				s.ChannelMessageEdit(message.ChannelID, message.ID, "No luck... the mapper line I parsed was ** "+mapperName+" ** and the title line I parsed was ** "+title+" **")
			} else {
				s.ChannelMessageDelete(message.ChannelID, message.ID)
			}
			deleteFile("./" + name + ".png")
			deleteFile("./" + name + ".txt")
			return
		}

		// Reorder the maps so that it returns the highest difficulty in the set
		sort.Slice(beatmaps, func(i, j int) bool {
			return beatmaps[i].DifficultyRating > beatmaps[j].DifficultyRating
		})

		// Look for the beatmap in the results
		for _, b := range beatmaps {
			if b.Title == title {
				beatmap = b
				break
			}
		}

		// Check if anything was actually found
		if beatmap == (osuapi.Beatmap{}) || len(beatmaps) == 0 {
			if diagnosisRegex.MatchString(m.Message.Content) {
				s.ChannelMessageEdit(message.ChannelID, message.ID, "No luck... the mapper line I parsed was ** "+mapperName+" ** and the title line I parsed was ** "+title+" **")
			} else {
				s.ChannelMessageDelete(message.ChannelID, message.ID)
			}
			deleteFile("./" + name + ".png")
			deleteFile("./" + name + ".txt")
			return
		}

		// Download the .osu file for the map
		replacer, _ := regexp.Compile(`[^a-zA-Z0-9\s\(\)]`)
		tools.DownloadFile(
			"./data/osuFiles/"+
				strconv.Itoa(beatmap.BeatmapID)+
				" "+
				replacer.ReplaceAllString(beatmap.Artist, "")+
				" - "+
				replacer.ReplaceAllString(beatmap.Title, "")+
				".osu",
			"https://osu.ppy.sh/osu/"+
				strconv.Itoa(beatmap.BeatmapID))
		// Assign embed colour for different modes
		Color := tools.ModeColour(beatmap.Mode)

		// Temporary method to obtain mapper user id, once creator id is available, actual user avatars will be used for banned users
		mapper, err := osu.GetUser(osuapi.GetUserOpts{
			Username: beatmap.Creator,
		})
		if err != nil {
			mapper, err = osu.GetUser(osuapi.GetUserOpts{
				UserID: 3,
			})
			mapper.Username = beatmap.Creator
		}

		// Obtain whole set
		beatmaps, err = osu.GetBeatmaps(osuapi.GetBeatmapsOpts{
			BeatmapSetID: beatmap.BeatmapSetID,
		})
		tools.ErrRead(err)

		// Assign variables for map specs
		totalMinutes := math.Floor(float64(beatmap.TotalLength / 60))
		totalSeconds := math.Mod(float64(beatmap.TotalLength), float64(60))
		hitMinutes := math.Floor(float64(beatmap.HitLength / 60))
		hitSeconds := math.Mod(float64(beatmap.HitLength), float64(60))

		starRating := "**SR:** " + fmt.Sprintf("%.2f", beatmap.DifficultyRating) + " "
		length := "**Length:** " + fmt.Sprint(totalMinutes) + ":" + fmt.Sprint(totalSeconds) + " (" + fmt.Sprint(hitMinutes) + ":" + fmt.Sprint(hitSeconds) + ") "
		bpm := "**BPM:** " + fmt.Sprint(beatmap.BPM) + " "
		combo := "**FC:** " + strconv.Itoa(beatmap.MaxCombo) + "x"

		status := "**Rank Status:** " + beatmap.Approved.String()

		download := "**Download:** [osz link](https://osu.ppy.sh/d/" + strconv.Itoa(beatmap.BeatmapSetID) + ")" + " | <osu://dl/" + strconv.Itoa(beatmap.BeatmapSetID) + ">"
		diffs := "**" + strconv.Itoa(len(beatmaps)) + `** difficulties <:ahFuck:550808614202245131>`

		// Calculate SR
		//aimRating, speedRating, totalRating := SRCalc(beatmap, "NM")

		// Calculate pp
		var (
			ppSS string
			pp99 string
			pp98 string
			pp97 string
			pp95 string
		)
		values, check := hash[beatmap.BeatmapID]
		if check {
			ppValueArray := values
			ppSS = "**100%:** " + strconv.Itoa(ppValueArray[0]) + "pp | "
			pp99 = "**99%:** " + strconv.Itoa(ppValueArray[1]) + "pp | "
			pp98 = "**98%:** " + strconv.Itoa(ppValueArray[2]) + "pp | "
			pp97 = "**97%:** " + strconv.Itoa(ppValueArray[3]) + "pp | "
			pp95 = "**95%:** " + strconv.Itoa(ppValueArray[4]) + "pp"
		} else {
			if beatmap.Mode != osuapi.ModeCatchTheBeat {
				s.ChannelMessageEdit(message.ChannelID, message.ID, "Calculating pp...")
				ppValues := make(chan int, 5)
				var ppValueArray [5]int
				go PPCalc(beatmap, 100.0, "NM", ppValues)
				go PPCalc(beatmap, 99.0, "NM", ppValues)
				go PPCalc(beatmap, 98.0, "NM", ppValues)
				go PPCalc(beatmap, 97.0, "NM", ppValues)
				go PPCalc(beatmap, 95.0, "NM", ppValues)
				for v := 0; v < 5; v++ {
					ppValueArray[v] = <-ppValues
				}
				sort.Slice(ppValueArray[:], func(i, j int) bool {
					return ppValueArray[i] > ppValueArray[j]
				})
				ppSS = "**100%:** " + strconv.Itoa(ppValueArray[0]) + "pp | "
				pp99 = "**99%:** " + strconv.Itoa(ppValueArray[1]) + "pp | "
				pp98 = "**98%:** " + strconv.Itoa(ppValueArray[2]) + "pp | "
				pp97 = "**97%:** " + strconv.Itoa(ppValueArray[3]) + "pp | "
				pp95 = "**95%:** " + strconv.Itoa(ppValueArray[4]) + "pp"
				hash[beatmap.BeatmapID] = ppValueArray
			} else {
				ppSS = "pp is not available for ctb yet"
				pp99 = ""
				pp98 = ""
				pp97 = ""
				pp95 = ""
			}
		}

		embed := &discordgo.MessageEmbed{
			Author: &discordgo.MessageEmbedAuthor{
				URL:     "https://osu.ppy.sh/beatmaps/" + strconv.Itoa(beatmap.BeatmapID),
				Name:    beatmap.Artist + " - " + beatmap.Title + " by " + mapper.Username,
				IconURL: "https://a.ppy.sh/" + strconv.Itoa(mapper.UserID),
			},
			Color: Color,
			Description: starRating + length + bpm + combo + "\n" +
				status + "\n" +
				download + "\n" +
				diffs + "\n" + "\n" +
				"**[" + beatmap.DiffName + "]**\n" +
				//aimRating + speedRating + totalRating + "\n" +
				ppSS + pp99 + pp98 + pp97 + pp95,
			Thumbnail: &discordgo.MessageEmbedThumbnail{
				URL: "https://b.ppy.sh/thumb/" + strconv.Itoa(beatmap.BeatmapSetID) + "l.jpg",
			},
		}
		s.ChannelMessageEdit(message.ChannelID, message.ID, "")
		s.ChannelMessageEditEmbed(message.ChannelID, message.ID, embed)

	} else {
		deleteFile("./" + name + ".png")
		deleteFile("./" + name + ".txt")
		return
	}

	// Close files
	deleteFile("./" + name + ".png")
	deleteFile("./" + name + ".txt")
	return
}

func deleteFile(path string) {
	var err = os.Remove(path)
	tools.ErrRead(err)
	return
}
