package main

import (
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	handlers "./handlers"
	tools "./tools"

	"github.com/bwmarrin/discordgo"
)

func main() {
	discord, err := discordgo.New("Bot " + os.Getenv("DISCORD_BOT_TOKEN"))
	tools.ErrRead(err, "17", "core.go")

	// Register the messageCreate func as a callback for MessageCreate events
	discord.AddHandler(handlers.MessageHandler)

	// Open a websocket connection to Discord and begin listening
	err = discord.Open()
	tools.ErrRead(err, "24", "core.go")
	fmt.Println("Bot is now running in " + strconv.Itoa(len(discord.State.Guilds)) + " servers.")

	// Create a channel to keep the bot running until a prompt is given to close
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Kill)
	<-sc

	// Close the Discord Session
	discord.Close()
}
