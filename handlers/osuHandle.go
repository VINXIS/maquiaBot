package handlers

import (
	osuapi "../osu-api"
	structs "../structs"
	admincommands "./admin-commands"
	osucommands "./osu-commands"
	"github.com/bwmarrin/discordgo"
)

// OsuHandle handles commands that are regarding osu!
func OsuHandle(s *discordgo.Session, m *discordgo.MessageCreate, args []string, osuAPI *osuapi.Client, playerCache []structs.PlayerData, mapCache []structs.MapData, mapperData []structs.MapperData, serverPrefix string) {
	// Check if any args were even given
	if len(args) > 1 {
		mainArg := args[1]
		switch mainArg {
		// Admin specific
		case "tr", "track":
			go admincommands.Track(s, m, osuAPI, mapCache)
		case "tt", "trackt", "ttoggle", "tracktoggle":
			go admincommands.TrackToggle(s, m, mapCache)
		case "toggle":
			go admincommands.OsuToggle(s, m)

		// non-Admin specific
		case "bfarm", "bottomfarm":
			go osucommands.BottomFarm(s, m, osuAPI, playerCache, serverPrefix)
		case "c", "compare":
			go osucommands.Compare(s, m, args, osuAPI, playerCache, serverPrefix, mapCache)
		case "farm":
			go osucommands.Farmerdog(s, m, osuAPI, playerCache)
		case "link", "set":
			go osucommands.Link(s, m, args, osuAPI, playerCache)
		case serverPrefix + "mt", serverPrefix + "mtrack", serverPrefix + "maptrack", serverPrefix + "mappertrack":
			go osucommands.TrackMapper(s, m, osuAPI, mapperData)
		case serverPrefix + "mti", serverPrefix + "mtinfo", serverPrefix + "mtrackinfo", serverPrefix + "maptracking", serverPrefix + "mappertracking", serverPrefix + "mappertrackinfo":
			go osucommands.TrackMapperInfo(s, m, mapperData)
		case "ppadd":
			go osucommands.PPAdd(s, m, osuAPI, playerCache)
		case "recent", "r", "rs":
			go osucommands.Recent(s, m, osuAPI, "recent", playerCache, mapCache)
		case "recentb", "rb", "recentbest":
			go osucommands.Recent(s, m, osuAPI, "best", playerCache, mapCache)
		case "t", "top":
			go osucommands.Top(s, m, osuAPI, playerCache, mapCache)
		case "tfarm", "topfarm":
			go osucommands.TopFarm(s, m, osuAPI, playerCache, serverPrefix)
		case "ti", "tinfo", "tracking", "trackinfo":
			go osucommands.TrackInfo(s, m)
		}
	} else {
		s.ChannelMessageSend(m.ChannelID, "Please specify a command! Check `"+serverPrefix+"help` for more details!")
	}
}
