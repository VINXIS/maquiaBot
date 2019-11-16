package tools

import (
	"encoding/json"
	"io/ioutil"
	"os"

	structs "../structs"
	"github.com/bwmarrin/discordgo"
)

// GetServer obtains a server using its guild ID
func GetServer(server discordgo.Guild) structs.ServerData {
	serverData := structs.NewServer(server)
	_, err := os.Stat("./data/serverData/" + server.ID + ".json")
	if err == nil {
		f, err := ioutil.ReadFile("./data/serverData/" + server.ID + ".json")
		ErrRead(err)
		_ = json.Unmarshal(f, &serverData)
	}
	return serverData
}
