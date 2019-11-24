package admincommands

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	tools "../../tools"
	"github.com/bwmarrin/discordgo"
)

// Purge lets admins purge messages including their purge command
func Purge(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Check if server exists
	server, err := s.Guild(m.GuildID)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "This is not a server so custom prefixes are unavailable! Please use `$` instead for commands!")
		return
	}

	if !tools.AdminCheck(s, m, *server) {
		s.ChannelMessageSend(m.ChannelID, "You must be an admin, server manager, or server owner!")
		return
	}

	// Get messages
	messages, err := s.ChannelMessages(m.ChannelID, -1, "", "", "")
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "Error obtaining messages!")
		return
	}

	// Get username(s) and number of messages
	userRegex, _ := regexp.Compile(`purge\s+(.+)`)

	users := m.Mentions
	userText := "!"
	num := 4
	var usernames []string
	if len(users) > 0 {
		for _, user := range users {
			usernames = append(usernames, user.Username)
		}
	}
	if userRegex.MatchString(m.Content) {
		userNum := userRegex.FindStringSubmatch(m.ContentWithMentionsReplaced())[1]
		args := strings.Split(userNum, " ")
		for _, arg := range args {
			if i, err := strconv.Atoi(arg); err == nil && i > 0 && i <= 100 {
				userNum = strings.TrimSpace(strings.Replace(userNum, arg, "", 1))
				num = i + 1
				break
			}
		}
		usernames = append(usernames, strings.Split(userNum, " ")...)
	}
	fmt.Println(usernames)
	fmt.Println(num)
	if len(usernames) != 0 {
		userText = " from the following people: "
		for _, username := range usernames {
			userText += "**" + username + "** "
		}
	}

	// Get messages and delete them
	var messageIDs []string
	for _, msg := range messages {
		if len(usernames) == 0 {
			messageIDs = append(messageIDs, msg.ID)
		} else {
			for _, username := range usernames {
				if strings.HasPrefix(strings.ToLower(msg.Author.Username), username) || strings.HasPrefix(strings.ToLower(msg.Author.Username), strings.Replace(username, "@", "", -1)) {
					messageIDs = append(messageIDs, msg.ID)
					break
				}
			}
		}
		if len(messageIDs) == num {
			break
		}
	}

	if len(messageIDs) == 0 {
		s.ChannelMessageSend(m.ChannelID, "No messages found with the given usernames!")
		return
	}

	err = s.ChannelMessagesBulkDelete(m.ChannelID, messageIDs)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "Could not delete messages! Please make sure I have the proper permissions!")
		return
	}

	// Send confirmation message and then delete it after
	msg, _ := s.ChannelMessageSend(m.ChannelID, "Removed the "+strconv.Itoa(num-1)+" latest messages"+userText)
	timer := time.NewTimer(5 * time.Second)
	<-timer.C
	s.ChannelMessageDelete(msg.ChannelID, msg.ID)
}
