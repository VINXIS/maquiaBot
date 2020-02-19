package tools

import (
	"fmt"
	"log"
	"runtime"

	config "../config"
	"github.com/bwmarrin/discordgo"
)

// ErrRead will check to see if there is an error; it will print the error and kill the bot if there is any
func ErrRead(s *discordgo.Session, err error) {
	if err != nil {
		pc, fn, line, _ := runtime.Caller(1)
		dm, err := s.UserChannelCreate(config.Conf.BotHoster.UserID)
		if err != nil {
			log.Printf("[error] in %s[%s:%d] %v\n", runtime.FuncForPC(pc).Name(), fn, line, err)
			return
		}
		s.ChannelMessageSend(dm.ID, fmt.Sprintf("[error] in %s[%s:%d] %v\n", runtime.FuncForPC(pc).Name(), fn, line, err))
	}
}
