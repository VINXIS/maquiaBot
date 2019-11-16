package gencommands

import (
	cryptoRand "crypto/rand"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"math/big"
	"math/rand"
	"regexp"
	"strconv"
	"strings"
	"time"

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
	serverData := tools.GetServer(*server)

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

// Adjectives allows users to add/change/see their custom adjectives
func Adjectives(s *discordgo.Session, m *discordgo.MessageCreate) {
	adjectivesRegex, _ := regexp.Compile(`(adj|adjectives?)\s*(add|remove)?\s+(.+)`)

	server, err := s.Guild(m.GuildID)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "This is not a server!")
		return
	}

	// Obtain server data
	serverData := tools.GetServer(*server)

	if m.Author.ID != server.OwnerID || !serverData.AllowAnyoneStats {
		member, _ := s.GuildMember(server.ID, m.Author.ID)
		admin := false
		for _, roleID := range member.Roles {
			role, _ := s.State.Role(m.GuildID, roleID)
			if role.Permissions&discordgo.PermissionAdministrator != 0 || role.Permissions&discordgo.PermissionManageServer != 0 {
				admin = true
				break
			}
		}
		if !admin && len(m.Mentions) >= 1 {
			s.ChannelMessageSend(m.ChannelID, "You must be an admin, server manager, or server owner, or one of them must allow anyone to add words.")
			return
		}
	}

	// Obtain word and if they want to add/remove it or just see the current ones
	mode := "add"
	word := ""
	if adjectivesRegex.MatchString(m.Content) {
		matches := adjectivesRegex.FindStringSubmatch(m.Content)
		if matches[2] != "" {
			mode = matches[2]
		}
		word = strings.ReplaceAll(matches[3], "`", "")
	} else if len(serverData.Adjectives) != 0 {
		text := "There are " + strconv.Itoa(len(serverData.Adjectives)) + " adjectives listed for this server! The adjectives are:\n"
		for i, adjective := range serverData.Adjectives {
			if i != len(serverData.Adjectives)-1 {
				text = text + adjective + ", "
			} else {
				text = text + adjective
			}
		}
		s.ChannelMessageSend(m.ChannelID, text)
		return
	}

	if word == "" {
		s.ChannelMessageSend(m.ChannelID, "No word to add/remove from the server's adjective list found!")
		return
	}

	// Commence operation
	err = serverData.Word(word, mode, "adjective")
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, err.Error())
		return
	}

	// Set new information in server data
	serverData.Time = time.Now()

	jsonCache, err := json.Marshal(serverData)
	tools.ErrRead(err)

	err = ioutil.WriteFile("./data/serverData/"+m.GuildID+".json", jsonCache, 0644)
	tools.ErrRead(err)
	if mode == "add" {
		s.ChannelMessageSend(m.ChannelID, "`"+word+"` is now added to the server's adjective list!")
	} else if mode == "remove" {
		s.ChannelMessageSend(m.ChannelID, "`"+word+"` is now removed from the server's adjective list!")
	}
	return
}

// Nouns allows users to add/change/see their custom nouns
func Nouns(s *discordgo.Session, m *discordgo.MessageCreate) {
	nounsRegex, _ := regexp.Compile(`(nouns?)\s*(add|remove)?\s+(.+)`)

	server, err := s.Guild(m.GuildID)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "This is not a server!")
		return
	}

	// Obtain server data
	serverData := tools.GetServer(*server)

	if m.Author.ID != server.OwnerID || !serverData.AllowAnyoneStats {
		member, _ := s.GuildMember(server.ID, m.Author.ID)
		admin := false
		for _, roleID := range member.Roles {
			role, _ := s.State.Role(m.GuildID, roleID)
			if role.Permissions&discordgo.PermissionAdministrator != 0 || role.Permissions&discordgo.PermissionManageServer != 0 {
				admin = true
				break
			}
		}
		if !admin && len(m.Mentions) >= 1 {
			s.ChannelMessageSend(m.ChannelID, "You must be an admin, server manager, or server owner, or one of them must allow anyone to add words.")
			return
		}
	}

	// Obtain word and if they want to add/remove it or just see the current ones
	mode := "add"
	word := ""
	if nounsRegex.MatchString(m.Content) {
		matches := nounsRegex.FindStringSubmatch(m.Content)
		if matches[2] != "" {
			mode = matches[2]
		}
		word = strings.ReplaceAll(matches[3], "`", "")
	} else if len(serverData.Nouns) != 0 {
		text := "There are " + strconv.Itoa(len(serverData.Nouns)) + " nouns listed for this server! The nouns are:\n"
		for i, noun := range serverData.Nouns {
			if i != len(serverData.Nouns)-1 {
				text = text + noun + ", "
			} else {
				text = text + noun
			}
		}
		s.ChannelMessageSend(m.ChannelID, text)
		return
	}

	if word == "" {
		s.ChannelMessageSend(m.ChannelID, "No word to add/remove from the server's noun list found!")
		return
	}

	// Commence operation
	err = serverData.Word(word, mode, "noun")
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, err.Error())
		return
	}

	// Set new information in server data
	serverData.Time = time.Now()

	jsonCache, err := json.Marshal(serverData)
	tools.ErrRead(err)

	err = ioutil.WriteFile("./data/serverData/"+m.GuildID+".json", jsonCache, 0644)
	tools.ErrRead(err)
	if mode == "add" {
		s.ChannelMessageSend(m.ChannelID, "`"+word+"` is now added to the server's noun list!")
	} else if mode == "remove" {
		s.ChannelMessageSend(m.ChannelID, "`"+word+"` is now removed from the server's noun list!")
	}
	return
}

