package handlers

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/json"
	"io/ioutil"
	"math/big"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	admincommands "maquiaBot/handlers/admin-commands"
	botcreatorcommands "maquiaBot/handlers/bot-creator-commands"
	gencommands "maquiaBot/handlers/general-commands"
	mathcommands "maquiaBot/handlers/math-commands"
	osucommands "maquiaBot/handlers/osu-commands"
	pokemoncommands "maquiaBot/handlers/pokemon-commands"
	"maquiaBot/structs"
	"maquiaBot/tools"

	"github.com/bwmarrin/discordgo"
)

// MessageHandler handles any incoming messages
func MessageHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Ignore all messages created by the bot itself or by other bots
	if m.Author.ID == s.State.User.ID || m.Author.Bot {
		return
	}

	if strings.Contains(m.Content, "@everyone") {
		m.Content = strings.TrimSpace(strings.ReplaceAll(m.Content, "@everyone", ""))
	}
	if strings.Contains(m.Content, "@here") {
		m.Content = strings.TrimSpace(strings.ReplaceAll(m.Content, "@here", ""))
	}

	emojiRegex, _ := regexp.Compile(`(?i)<a?(:.+:)\d+>`)
	noEmoji := m.Content
	if emojiRegex.MatchString(m.Content) {
		noEmoji = emojiRegex.ReplaceAllString(m.Content, emojiRegex.FindStringSubmatch(m.Content)[1])
	}

	// Obtain server data
	server, err := s.Guild(m.GuildID)
	if err != nil {
		server = &discordgo.Guild{}
	}
	serverData, _ := tools.GetServer(*server, s)
	serverPrefix := serverData.Prefix

	// Obtain channel data
	channel, err := s.Channel(m.ChannelID)
	if err != nil {
		channel = &discordgo.Channel{}
	}
	channelData, _ := tools.GetChannel(*channel, s)

	// Generate regexes for message parsing
	profileRegex, _ := regexp.Compile(`(?i)(osu|old)\.ppy\.sh\/(u|users)\/(\S+)`)
	beatmapRegex, _ := regexp.Compile(`(?i)(https:\/\/)?(osu|old)\.ppy\.sh\/(s|b|beatmaps|beatmapsets)\/(\d+)(#(osu|taiko|fruits|mania)\/(\d+))?`)
	linkRegex, _ := regexp.Compile(`(?i)https?:\/\/\S*`)
	timestampRegex, _ := regexp.Compile(`(?i)(\d+):(\d{2}):(\d{3})\s*(\(((\d\,?)+)\))?`)

	// Timestamp conversions
	if timestampRegex.MatchString(noEmoji) && (serverData.TimestampToggle || channelData.TimestampToggle) {
		go osucommands.TimestampMessage(s, m, timestampRegex)
	}

	// Vibe check (1/100000 chance if vibe is on in the server)
	if serverData.Vibe || channelData.Vibe {
		roll, _ := rand.Int(rand.Reader, big.NewInt(100000))
		number := roll.Int64()
		if number == 0 {
			go gencommands.Vibe(s, m, "")
		}
	}

	// Role checks
	for _, roleAuto := range serverData.RoleAutomation {
		roleAuto.Text = `(?i)` + roleAuto.Text
		reg, err := regexp.Compile(roleAuto.Text)

		match := false
		if err != nil {
			if strings.Contains(strings.ToLower(m.Content), roleAuto.Text) {
				match = true
			}
		} else if reg.MatchString(strings.ToLower(m.Content)) {
			match = true
		}

		if match {
			for _, role := range roleAuto.Roles {
				hasRole := false
				if m.Member != nil {
					for _, memberRole := range m.Member.Roles {
						if memberRole == role {
							hasRole = true
							break
						}
					}
				}
				if hasRole {
					s.GuildMemberRoleRemove(m.GuildID, m.Author.ID, role)
				} else {
					s.GuildMemberRoleAdd(m.GuildID, m.Author.ID, role)
				}
			}
		}
	}

	// Trigger checks
	for _, trigger := range serverData.Triggers {
		trigger.Cause = `(?i)` + trigger.Cause
		reg, err := regexp.Compile(trigger.Cause)
		send := false
		if (err != nil && strings.Contains(strings.ToLower(m.Content), trigger.Cause)) || (err == nil && reg.MatchString(strings.ToLower(m.Content))) {
			send = true
		}

		if send {
			allowedFormats := []string{".jpg", ".jpeg", ".png", ".gif", ".bmp", ".tiff", ".svg", ".mp4", ".avi", ".mov", ".webm", ".flv"}
			format := ""
			for _, allowedFormat := range allowedFormats {
				if strings.Contains(trigger.Result, allowedFormat) {
					format = allowedFormat
				}
			}

			if format == "" {
				s.ChannelMessageSend(m.ChannelID, trigger.Result)
			} else {
				response, err := http.Get(trigger.Result)
				if err == nil {
					s.ChannelMessageSendComplex(m.ChannelID, &discordgo.MessageSend{
						File: &discordgo.File{
							Name:   "trigger" + format,
							Reader: response.Body,
						},
					})
				}
			}
		}
	}

	// Counter checks
	change := false
	for i, counter := range serverData.Counters {
		counter.Text = `(?i)` + counter.Text
		reg, err := regexp.Compile(counter.Text)
		count := false
		if (err != nil && strings.Contains(strings.ToLower(m.Content), counter.Text)) || reg.MatchString(strings.ToLower(m.Content)) {
			count = true
			change = true
		}

		if count {
			exists := false
			for j, countUser := range counter.Users {
				if countUser.UserID == m.Author.ID {
					serverData.Counters[i].Users[j].Count++
					exists = true
					break
				}
			}

			if !exists {
				serverData.Counters[i].Users = append(serverData.Counters[i].Users, structs.CounterTrack{
					Username: m.Author.Username,
					UserID:   m.Author.ID,
					Count:    1,
				})
			}
		}
	}

	if change {
		jsonCache, err := json.Marshal(serverData)
		tools.ErrRead(s, err)

		err = ioutil.WriteFile("./data/serverData/"+m.GuildID+".json", jsonCache, 0644)
		tools.ErrRead(s, err)
	}

	// Command checks
	if strings.HasPrefix(m.Content, "maquia") {
		args := strings.Split(strings.Split(m.Content, "\n")[0], " ")
		switch strings.ToLower(strings.Replace(args[0], "maquia", "", 1)) {
		case "clean":
			go botcreatorcommands.Clean(s, m)
		case "cleanf", "cleanfarm":
			go botcreatorcommands.CleanFarm(s, m)

		case "prefix", "newprefix":
			prefixArgs := strings.Split(m.Content, " ")
			if len(prefixArgs) < 2 {
				go s.ChannelMessageSend(m.ChannelID, "The server prefix is "+serverPrefix)
			} else {
				go admincommands.Prefix(s, m)
			}
		case "help":
			go HelpHandle(s, m, serverPrefix)
		}
	} else if strings.HasPrefix(m.Content, serverPrefix) {
		args := strings.Split(strings.Split(m.Content, "\n")[0], " ")
		switch strings.ToLower(strings.Replace(args[0], serverPrefix, "", 1)) {
		// Commands without functions
		case "complain":
			go s.ChannelMessageSend(m.ChannelID, "Shut up hoe")
		case "dubs", "doubles", "trips", "triples", "quads", "quadruples", "quints", "quintuples", "sexts", "sextuples", "septs", "septuples", "octs", "octuples", "nons", "nontuples":
			go s.ChannelMessageSend(m.ChannelID, "Fuck you")
		case "k", "key":
			go s.ChannelMessageSend(m.ChannelID, "``` Default AES encryption key: Nxb]^NSc;L*qn3K(/tN{6N7%4n32fF#@```\n This key is given out publicly and I use it for all of my encryption tools, so please do not use me for sensitive data.\n To use your own key, make sure you add a `-k` flag!")
		case "noncesize", "nsize":
			key := []byte("Nxb]^NSc;L*qn3K(/tN{6N7%4n32fF#@")
			block, _ := aes.NewCipher(key)
			gcm, _ := cipher.NewGCM(block)
			go s.ChannelMessageSend(m.ChannelID, "The nonce size using the default AES encryption key is "+strconv.Itoa(gcm.NonceSize()))
		case "src", "source":
			go s.ChannelMessageSend(m.ChannelID, "https://github.com/VINXIS/maquiaBot")

		// Bot owner commands
		case "announce":
			go botcreatorcommands.Announce(s, m)
		case "clean":
			go botcreatorcommands.Clean(s, m)
		case "cleane", "cleanempty":
			go botcreatorcommands.CleanEmpty(s, m)
		case "cleanf", "cleanfarm":
			go botcreatorcommands.CleanFarm(s, m)
		case "up", "update":
			go botcreatorcommands.Update(s, m)
		case "servers":
			go botcreatorcommands.Servers(s, m)

		// Sub-handles for pokemon and osu!
		case "h", "help":
			go HelpHandle(s, m, serverPrefix)
		case "o", "osu":
			go OsuHandle(s, m, args)
		case "pokemon":
			go PokemonHandle(s, m, args, serverPrefix)
		case "math":
			go MathHandle(s, m, args)

		// Admin commands
		case "counter":
			go admincommands.Counter(s, m)
		case "dlch", "dlchannel", "downloadch", "downloadchannel":
			go admincommands.DownloadChannel(s, m)
		case "dlsv", "dlserver", "downloadsv", "downloadserver":
			go admincommands.DownloadServer(s, m)
		case "prefix", "newprefix":
			go admincommands.Prefix(s, m)
		case "purge":
			go admincommands.Purge(s, m)
		case "rmch", "rmchannel", "removech", "removechannel":
			go admincommands.RemoveChannel(s, m)
		case "rmsv", "rmserver", "removesv", "removeserver":
			go admincommands.RemoveServer(s, m)
		case "rolea", "roleauto", "roleautomation":
			go admincommands.RoleAutomation(s, m)
		case "toggle":
			go admincommands.Toggle(s, m)
		case "tr", "track":
			go admincommands.Track(s, m)
		case "trigger":
			go admincommands.Trigger(s, m)
		case "tt", "trackt", "ttoggle", "tracktoggle":
			go admincommands.TrackToggle(s, m)

		// General commands
		case "adj", "adjective", "adjectives":
			go gencommands.Adjectives(s, m)
		case "a", "ava", "avatar":
			go gencommands.Avatar(s, m)
		case "aq", "avaquote", "avatarquote", "quoteava", "quoteavatar":
			if len(strings.Split(m.Content, " ")) < 2 {
				m.Content += " " + m.Author.Username
			}
			go gencommands.Avatar(s, m)
			go gencommands.Quote(s, m)
		case "cap", "caps", "upper":
			go gencommands.TextManipulation(s, m, "allCaps")
		case "conv", "convo", "conversation":
			go gencommands.Conversation(s, m)
		case "cp", "comparep", "comparepenis":
			if serverData.Daily || channelData.Daily {
				go gencommands.PenisCompare(s, m)
			}
		case "cv", "comparev", "comparevagina":
			if serverData.Daily || channelData.Daily {
				go gencommands.VaginaCompare(s, m)
			}
		case "cinfo", "chinfo", "channelinfo":
			go gencommands.ChannelInfo(s, m)
		case "ch", "choose":
			go gencommands.Choose(s, m)
		case "cheers":
			go gencommands.FunnyMedia(s, m, "cheers")
		case "col", "color", "colour":
			go gencommands.Colour(s, m)
		case "cr", "crank", "countr", "countrank":
			go gencommands.CountRank(s, m)
		case "cs", "counters":
			go gencommands.Counters(s, m)
		case "crab":
			go gencommands.FunnyMedia(s, m, "crab")
		case "daily":
			go gencommands.Daily(s, m)
		case "decrypt":
			go gencommands.Decrypt(s, m)
		case "dl", "download":
			go gencommands.Download(s, m)
		case "e", "emoji", "emote":
			go gencommands.Emoji(s, m)
		case "encrypt":
			go gencommands.Encrypt(s, m)
		case "face":
			go gencommands.Face(s, m)
		case "history":
			if serverData.Daily || channelData.Daily {
				go gencommands.History(s, m)
			}
		case "idea", "niceidea":
			go gencommands.FunnyMedia(s, m, "idea")
		case "info":
			go gencommands.Info(s, m)
		case "kanye":
			go gencommands.Kanye(s, m)
		case "late", "old", "ancient":
			go gencommands.FunnyMedia(s, m, "late")
		case "leven", "levenshtein":
			go gencommands.Levenshtein(s, m)
		case "list":
			go gencommands.List(s, m)
		case "lower":
			go gencommands.TextManipulation(s, m, "allLower")
		case "meme":
			go gencommands.Meme(s, m)
		case "merge":
			go gencommands.Merge(s, m)
		case "noun", "nouns":
			go gencommands.Nouns(s, m)
		case "ntw", "numtw", "ntow", "ntword", "numtow", "numtword", "ntoword", "numtoword", "numbertoword":
			go gencommands.Ntow(s, m)
		case "ocr":
			go gencommands.OCR(s, m)
		case "over":
			go gencommands.FunnyMedia(s, m, "over")
		case "p", "per", "percent", "percentage":
			go gencommands.Percentage(s, m)
		case "parse":
			go gencommands.Parse(s, m)
		case "penis":
			if serverData.Daily || channelData.Daily {
				go gencommands.Penis(s, m)
			}
		case "ping":
			go gencommands.Ping(s, m)
		case "q", "quote":
			go gencommands.Quote(s, m)
		case "qa", "qadd", "quotea", "quoteadd":
			go gencommands.QuoteAdd(s, m)
		case "qd", "qr", "qdelete", "qremove", "quotedelete", "quoteremove":
			go gencommands.QuoteRemove(s, m)
		case "qs", "quotes":
			go gencommands.Quotes(s, m)
		case "rcap", "rcaps", "rupper", "rlower", "randomcap", "randomcaps", "randomupper", "randomlower":
			go gencommands.TextManipulation(s, m, "random")
		case "rp", "rankp", "rankpenis":
			if serverData.Daily || channelData.Daily {
				go gencommands.PenisRank(s, m)
			}
		case "remind", "reminder":
			go gencommands.Remind(s, m)
		case "reminders":
			go gencommands.Reminders(s, m)
		case "remindremove", "rremove":
			go gencommands.RemoveReminder(s, m)
		case "rinfo", "roleinfo":
			go gencommands.RoleInfo(s, m)
		case "roll":
			go gencommands.Roll(s, m)
		case "rv", "rankv", "rankvagina":
			if serverData.Daily || channelData.Daily {
				go gencommands.VaginaRank(s, m)
			}
		case "sinfo", "serverinfo":
			go gencommands.ServerInfo(s, m)
		case "skill", "skills":
			go gencommands.Skills(s, m)
		case "stats", "class":
			go gencommands.Stats(s, m)
		case "swap":
			go gencommands.TextManipulation(s, m, "swap")
		case "task":
			go gencommands.Task(s, m)
		case "tasks":
			go gencommands.Tasks(s, m)
		case "taskremove", "tremove":
			go gencommands.RemoveTask(s, m)
		case "title":
			go gencommands.TextManipulation(s, m, "title")
		case "transfer":
			go gencommands.Transfer(s, m)
		case "triggers":
			go gencommands.Triggers(s, m)
		// case "tts":
		// 	go gencommands.TTS(s, m)
		case "twitch", "twitchdl":
			go gencommands.Twitch(s, m)
		case "twitter", "twitterdl":
			go gencommands.Twitter(s, m)
		case "u", "unlink":
			go gencommands.Unlink(s, m)
		case "vagina":
			if serverData.Daily || channelData.Daily {
				go gencommands.Vagina(s, m)
			}
		case "vibe", "vibec", "vibecheck":
			go gencommands.Vibe(s, m, "notRandom")
		case "w", "weather":
			go gencommands.Weather(s, m)
		case "wtn", "wtnum", "wton", "wordtn", "wtonum", "wordtnum", "wordton", "wordtonum", "wordtonumber":
			go gencommands.Wton(s, m)

		// Math commands
		case "ave", "average", "mean":
			go mathcommands.Average(s, m)
		case "d", "dist", "distance", "dir", "direction":
			go mathcommands.DistanceDirection(s, m)
		case "dr", "degrad", "degreesradians":
			go mathcommands.DegreesRadians(s, m)
		case "rd", "raddeg", "radiansdegrees":
			go mathcommands.RadiansDegrees(s, m)
		case "stddev", "standarddev", "stddeviation", "standarddeviation":
			go mathcommands.StandardDeviation(s, m)
		case "va", "vadd", "vectora", "vectoradd":
			go mathcommands.VectorAdd(s, m)
		case "vc", "vcross", "vectorc", "vectorcross":
			go mathcommands.VectorCross(s, m)
		case "vd", "vdiv", "vdivide", "vectord", "vectordiv", "vectordivide":
			go mathcommands.VectorDivide(s, m)
		case "vdot", "vectordot":
			go mathcommands.VectorDot(s, m)
		case "vm", "vmult", "vmultiply", "vectorm", "vectormult", "vectormultiply":
			go mathcommands.VectorMultiply(s, m)
		case "vs", "vsub", "vsubtract", "vectors", "vectorsub", "vectorsubtract":
			go mathcommands.VectorSubtract(s, m)

		// osu! commands
		case "bfarm", "bottomfarm":
			go osucommands.BottomFarm(s, m)
		case "bpm":
			if serverData.Daily || channelData.Daily {
				go osucommands.BPM(s, m)
			}
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
		case "osutop", "osudetail":
			go osucommands.ProfileMessage(s, m, profileRegex)
		case "ppadd", "addpp":
			go osucommands.PPAdd(s, m)
		case "profile":
			go osucommands.ProfileMessage(s, m, profileRegex)
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

		// Pokemon commands
		case "b", "berry":
			go pokemoncommands.Berry(s, m)
		}
	} else if beatmapRegex.MatchString(m.Content) && (serverData.OsuToggle || channelData.OsuToggle) { // If a beatmap was linked
		go osucommands.BeatmapMessage(s, m, "", beatmapRegex)
	} else if profileRegex.MatchString(m.Content) && (serverData.OsuToggle || channelData.OsuToggle) { // If a profile was linked
		go osucommands.ProfileMessage(s, m, profileRegex)
	}

	// Dont mention me mate. Ill fuck u up
	if len(m.Mentions) > 0 {
		for _, mention := range m.Mentions {
			if mention.ID == s.State.User.ID {
				roll, _ := rand.Int(rand.Reader, big.NewInt(100))
				number := roll.Int64()
				if number == 51 {
					go gencommands.Vibe(s, m, "")
				} else if number%17 == 0 {
					s.ChannelMessageSend(m.ChannelID, "Dude I'm serious. Stop pinging me or there will be consequences.")
				} else if number%11 == 0 {
					s.ChannelMessageSend(m.ChannelID, "lol what do u want dude i bet u havent even watched the Maquia movie stop pinging me .")
				} else if number%5 == 0 {
					s.ChannelMessageSend(m.ChannelID, "what!!!! what do u want!!!!!")
				} else {
					s.ChannelMessageSend(m.ChannelID, "what do u want dude lol")
				}
				break
			}
		}
	}

	// Check if an image was linked
	if len(m.Attachments) > 0 || linkRegex.MatchString(m.Content) || (len(m.Embeds) > 0 && m.Embeds[0].Image != nil) {
		if serverData.OsuToggle || channelData.OsuToggle {
			go osucommands.OsuImageParse(s, m, linkRegex)
		}
		go osucommands.ReplayMessage(s, m, linkRegex)
	}
}
