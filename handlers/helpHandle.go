package handlers

import (
	"math/rand"
	"regexp"
	"strings"

	osuapi "../osu-api"
	osutools "../osu-tools"
	helpcommands "./help-commands"
	"github.com/bwmarrin/discordgo"
)

// HelpHandle lets you know the commands available
func HelpHandle(s *discordgo.Session, m *discordgo.MessageCreate, prefix string) {
	embed := &discordgo.MessageEmbed{
		Author: &discordgo.MessageEmbedAuthor{
			URL:     "https://discordapp.com/oauth2/authorize?&client_id=" + s.State.User.ID + "&scope=bot&permissions=0",
			Name:    "Click here to invite MaquiaBot!",
			IconURL: s.State.User.AvatarURL("2048"),
		},
		Description: "**Most commands have other forms as well for convenience!**" + "\n\n" +
			"**Please do `" + prefix + "help <command>` for more information about the command!** \n" +
			"Help information format: `(cmd|names) <args> [optional args]`",
		Fields: []*discordgo.MessageEmbedField{
			&discordgo.MessageEmbedField{
				Name: "Admin commands:",
				Value: "`maquiaprefix`, " +
					"`purge`, " +
					"`roleautomation`, " +
					"`toggle`, " +
					"`track`, " +
					"`trigger`, " +
					"`tracktoggle`",
			},
			&discordgo.MessageEmbedField{
				Name: "General commands:",
				Value: "`adjectives`, " +
					"`avatar`, " +
					"`caps`, " +
					"`comparepenis`, " +
					"`comparevagina`, " +
					"`choose`, " +
					"`cheers`, " +
					"`colour`, " +
					"`crab`, " +
					"`decrypt`, " +
					"`emote`, " +
					"`encrypt`, " +
					"`face`, " +
					"`history`, " +
					"`idea`, " +
					"`info`, " +
					"`kanye`, " +
					"`late`, " +
					"`levenshtein`, " +
					"`list`, " +
					"`lower`, " +
					"`meme`, " +
					"`nouns`, " +
					"`ocr`, " +
					"`over`, " +
					"`percentage`, " +
					"`parse`, " +
					"`penis`, " +
					"`ping`, " +
					"`quote`, " +
					"`quoteadd`, " +
					"`quotedelete`, " +
					"`quotes`, " +
					"`randomcaps`, " +
					"`rankpenis`, " +
					"`rankvagina`, " +
					"`reminder`, " +
					"`reminders`, " +
					"`remindremove`, " +
					"`roleinfo`, " +
					"`roll`, " +
					"`serverinfo`, " +
					"`skills`, " +
					"`stats`, " +
					"`title`, " +
					"`triggers`, " +
					"`twitchdl`, " +
					"`twitterdl`, " +
					"`vagina`, " +
					"`vibecheck`, " +
					"`weather`",
			},
			&discordgo.MessageEmbedField{
				Name: "math commands:",
				Value: "`average`, " +
					"`distance`, " +
					"`degreesradians`, " +
					"`radiansdegrees`, " +
					"`standarddeviation`, " +
					"`vectoradd`, " +
					"`vectorcross`, " +
					"`vectordivide`, " +
					"`vectordot`, " +
					"`vectormultiply`, " +
					"`vectorsubtract`",
			},
			&discordgo.MessageEmbedField{
				Name: "osu! commands:",
				Value: "`bottomfarm`, " +
					"`bpm`, " +
					"`compare`, " +
					"`farm`, " +
					"`leaderboard`, " +
					"`link`, " +
					"`map`, " +
					"`osu`, " +
					"`osudetail`, " +
					"`osutop`, " +
					"`ppadd`, " +
					"`recent`, " +
					"`recentbest`, " +
					"`scorepost`, " +
					"`top`, " +
					"`topfarm`, " +
					"`trackinfo`",
			},
			&discordgo.MessageEmbedField{
				Name: "Pokemon commands:",
				Value: "`berry`, " +
					"`pokemon`",
			},
		},
		Color: osutools.ModeColour(osuapi.ModeOsu),
	}

	argRegex, _ := regexp.Compile(`help\s+(.+)`)
	if argRegex.MatchString(m.Content) {
		arg := argRegex.FindStringSubmatch(m.Content)[1]
		args := strings.Split(arg, " ")
		if (args[0] == "pokemon" || args[0] == "osu") && len(args) > 1 {
			arg = args[1]
		}
		switch arg {
		// Admin commands
		case "prefix", "maquiaprefix", "newprefix":
			embed = helpcommands.Prefix(embed)
		case "purge":
			embed = helpcommands.Purge(embed)
		case "rolea", "roleauto", "roleautomation":
			embed = helpcommands.RoleAutomation(embed)
		case "toggle":
			embed = helpcommands.Toggle(embed)
		case "tr", "track":
			embed = helpcommands.Track(embed)
		case "trigger":
			embed = helpcommands.Trigger(embed)
		case "tt", "trackt", "ttoggle", "tracktoggle":
			embed = helpcommands.TrackToggle(embed)

		// General commands
		case "adj", "adjective", "adjectives":
			embed = helpcommands.Adjectives(embed)
		case "avatar", "ava", "a":
			embed = helpcommands.Avatar(embed)
		case "cap", "caps", "upper":
			embed = helpcommands.AllCaps(embed)
		case "cp", "comparep", "comparepenis":
			embed = helpcommands.PenisCompare(embed)
		case "cv", "comparev", "comparevagina":
			embed = helpcommands.VaginaCompare(embed)
		case "ch", "choose":
			embed = helpcommands.Choose(embed)
		case "cheers":
			embed = helpcommands.Cheers(embed)
		case "col", "color", "colour":
			embed = helpcommands.Colour(embed)
		case "crab":
			embed = helpcommands.Crab(embed)
		case "decrypt":
			embed = helpcommands.Decrypt(embed)
		case "e", "emoji", "emote":
			embed = helpcommands.Emoji(embed)
		case "encrypt":
			embed = helpcommands.Encrypt(embed)
		case "face":
			embed = helpcommands.Face(embed)
		case "history":
			embed = helpcommands.History(embed)
		case "idea", "niceidea":
			embed = helpcommands.NiceIdea(embed)
		case "info":
			embed = helpcommands.Info(embed)
		case "kanye":
			embed = helpcommands.Kanye(embed)
		case "late", "old", "ancient":
			embed = helpcommands.Late(embed)
		case "leven", "levenshtein":
			embed = helpcommands.Levenshtein(embed)
		case "list":
			embed = helpcommands.List(embed)
		case "lower":
			embed = helpcommands.AllLower(embed)
		case "meme":
			embed = helpcommands.Meme(embed)
		case "noun", "nouns":
			embed = helpcommands.Nouns(embed)
		case "ocr":
			embed = helpcommands.OCR(embed)
		case "over":
			embed = helpcommands.OverIt(embed)
		case "p", "per", "percent", "percentage":
			embed = helpcommands.Percentage(embed)
		case "parse":
			embed = helpcommands.Parse(embed)
		case "penis":
			embed = helpcommands.Penis(embed)
		case "ping":
			embed = helpcommands.Ping(embed)
		case "q", "quote":
			embed = helpcommands.Quote(embed)
		case "qa", "qadd", "quotea", "quoteadd":
			embed = helpcommands.QuoteAdd(embed)
		case "qd", "qr", "qremove", "qdelete", "quoter", "quoted", "quoteremove", "quotedelete":
			embed = helpcommands.QuoteRemove(embed)
		case "qs", "quotes":
			embed = helpcommands.Quotes(embed)
		case "rcap", "rcaps", "rupper", "rlower", "randomcap", "randomcaps", "randomupper", "randomlower":
			embed = helpcommands.RandomCaps(embed)
		case "rp", "rankp", "rankpenis":
			embed = helpcommands.PenisRank(embed)
		case "rv", "rankv", "rankvagina":
			embed = helpcommands.VaginaRank(embed)
		case "remind", "reminder":
			embed = helpcommands.Remind(embed)
		case "reminders":
			embed = helpcommands.Reminders(embed)
		case "remindremove", "rremove":
			embed = helpcommands.RemindRemove(embed)
		case "rinfo", "roleinfo":
			embed = helpcommands.RoleInfo(embed)
		case "roll":
			embed = helpcommands.Roll(embed)
		case "sinfo", "serverinfo":
			embed = helpcommands.ServerInfo(embed)
		case "skill", "skills":
			embed = helpcommands.Skills(embed)
		case "stats", "class":
			embed = helpcommands.Stats(embed)
		case "title":
			embed = helpcommands.Title(embed)
		case "triggers":
			embed = helpcommands.Triggers(embed)
		case "twitch", "twitchdl":
			embed = helpcommands.Twitch(embed)
		case "twitter", "twitterdl":
			embed = helpcommands.Twitter(embed)
		case "vagina":
			embed = helpcommands.Vagina(embed)
		case "vibe", "vibec", "vibecheck":
			embed = helpcommands.Vibe(embed)
		case "w", "weather":
			embed = helpcommands.Weather(embed)

		// Math commands
		case "ave", "average", "mean":
			embed = helpcommands.Average(embed)
		case "d", "dist", "distance", "dir", "direction":
			embed = helpcommands.DistanceDirection(embed)
		case "dr", "degrad", "degreesradians":
			embed = helpcommands.DegreesRadians(embed)
		case "rd", "raddeg", "radiansdegrees":
			embed = helpcommands.RadiansDegrees(embed)
		case "stddev", "standarddev", "stddeviation", "standarddeviation":
			embed = helpcommands.StandardDeviation(embed)
		case "va", "vadd", "vectora", "vectoradd":
			embed = helpcommands.VectorAdd(embed)
		case "vc", "vcross", "vectorc", "vectorcross":
			embed = helpcommands.VectorCross(embed)
		case "vd", "vdiv", "vdivide", "vectord", "vectordiv", "vectordivide":
			embed = helpcommands.VectorDivide(embed)
		case "vdot", "vectordot":
			embed = helpcommands.VectorDot(embed)
		case "vm", "vmult", "vmultiply", "vectorm", "vectormult", "vectormultiply":
			embed = helpcommands.VectorMultiply(embed)
		case "vs", "vsub", "vsubtract", "vectors", "vectorsub", "vectorsubtract":
			embed = helpcommands.VectorSubtract(embed)

		// osu! commands
		case "bfarm", "bottomfarm":
			embed = helpcommands.BottomFarm(embed)
		case "bpm":
			embed = helpcommands.BPM(embed)
		case "c", "compare":
			embed = helpcommands.Compare(embed)
		case "farm":
			embed = helpcommands.Farm(embed)
		case "l", "leader", "leaderboard":
			embed = helpcommands.Leaderboard(embed)
		case "link", "set":
			embed = helpcommands.Link(embed)
		case "m", "map":
			embed = helpcommands.Map(embed)
		case "osu", "profile":
			embed = helpcommands.Profile(embed)
		case "osudetail":
			embed = helpcommands.ProfileDetail(embed)
		case "osutop":
			embed = helpcommands.ProfileTop(embed)
		case "ppadd":
			embed = helpcommands.PPAdd(embed)
		case "r", "rs", "recent":
			embed = helpcommands.Recent(embed)
		case "rb", "recentb", "recentbest":
			embed = helpcommands.RecentBest(embed)
		case "s", "sc", "scorepost":
			embed = helpcommands.ScorePost(embed)
		case "t", "top":
			embed = helpcommands.Top(embed)
		case "tfarm", "topfarm":
			embed = helpcommands.TopFarm(embed)
		case "ti", "tinfo", "tracking", "trackinfo":
			embed = helpcommands.TrackInfo(embed)

		// Pokemon commands
		case "b", "berry":
			embed = helpcommands.Berry(embed)
		case "pokemon":
			embed = helpcommands.Pokemon(embed)
		}
	}

	if !strings.HasPrefix(embed.Description, "**Most") && embed.Fields[0].Name == "Admin commands:" {
		embed.Fields = []*discordgo.MessageEmbedField{}
	}

	switch rand.Intn(11) {
	case 0:
		embed.Image = &discordgo.MessageEmbedImage{
			URL: "https://cdn.discordapp.com/attachments/555493588465877012/555994312760885248/epicAnimeScene.gif",
		}
	case 1:
		embed.Image = &discordgo.MessageEmbedImage{
			URL: "https://cdn.discordapp.com/attachments/555493588465877012/555996915884490752/epicAnimeGifTWO.gif",
		}
	case 2:
		embed.Image = &discordgo.MessageEmbedImage{
			URL: "https://cdn.discordapp.com/attachments/555493588465877012/556000178406948875/epicAnimeGif5.gif",
		}
	case 3, 4:
		embed.Image = &discordgo.MessageEmbedImage{
			URL: "https://cdn.discordapp.com/attachments/555493588465877012/555998627940532237/tumblr_phjkel3lgn1xlyyvto1_1280.png",
		}
	case 5, 6:
		embed.Image = &discordgo.MessageEmbedImage{
			URL: "https://cdn.discordapp.com/attachments/555493588465877012/555998640142024704/tumblr_phjkel3lgn1xlyyvto2_1280.png",
		}
	case 7, 8:
		embed.Image = &discordgo.MessageEmbedImage{
			URL: "https://cdn.discordapp.com/attachments/555493588465877012/555998669761937418/tumblr_phjkel3lgn1xlyyvto3_1280.png",
		}
	case 9, 10:
		embed.Image = &discordgo.MessageEmbedImage{
			URL: "https://cdn.discordapp.com/attachments/555493588465877012/555998681375965194/tumblr_phjkel3lgn1xlyyvto5_1280.png",
		}
	}
	s.ChannelMessageSendComplex(m.ChannelID, &discordgo.MessageSend{
		Content: "All commands in PM will use the bot's default prefix `$` instead!",
		Embed:   embed,
	})
}