// Skills allows users to add/change/see their custom skills
func Skills(s *discordgo.Session, m *discordgo.MessageCreate) {
	skillsRegex, _ := regexp.Compile(`(skills?)\s*(add|remove)?\s+(.+)`)

	server, err := s.Guild(m.GuildID)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "This is not a server!")
		return
	}

	// Obtain server data
	serverData := tools.GetServer(*server)

	if m.Author.ID != server.OwnerID || !serverData.AllowAnyoneStats {
		member, _ := s.GuildMember(server.ID, m.Author.ID)
		admin := false
		for _, roleID := range member.Roles {
			role, _ := s.State.Role(m.GuildID, roleID)
			if role.Permissions&discordgo.PermissionAdministrator != 0 || role.Permissions&discordgo.PermissionManageServer != 0 {
				admin = true
				break
			}
		}
		if !admin && len(m.Mentions) >= 1 {
			s.ChannelMessageSend(m.ChannelID, "You must be an admin, server manager, or server owner, or one of them must allow anyone to add words.")
			return
		}
	}

	// Obtain word and if they want to add/remove it or just see the current ones
	mode := "add"
	word := ""
	if skillsRegex.MatchString(m.Content) {
		matches := skillsRegex.FindStringSubmatch(m.Content)
		if matches[2] != "" {
			mode = matches[2]
		}
		word = strings.ReplaceAll(matches[3], "`", "")
	} else if len(serverData.Skills) != 0 {
		text := "There are " + strconv.Itoa(len(serverData.Skills)) + " skills listed for this server! The skills are:\n"
		for i, skill := range serverData.Skills {
			if i != len(serverData.Skills)-1 {
				text = text + skill + ", "
			} else {
				text = text + skill
			}
		}
		s.ChannelMessageSend(m.ChannelID, text)
		return
	}

	if word == "" {
		s.ChannelMessageSend(m.ChannelID, "No word to add/remove from the server's skill list found!")
		return
	}

	// Commence operation
	err = serverData.Word(word, mode, "skill")
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, err.Error())
		return
	}

	// Set new information in server data
	serverData.Time = time.Now()

	jsonCache, err := json.Marshal(serverData)
	tools.ErrRead(err)

	err = ioutil.WriteFile("./data/serverData/"+m.GuildID+".json", jsonCache, 0644)
	tools.ErrRead(err)
	if mode == "add" {
		s.ChannelMessageSend(m.ChannelID, "`"+word+"` is now added to the server's skill list!")
	} else if mode == "remove" {
		s.ChannelMessageSend(m.ChannelID, "`"+word+"` is now removed from the server's skill list!")
	}
	return
}

// StatsToggle allows users to add/remove nouns, adjectives, and skills
func StatsToggle(s *discordgo.Session, m *discordgo.MessageCreate) {
	server, err := s.Guild(m.GuildID)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "This is not a server!")
		return
	}

	if m.Author.ID != server.OwnerID {
		member, _ := s.GuildMember(server.ID, m.Author.ID)
		admin := false
		for _, roleID := range member.Roles {
			role, _ := s.State.Role(m.GuildID, roleID)
			if role.Permissions&discordgo.PermissionAdministrator != 0 || role.Permissions&discordgo.PermissionManageServer != 0 {
				admin = true
				break
			}
		}
		if !admin && len(m.Mentions) >= 1 {
			s.ChannelMessageSend(m.ChannelID, "You must be an admin, server manager, or server owner!")
			return
		}
	}

	// Obtain server data
	serverData := tools.GetServer(*server)

	// Set new information in server data
	serverData.Time = time.Now()
	serverData.AllowAnyoneStats = !serverData.AllowAnyoneStats

	jsonCache, err := json.Marshal(serverData)
	tools.ErrRead(err)

	err = ioutil.WriteFile("./data/serverData/"+m.GuildID+".json", jsonCache, 0644)
	tools.ErrRead(err)

	if serverData.Crab {
		s.ChannelMessageSend(m.ChannelID, "Anyone may now add nouns/adjectives/skills for stats now.")
	} else {
		s.ChannelMessageSend(m.ChannelID, "Only admins, server managers, or the server owner may add nouns/adjectives/skills for stats now.")
	}
	return
}
