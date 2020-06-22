package tools

import (
	"encoding/json"
	"io/ioutil"
	"os"

	"github.com/bwmarrin/discordgo"
	structs "maquiaBot/structs"
)

// GetGenitalRecord obtains the penis records
func GetGenitalRecord(s *discordgo.Session) structs.GenitalRecordData {
	genitalRecords := structs.GenitalRecordData{
		Penis: struct {
			Largest  structs.GenitalData
			Smallest structs.GenitalData
		}{
			Smallest: structs.GenitalData{
				Size: 1e308,
			},
		},
		Vagina: struct {
			Largest  structs.GenitalData
			Smallest structs.GenitalData
		}{
			Smallest: structs.GenitalData{
				Size: 1e308,
			},
		},
	}
	_, err := os.Stat("./data/genitalRecords.json")
	if err == nil {
		f, err := ioutil.ReadFile("./data/genitalRecords.json")
		ErrRead(s, err)
		_ = json.Unmarshal(f, &genitalRecords)
	}
	return genitalRecords
}
