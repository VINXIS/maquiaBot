package handlers

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"regexp"
	"strings"

	structs "../structs"
	tools "../tools"
	gencommands "./general-commands"
	osucommands "./osu-commands"

	"github.com/bwmarrin/discordgo"
	"github.com/thehowl/go-osuapi"
)

// MessageHandler handles any incoming messages
func MessageHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	negateRegex, _ := regexp.Compile(`-n\W`)

	// Ignore all messages created by the bot itself
	if m.Author.ID == s.State.User.ID || negateRegex.MatchString(m.Content) {
		return
	}

	go ML(s, m)

	osuAPI := osuapi.NewClient(os.Getenv("OSU_API"))

	// Obtain map cache data
	mapCache := []structs.MapData{}
	f, err := ioutil.ReadFile("./data/osuData/mapCache.json")
	tools.ErrRead(err)
	_ = json.Unmarshal(f, &mapCache)

	// Obtain profile cache data
	profileCache := []structs.PlayerData{}
	f, err = ioutil.ReadFile("./data/osuData/profileCache.json")
	tools.ErrRead(err)
	_ = json.Unmarshal(f, &profileCache)

	// Obtain server data
	serverData := structs.ServerData{}
	_, err = os.Stat("./data/serverData/" + m.GuildID + ".json")
	if err == nil {
		f, err = ioutil.ReadFile("./data/serverData/" + m.GuildID + ".json")
		tools.ErrRead(err)
		_ = json.Unmarshal(f, &serverData)
	}

	// Check for custom prefix
	serverPrefix := `$`
	serverRegex := `^\` + serverPrefix + `(\S+)`
	emptyServerData := structs.ServerData{}
	if serverData.Server.ID != emptyServerData.Server.ID {
		serverPrefix = serverData.Prefix
		charCheck, _ := regexp.Compile(`[^a-zA-Z0-9\s\(\)]`)
		if charCheck.MatchString(serverPrefix) {
			serverRegex = `^` + charCheck.ReplaceAllString(serverPrefix, `\`+serverPrefix) + `(\S+)`
		} else {
			serverRegex = `^` + serverPrefix + `(\S+)`
		}
	}

	// Generate regexes for message parsing
	profileRegex, _ := regexp.Compile(`(osu|old)\.ppy\.sh/u(sers)?/(\d+|\S+)`)
	beatmapRegex, _ := regexp.Compile(`(osu|old)\.ppy\.sh/(s|b|beatmaps(ets)?)/(\d+)(#(osu|taiko|fruits|mania)/(\d+)|\S+)?(\s)*(-n)?(\s)*(-m (\S+))?`)
	commandRegex, _ := regexp.Compile(serverRegex)
	linkRegex, _ := regexp.Compile(`https?:\/\/\S*`)

	if strings.HasPrefix(m.Content, "maquiaprefix") {
		args := strings.Split(m.Content, " ")
		go gencommands.NewPrefix(s, m, args, serverPrefix)
	} else if beatmapRegex.MatchString(m.Content) { // If a beatmap is linked
		go osucommands.BeatmapMessage(s, m, beatmapRegex, osuAPI, mapCache)
		return
	} else if profileRegex.MatchString(m.Content) { // if a profile was linked
		go osucommands.ProfileMessage(s, m, profileRegex, osuAPI, profileCache)
		return
	} else if commandRegex.MatchString(m.Content) { // If a command was declared
		args := strings.Split(m.Content, " ")
		command := args[0]
		switch command {
		case serverPrefix + "osu", serverPrefix + "o":
			go OsuHandle(s, m, args, osuAPI, profileCache, mapCache, serverPrefix)
		case serverPrefix + "pokemon", serverPrefix + "p":
			go PokemonHandle(s, m, args, serverPrefix)
		case serverPrefix + "avatar", serverPrefix + "ava", serverPrefix + "a":
			go gencommands.Avatar(s, m)
		case serverPrefix + "help":
			go gencommands.Help(s, m, serverPrefix)
		case serverPrefix + "src", serverPrefix + "source":
			s.ChannelMessageSend(m.ChannelID, "https://github.com/VINXIS/maquiaBot")
		case serverPrefix + "prefix", serverPrefix + "newprefix":
			go gencommands.NewPrefix(s, m, args, serverPrefix)
		case serverPrefix + "link", serverPrefix + "set":
			go osucommands.Link(s, m, args, osuAPI, profileCache)
		case serverPrefix + "r", serverPrefix + "rs", serverPrefix + "recent":
			go osucommands.Recent(s, m, args, osuAPI, profileCache, "recent", serverPrefix, mapCache)
		case serverPrefix + "rb", serverPrefix + "recentb", serverPrefix + "recentbest":
			go osucommands.Recent(s, m, args, osuAPI, profileCache, "best", serverPrefix, mapCache)
		case serverPrefix + "c", serverPrefix + "compare":
			go osucommands.Compare(s, m, args, osuAPI, profileCache, serverPrefix, mapCache)
		case serverPrefix + "t", serverPrefix + "track":
			go osucommands.Track(s, m, args, osuAPI, mapCache)
		case serverPrefix + "tt", serverPrefix + "trackt", serverPrefix + "tracktoggle":
			go osucommands.TrackToggle(s, m, mapCache)
		case serverPrefix + "tinfo", serverPrefix + "tracking", serverPrefix + "trackinfo":
			go osucommands.TrackInfo(s, m)
		}
		return
	}

	// Check if an image was linked
	if len(m.Attachments) > 0 || (linkRegex.MatchString(m.Content) && !beatmapRegex.MatchString(m.Content)) {
		go osucommands.OsuImageParse(s, m, linkRegex, osuAPI, mapCache)
		return
	}
}
