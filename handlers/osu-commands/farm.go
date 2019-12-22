package osucommands

import (
	"encoding/json"
	"io/ioutil"
	"math"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	osuapi "../../osu-api"
	structs "../../structs"
	tools "../../tools"
	"github.com/bwmarrin/discordgo"
)

// Farm gives a player's farmerdog rating
func Farm(s *discordgo.Session, m *discordgo.MessageCreate, cache []structs.PlayerData) {
	username := ""
	user := structs.PlayerData{}
	cached := false
	playerIndex := 0
	amount := 5

	// Obtain user if any user was stated
	userRegex, _ := regexp.Compile(`(.+)farm\s*(.+)?`)
	amountRegex, _ := regexp.Compile(`-a\s*(\d*)`)
	prefix := userRegex.FindStringSubmatch(m.Content)[1]
	username = userRegex.FindStringSubmatch(m.Content)[2]
	if amountRegex.MatchString(m.Content) {
		amount, _ = strconv.Atoi(amountRegex.FindStringSubmatch(m.Content)[1])
		username = strings.TrimSpace(strings.Replace(username, amountRegex.FindStringSubmatch(m.Content)[0], "", -1))
	}

	// Check for user
	if username == "" {
		for i, player := range cache {
			if m.Author.ID == player.Discord.ID {
				if player.Osu.Username == "" {
					s.ChannelMessageSend(m.ChannelID, "No user linked to your discord account! Use "+prefix+"link to link your account!")
					return
				}
				playerIndex = i
				cached = true
				user = player
				username = player.Osu.Username
				break
			}
		}
	} else {
		for i, player := range cache {
			if username == strings.ToLower(player.Osu.Username) {
				playerIndex = i
				cached = true
				user = player
				break
			}
		}
	}

	if username == "" {
		s.ChannelMessageSend(m.ChannelID, "No user given!")
		return
	}

	// Get User and see if user exists
	osuUser, err := OsuAPI.GetUser(osuapi.GetUserOpts{
		Username: username,
	})
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "User: **"+user.Osu.Username+"** may not exist!")
		return
	}
	user.Osu = *osuUser

	// Obtain farm data
	farmData := structs.FarmData{}
	f, err := ioutil.ReadFile("./data/osuData/mapFarm.json")
	tools.ErrRead(err)
	_ = json.Unmarshal(f, &farmData)

	// Add the new information to the full data
	user.FarmCalc(OsuAPI, farmData)
	user.Time = time.Now()
	if cached {
		cache[playerIndex] = user
	} else {
		cache = append(cache, user)
	}

	// Save info
	jsonCache, err := json.Marshal(cache)
	tools.ErrRead(err)

	err = ioutil.WriteFile("./data/osuData/profileCache.json", jsonCache, 0644)
	tools.ErrRead(err)

	// Get the list of scores
	sort.Slice(user.Farm.List, func(i, j int) bool { return user.Farm.List[i].FarmScore > user.Farm.List[j].FarmScore })
	list := ""
	if amount > len(user.Farm.List) {
		amount = len(user.Farm.List)
	}
	for n := 0; n < amount; n++ {
		list += "`" + user.Farm.List[n].Name + "`: " + strconv.FormatFloat(user.Farm.List[n].FarmScore, 'f', 2, 64) + " Farmerdog rating (" + strconv.FormatFloat(user.Farm.List[n].PP, 'f', 2, 64) + ")\n"
	}

	s.ChannelMessageSend(m.ChannelID, "Farmerdog rating for **"+user.Osu.Username+":** "+strconv.FormatFloat(user.Farm.Rating, 'f', 2, 64)+"\n**Your top "+strconv.Itoa(amount)+" farmerdog scores:** \n"+list)
}

