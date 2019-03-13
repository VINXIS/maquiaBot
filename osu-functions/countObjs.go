package osutools

import (
	"io/ioutil"
	"regexp"
	"strconv"
	"strings"

	tools "../tools"
	"github.com/thehowl/go-osuapi"
)

// CountObjs counts the amount of objects in a beatmap
func CountObjs(beatmap osuapi.Beatmap) (objs int) {
	replacer, _ := regexp.Compile(`[^a-zA-Z0-9\s\(\)]`)

	res, err := ioutil.ReadFile("./data/osuFiles/" + strconv.Itoa(beatmap.BeatmapID) + " " + replacer.ReplaceAllString(beatmap.Artist, "") + " - " + replacer.ReplaceAllString(beatmap.Title, "") + ".osu")
	tools.ErrRead(err)
	str := string(res)
	lines := strings.Split(str, "\n")

	objs = -1
	for _, line := range lines {
		if line == "[HitObjects]\r" {
			objs = 0
		}
		if objs != -1 && line != "" {
			objs++
		}
	}
	objs--

	return objs
}
