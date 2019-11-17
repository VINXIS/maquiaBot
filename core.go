package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/signal"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
	"time"

	handlers "./handlers"
	gencommands "./handlers/general-commands"
	osutools "./osu-functions"
	structs "./structs"
	tools "./tools"

	"github.com/bwmarrin/discordgo"
)

func main() {
	discord, err := discordgo.New("Bot " + os.Getenv("DISCORD_BOT_TOKEN"))
	tools.ErrRead(err)

	// Handle farm data
	go osutools.FarmUpdate()

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

	// Resume all reminder timers
	reminders := []structs.Reminder{}
	_, err = os.Stat("./data/reminders.json")
	if err == nil {
		f, err := ioutil.ReadFile("./data/reminders.json")
		tools.ErrRead(err)
		_ = json.Unmarshal(f, &reminders)
	} else {
		tools.ErrRead(err)
	}
	reminderTimers := []structs.ReminderTimer{}
	for _, reminder := range reminders {
		reminderTimer := structs.ReminderTimer{
			Reminder: reminder,
			Timer:    *time.NewTimer(reminder.Target.Sub(time.Now().UTC())),
		}
		reminderTimers = append(reminderTimers, reminderTimer)
		go gencommands.RunReminder(discord, reminderTimer)
	}
	gencommands.ReminderTimers = reminderTimers

	// Get osu! tracking data for channels
	var channels []string

	err = filepath.Walk("./data/channelData", func(path string, info os.FileInfo, err error) error {
		tools.ErrRead(err)
		channels = append(channels, path)
		return nil
	})
	tools.ErrRead(err)
	for _, channel := range channels {
		if strings.HasSuffix(channel, ".json") {
			chID := strings.Replace(strings.Replace(channel, "data\\channelData\\", "", -1), ".json", "", -1)
			ch, err := discord.Channel(chID)
			if err == nil {
				go osutools.TrackPost(*ch, discord, mapCache)
			}
		}
	}

	// Get osu! mapper tracking data
	// go osutools.TrackMapperPost(discord) Commented until a solution is found for its issues

	// Create a channel to keep the bot running until a prompt is given to close
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Kill)
	<-sc

	// Close the Discord Session
	discord.Close()
}
