package gencommands

import (
	"log"
	"maquiaBot/structs"
	"time"

	"github.com/bwmarrin/discordgo"
)

// Roulette
func Roulette(s *discordgo.Session, m *discordgo.MessageCreate, cache []structs.PlayerData) {
	player := structs.PlayerData{
		Time:    time.Now(),
		Discord: m.Author.ID,
	}
	exists := -1

	// Run through the player cache to find the user using discord ID.
	for i, cachePlayer := range cache {
		if cachePlayer.Discord == m.Author.ID {
			player = cachePlayer
			exists = i
			break
		}
	}

	log.Println(player)
	log.Println(exists)

}
