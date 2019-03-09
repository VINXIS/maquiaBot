package handle

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strings"

	structs "../structs"
	tools "../tools"
	osucommands "./osu-commands"

	"github.com/bwmarrin/discordgo"
	"github.com/thehowl/go-osuapi"
)

// MessageHandler handles any incoming messages
func MessageHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	osuAPI := osuapi.NewClient(os.Getenv("OSU_API"))

	// Obtain map cache data
	mapCache := []structs.MapData{}
	f, err := ioutil.ReadFile("./data/osuData/mapCache.json")
	tools.ErrRead(err)
	_ = json.Unmarshal(f, &mapCache)

	// Obtain profile cache data
	profileCache := []structs.MapData{}
	f, err = ioutil.ReadFile("./data/osuData/profileCache.json")
	tools.ErrRead(err)
	_ = json.Unmarshal(f, &profileCache)

	// Ignore all messages created by the bot itself
	if m.Author.ID == s.State.User.ID {
		return
	}

	profileRegex, _ := regexp.Compile(`(osu|old)\.ppy\.sh/u(sers)?/(\d+|\S+)`)
	beatmapRegex, _ := regexp.Compile(`(osu|old)\.ppy\.sh/(s|b|beatmaps(ets)?)/(\d+)(#(osu|taiko|fruits|mania)/(\d+)|\S+)?(\s)*(-n)?(\s)*(-m (\S+))?`)
	commandRegex, _ := regexp.Compile(`^\$(\S+)`)
	linkRegex, _ := regexp.Compile(`https?:\/\/\S*`)
	negateRegex, _ := regexp.Compile(`-n`)

	if negateRegex.MatchString(m.Content) {
		return
	}

	// If message linked beatmap(s) TODO: Multiple maps linked in a message
	if beatmapRegex.MatchString(m.Content) {
		go osucommands.BeatmapMessage(s, m, beatmapRegex, osuAPI, mapCache)
		return
	} else if commandRegex.MatchString(m.Content) {
		fmt.Println(strings.Split(m.Content, " -"))
		return
	} else if profileRegex.MatchString(m.Content) {
		go osucommands.ProfileMessage(s, m, profileRegex, osuAPI, profileCache)
		return
	}

	if len(m.Attachments) > 0 || (linkRegex.MatchString(m.Content) && !beatmapRegex.MatchString(m.Content)) {
		go osucommands.OsuImageParse(s, m, linkRegex, osuAPI, mapCache)
		return
	}
}
