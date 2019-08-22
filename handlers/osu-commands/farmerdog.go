package osucommands

import (
	"encoding/json"
	"io/ioutil"
	"math"
	"sort"
	"strconv"
	"strings"
	"time"

	structs "../../structs"
	tools "../../tools"
	"github.com/bwmarrin/discordgo"
	"github.com/thehowl/go-osuapi"
)

// Farmerdog gives a player's farmerdog rating
func Farmerdog(s *discordgo.Session, m *discordgo.MessageCreate, args []string, osuAPI *osuapi.Client, cache []structs.PlayerData, serverPrefix string) {
	username := ""
	user := structs.PlayerData{}
	FarmRating := 0.00
	cached := false
	playerIndex := 0

	// Obtain user from args
	if len(args) > 1 {
		if args[0] == serverPrefix+"osu" && len(args) > 2 {
			username = args[2]
		} else {
			username = args[1]
		}
	} else {
		for i, player := range cache {
			if m.Author.ID == player.Discord.ID && player.Osu.Username != user.Osu.Username {
				if time.Since(player.Farm.Time) < time.Hour {
					s.ChannelMessageSend(m.ChannelID, "Farmerdog rating for **"+player.Osu.Username+":** "+strconv.FormatFloat(player.Farm.Rating, 'f', 2, 64)+"\n**Your top 5 farmerdog scores:** \n"+player.Farm.List)
					return
				}
				playerIndex = i
				cached = true
				user = player
				break
			}
		}
	}

	// Check
	if username != "" {
		for i, player := range cache {
			if username == strings.ToLower(player.Osu.Username) {
				playerIndex = i
				cached = true
				user = player
				break
			}
		}
		if user.Osu.Username == "" {
			user.Osu.Username = username
		}
	}

	if user.Osu.Username == "" && username == "" {
		s.ChannelMessageSend(m.ChannelID, "No user found!")
		return
	}

	// Get best scores
	scoreList, err := osuAPI.GetUserBest(osuapi.GetUserScoresOpts{
		Username: user.Osu.Username,
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
	list := ""
	FarmValues := make(map[float64]string)

	for i, score := range scoreList {
		var HDVer osuapi.Mods
		var weight float64
		var artist string
		var title string
		var diffName string

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
				weight = math.Max(weight, math.Pow(0.95, float64(i))*farmMap.Overweightness)
				artist = farmMap.Artist
				title = farmMap.Title
				diffName = farmMap.DiffName
			}
		}
		if weight != 0 {
			FarmRating = FarmRating + weight
			FarmValues[weight] = "`" + artist + " - " + title + " [" + diffName + "]`: " + strconv.FormatFloat(weight, 'f', 2, 64) + " Farmerdog rating (" + strconv.FormatFloat(score.PP, 'f', 2, 64) + "pp) \n"
		}
	}

	keys := []float64{}
	for FarmVal := range FarmValues {
		keys = append(keys, FarmVal)
	}
	sort.Slice(keys, func(i, j int) bool {
		return keys[i] > keys[j]
	})

	for _, FarmVal := range keys {
		lines := strings.Split(list, "\n")
		if len(lines) > 5 {
			break
		}

		list = list + FarmValues[FarmVal]
	}

	cacheUser, err := osuAPI.GetUser(osuapi.GetUserOpts{
		Username: user.Osu.Username,
	})
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "User: **"+user.Osu.Username+"** may not exist!")
		return
	}

	user.Time = time.Now()
	user.Osu = *cacheUser
	user.Farm.Rating = FarmRating
	user.Farm.Time = time.Now()
	user.Farm.List = list

	if cached {
		cache[playerIndex] = user
	} else {
		cache = append(cache, user)
	}

	jsonCache, err := json.Marshal(cache)
	tools.ErrRead(err)

	err = ioutil.WriteFile("./data/osuData/profileCache.json", jsonCache, 0644)
	tools.ErrRead(err)

	s.ChannelMessageSend(m.ChannelID, "Farmerdog rating for **"+user.Osu.Username+":** "+strconv.FormatFloat(FarmRating, 'f', 2, 64)+"\n**Your top 5 farmerdog scores:** \n"+list)
}
