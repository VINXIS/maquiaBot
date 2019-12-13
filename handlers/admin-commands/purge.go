package admincommands

import (
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

	// Get username(s) and number of messages
	userRegex, _ := regexp.Compile(`purge\s+(.+)`)

	userText := ""
	num := 4
	var usernames []string
	if userRegex.MatchString(m.Content) {
		userNum := userRegex.FindStringSubmatch(strings.Replace(m.ContentWithMentionsReplaced(), "@", "", -1))[1]
		args := strings.Split(userNum, " ")
		for _, arg := range args {
			if i, err := strconv.Atoi(arg); err == nil && i > 0 {
				userNum = strings.TrimSpace(strings.Replace(userNum, arg, "", 1))
				num = i + 1
				break
			}
		}
		if userNum != "" {
			usernames = append(usernames, strings.Split(userNum, " ")...)
		}
	}
	if len(usernames) != 0 {
		for _, username := range usernames {
			userText += "**" + username + "** "
		}
	}

	// Get messages
	messages, err := s.ChannelMessages(m.ChannelID, -1, "", "", "")
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "Error obtaining messages!")
		return
	}

	// Get messages and delete them
	var messageIDs []string
	lastID := ""
	prevLength := 0
	recurring := 0
	for {
		messages, err = s.ChannelMessages(m.ChannelID, -1, lastID, "", "")
		if err != nil {
			break
		}
		for _, msg := range messages {
			if len(usernames) == 0 {
				messageIDs = append(messageIDs, msg.ID)
			} else {
				for _, username := range usernames {
					if strings.HasPrefix(strings.ToLower(msg.Author.Username), strings.ToLower(username)) || (msg.Member != nil && strings.HasPrefix(strings.ToLower(msg.Member.Nick), strings.ToLower(username))) {
						messageIDs = append(messageIDs, msg.ID)
						break
					}
				}
			}
			if len(messageIDs) == num {
				break
			}
			lastID = msg.ID
		}
		if len(messageIDs) == num {
			break
		}
		if prevLength == len(messageIDs) {
			recurring++
		} else {
			prevLength = len(messageIDs)
			recurring = 1
		}
		if recurring == 5 {
			num = len(messageIDs)
			break
		}
	}

	if len(messageIDs) == 0 {
		s.ChannelMessageSend(m.ChannelID, "No messages found with the given usernames: "+userText)
		return
	}

	if len(messageIDs) > 100 {
		i := 100
		for {
			if i >= len(messageIDs) {
				i = len(messageIDs) - 1
			}
			err = s.ChannelMessagesBulkDelete(m.ChannelID, messageIDs[i-100:i])
			if err != nil {
				s.ChannelMessageSend(m.ChannelID, "Could not delete messages! Please make sure I have the proper permissions!")
				return
			}
			if i == len(messageIDs)-1 {
				break
			}
			i += 100
		}
	}

	// Send confirmation message and then delete it after
	msg, err := s.ChannelMessageSend(m.ChannelID, "Removed the latest "+strconv.Itoa(num-1)+" messages from the following people: "+userText)
	if err != nil {
		return
	}
	timer := time.NewTimer(5 * time.Second)
	<-timer.C
	s.ChannelMessageDelete(msg.ChannelID, msg.ID)
}
