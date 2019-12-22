package handlers

import (
	"regexp"

	structs "../structs"
	admincommands "./admin-commands"
	osucommands "./osu-commands"
	"github.com/bwmarrin/discordgo"
)

// OsuHandle handles commands that are regarding osu!
func OsuHandle(s *discordgo.Session, m *discordgo.MessageCreate, args []string, playerCache []structs.PlayerData, mapCache []structs.MapData, mapperData []structs.MapperData) {
	profileRegex, _ := regexp.Compile(`(osu|old)\.ppy\.sh\/(u|users)\/(\S+)`)
	beatmapRegex, _ := regexp.Compile(`(osu|old)\.ppy\.sh\/(s|b|beatmaps|beatmapsets)\/(\d+)(#(osu|taiko|fruits|mania)\/(\d+))?`)
	// Check if any args were even given
	if len(args) == 1 {
		go osucommands.ProfileMessage(s, m, profileRegex, playerCache)
	} else if len(args) > 1 {
		mainArg := args[1]
		switch mainArg {
		// Admin specific
		case "tr", "track":
			go admincommands.Track(s, m, mapCache)
		case "tt", "trackt", "ttoggle", "tracktoggle":
			go admincommands.TrackToggle(s, m, mapCache)

		// non-Admin specific
		case "bfarm", "bottomfarm":
			go osucommands.BottomFarm(s, m, playerCache)
		case "bpm":
			go osucommands.BPM(s, m, playerCache)
		case "c", "compare":
			go osucommands.Compare(s, m, args, playerCache, mapCache)
		case "farm":
			go osucommands.Farm(s, m, playerCache)
		case "link", "set":
			go osucommands.Link(s, m, args, playerCache)
		case "m", "map":
			go osucommands.BeatmapMessage(s, m, beatmapRegex, mapCache)
		case "mt", "mtrack", "maptrack", "mappertrack":
			go osucommands.TrackMapper(s, m, mapperData)
		case "mti", "mtinfo", "mtrackinfo", "maptracking", "mappertracking", "mappertrackinfo":
			go osucommands.TrackMapperInfo(s, m, mapperData)
		case "ppadd":
			go osucommands.PPAdd(s, m, playerCache)
		case "r", "rs", "recent":
			go osucommands.Recent(s, m, "recent", playerCache, mapCache)
		case "rb", "recentb", "recentbest":
			go osucommands.Recent(s, m, "best", playerCache, mapCache)
		case "t", "top":
			go osucommands.Top(s, m, playerCache, mapCache)
		case "tfarm", "topfarm":
			go osucommands.TopFarm(s, m, playerCache)
		case "ti", "tinfo", "tracking", "trackinfo":
			go osucommands.TrackInfo(s, m)
		default:
			go osucommands.ProfileMessage(s, m, profileRegex, playerCache)
		}
	} else {
		s.ChannelMessageSend(m.ChannelID, "Please specify a command! Check `help` for more details!")
	}
}