// TopFarm gives the top farmerdogs in the game based on who's been run
func TopFarm(s *discordgo.Session, m *discordgo.MessageCreate, cache []structs.PlayerData) {
	farmCountRegex, _ := regexp.Compile(`.+(tfarm|topfarm)\s*(-s)?\s*(\d+)?`)

	farmAmount := 1

	if strings.Contains(m.Content, "-s") {
		members, err := s.GuildMembers(m.GuildID, "", 1000)
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, "This is not a server!")
			return
		}
		trueCache := []structs.PlayerData{}

		for _, player := range cache {
			for _, member := range members {
				if player.Discord.ID == member.User.ID && math.Round(player.Farm.Rating*100)/100 != 0.00 {
					trueCache = append(trueCache, player)
					break
				}
			}
		}

		cache = trueCache
	}

	sort.Slice(cache, func(i, j int) bool {
		return cache[i].Farm.Rating > cache[j].Farm.Rating
	})

	farmCount := farmCountRegex.FindStringSubmatch(m.Content)[3]

	if farmCount != "" {
		farmAmount, _ = strconv.Atoi(farmCount)
	}

	if farmAmount == 1 {
		if strings.Contains(m.Content, "-s") {
			s.ChannelMessageSend(m.ChannelID, "The worst farmerdog in the server is **"+cache[0].Osu.Username+"** with a farmerdog rating of "+strconv.FormatFloat(cache[0].Farm.Rating, 'f', 2, 64))
			return
		}
		s.ChannelMessageSend(m.ChannelID, "The worst farmerdog is **"+cache[0].Osu.Username+"** with a farmerdog rating of "+strconv.FormatFloat(cache[0].Farm.Rating, 'f', 2, 64))
		return
	} else if farmAmount > len(cache) {
		farmAmount = len(cache)
	}

	msg := "**Top farmerdogs:** \n"
	if strings.Contains(m.Content, "-s") {
		msg = "**Top farmerdogs in this server:** \n"
	}
	max := 0

	for i := 0; i < farmAmount; i++ {
		if len(msg) >= 2000 {
			max = i + 1
			break
		}

		msg = msg + "#" + strconv.Itoa(i+1) + ": **" + cache[i].Osu.Username + "** - " + strconv.FormatFloat(cache[i].Farm.Rating, 'f', 2, 64) + " farmerdog rating \n"
	}

	if len(msg) > 2000 {
		for {
			lines := strings.Split(msg, "\n")
			lines = lines[:len(lines)-1]
			msg = strings.Join(lines, "\n")
			if len(msg) <= 2000 {
				break
			}
		}
	}

	if max == 0 {
		s.ChannelMessageSend(m.ChannelID, msg)
	} else {
		s.ChannelMessageSend(m.ChannelID, "Only showing top "+strconv.Itoa(max)+" farmerdogs")
		s.ChannelMessageSend(m.ChannelID, msg)
	}
}

// BottomFarm gives the top farmerdogs in the game based on who's been run
func BottomFarm(s *discordgo.Session, m *discordgo.MessageCreate, cache []structs.PlayerData) {
	farmCountRegex, _ := regexp.Compile(`.+(bfarm|bottomfarm)\s*(-s)?\s*(\d+)?`)

	farmAmount := 1

	if strings.Contains(m.Content, "-s") {
		members, err := s.GuildMembers(m.GuildID, "", 1000)
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, "This is not a server!")
			return
		}
		trueCache := []structs.PlayerData{}

		for _, player := range cache {
			for _, member := range members {
				if player.Discord.ID == member.User.ID && player.Farm.Rating != 0.00 {
					trueCache = append(trueCache, player)
					break
				}
			}
		}

		cache = trueCache
	} else {
		trueCache := []structs.PlayerData{}

		for _, player := range cache {
			if player.Farm.Rating != 0.00 {
				trueCache = append(trueCache, player)
			}
		}

		cache = trueCache
	}

	sort.Slice(cache, func(i, j int) bool {
		return cache[i].Farm.Rating < cache[j].Farm.Rating
	})

	farmCount := farmCountRegex.FindStringSubmatch(m.Content)[3]

	if farmCount != "" {
		farmAmount, _ = strconv.Atoi(farmCount)
	}

	if farmAmount == 1 {
		if strings.Contains(m.Content, "-s") {
			s.ChannelMessageSend(m.ChannelID, "The best farmerdog in this server is **"+cache[0].Osu.Username+"** with a farmerdog rating of "+strconv.FormatFloat(cache[0].Farm.Rating, 'f', 2, 64))
		}
		s.ChannelMessageSend(m.ChannelID, "The best farmerdog is **"+cache[0].Osu.Username+"** with a farmerdog rating of "+strconv.FormatFloat(cache[0].Farm.Rating, 'f', 2, 64))
		return
	} else if farmAmount > len(cache) {
		farmAmount = len(cache)
	}

	msg := "**Lowest farmerdogs (excluding anyone with 0.00 rating):** \n"
	if strings.Contains(m.Content, "-s") {
		msg = "**Lowest farmerdogs in this server (excluding anyone with 0.00 rating):** \n"
	}
	max := 0

	for i := 0; i < farmAmount; i++ {
		if len(msg) >= 2000 {
			max = i + 1
			break
		}

		msg = msg + "#" + strconv.Itoa(i+1) + ": **" + cache[i].Osu.Username + "** - " + strconv.FormatFloat(cache[i].Farm.Rating, 'f', 2, 64) + " farmerdog rating \n"
	}

	if len(msg) > 2000 {
		for {
			lines := strings.Split(msg, "\n")
			lines = lines[:len(lines)-1]
			msg = strings.Join(lines, "\n")
			if len(msg) <= 2000 {
				break
			}
		}
	}

	if max == 0 {
		s.ChannelMessageSend(m.ChannelID, msg)
	} else {
		s.ChannelMessageSend(m.ChannelID, "Only showing lowest "+strconv.Itoa(max)+" farmerdogs")
		s.ChannelMessageSend(m.ChannelID, msg)
	}
}
