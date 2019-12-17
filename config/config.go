package config

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

// Conf is the configuration
var Conf Config

// Config is the main part of the config
type Config struct {
	BotHoster    BotHoster
	Database     Database
	Server       string
	DiscordToken string
	OsuToken     string
}

// BotHoster holds info about who is hosting the bot (also known as Bot Creator)
type BotHoster struct {
	Username string
	UserID   string
}

// Database holds info about database login
type Database struct {
	Username string
	Password string
	Name     string
}

// NewConfig creates the new configuration from the JSON file
func NewConfig() {
	config := Config{}
	f, err := ioutil.ReadFile("./config/config.json")
	if err != nil {
		log.Fatalln("Error obtaining config information: " + err.Error())
	}
	_ = json.Unmarshal(f, &config)
	Conf = config
}
