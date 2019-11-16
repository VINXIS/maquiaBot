package osucommands

import (
	"encoding/json"
	"io/ioutil"
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

// Farmerdog gives a player's farmerdog rating
func Farmerdog(s *discordgo.Session, m *discordgo.MessageCreate, osuAPI *osuapi.Client, cache []structs.PlayerData) {
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

	// Get User and see if user exists
	osuUser, err := osuAPI.GetUser(osuapi.GetUserOpts{
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

	// Calc stuff
	user.FarmCalc(osuAPI, farmData)

	// Add the new information to the full data
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

	_, err = s.ChannelMessageSend(m.ChannelID, "Farmerdog rating for **"+user.Osu.Username+":** "+strconv.FormatFloat(user.Farm.Rating, 'f', 2, 64)+"\n**Your top "+strconv.Itoa(amount)+" farmerdog scores:** \n"+list)
}
