package gencommands

import (
	"encoding/json"
	"io/ioutil"
	"maquiaBot/structs"
	"maquiaBot/tools"
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

// Daily lets you obtain 10 + 1% of current points everyday
func Daily(s *discordgo.Session, m *discordgo.MessageCreate) {
	player := structs.PlayerData{
		Time:    time.Now(),
		Discord: m.Author.ID,
	}
	exists := -1

	// Obtain profile cache data
	var cache []structs.PlayerData
	f, err := ioutil.ReadFile("./data/osuData/profileCache.json")
	tools.ErrRead(s, err)
	_ = json.Unmarshal(f, &cache)

	// Run through the player cache to find the user using discord ID.
	for i, cachePlayer := range cache {
		if cachePlayer.Discord == m.Author.ID {
			player = cachePlayer
			exists = i
			break
		}
	}

	year, month, day := player.Currency.LastDaily.Date()
	if exists != -1 && year == time.Now().UTC().Year() && month == time.Now().UTC().Month() && day == time.Now().UTC().Day() {
		s.ChannelMessageSend(m.ChannelID, "You have already received your daily amount of points today!")
		return
	}

	newAmount := math.Max(10, 0.01*player.Currency.Amount)
	if player.Currency.Amount > 0 {
		switch {
		case int64(player.Currency.Amount)%445 == 0:
			newAmount = 1.1 * player.Currency.Amount
		case int64(player.Currency.Amount)%727 == 0:
			newAmount = 1.01 * player.Currency.Amount
		case int64(player.Currency.Amount)%2 == 1:
			newAmount = -0.01 * player.Currency.Amount
		}
	}
	player.Currency.Amount += newAmount
	player.Currency.LastDaily = time.Now()
	player.Time = time.Now()

	// Save player
	if exists != -1 {
		cache[exists] = player
	} else {
		cache = append(cache, player)
	}
	jsonCache, err := json.Marshal(cache)
	tools.ErrRead(s, err)

	err = ioutil.WriteFile("./data/osuData/profileCache.json", jsonCache, 0644)
	tools.ErrRead(s, err)
	s.ChannelMessageSend(m.ChannelID, "You have now obtained **"+strconv.FormatFloat(newAmount, 'f', 2, 64)+"** points for today! Your new balance is **"+strconv.FormatFloat(player.Currency.Amount, 'f', 2, 64)+"** points.")
}

// Transfer lets users transfer points between the author and a target user
func Transfer(s *discordgo.Session, m *discordgo.MessageCreate) {
	var err error

	// Mention check
	if len(m.Mentions) == 0 {
		s.ChannelMessageSend(m.ChannelID, "Please mention someone to transfer points to!")
		return
	}
	if m.Mentions[0].ID == m.Author.ID {
		s.ChannelMessageSend(m.ChannelID, "Ur really funny mate")
		return
	}
	target := m.Mentions[0]

	// Get number of points to transfer
	val := 0.00
	if val, err = strconv.ParseFloat(strings.Split(m.Content, " ")[len(strings.Split(m.Content, " "))-1], 64); err != nil || val < 0 || math.IsInf(val, 1) || math.IsInf(val, -1) || math.IsNaN(val) {
		s.ChannelMessageSend(m.ChannelID, "Please provide a number at the end of your message!")
		return
	}

	// Obtain profile cache data
	var cache []structs.PlayerData
	f, err := ioutil.ReadFile("./data/osuData/profileCache.json")
	tools.ErrRead(s, err)
	_ = json.Unmarshal(f, &cache)

	// Find message author user and target user in player cache
	playerAuthor := structs.PlayerData{
		Time:     time.Now(),
		Discord:  m.Author.ID,
		Currency: structs.CurrencyData{10, time.Now()},
	}
	authorExists := -1

	playerTarget := structs.PlayerData{
		Time:     time.Now(),
		Discord:  target.ID,
		Currency: structs.CurrencyData{10, time.Now()},
	}
	targetExists := -1
	for i, cachePlayer := range cache {
		if cachePlayer.Discord == m.Author.ID {
			playerAuthor = cachePlayer
			authorExists = i
			if val > playerAuthor.Currency.Amount {
				s.ChannelMessageSend(m.ChannelID, "Basic math dictates your current balance of "+strconv.FormatFloat(playerAuthor.Currency.Amount, 'f', 2, 64)+" points is less than fucking "+strconv.FormatFloat(val, 'f', 2, 64)+" dude.")
				return
			}
		} else if cachePlayer.Discord == target.ID {
			playerTarget = cachePlayer
			targetExists = i
		}

		if authorExists != -1 && targetExists != -1 {
			break
		}
	}

	// Transfer credits
	playerAuthor.Currency.Amount -= val
	playerTarget.Currency.Amount += val

	// Save author
	if authorExists != -1 {
		cache[authorExists] = playerAuthor
	} else {
		cache = append(cache, playerAuthor)
	}

	// Save target
	if targetExists != -1 {
		cache[targetExists] = playerTarget
	} else {
		cache = append(cache, playerTarget)
	}

	jsonCache, err := json.Marshal(cache)
	tools.ErrRead(s, err)

	err = ioutil.WriteFile("./data/osuData/profileCache.json", jsonCache, 0644)
	tools.ErrRead(s, err)

	s.ChannelMessageSend(m.ChannelID, "Transferred **"+strconv.FormatFloat(val, 'f', 2, 64)+"** points to "+target.Username+"\nYour new balance: **"+strconv.FormatFloat(playerAuthor.Currency.Amount, 'f', 2, 64)+"**\n"+target.Username+"'s new balance: **"+strconv.FormatFloat(playerTarget.Currency.Amount, 'f', 2, 64)+"**")
}
