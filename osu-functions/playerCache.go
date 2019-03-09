package osutools

import (
	"encoding/json"
	"io/ioutil"
	"time"

	"github.com/thehowl/go-osuapi"

	structs "../structs"
	tools "../tools"
)

// PlayerCache checks to see if the latest user information is already saved, otherwise it will update as necessary
func PlayerCache(user osuapi.User, cache []structs.PlayerData) {
	exists := false
	for _, player := range cache {
		if player.Osu.UserID == user.UserID {
			exists = true
			player.Osu = user
		}
	}
	if !exists {
		cache = append(cache, structs.PlayerData{
			Time: time.Now(),
			Osu:  user,
		})
	}
	jsonCache, err := json.Marshal(cache)
	tools.ErrRead(err)
	err = ioutil.WriteFile("./data/osuData/profileCache.json", jsonCache, 0644)
	tools.ErrRead(err)
}
