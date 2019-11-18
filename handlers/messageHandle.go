package handlers

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/json"
	"io/ioutil"
	"math/big"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"

	osuapi "../osu-api"
	structs "../structs"
	tools "../tools"
	gencommands "./general-commands"
	osucommands "./osu-commands"
	pokemoncommands "./pokemon-commands"

	"github.com/bwmarrin/discordgo"
)

// MessageHandler handles any incoming messages
func MessageHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	negateRegex, _ := regexp.Compile(`\s(-n\s|-n$)`)

	// Ignore all messages created by the bot itself or if the negate command was stated
	if m.Author.ID == s.State.User.ID || negateRegex.MatchString(m.Content) {
		return
	}

	m.Content = strings.ToLower(m.Content)
	if strings.Contains(m.Content, "@everyone") {
		m.Content = strings.TrimSpace(strings.ReplaceAll(m.Content, "@everyone", ""))
	}
	if strings.Contains(m.Content, "@here") {
		m.Content = strings.TrimSpace(strings.ReplaceAll(m.Content, "@here", ""))
	}

	emojiRegex, _ := regexp.Compile(`<(:.+:)\d+>`)
	noEmoji := m.Content
	if emojiRegex.MatchString(m.Content) {
		noEmoji = emojiRegex.ReplaceAllString(m.Content, emojiRegex.FindStringSubmatch(m.Content)[1])
	}

	osuAPI := osuapi.NewClient(os.Getenv("OSU_API"))

	// Obtain map cache data
	var mapCache []structs.MapData
	f, err := ioutil.ReadFile("./data/osuData/mapCache.json")
	tools.ErrRead(err)
	_ = json.Unmarshal(f, &mapCache)

	// Obtain profile cache data
	var profileCache []structs.PlayerData
	f, err = ioutil.ReadFile("./data/osuData/profileCache.json")
	tools.ErrRead(err)
	_ = json.Unmarshal(f, &profileCache)

	// Obtain mapper data
	var mapperData []structs.MapperData
	f, err = ioutil.ReadFile("./data/osuData/mapperData.json")
	tools.ErrRead(err)
	_ = json.Unmarshal(f, &mapperData)

	// Obtain server data
	var serverData structs.ServerData
	_, err = os.Stat("./data/serverData/" + m.GuildID + ".json")
	if err == nil {
		f, err = ioutil.ReadFile("./data/serverData/" + m.GuildID + ".json")
		tools.ErrRead(err)
		_ = json.Unmarshal(f, &serverData)
	}

	// Check for custom prefix
	serverPrefix := `$`
	crab := true
	osuToggle := true
	vibe := false
	if serverData.Server.ID != "" {
		serverPrefix = serverData.Prefix
		crab = serverData.Crab
		osuToggle = serverData.OsuToggle
		vibe = serverData.Vibe
	}

	// CRAB RAVE
	if crab && (strings.Contains(m.Content, "crab") || strings.Contains(m.Content, "rave")) && m.Content != serverPrefix+"crab" {
		response, err := http.Get("https://cdn.discordapp.com/emojis/510169818893385729.gif")
		if err != nil {
			return
		}

		message := &discordgo.MessageSend{
			File: &discordgo.File{
				Name:   "crab.gif",
				Reader: response.Body,
			},
		}
		s.ChannelMessageSendComplex(m.ChannelID, message)
		response.Body.Close()
	}

	// Generate regexes for message parsing
	profileRegex, _ := regexp.Compile(`(osu|old)\.ppy\.sh\/(u|users)\/(\S+)`)
	beatmapRegex, _ := regexp.Compile(`(osu|old)\.ppy\.sh\/(s|b|beatmaps|beatmapsets)\/(\d+)(#(osu|taiko|fruits|mania)\/(\d+))?`)
	linkRegex, _ := regexp.Compile(`https?:\/\/\S*`)
	timestampRegex, _ := regexp.Compile(`(\d+):(\d{2}):(\d{3})\s*(\(((\d\,?)+)\))?`)

	if timestampRegex.MatchString(noEmoji) && osuToggle {
		go osucommands.TimestampMessage(s, m, timestampRegex)
	}

	if vibe {
		roll, _ := rand.Int(rand.Reader, big.NewInt(int64(100000)))
		number := roll.Int64() + 1
		if number <= 1 {
			go gencommands.VibeCheck(s, m, "")
		}
	}

	// Command checks
	if strings.HasPrefix(m.Content, "maquiaprefix") {
		go gencommands.Prefix(s, m)
		return
	} else if strings.HasPrefix(m.Content, "maquiacleanf") || strings.Contains(m.Content, "maquiacleanfarm") {
		go gencommands.CleanFarm(s, m, profileCache, osuAPI)
		return
	} else if strings.HasPrefix(m.Content, "maquiaclean") {
		go gencommands.Clean(s, m, profileCache)
		return
	} else if strings.HasPrefix(m.Content, serverPrefix) {
		args := strings.Split(m.Content, " ")
		switch args[0] {
		case serverPrefix + "k", serverPrefix + "key":
			go s.ChannelMessageSend(m.ChannelID, "``` Default AES encryption key: Nxb]^NSc;L*qn3K(/tN{6N7%4n32fF#@```\n This key is given out publicly and I use it for all of my encryption tools, so please do not use me for sensitive data.\n I will have custom key functionality later.")
		case serverPrefix + "noncesize", serverPrefix + "nsize":
			key := []byte("Nxb]^NSc;L*qn3K(/tN{6N7%4n32fF#@")
			block, _ := aes.NewCipher(key)
			gcm, _ := cipher.NewGCM(block)
			go s.ChannelMessageSend(m.ChannelID, "The nonce size using the default AES encryption key is "+strconv.Itoa(gcm.NonceSize()))
		case serverPrefix + "src", serverPrefix + "source":
			go s.ChannelMessageSend(m.ChannelID, "https://github.com/VINXIS/maquiaBot")
		case serverPrefix + "osu", serverPrefix + "o":
			go OsuHandle(s, m, args, osuAPI, profileCache, mapCache, serverPrefix)
		case serverPrefix + "pokemon":
			go PokemonHandle(s, m, args, serverPrefix)
		case serverPrefix + "avatar", serverPrefix + "ava", serverPrefix + "a":
			go gencommands.Avatar(s, m)
		case serverPrefix + "ping":
			go gencommands.Ping(s, m)
		case serverPrefix + "up", serverPrefix + "update":
			go gencommands.Update(s, m)
		case serverPrefix + "parse":
			go gencommands.ParseID(s, m)
		case serverPrefix + "info":
			go gencommands.Info(s, m, profileCache)
		case serverPrefix + "sinfo", serverPrefix + "serverinfo":
			go gencommands.ServerInfo(s, m)
		case serverPrefix + "h", serverPrefix + "help":
			go gencommands.Help(s, m, serverPrefix, args)
		case serverPrefix + "crab":
			go gencommands.Crab(s, m)
		case serverPrefix + "vibet", serverPrefix + "vibetoggle":
			go gencommands.Vibe(s, m)
		case serverPrefix + "vibe", serverPrefix + "vibec", serverPrefix + "vibecheck":
			go gencommands.VibeCheck(s, m, "notRandom")
		case serverPrefix + "prefix", serverPrefix + "newprefix":
			go gencommands.Prefix(s, m)
		case serverPrefix + "remind":
			go gencommands.Remind(s, m)
		case serverPrefix + "rremove", serverPrefix + "remindremove":
			go gencommands.RemoveReminder(s, m)
		case serverPrefix + "reminders":
			go gencommands.Reminders(s, m)
		case serverPrefix + "l", serverPrefix + "leven", serverPrefix + "levenshtein":
			go gencommands.Levenshtein(s, m)
		case serverPrefix + "penis":
			go gencommands.Penis(s, m)
		case serverPrefix + "funny":
			go gencommands.Funny(s, m)
		case serverPrefix + "kanye":
			go gencommands.Kanye(s, m)
		case serverPrefix + "encrypt":
			go gencommands.Encrypt(s, m)
		case serverPrefix + "decrypt":
			go gencommands.Decrypt(s, m)
		case serverPrefix + "osutoggle", serverPrefix + "osut":
			go osucommands.OsuToggle(s, m)
		case serverPrefix + "link", serverPrefix + "set":
			go osucommands.Link(s, m, args, osuAPI, profileCache)
		case serverPrefix + "tfarm", serverPrefix + "topfarm":
			go osucommands.TopFarm(s, m, osuAPI, profileCache, serverPrefix)
		case serverPrefix + "bfarm", serverPrefix + "bottomfarm":
			go osucommands.BottomFarm(s, m, osuAPI, profileCache, serverPrefix)
		case serverPrefix + "farm":
			go osucommands.Farmerdog(s, m, osuAPI, profileCache)
		case serverPrefix + "ppadd":
			go osucommands.PPAdd(s, m, osuAPI, profileCache)
		case serverPrefix + "r", serverPrefix + "rs", serverPrefix + "recent":
			go osucommands.Recent(s, m, osuAPI, "recent", profileCache, mapCache)
		case serverPrefix + "rb", serverPrefix + "recentb", serverPrefix + "recentbest":
			go osucommands.Recent(s, m, osuAPI, "best", profileCache, mapCache)
		case serverPrefix + "c", serverPrefix + "compare":
			go osucommands.Compare(s, m, args, osuAPI, profileCache, serverPrefix, mapCache)
		case serverPrefix + "t", serverPrefix + "top":
			go osucommands.Top(s, m, osuAPI, profileCache, mapCache)
		case serverPrefix + "tr", serverPrefix + "track":
			go osucommands.Track(s, m, osuAPI, mapCache)
		case serverPrefix + "tt", serverPrefix + "trackt", serverPrefix + "ttoggle", serverPrefix + "tracktoggle":
			go osucommands.TrackToggle(s, m, mapCache)
		case serverPrefix + "ti", serverPrefix + "tinfo", serverPrefix + "tracking", serverPrefix + "trackinfo":
			go osucommands.TrackInfo(s, m)
		case serverPrefix + "mt", serverPrefix + "mtrack", serverPrefix + "maptrack", serverPrefix + "mappertrack":
			go osucommands.TrackMapper(s, m, osuAPI, mapperData)
		case serverPrefix + "mti", serverPrefix + "mtinfo", serverPrefix + "mtrackinfo", serverPrefix + "maptracking", serverPrefix + "mappertracking", serverPrefix + "mappertrackinfo":
			go osucommands.TrackMapperInfo(s, m, mapperData)
		case serverPrefix + "b", serverPrefix + "berry":
			go pokemoncommands.Berry(s, m)
		case serverPrefix + "p", serverPrefix + "percentage", serverPrefix + "per", serverPrefix + "percent":
			go gencommands.Percentage(s, m)
		case serverPrefix + "roll":
			go gencommands.Roll(s, m)
		case serverPrefix + "dubs":
			go s.ChannelMessageSend(m.ChannelID, "Ur retarded")
		case serverPrefix + "kanye":
			go gencommands.Kanye(s, m)
		case serverPrefix + "ch", serverPrefix + "choose":
			go gencommands.Choose(s, m)
		case serverPrefix + "stats":
			go gencommands.Stats(s, m)
		case serverPrefix + "adj", serverPrefix + "adjective", serverPrefix + "adjectives":
			go gencommands.Adjectives(s, m)
		case serverPrefix + "noun", serverPrefix + "nouns":
			go gencommands.Nouns(s, m)
		case serverPrefix + "skill", serverPrefix + "skills":
			go gencommands.Skills(s, m)
		case serverPrefix + "statst", serverPrefix + "statstoggle":
			go gencommands.StatsToggle(s, m)
		case serverPrefix + "clean":
			go gencommands.Clean(s, m, profileCache)
		case serverPrefix + "cleanf", serverPrefix + "cleanfarm":
			go gencommands.CleanFarm(s, m, profileCache, osuAPI)
		case serverPrefix + "ocr":
			go gencommands.OCR(s, m)
		}
		return
	} else if beatmapRegex.MatchString(m.Content) && osuToggle { // If a beatmap was linked
		go osucommands.BeatmapMessage(s, m, beatmapRegex, osuAPI, mapCache)
		return
	} else if profileRegex.MatchString(m.Content) && osuToggle { // if a profile was linked
		go osucommands.ProfileMessage(s, m, profileRegex, osuAPI, profileCache)
		return
	}

	if len(m.Mentions) > 0 {
		for _, mention := range m.Mentions {
			if mention.ID == s.State.User.ID {
				s.ChannelMessageSend(m.ChannelID, "what do u want dude lol")
				break
			}
		}
	}

	// Check if an image was linked
	if len(m.Attachments) > 0 || (linkRegex.MatchString(m.Content)) || (len(m.Embeds) > 0 && m.Embeds[0].Image != nil) {
		go osucommands.OsuImageParse(s, m, linkRegex, osuAPI, mapCache)
		return
	}
}
