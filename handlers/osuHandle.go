package handlers

import (
	structs "../structs"
	osucommands "./osu-commands"
	"github.com/bwmarrin/discordgo"
	"github.com/thehowl/go-osuapi"
)

// OsuHandle handles commands that are regarding osu!
func OsuHandle(s *discordgo.Session, m *discordgo.MessageCreate, args []string, osuAPI *osuapi.Client, playerCache []structs.PlayerData, mapCache []structs.MapData, serverPrefix string) {
	// Check if any args were even given
	if len(args) > 1 {
		mainArg := args[1]
		switch mainArg {
		case "link", "set":
			go osucommands.Link(s, m, args, osuAPI, playerCache)
		case "recent", "r", "rs":
			go osucommands.Recent(s, m, args, osuAPI, playerCache, "recent", serverPrefix, mapCache)
		case "recentb", "rb", "recentbest":
			go osucommands.Recent(s, m, args, osuAPI, playerCache, "best", serverPrefix, mapCache)
		case "t", "top":
			go osucommands.Top(s, m, args, osuAPI, playerCache, serverPrefix, mapCache)
		case "tfarm", "topfarm", "tfarmerdog", "topfarmerdog":
			go osucommands.TopFarm(s, m, args, osuAPI, playerCache, serverPrefix)
		case "bfarm", "bottomfarm", "bfarmerdog", "bottomfarmerdog":
			go osucommands.BottomFarm(s, m, args, osuAPI, playerCache, serverPrefix)
		case "farm", "farmerdog", "f":
			go osucommands.Farmerdog(s, m, args, osuAPI, playerCache, serverPrefix)
		case "tr", "track":
			go osucommands.Track(s, m, args, osuAPI, mapCache)
		case "ti", "tinfo", "tracking", "trackinfo":
			go osucommands.TrackInfo(s, m)
		case "tt", "trackt", "tracktoggle":
			go osucommands.TrackToggle(s, m, mapCache)
		case "c", "compare":
			go osucommands.Compare(s, m, args, osuAPI, playerCache, serverPrefix, mapCache)
		}
	} else {
		s.ChannelMessageSend(m.ChannelID, "Please specify a command! Check `"+serverPrefix+"help` for more details!")
	}
}
