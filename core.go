package main

import (
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"strconv"
	"syscall"

	handlers "./handlers"
	tools "./tools"

	"github.com/bwmarrin/discordgo"
)

func main() {
	discord, err := discordgo.New("Bot " + os.Getenv("DISCORD_BOT_TOKEN"))
	tools.ErrRead(err)

	// Register the messageCreate func as a callback for MessageCreate events
	discord.AddHandler(handlers.MessageHandler)

	// Open a websocket connection to Discord and begin listening
	err = discord.Open()
	tools.ErrRead(err)
	fmt.Println("Bot is now running in " + strconv.Itoa(len(discord.State.Guilds)) + " servers.")

	var servers []string

	err = filepath.Walk("./data/channelData", func(path string, info os.FileInfo, err error) error {
		tools.ErrRead(err)
		servers = append(servers, path)
		return nil
	})
	tools.ErrRead(err)
	for _, server := range servers {
		fmt.Println(server)
	}

	// Create a channel to keep the bot running until a prompt is given to close
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Kill)
	<-sc

	// Close the Discord Session
	discord.Close()
}
