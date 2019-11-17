package handlers

import (
	osuapi "../osu-api"
	structs "../structs"
	osucommands "./osu-commands"
	"github.com/bwmarrin/discordgo"
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
			go osucommands.Recent(s, m, osuAPI, "recent", playerCache, mapCache)
		case "recentb", "rb", "recentbest":
			go osucommands.Recent(s, m, osuAPI, "best", playerCache, mapCache)
		case "t", "top":
			go osucommands.Top(s, m, osuAPI, playerCache, mapCache)
		case "tfarm", "topfarm":
			go osucommands.TopFarm(s, m, osuAPI, playerCache, serverPrefix)
		case "bfarm", "bottomfarm":
			go osucommands.BottomFarm(s, m, osuAPI, playerCache, serverPrefix)
		case "farm":
			go osucommands.Farmerdog(s, m, osuAPI, playerCache)
		case "ppadd":
			go osucommands.PPAdd(s, m, osuAPI, playerCache)
		case "tr", "track":
			go osucommands.Track(s, m, osuAPI, mapCache)
		case "ti", "tinfo", "tracking", "trackinfo":
			go osucommands.TrackInfo(s, m)
		case "tt", "trackt", "ttoggle", "tracktoggle":
			go osucommands.TrackToggle(s, m, mapCache)
		case "toggle":
			go osucommands.OsuToggle(s, m)
		case "c", "compare":
			go osucommands.Compare(s, m, args, osuAPI, playerCache, serverPrefix, mapCache)
		}
	} else {
		s.ChannelMessageSend(m.ChannelID, "Please specify a command! Check `"+serverPrefix+"help` for more details!")
	}
}
