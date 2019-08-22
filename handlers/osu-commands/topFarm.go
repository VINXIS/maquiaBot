package osucommands

import (
	"sort"
	"strconv"
	"strings"

	structs "../../structs"
	tools "../../tools"
	"github.com/bwmarrin/discordgo"
	"github.com/thehowl/go-osuapi"
)

//TopFarm gives the top farmerdogs in the game based on who's been run
func TopFarm(s *discordgo.Session, m *discordgo.MessageCreate, args []string, osuAPI *osuapi.Client, cache []structs.PlayerData, serverPrefix string) {
	if strings.Contains(m.Content, "-s") {
		guild, err := s.Guild(m.GuildID)
		tools.ErrRead(err)
		trueCache := []structs.PlayerData{}

		for _, member := range guild.Members {
			for _, player := range cache {
				if player.Discord.ID == member.User.ID {
					trueCache = append(trueCache, player)
				}
			}
		}

		cache = trueCache
	}

	sort.Slice(cache, func(i, j int) bool {
		return cache[i].Farm.Rating > cache[j].Farm.Rating
	})

	farmString := "1"

	if len(args) > 1 {
		if args[0] == serverPrefix+"osu" && len(args) > 2 {
			if strings.Contains(m.Content, "-s") && len(args) > 3 {
				farmString = args[3]
			} else if !strings.Contains(m.Content, "-s") {
				farmString = args[2]
			}
		} else {
			if strings.Contains(m.Content, "-s") && len(args) > 2 {
				farmString = args[2]
			} else if !strings.Contains(m.Content, "-s") {
				farmString = args[1]
			}
		}
	}

	farmAmount, err := strconv.Atoi(farmString)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "Please state an actual number or do not state anything at all after the command!")
		return
	}

	if farmAmount == 1 {
		if strings.Contains(m.Content, "-s") {
			s.ChannelMessageSend(m.ChannelID, "The worst farmerdog in the server is **"+cache[0].Osu.Username+"** with a farmerdog rating of "+strconv.FormatFloat(cache[0].Farm.Rating, 'f', 2, 64))
			return
		}
		s.ChannelMessageSend(m.ChannelID, "The worst farmerdog is **"+cache[0].Osu.Username+"** with a farmerdog rating of "+strconv.FormatFloat(cache[0].Farm.Rating, 'f', 2, 64))
		return
	}

	msg := "**Top farmerdogs:** \n"
	if strings.Contains(m.Content, "-s") {
		msg = "**Top farmerdogs in this server:** \n"
	}
	max := 0

	if farmAmount > len(cache)-1 {
		s.ChannelMessageSend(m.ChannelID, "Not enough players in the data set!")
		return
	}

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
