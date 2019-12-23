package tools

import (
	"encoding/json"
	"io/ioutil"
	"os"

	structs "../structs"
)

// GetPenisRecord obtains the penis records
func GetPenisRecord() structs.PenisRecordData {
	penisRecords := structs.PenisRecordData{
		Smallest: structs.PenisData{
			Size: 1e308,
		},
	}
	_, err := os.Stat("./data/penisRecords.json")
	if err == nil {
		f, err := ioutil.ReadFile("./data/penisRecords.json")
		ErrRead(err)
		_ = json.Unmarshal(f, &penisRecords)
	}
	return penisRecords
}
