package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"

	handlers "./handlers"
	osutools "./osu-functions"
	structs "./structs"
	tools "./tools"

	"github.com/bwmarrin/discordgo"
)

func main() {
	_, err := exec.Command("dotnet", "build", "./osu-tools/PerformanceCalculator").Output()
	tools.ErrRead(err)

	discord, err := discordgo.New("Bot " + os.Getenv("DISCORD_BOT_TOKEN"))
	tools.ErrRead(err)

	// Obtain map cache data
	mapCache := []structs.MapData{}
	f, err := ioutil.ReadFile("./data/osuData/mapCache.json")
	tools.ErrRead(err)
	_ = json.Unmarshal(f, &mapCache)

	// Register the messageCreate func as a callback for MessageCreate events
	discord.AddHandler(handlers.MessageHandler)

	// Open a websocket connection to Discord and begin listening
	err = discord.Open()
	tools.ErrRead(err)
	fmt.Println("Bot is now running in " + strconv.Itoa(len(discord.State.Guilds)) + " servers.")

	var channels []string

	err = filepath.Walk("./data/channelData", func(path string, info os.FileInfo, err error) error {
		tools.ErrRead(err)
		channels = append(channels, path)
		return nil
	})
	tools.ErrRead(err)
	for _, channel := range channels {
		if strings.HasSuffix(channel, ".json") {
			go osutools.TrackPost(channel, discord, mapCache)
		}
	}

	// Create a channel to keep the bot running until a prompt is given to close
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Kill)
	<-sc

	// Close the Discord Session
	discord.Close()
}
