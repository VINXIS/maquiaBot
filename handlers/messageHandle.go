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
		serverRegex = `^\` + serverPrefix + `(\S+)`
	}

	profileRegex, _ := regexp.Compile(`(osu|old)\.ppy\.sh/u(sers)?/(\d+|\S+)`)
	beatmapRegex, _ := regexp.Compile(`(osu|old)\.ppy\.sh/(s|b|beatmaps(ets)?)/(\d+)(#(osu|taiko|fruits|mania)/(\d+)|\S+)?(\s)*(-n)?(\s)*(-m (\S+))?`)
	commandRegex, _ := regexp.Compile(serverRegex)
	linkRegex, _ := regexp.Compile(`https?:\/\/\S*`)

	if beatmapRegex.MatchString(m.Content) { // If a beatmap is linked
		go osucommands.BeatmapMessage(s, m, beatmapRegex, osuAPI, mapCache)
		return
	} else if profileRegex.MatchString(m.Content) { // if a profile was linked
		go osucommands.ProfileMessage(s, m, profileRegex, osuAPI, profileCache)
		return
	} else if commandRegex.MatchString(m.Content) { // If a command was written
		args := strings.Split(m.Content, " ")
		command := args[0]
		switch command {
		case serverPrefix + "osu", serverPrefix + "o":
			go OsuHandle(s, m, args, osuAPI, profileCache, mapCache, serverPrefix)
		case serverPrefix + "avatar":
			go gencommands.Avatar(s, m)
		case serverPrefix + "help":
			go gencommands.Help(s, m, serverPrefix)
		case serverPrefix + "prefix":
			go gencommands.NewPrefix(s, m, args)
		case serverPrefix + "r", serverPrefix + "rs", serverPrefix + "recent":
			go osucommands.Recent(s, m, args, osuAPI, profileCache, "recent", mapCache)
		case serverPrefix + "rb", serverPrefix + "recentb", serverPrefix + "recentbest":
			go osucommands.Recent(s, m, args, osuAPI, profileCache, "best", mapCache)
		}
		return
	}
	if len(m.Attachments) > 0 || (linkRegex.MatchString(m.Content) && !beatmapRegex.MatchString(m.Content)) {
		go osucommands.OsuImageParse(s, m, linkRegex, osuAPI, mapCache)
		return
	}
}
