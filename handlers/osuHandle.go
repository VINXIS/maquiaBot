package handlers

import (
	structs "../structs"
	"github.com/bwmarrin/discordgo"
	"github.com/thehowl/go-osuapi"
)

// OsuHandle handles commands that are regarding osu!
func OsuHandle(s *discordgo.Session, m *discordgo.MessageCreate, args []string, osuAPI *osuapi.Client, cache []structs.PlayerData) {

}
