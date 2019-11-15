package osucommands

import (
	"math"
	"regexp"
	"sort"
	"strconv"
	"strings"

	osuapi "../../osu-api"
	structs "../../structs"
	tools "../../tools"
	"github.com/bwmarrin/discordgo"
)

//BottomFarm gives the top farmerdogs in the game based on who's been run
func BottomFarm(s *discordgo.Session, m *discordgo.MessageCreate, osuAPI *osuapi.Client, cache []structs.PlayerData, serverPrefix string) {
	farmCountRegex, _ := regexp.Compile(`.+(bfarm|bottomfarm)\s*(-s)?\s*(\d+)?`)

	farmAmount := 1

	if strings.Contains(m.Content, "-s") {
		guild, err := s.Guild(m.GuildID)
		tools.ErrRead(err)
		trueCache := []structs.PlayerData{}

		for _, player := range cache {
			for _, member := range guild.Members {
				if player.Discord.ID == member.User.ID {
					trueCache = append(trueCache, player)
					break
				}
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
		i := 0
		for {
			if math.Round(cache[i].Farm.Rating*100)/100 != 0.00 {
				if strings.Contains(m.Content, "-s") {
					s.ChannelMessageSend(m.ChannelID, "The best farmerdog in this server is **"+cache[i].Osu.Username+"** with a farmerdog rating of "+strconv.FormatFloat(cache[i].Farm.Rating, 'f', 2, 64))
					break
				}
				s.ChannelMessageSend(m.ChannelID, "The best farmerdog is **"+cache[i].Osu.Username+"** with a farmerdog rating of "+strconv.FormatFloat(cache[i].Farm.Rating, 'f', 2, 64))
				break
			} else {
				i++
			}
		}
		return
	} else if farmAmount > len(cache) {
		s.ChannelMessageSend(m.ChannelID, "Not enough players in the data set!")
		return
	}

	msg := "**Lowest farmerdogs (excluding anyone with 0.00 rating):** \n"
	if strings.Contains(m.Content, "-s") {
		msg = "**Lowest farmerdogs in this server (excluding anyone with 0.00 rating):** \n"
	}
	max := 0

	for i := 0; i < farmAmount; i++ {
		if math.Round(cache[i].Farm.Rating*100)/100 == 0.00 {
			farmAmount++
			continue
		}

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
