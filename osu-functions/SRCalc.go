package osutools

import (
	"github.com/thehowl/go-osuapi"
)

// SRCalc calculates the aim, speed, and total SR for a beatmap
func SRCalc(beatmap osuapi.Beatmap, mods string) (aim, speed, total string) {
	/*replacer, _ := regexp.Compile(`[^a-zA-Z0-9\s\(\)]`)
	//fileName := strconv.Itoa(rand.Intn(10000000))

	var commands []string
	commands = append(commands, "/C", "start", "dotnet", "run", "-p", "./osu-tools/PerformanceCalculator", "difficulty", "./data/osuFiles/"+strconv.Itoa(beatmap.BeatmapID)+" "+replacer.ReplaceAllString(beatmap.Artist, "")+" - "+replacer.ReplaceAllString(beatmap.Title, "")+".osu")
	cmd := exec.Command("cmd", commands[:]...)
	err := cmd.Start()
	tools.ErrRead(err)
	fmt.Println(string(res))
	for {
		if cmd.ProcessState.ExitCode() == 0 {
			text, err := ioutil.ReadFile(fileName + ".txt")
			tools.ErrRead(err)
			fmt.Println(string(text))
			tools.DeleteFile("./" + fileName + ".txt")
			break
		}
	}*/

	aim = "**Aim SR**: 0 "
	speed = "**Speed SR**: 0 "
	total = "**Total SR**: 0"
	return aim, speed, total
}
