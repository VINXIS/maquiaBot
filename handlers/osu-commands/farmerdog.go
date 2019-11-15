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

	// Check
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
	user.Farm = structs.FarmerdogData{}

	// Get best scores
	scoreList, err := osuAPI.GetUserBest(osuapi.GetUserScoresOpts{
		Username: username,
		Limit:    100,
	})
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "User: **"+user.Osu.Username+"** may not exist!")
		return
	}

	// Obtain farm data
	farmData := structs.FarmData{}
	f, err := ioutil.ReadFile("./data/osuData/mapFarm.json")
	tools.ErrRead(err)
	_ = json.Unmarshal(f, &farmData)

	// Get the farmerdog rating
	for j, score := range scoreList {
		var HDVer osuapi.Mods
		var playerFarmScore = structs.PlayerScore{}

		if strings.Contains(score.Mods.String(), "NC") {
			stringMods := strings.Replace(score.Mods.String(), "NC", "", 1)
			score.Mods = osuapi.ParseMods(stringMods)
		}
		if strings.Contains(score.Mods.String(), "HD") {
			HDVer = score.Mods
			stringMods := strings.Replace(score.Mods.String(), "HD", "", 1)
			score.Mods = osuapi.ParseMods(stringMods)
		} else {
			stringMods := score.Mods.String() + "HD"
			HDVer = osuapi.ParseMods(stringMods)
		}
		for _, farmMap := range farmData.Maps {
			if score.BeatmapID == farmMap.BeatmapID && (score.Mods == farmMap.Mods || HDVer == farmMap.Mods) {
				playerFarmScore.BeatmapSet = score.BeatmapID
				playerFarmScore.PP = score.PP
				playerFarmScore.FarmScore = math.Max(playerFarmScore.FarmScore, math.Pow(0.95, float64(j))*farmMap.Overweightness)
				playerFarmScore.Name = farmMap.Artist + " - " + farmMap.Title + " [" + farmMap.DiffName + "]"
			}
		}
		user.Farm.List = append(user.Farm.List, playerFarmScore)
		user.Farm.Rating += playerFarmScore.FarmScore
	}

	cacheUser, err := osuAPI.GetUser(osuapi.GetUserOpts{
		Username: username,
	})
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "User: **"+user.Osu.Username+"** may not exist!")
		return
	}

	user.Time = time.Now()
	user.Osu = *cacheUser

	if cached {
		cache[playerIndex] = user
	} else {
		cache = append(cache, user)
	}

	jsonCache, err := json.Marshal(cache)
	tools.ErrRead(err)

	err = ioutil.WriteFile("./data/osuData/profileCache.json", jsonCache, 0644)
	tools.ErrRead(err)

	sort.Slice(user.Farm.List, func(i, j int) bool { return user.Farm.List[i].FarmScore > user.Farm.List[j].FarmScore })
	list := ""
	for n := 0; n < amount; n++ {
		list += "`" + user.Farm.List[n].Name + "`: " + strconv.FormatFloat(user.Farm.List[n].FarmScore, 'f', 2, 64) + " Farmerdog rating (" + strconv.FormatFloat(user.Farm.List[n].PP, 'f', 2, 64) + ")\n"
	}

	_, err = s.ChannelMessageSend(m.ChannelID, "Farmerdog rating for **"+user.Osu.Username+":** "+strconv.FormatFloat(user.Farm.Rating, 'f', 2, 64)+"\n**Your top "+strconv.Itoa(amount)+" farmerdog scores:** \n"+list)
}
