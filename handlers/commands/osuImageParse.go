package commands

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"

	tools "../../tools"
	"github.com/bwmarrin/discordgo"
	"github.com/thehowl/go-osuapi"
)

// OsuImageParse detects for an osu image
func OsuImageParse(s *discordgo.Session, m *discordgo.MessageCreate, osu *osuapi.Client) {
	log.Println("Someone linked an image! the image URL is: " + m.Attachments[0].URL)
	message, err := s.ChannelMessageSend(m.ChannelID, "Checking if image is from osu!...")

	name := m.Attachments[0].Filename
	url := m.Attachments[0].URL

	// Fetch the image data
	response, err := http.Get(url)
	tools.ErrRead(err, "30", "osuImageParse.go")

	// Check if name already exists, create a seperate name instead if name is currently in use
	_, err1 := os.Stat("./" + name + ".jpg")
	_, err2 := os.Stat("./" + name + ".txt")
	if err1 != nil || err2 != nil {
		if !os.IsNotExist(err1) || !os.IsNotExist(err2) {
			i := 1
			for {
				s.ChannelMessageEdit(message.ChannelID, message.ID, name+strconv.Itoa(i))
				_, err1 := os.Stat("./" + name + strconv.Itoa(i) + ".jpg")
				_, err2 := os.Stat("./" + name + strconv.Itoa(i) + ".txt")
				if err1 != nil || err2 != nil {
					if os.IsNotExist(err1) && os.IsNotExist(err2) {
						name = name + strconv.Itoa(i)
						break
					} else {
						i = i + 1
					}
				} else {
					name = name + strconv.Itoa(i)
					break
				}
			}
		}
	}

	// Create the file to write in
	file, err := os.Create("./" + name + ".jpg")
	tools.ErrRead(err, "59", "osuImageParse.go")

	// Dump the image data into the file
	_, err = io.Copy(file, response.Body)
	tools.ErrRead(err, "63", "osuImageParse.go")
	response.Body.Close()
	file.Close()

	// Run tesseract to parse the image
	_, err = exec.Command("tesseract", "./"+name+".jpg", name).Output()
	tools.ErrRead(err, "69", "osuImageParse.go")

	// Read result and parse it
	text, err := ioutil.ReadFile(name + ".txt")
	tools.ErrRead(err, "73", "osuImageParse.go")

	// Parse result
	raw := string(text)
	str := strings.Split(raw, "\n")

	// Create regexps for checks
	mapperRegex, _ := regexp.Compile(`B?e?a?t?map by (\S*)`)
	titleRegex, _ := regexp.Compile(`\- (.*) \[`)
	diagnosisRegex, _ := regexp.Compile(`(-v)`)

	if len(str) < 2 {
		if diagnosisRegex.MatchString(m.Content) {
			s.ChannelMessageEdit(message.ChannelID, message.ID, "No luck... could not find a mapper not a title in the image!")
		} else {
			s.ChannelMessageDelete(message.ChannelID, message.ID)
		}
		return
	}
	title := str[0]
	mapper := str[1]

	// Check if the result was clean
	if mapperRegex.MatchString(mapper) && titleRegex.MatchString(title) {
		var beatmap osuapi.Beatmap
		r := mapperRegex.FindStringSubmatch(mapper)
		t := titleRegex.FindStringSubmatch(title)
		s.ChannelMessageEdit(message.ChannelID, message.ID, "Possible beatmap match found! Doing an API call with `"+r[1]+"` as mapper, and `"+t[1]+"` as title...")

		beatmaps, err := osu.GetBeatmaps(osuapi.GetBeatmapsOpts{
			Username: r[1],
		})
		if err != nil {
			if diagnosisRegex.MatchString(m.Content) {
				s.ChannelMessageEdit(message.ChannelID, message.ID, "No luck... the mapper line I parsed was `"+mapper+"` and the title line I parsed was `"+title+"`")
			} else {
				s.ChannelMessageDelete(message.ChannelID, message.ID)
			}
			return
		}

		// Look for the beatmap in the results
		for _, b := range beatmaps {
			if b.Title == t[1] {
				beatmap = b
				s.ChannelMessageEdit(message.ChannelID, message.ID, "Beatmap found! Calculating specs...")
				break
			}
		}

		// Check if anything was actually found
		if beatmap == (osuapi.Beatmap{}) || len(beatmaps) == 0 {
			if diagnosisRegex.MatchString(m.Content) {
				s.ChannelMessageEdit(message.ChannelID, message.ID, "No luck... the mapper line I parsed was `"+mapper+"` and the title line I parsed was `"+title+"`")
			} else {
				s.ChannelMessageDelete(message.ChannelID, message.ID)
			}
			return
		}

		// Download the .osu file for the map
		err = tools.DownloadFile(
			"./data/osuFiles/"+
				strconv.Itoa(beatmap.BeatmapID)+
				" "+
				beatmap.Artist+
				" - "+
				beatmap.Title+
				".osu",
			"https://osu.ppy.sh/osu/"+
				strconv.Itoa(beatmap.BeatmapID))
		tools.ErrRead(err, "135", "osuImageParse.go")

		// Assign embed colour for different modes
		var Color int
		switch beatmap.Mode {
		case osuapi.ModeOsu:
			Color = 0xD65288
		case osuapi.ModeTaiko:
			Color = 0xFF0000
		case osuapi.ModeCatchTheBeat:
			Color = 0x007419
		case osuapi.ModeOsuMania:
			Color = 0xff6200
		}

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
		tools.ErrRead(err, "172", "osuImageParse.go")

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
		ppSS := "**100%:** " + PPCalc(beatmap, 100.0, "NM") + " | "
		pp99 := "**99%:** " + PPCalc(beatmap, 99.0, "NM") + " | "
		pp98 := "**98%:** " + PPCalc(beatmap, 98.0, "NM") + " | "
		pp97 := "**97%:** " + PPCalc(beatmap, 97.0, "NM") + " | "
		pp95 := "**95%:** " + PPCalc(beatmap, 95.0, "NM")

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

	} else if diagnosisRegex.MatchString(m.Content) {
		s.ChannelMessageEdit(message.ChannelID, message.ID, "No luck... the mapper line I parsed was `"+mapper+"` and the title line I parsed was `"+title+"`")
	} else {
		s.ChannelMessageDelete(message.ChannelID, message.ID)
	}

	// Close files
	deleteFile("./" + name + ".jpg")
	deleteFile("./" + name + ".txt")
}

func deleteFile(path string) {
	var err = os.Remove(path)
	tools.ErrRead(err, "235", "osuImageParse.go")
	return
}
