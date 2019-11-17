package osutools

import (
	"os/exec"
	"regexp"
	"strconv"
	"strings"

	osuapi "../osu-api"
	tools "../tools"
)

// PPCalc calculates the pp given by the beatmap with specified acc and mods TODO: More args
func PPCalc(beatmap osuapi.Beatmap, acc float64, combo string, misses string, mods string, store chan<- string) {
	replacer, _ := regexp.Compile(`[^a-zA-Z0-9\s\(\)]`)

	regex, err := regexp.Compile(`pp\s+:\s+(.+)`)
	tools.ErrRead(err)

	var data []string
	var args []string
	var mode string
	switch beatmap.Mode {
	case osuapi.ModeOsu:
		mode = "osu"
	case osuapi.ModeTaiko:
		mode = "taiko"
	case osuapi.ModeOsuMania:
		mode = "mania"
	}
	args = append(args, "./osu-tools/PerformanceCalculator/bin/Debug/netcoreapp2.0/PerformanceCalculator.dll", "simulate", mode, "./data/osuFiles/"+strconv.Itoa(beatmap.BeatmapID)+" "+replacer.ReplaceAllString(beatmap.Artist, "")+" - "+replacer.ReplaceAllString(beatmap.Title, "")+".osu", "-a", strconv.FormatFloat(acc, 'f', 2, 64))

	// Check combo and misses
	if combo != "" {
		args = append(args, "-c", combo)
	}

	if misses != "" {
		args = append(args, "-X", misses)
	}

	// Check mods
	if len(mods) > 0 && mods != "NM" {
		modList := tools.StringSplit(mods, 2)
		for i := range modList {
			args = append(args, "-m", strings.ToLower(modList[i]))
		}
	}

	out, err := exec.Command("dotnet", args...).Output()
	if err != nil {
		store <- "N/A"
		return
	}
	data = strings.Split(string(out), "\n")

	var res []string
	for _, line := range data {
		if regex.MatchString(line) {
			res = regex.FindStringSubmatch(line)
		}
	}
	if len(res) == 0 {
		store <- "N/A"
		return
	}
	ppValue, err := strconv.ParseFloat(strings.Replace(res[1], "\r", "", -1), 64)
	tools.ErrRead(err)

	value := strconv.FormatFloat(ppValue, 'f', 2, 64)
	store <- value
}
