package handlers

import (
	"regexp"

	admincommands "maquiaBot/handlers/admin-commands"
	osucommands "maquiaBot/handlers/osu-commands"

	"github.com/bwmarrin/discordgo"
)

// OsuHandle handles commands that are regarding osu!
func OsuHandle(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
	profileRegex, _ := regexp.Compile(`(?i)(osu|old)\.ppy\.sh\/(u|users)\/(\S+)`)
	beatmapRegex, _ := regexp.Compile(`(?i)(osu|old)\.ppy\.sh\/(s|b|beatmaps|beatmapsets)\/(\d+)(#(osu|taiko|fruits|mania)\/(\d+))?`)
	// Check if any args were even given
	if len(args) == 1 {
		go osucommands.ProfileMessage(s, m, profileRegex)
	} else if len(args) > 1 {
		mainArg := args[1]
		switch mainArg {
		// Admin specific
		case "tr", "track":
			go admincommands.Track(s, m)
		case "tt", "trackt", "ttoggle", "tracktoggle":
			go admincommands.TrackToggle(s, m)

		// non-Admin specific
		case "bfarm", "bottomfarm":
			go osucommands.BottomFarm(s, m)
		case "bpm":
			go osucommands.BPM(s, m)
		case "c", "compare":
			go osucommands.Compare(s, m)
		case "farm":
			go osucommands.Farm(s, m)
		case "l", "leader", "leaderboard":
			go osucommands.Leaderboard(s, m, beatmapRegex)
		case "link", "set":
			go osucommands.Link(s, m, args)
		case "m", "map":
			go osucommands.BeatmapMessage(s, m, "", beatmapRegex)
		case "ppadd", "addpp":
			go osucommands.PPAdd(s, m)
		case "r", "rs", "recent":
			go osucommands.Recent(s, m, "recent")
		case "rb", "recentb", "recentbest":
			go osucommands.Recent(s, m, "best")
		case "s", "sc", "scorepost":
			go osucommands.ScorePost(s, m, "scorePost", "")
		case "similar":
			go osucommands.Similar(s, m)
		case "t", "top":
			go osucommands.Top(s, m)
		case "tfarm", "topfarm":
			go osucommands.TopFarm(s, m)
		case "ti", "tinfo", "tracking", "trackinfo":
			go osucommands.TrackInfo(s, m)
		default:
			go osucommands.ProfileMessage(s, m, profileRegex)
		}
	} else {
		s.ChannelMessageSend(m.ChannelID, "Please specify a command! Check `help` for more details!")
	}
}
