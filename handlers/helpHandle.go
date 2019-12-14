package handlers

import (
	"math/rand"
	"regexp"
	"strings"

	osuapi "../osu-api"
	osutools "../osu-functions"
	helpcommands "./help-commands"
	"github.com/bwmarrin/discordgo"
)

// HelpHandle lets you know the commands available
func HelpHandle(s *discordgo.Session, m *discordgo.MessageCreate, prefix string) {
	embed := &discordgo.MessageEmbed{
		Author: &discordgo.MessageEmbedAuthor{
			URL:     "https://discordapp.com/oauth2/authorize?&client_id=551667572723023893&scope=bot&permissions=0",
			Name:    "Click here to invite MaquiaBot!",
			IconURL: s.State.User.AvatarURL("2048"),
		},
		Description: "Detailed version of the commands list [here](https://docs.google.com/spreadsheets/d/12VzMXGoxliSVv6Rrr6tEy_-Qe9oJ0TNF4MoPGcxIpcU/edit?usp=sharing). **Most commands have other forms as well for convenience!**" + "\n\n" +
			"**Please do `" + prefix + "help <command>` for more information about the command!** \n" +
			"Help information format: `(cmd|names) <args> [optional args]`",
		Fields: []*discordgo.MessageEmbedField{
			&discordgo.MessageEmbedField{
				Name: "Admin commands:",
				Value: "`(ct|crabt|ctoggle|crabtoggle)`, " +
					"`(it|ideat|itoggle|ideatoggle)`, " +
					"`(ot|osut|otoggle|osutoggle)`, " +
					"`(prefix|maquiaprefix|newprefix)`, " +
					"`purge`, " +
					"`(st|statst|stoggle|statstoggle)`, " +
					"`(tr|track)`, " +
					"`(tt|trackt|ttoggle|tracktoggle)`, " +
					"`(vt|vibet|vtoggle|vibetoggle)`",
			},
			&discordgo.MessageEmbedField{
				Name: "General commands:",
				Value: "`(adj|adjective|adjectives)`, " +
					"`(a|ava|avatar)`, " +
					"`(cc|cp|comparec|comparep|comparecock|comparepenis)`, " +
					"`(ch|choose)`, " +
					"`crab`, " +
					"`decrypt`, " +
					"`(e|emoji|emote)`, " +
					"`encrypt`, " +
					"`face`, " +
					"`funny`, " +
					"`(idea|niceidea)`, " + 
					"`info`, " +
					"`kanye`, " +
					"`(l|leven|levenshtein)`, " +
					"`(noun|nouns)`, " +
					"`ocr`, " +
					"`(p|per|percent|percentage)`, " +
					"`parse`, " +
					"`(penis|cock)`, " +
					"`ping`, " +
					"`(q|quote)`, " +
					"`(qa|qadd|quotea|quoteadd)`, " +
					"`(qd|qr|qremove|qdelete|quoter|quoted|quoteremove|quotedelete)`, " +
					"`(qs|quotes)`, " +
					"`(rc|rp|rankc|rankp|rankcock|rankpenis)`, " +
					"`(remind|reminder)`, " +
					"`reminders`, " +
					"`(remindremove|rremove)`, " +
					"`(rinfo|roleinfo)`, " +
					"`roll`, " +
					"`(sinfo|serverinfo)`, " +
					"`(skill|skills)`, " +
					"`(stats|class)`, " +
					"`(twitter|twitterdl)`, " +
					"`(vibe|vibec|vibecheck)`",
			},
			&discordgo.MessageEmbedField{
				Name: "osu! commands:",
				Value: "`(bfarm|bottomfarm)`, " +
					"`bpm`, " +
					"`farm`, " +
					"`(link|set)`, " +
					"`(tfarm|topfarm)`, " +
					"`(ti|tinfo|tracking|trackinfo)`, ",
			},
		},
		Color: osutools.ModeColour(osuapi.ModeOsu),
	}

	argRegex, _ := regexp.Compile(`help\s+(.+)`)
	if argRegex.MatchString(m.Content) {
		arg := argRegex.FindStringSubmatch(m.Content)[1]
		args := strings.Split(arg, " ")
		if (args[0] == "pokemon" || args[0] == "osu") && len(args) > 1 {
			args := strings.Split(argRegex.FindStringSubmatch(m.Content)[1], " ")
			if len(args) > 1 {
				arg = args[1]
			}
		}
		switch arg {
		// Admin commands
		case "ct", "crabt", "ctoggle", "crabtoggle":
			embed = helpcommands.CrabToggle(embed)
		case "it", "ideat", "itoggle", "ideatoggle":
			embed = helpcommands.NiceIdeaToggle(embed)
		case "ot", "osut", "otoggle", "osutoggle":
			embed = helpcommands.OsuToggle(embed)
		case "prefix", "maquiaprefix", "newprefix":
			embed = helpcommands.Prefix(embed)
		case "purge":
			embed = helpcommands.Purge(embed)
		case "st", "statst", "stoggle", "statstoggle":
			embed = helpcommands.StatsToggle(embed)
		case "tr", "track":
			embed = helpcommands.Track(embed)
		case "tt", "trackt", "ttoggle", "tracktoggle":
			embed = helpcommands.TrackToggle(embed)
		case "vt", "vibet", "vtoggle", "vibetoggle":
			embed = helpcommands.VibeToggle(embed)

		// General commands
		case "adj", "adjective", "adjectives":
			embed = helpcommands.Adjectives(embed)
		case "avatar", "ava", "a":
			embed = helpcommands.Avatar(embed)
		case "cc", "cp", "comparec", "comparep", "comparecock", "comparepenis":
			embed = helpcommands.PenisCompare(embed)
		case "ch", "choose":
			embed = helpcommands.Choose(embed)
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
		case "funny":
			embed = helpcommands.Funny(embed)
		case "idea", "niceidea":
			embed = helpcommands.NiceIdea(embed)
		case "info":
			embed = helpcommands.Info(embed)
		case "kanye":
			embed = helpcommands.Kanye(embed)
		case "l", "leven", "levenshtein":
			embed = helpcommands.Levenshtein(embed)
		case "noun", "nouns":
			embed = helpcommands.Nouns(embed)
		case "ocr":
			embed = helpcommands.OCR(embed)
		case "p", "per", "percent", "percentage":
			embed = helpcommands.Percentage(embed)
		case "parse":
			embed = helpcommands.Parse(embed)
		case "penis", "cock":
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
		case "rc", "rp", "rankc", "rankp", "rankcock", "rankpenis":
			embed = helpcommands.PenisRank(embed)
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
		case "twitter", "twitterdl":
			embed = helpcommands.Twitter(embed)
		case "vibe", "vibec", "vibecheck":
			embed = helpcommands.Vibe(embed)

		// osu! commands
		case "bfarm", "bottomfarm":
			embed = helpcommands.BottomFarm(embed)
		case "bpm":
			embed = helpcommands.BPM(embed)
		case "farm":
			embed = helpcommands.Farm(embed)
		case "link", "set":
			embed = helpcommands.Link(embed)
		case "tfarm", "topfarm":
			embed = helpcommands.TopFarm(embed)
		case "ti", "tinfo", "tracking", "trackinfo":
			embed = helpcommands.TrackInfo(embed)
		}
	}

	if !strings.HasPrefix(embed.Description, "Detailed") && embed.Fields[0].Name == "Admin commands:" {
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
