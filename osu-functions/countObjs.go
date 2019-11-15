package osutools

import (
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strconv"
	"strings"

	osuapi "../osu-api"
)

// CountObjs counts the amount of objects in a beatmap
func CountObjs(beatmap osuapi.Beatmap) (objs int) {
	replacer, _ := regexp.Compile(`[^a-zA-Z0-9\s\(\)]`)

	res, err := ioutil.ReadFile("./data/osuFiles/" + strconv.Itoa(beatmap.BeatmapID) + " " + replacer.ReplaceAllString(beatmap.Artist, "") + " - " + replacer.ReplaceAllString(beatmap.Title, "") + ".osu")

	for {
		if err != nil {
			osuAPI := osuapi.NewClient(os.Getenv("OSU_API"))
			fmt.Print("An error occured trying to fetch the beatmap file: ")
			fmt.Println(err)
			fmt.Println("Trying again...")
			BeatmapParse(strconv.Itoa(beatmap.BeatmapID), "map", 0, osuAPI)
			res, err = ioutil.ReadFile("./data/osuFiles/" + strconv.Itoa(beatmap.BeatmapID) + " " + replacer.ReplaceAllString(beatmap.Artist, "") + " - " + replacer.ReplaceAllString(beatmap.Title, "") + ".osu")
		} else if err == nil {
			break
		}
	}
	str := string(res)
	lines := strings.Split(str, "\n")

	objs = -1
	for _, line := range lines {
		if strings.Contains(line, "[HitObjects]") {
			objs = 0
		}
		if objs != -1 && line != "" {
			objs++
		}
	}
	objs--

	return objs
}
