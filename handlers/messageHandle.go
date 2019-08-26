package handlers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"strings"

	structs "../structs"
	tools "../tools"
	gencommands "./general-commands"
	osucommands "./osu-commands"
	pokemoncommands "./pokemon-commands"

	"github.com/bwmarrin/discordgo"
	"github.com/thehowl/go-osuapi"
)

// MessageHandler handles any incoming messages
func MessageHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	negateRegex, _ := regexp.Compile(`-n\s`)

	// Ignore all messages created by the bot itself or if the negate command was stated
	if m.Author.ID == s.State.User.ID || negateRegex.MatchString(m.Content) {
		return
	}

	m.Content = strings.ToLower(m.Content)

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
	crab := true
	emptyServerData := structs.ServerData{}
	if serverData.Server.ID != emptyServerData.Server.ID {
		serverPrefix = serverData.Prefix
		crab = serverData.Crab
	}

	// CRAB RAVE
	if crab && (strings.Contains(m.Content, "crab") || strings.Contains(m.Content, "rave")) {
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
	profileRegex, _ := regexp.Compile(`(osu|old)\.ppy\.sh/u(sers)?/(\d+|\S+)`)
	beatmapRegex, _ := regexp.Compile(`(osu|old)\.ppy\.sh/(s|b|beatmaps(ets)?)/(\d+)(#(osu|taiko|fruits|mania)/(\d+)|\S+)?(\s)*(-n)?(\s)*(-m (\S+))?`)
	linkRegex, _ := regexp.Compile(`https?:\/\/\S*`)

	if strings.HasPrefix(m.Content, "maquiaprefix") {
		args := strings.Split(m.Content, " ")
		go gencommands.NewPrefix(s, m, args, serverPrefix)
		return
	} else if strings.Contains(m.Content, "maquiacleanf") || strings.Contains(m.Content, "maquiacleanfarm") {
		go gencommands.CleanFarm(s, m, profileCache)
		return
	} else if strings.Contains(m.Content, "maquiaclean") {
		go gencommands.Clean(s, m, profileCache)
		return
	} else if beatmapRegex.MatchString(m.Content) { // If a beatmap is linked
		go osucommands.BeatmapMessage(s, m, beatmapRegex, osuAPI, mapCache)
	} else if profileRegex.MatchString(m.Content) { // if a profile was linked
		go osucommands.ProfileMessage(s, m, profileRegex, osuAPI, profileCache)
	} else if strings.HasPrefix(m.Content, serverPrefix) { // If a command was declared
		args := strings.Split(m.Content, " ")
		switch args[0] {
		case serverPrefix + "osu", serverPrefix + "o":
			go OsuHandle(s, m, args, osuAPI, profileCache, mapCache, serverPrefix)
		case serverPrefix + "pokemon", serverPrefix + "p":
			go PokemonHandle(s, m, args, serverPrefix)
		case serverPrefix + "avatar", serverPrefix + "ava", serverPrefix + "a":
			go gencommands.Avatar(s, m)
		case serverPrefix + "h", serverPrefix + "help":
			go gencommands.Help(s, m, serverPrefix, args)
		case serverPrefix + "src", serverPrefix + "source":
			s.ChannelMessageSend(m.ChannelID, "https://github.com/VINXIS/maquiaBot")
		case serverPrefix + "prefix", serverPrefix + "newprefix":
			go gencommands.NewPrefix(s, m, args, serverPrefix)
		case serverPrefix + "link", serverPrefix + "set":
			go osucommands.Link(s, m, args, osuAPI, profileCache)
		case serverPrefix + "tfarm", serverPrefix + "topfarm", serverPrefix + "tfarmerdog", serverPrefix + "topfarmerdog":
			go osucommands.TopFarm(s, m, args, osuAPI, profileCache, serverPrefix)
		case serverPrefix + "bfarm", serverPrefix + "bottomfarm", serverPrefix + "bfarmerdog", serverPrefix + "bottomfarmerdog":
			go osucommands.BottomFarm(s, m, args, osuAPI, profileCache, serverPrefix)
		case serverPrefix + "farm", serverPrefix + "farmerdog", serverPrefix + "f":
			go osucommands.Farmerdog(s, m, args, osuAPI, profileCache, serverPrefix)
		case serverPrefix + "r", serverPrefix + "rs", serverPrefix + "recent":
			go osucommands.Recent(s, m, args, osuAPI, profileCache, "recent", serverPrefix, mapCache)
		case serverPrefix + "rb", serverPrefix + "recentb", serverPrefix + "recentbest":
			go osucommands.Recent(s, m, args, osuAPI, profileCache, "best", serverPrefix, mapCache)
		case serverPrefix + "c", serverPrefix + "compare":
			go osucommands.Compare(s, m, args, osuAPI, profileCache, serverPrefix, mapCache)
		case serverPrefix + "t", serverPrefix + "top":
			go osucommands.Top(s, m, args, osuAPI, profileCache, serverPrefix, mapCache)
		case serverPrefix + "tr", serverPrefix + "track":
			go osucommands.Track(s, m, args, osuAPI, mapCache)
		case serverPrefix + "tt", serverPrefix + "trackt", serverPrefix + "tracktoggle":
			go osucommands.TrackToggle(s, m, mapCache)
		case serverPrefix + "ti", serverPrefix + "tinfo", serverPrefix + "tracking", serverPrefix + "trackinfo":
			go osucommands.TrackInfo(s, m)
		case serverPrefix + "b", serverPrefix + "berry":
			go pokemoncommands.Berry(s, m, args)
		case serverPrefix + "clean":
			go gencommands.Clean(s, m, profileCache)
		case serverPrefix + "cleanf", serverPrefix + "cleanfarm":
			go gencommands.CleanFarm(s, m, profileCache)
		}
		return
	}

	// Check if an image was linked
	if len(m.Attachments) > 0 || (linkRegex.MatchString(m.Content) && !beatmapRegex.MatchString(m.Content)) {
		go osucommands.OsuImageParse(s, m, linkRegex, osuAPI, mapCache)
		return
	}
}
