package osutools

import (
	"fmt"
	"math"
	"os/exec"
	"regexp"
	"strconv"
	"strings"

	tools "../tools"

	"github.com/thehowl/go-osuapi"
)

// PPCalc calculates the pp given by the beatmap with specified acc and mods TODO: More args
func PPCalc(beatmap osuapi.Beatmap, acc float64, combo string, misses string, mods string, store chan<- string) {
	replacer, _ := regexp.Compile(`[^a-zA-Z0-9\s\(\)]`)

	regex, err := regexp.Compile(`pp             : (\d+)(\.\d+)?`)
	tools.ErrRead(err)

	var data []string
	var commands []string
	var mode string
	switch beatmap.Mode {
	case osuapi.ModeOsu:
		mode = "osu"
	case osuapi.ModeOsuMania:
		mode = "mania"
	case osuapi.ModeTaiko:
		mode = "taiko"
	}
	commands = append(commands, "run", "-p", "./osu-tools/PerformanceCalculator", "simulate", mode, "./data/osuFiles/"+strconv.Itoa(beatmap.BeatmapID)+" "+replacer.ReplaceAllString(beatmap.Artist, "")+" - "+replacer.ReplaceAllString(beatmap.Title, "")+".osu", "-a", fmt.Sprint(acc))

	var argResult strings.Builder

	// Check combo and misses
	if combo != "" {
		argResult.WriteString("-c " + combo + " ")
	}

	if misses != "" {
		argResult.WriteString("-X " + misses + " ")
	}

	// Check mods
	if len(mods) > 0 && mods != "NM" {
		modList := tools.StringSplit(mods, 2)
		for i := range modList {
			argResult.WriteString("-m " + strings.ToLower(modList[i]) + " ")
		}
	}

	commands = append(commands, strings.Split(argResult.String(), " ")[:]...)
	out, err := exec.Command("dotnet", commands[:]...).Output()
	tools.ErrRead(err)
	data = strings.Split(string(out), "\n")

	var res []string
	for _, line := range data {
		if regex.MatchString(line) {
			res = regex.FindStringSubmatch(line)
		}
	}
	ppValue, err := strconv.ParseFloat(res[1]+res[2], 64)
	tools.ErrRead(err)

	value := strconv.FormatFloat(math.Round(ppValue), 'f', 0, 64)
	store <- value
}
