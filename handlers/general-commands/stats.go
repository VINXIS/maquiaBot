package gencommands

import (
	cryptoRand "crypto/rand"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"math/big"
	"math/rand"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	structs "../../structs"
	tools "../../tools"
	"github.com/bwmarrin/discordgo"
)

// Stats creates and outputs randomized stats for the user in question
func Stats(s *discordgo.Session, m *discordgo.MessageCreate) {
	statsRegex, _ := regexp.Compile(`(stats)\s+(.+)?`)
	prefixRegex, _ := regexp.Compile(`(.+)(stats.*)`)

	server, err := s.Guild(m.GuildID)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "This is not a server!")
		return
	}

	prefix := prefixRegex.FindStringSubmatch(m.Content)[1]

	// Parse emssage to see if a skill count was given/object of reference
	text := "You have"
	skillCount := 4
	if statsRegex.MatchString(m.Content) {
		var err error
		skillCount, err = strconv.Atoi(statsRegex.FindStringSubmatch(m.Content)[2])
		if err != nil {
			text = strings.ReplaceAll(statsRegex.FindStringSubmatch(m.Content)[2]+" has", "`", "")
			list := strings.Split(statsRegex.FindStringSubmatch(m.Content)[2], " ")
			skillCount = 4
			if len(list) > 1 {
				skillCount, err = strconv.Atoi(list[len(list)-1])
				if err != nil {
					skillCount = 4
				} else {
					text = strings.ReplaceAll(strings.Join(list[:len(list)-1], " "), "`", "") + " has"
				}
			}
		}
	}

	// Obtain server data
	serverData := structs.ServerData{
		Prefix:    "$",
		Crab:      true,
		OsuToggle: true,
	}
	_, err = os.Stat("./data/serverData/" + m.GuildID + ".json")
	if err == nil {
		f, err := ioutil.ReadFile("./data/serverData/" + m.GuildID + ".json")
		tools.ErrRead(err)
		_ = json.Unmarshal(f, &serverData)
	} else if os.IsNotExist(err) {
		serverData.Server = *server
	} else {
		tools.ErrRead(err)
		return
	}

	// Check if the minimum amount of skills, nouns, and adjectives are there
	if len(serverData.Skills) < skillCount {
		s.ChannelMessageSend(m.ChannelID, "You need at least "+strconv.Itoa(skillCount)+" skills! Add skills with `"+prefix+"skills add`")
		return
	} else if len(serverData.Nouns) == 0 {
		s.ChannelMessageSend(m.ChannelID, "You have no nouns! Add skills with `"+prefix+"nouns add`")
		return
	} else if len(serverData.Adjectives) == 0 {
		s.ChannelMessageSend(m.ChannelID, "You need no adjectives! Add skills with `"+prefix+"adj(ectives) add`")
		return
	}

	// Obtain 4 skills
	var skills []string
	maxLength := float64(0)
	for len(skills) < skillCount {
		randNum, _ := cryptoRand.Int(cryptoRand.Reader, big.NewInt(int64(len(serverData.Skills))))
		newSkill := serverData.Skills[randNum.Int64()]
		contained := false
		for _, skill := range skills {
			if skill == newSkill {
				contained = true
				break
			}
		}
		if !contained {
			maxLength = math.Max(maxLength, float64(len(newSkill)))
			skills = append(skills, newSkill)
		}
	}

	// Obtain the percentages of the skills, alongside an adjective and noun
	fullText := "```"
	skillRang := rand.New(rand.NewSource(time.Now().UnixNano()))
	for _, skill := range skills {
		percent := math.Max(0, math.Min(100, skillRang.NormFloat64()*12.5+50))
		bar := tools.BarCreation(percent / 100)
		fullText = fullText + fmt.Sprintf("%-"+strconv.FormatFloat(maxLength, 'f', 0, 64)+"s: %3d%% %s", skill, int(percent), bar) + "\n"
	}
	randNum, _ := cryptoRand.Int(cryptoRand.Reader, big.NewInt(int64(len(serverData.Adjectives))))
	adjective := serverData.Adjectives[randNum.Int64()]

	randNum, _ = cryptoRand.Int(cryptoRand.Reader, big.NewInt(int64(len(serverData.Nouns))))
	noun := serverData.Nouns[randNum.Int64()]

	fullText = fullText + "\n" + text + " chosen the \"" + adjective + " " + noun + "\" class.```"
	_, err = s.ChannelMessageSend(m.ChannelID, fullText)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "Message probably went over the 2000 character limit!")
	}
	return
}
