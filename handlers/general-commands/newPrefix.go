package gencommands

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"time"

	structs "../../structs"
	tools "../../tools"

	"github.com/bwmarrin/discordgo"
)

// NewPrefix sets a new prefix for the bot
func NewPrefix(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
	// Obtain server data
	serverData := structs.ServerData{}
	_, err := os.Stat("./data/serverData/" + m.GuildID + ".json")
	if err == nil {
		f, err := ioutil.ReadFile("./data/serverData/" + m.GuildID + ".json")
		tools.ErrRead(err)
		_ = json.Unmarshal(f, &serverData)
	} else if os.IsNotExist(err) {
		server, err := s.Guild(m.GuildID)
		tools.ErrRead(err)

		serverData.Server = *server
	} else {
		tools.ErrRead(err)
		return
	}

	// Set new information in server data
	serverData.Time = time.Now()
	serverData.Prefix = args[1]

	jsonCache, err := json.Marshal(serverData)
	tools.ErrRead(err)

	err = ioutil.WriteFile("./data/serverData/"+m.GuildID+".json", jsonCache, 0644)
	tools.ErrRead(err)

	s.ChannelMessageSend(m.ChannelID, "Prefix changed from "+string([]rune(args[0])[0])+" to "+args[1])
	return
}
