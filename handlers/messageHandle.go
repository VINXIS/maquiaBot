package handle

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	commands "./commands"

	"github.com/bwmarrin/discordgo"
	"github.com/thehowl/go-osuapi"
)

// MessageHandler handles any incoming messages
func MessageHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	osu := osuapi.NewClient(os.Getenv("OSU_API"))

	// Ignore all messages created by the bot itself
	if m.Author.ID == s.State.User.ID {
		return
	}

	beatmapRegex, _ := regexp.Compile(`(osu|old).ppy.sh/(s|b|beatmaps(ets)?)/(\d+)(#(osu|taiko|fruits|mania)/(\d+)|\S+)?(\s)*(-n)?(\s)*(-m (\S+))?`)
	commandRegex, _ := regexp.Compile(`^\$(\S+)`)
	linkRegex, _ := regexp.Compile(`https?:\/\/\S*`)
	negateRegex, _ := regexp.Compile(`-n`)

	if negateRegex.MatchString(m.Content) {
		return
	}

	// If message linked beatmap(s) TODO: Multiple maps linked in a message
	if beatmapRegex.MatchString(m.Content) {
		go commands.BeatmapMessage(s, m, beatmapRegex, osu)
	} else if commandRegex.MatchString(m.Content) {
		fmt.Println(strings.Split(m.Content, " -"))
	}

	if len(m.Attachments) > 0 || (linkRegex.MatchString(m.Content) && !beatmapRegex.MatchString(m.Content)) {
		go commands.OsuImageParse(s, m, osu)
	}
}
