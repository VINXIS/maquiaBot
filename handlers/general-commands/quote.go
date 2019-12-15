package gencommands

import (
	"crypto/rand"
	"encoding/json"
	"io/ioutil"
	"math/big"
	"net/http"
	"regexp"
	"sort"
	"strings"
	"time"

	tools "../../tools"
	"github.com/bwmarrin/discordgo"
)

// Quote lets you get a quote from someone
func Quote(s *discordgo.Session, m *discordgo.MessageCreate) {
	quoteRegex, _ := regexp.Compile(`q(uote)?\s+(.+)`)
	linkRegex, _ := regexp.Compile(`https?:\/\/\S+`)
	extensionRegex, _ := regexp.Compile(`\.(\S{3,4})`)

	// Get server
	server, err := s.Guild(m.GuildID)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "This is not a server!")
		return
	}
	serverData := tools.GetServer(*server)
	if len(serverData.Quotes) == 0 {
		s.ChannelMessageSend(m.ChannelID, "No quotes saved for this server! Please use `quoteadd` to add a quote!")
		return
	}

	// Get username
	user := &discordgo.User{}
	username := ""
	userQuotes := serverData.Quotes
	if quoteRegex.MatchString(m.Content) {
		username = quoteRegex.FindStringSubmatch(m.Content)[2]

		// Get user
		members, _ := s.GuildMembers(m.GuildID, "", 1000)
		sort.Slice(members, func(i, j int) bool {
			time1, _ := members[i].JoinedAt.Parse()
			time2, _ := members[j].JoinedAt.Parse()
			return time1.Unix() < time2.Unix()
		})
		for _, member := range members {
			if strings.HasPrefix(strings.ToLower(member.User.Username), username) || strings.HasPrefix(strings.ToLower(member.Nick), username) {
				user, _ = s.User(member.User.ID)
				break
			}
		}

		if user.ID == "" {
			s.ChannelMessageSend(m.ChannelID, "No user with the name **"+username+"** found!")
		}

		userQuotes = []discordgo.Message{}
		for _, quote := range serverData.Quotes {
			if quote.Author.ID == user.ID {
				userQuotes = append(userQuotes, quote)
			}
		}
		if len(userQuotes) == 0 {
			s.ChannelMessageSend(m.ChannelID, "No quotes saved for **"+user.Username+"**! Please use `quoteadd` to add a quote!")
			return
		}
	}

	roll, _ := rand.Int(rand.Reader, big.NewInt(int64(len(userQuotes))))
	quote := userQuotes[roll.Int64()]
	timestamp, _ := quote.Timestamp.Parse()
	timestampString := strings.Replace(timestamp.Format(time.RFC822Z), "+0000", "UTC", -1)
	if user.ID == "" {
		user = quote.Author
	}

	embed := &discordgo.MessageEmbed{
		Author: &discordgo.MessageEmbedAuthor{
			URL:     "https://discordapp.com/channels/" + server.ID + "/" + quote.ChannelID + "/" + quote.ID,
			Name:    user.Username,
			IconURL: user.AvatarURL(""),
		},
		Description: quote.Content,
		Footer: &discordgo.MessageEmbedFooter{
			Text: timestampString,
		},
	}
	link := ""
	name := ""
	if len(quote.Attachments) != 0 {
		if !(strings.HasSuffix(quote.Attachments[0].URL, "png") || strings.HasSuffix(quote.Attachments[0].URL, "jpg") || strings.HasSuffix(quote.Attachments[0].URL, "gif")) {
			link = quote.Attachments[0].URL
			name = quote.Attachments[0].Filename
		}
		embed.Image = &discordgo.MessageEmbedImage{
			URL: quote.Attachments[0].URL,
		}
	} else if linkRegex.MatchString(quote.Content) {
		if !(strings.HasSuffix(linkRegex.FindStringSubmatch(quote.Content)[0], "png") || strings.HasSuffix(linkRegex.FindStringSubmatch(quote.Content)[0], "jpg") || strings.HasSuffix(linkRegex.FindStringSubmatch(quote.Content)[0], "gif")) {
			link = linkRegex.FindStringSubmatch(quote.Content)[0]
			res := extensionRegex.FindAllStringSubmatch(quote.Content, -1)
			name = "video." + res[len(res)-1][0]
		}
		embed.Image = &discordgo.MessageEmbedImage{
			URL: linkRegex.FindStringSubmatch(quote.Content)[0],
		}
	}
	if link != "" {
		response, err := http.Get(link)
		if err == nil {
			s.ChannelMessageSendComplex(m.ChannelID, &discordgo.MessageSend{
				Embed: embed,
				File: &discordgo.File{
					Name:   name,
					Reader: response.Body,
				},
			})
			return
		}
	}
	s.ChannelMessageSendEmbed(m.ChannelID, embed)
	return
}

// QuoteAdd lets you add quotes
func QuoteAdd(s *discordgo.Session, m *discordgo.MessageCreate) {
	quoteAddRegex, _ := regexp.Compile(`q(uote)?a(dd)?\s+(.+)`)
	randomRegex, _ := regexp.Compile(`-r`)
	channelRegex, _ := regexp.Compile(`https://discordapp.com/channels\/(\d+)\/(\d+)\/(\d+)`)

	// Get server
	server, err := s.Guild(m.GuildID)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "This is not a server!")
		return
	}
	serverData := tools.GetServer(*server)

	// Get message
	message := &discordgo.Message{}
	username := ""
	msgs, err := s.ChannelMessages(m.ChannelID, -1, m.ID, "", "")
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "Ur really funny mate")
		return
	}
	if channelRegex.MatchString(m.Content) {
		if channelRegex.FindStringSubmatch(m.Content)[1] != server.ID {
			s.ChannelMessageSend(m.ChannelID, "Please do not quote from other servers!")
			return
		}
		message, _ = s.ChannelMessage(channelRegex.FindStringSubmatch(m.Content)[2], channelRegex.FindStringSubmatch(m.Content)[3])
	} else if quoteAddRegex.MatchString(m.Content) {
		username = quoteAddRegex.FindStringSubmatch(m.Content)[3]
		message, err = s.ChannelMessage(m.ChannelID, username)
		if err != nil {
			random := false
			if randomRegex.MatchString(m.Content) {
				random = true
				username = strings.TrimSpace(strings.Replace(username, "-r", "", -1))
			}
			if random {
				msgList := []*discordgo.Message{}
				if username != "" {
					for _, msg := range msgs {
						if strings.HasPrefix(strings.ToLower(msg.Author.Username), username) || (msg.Member != nil && strings.HasPrefix(strings.ToLower(msg.Member.Nick), username)) {
							msgList = append(msgList, msg)
						}
					}
				} else {
					msgList = msgs
				}
				if len(msgList) == 0 {
					s.ChannelMessageSend(m.ChannelID, "No message found for the user **"+username+"** to randomly choose from!")
					return
				}
				roll, _ := rand.Int(rand.Reader, big.NewInt(int64(len(msgList))))
				message = msgList[roll.Int64()]
			} else {
				for _, msg := range msgs {
					if strings.HasPrefix(strings.ToLower(msg.Author.Username), username) || (msg.Member != nil && strings.HasPrefix(strings.ToLower(msg.Member.Nick), username)) {
						message = msg
						break
					}
				}
			}
		}
	} else {
		for _, msg := range msgs {
			if msg.Author.ID != m.Author.ID && msg.Author.ID != s.State.User.ID {
				message = msg
				break
			}
		}
	}

	if message == nil || message.ID == "" || (message.Content == "" && len(message.Attachments) == 0) {
		s.ChannelMessageSend(m.ChannelID, "No message found!")
		return
	}

	err = serverData.AddQuote(message)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "Quote: `"+message.Content+"` already exists for **"+message.Author.Username+"**!")
		return
	}

	serverData.Time = time.Now()

	jsonCache, err := json.Marshal(serverData)
	tools.ErrRead(err)

	err = ioutil.WriteFile("./data/serverData/"+m.GuildID+".json", jsonCache, 0644)
	tools.ErrRead(err)

	if message.Content != "" {
		s.ChannelMessageSend(m.ChannelID, "Quote: `"+message.Content+"` added for **"+message.Author.Username+"**!")
	} else if len(message.Attachments) != 0 {
		s.ChannelMessageSend(m.ChannelID, "Image / video quote added for **"+message.Author.Username+"**!")
	}
	return
}

// QuoteRemove lets you remove quotes
func QuoteRemove(s *discordgo.Session, m *discordgo.MessageCreate) {
	quoteRemoveRegex, _ := regexp.Compile(`q(uote)?(r(emove)?|d(elete)?)\s+(\d+)`)
	channelRegex, _ := regexp.Compile(`https://discordapp.com/channels\/\d+\/(\d+)\/(\d+)`)

	// Get server
	server, err := s.Guild(m.GuildID)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "This is not a server!")
		return
	}
	serverData := tools.GetServer(*server)

	mID := ""
	if quoteRemoveRegex.MatchString(m.Content) {
		mID = quoteRemoveRegex.FindStringSubmatch(m.Content)[5]
	} else if channelRegex.MatchString(m.Content) {
		mID = channelRegex.FindStringSubmatch(m.Content)[2]
	} else {
		s.ChannelMessageSend(m.ChannelID, "Please link a message / give a message ID to remove!")
		return
	}

	err = serverData.RemoveQuote(mID)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "No message found with that ID! Use `quotes` to see the list of quotes.")
		return
	}

	serverData.Time = time.Now()

	jsonCache, err := json.Marshal(serverData)
	tools.ErrRead(err)

	err = ioutil.WriteFile("./data/serverData/"+m.GuildID+".json", jsonCache, 0644)
	tools.ErrRead(err)

	s.ChannelMessageSend(m.ChannelID, "Quote removed!")
	return
}

// Quotes lets you see all the quotes of a user
func Quotes(s *discordgo.Session, m *discordgo.MessageCreate) {
	quotesRegex, _ := regexp.Compile(`q(uote)?s\s+(.+)`)

	// Get server
	server, err := s.Guild(m.GuildID)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "This is not a server!")
		return
	}
	serverData := tools.GetServer(*server)
	serverImg := "https://cdn.discordapp.com/icons/" + server.ID + "/" + server.Icon
	if strings.Contains(server.Icon, "a_") {
		serverImg += ".gif"
	} else {
		serverImg += ".png"
	}

	// Get username if there is any
	username := ""
	userQuotes := serverData.Quotes
	embed := &discordgo.MessageEmbed{
		Author: &discordgo.MessageEmbedAuthor{
			Name:    server.Name,
			IconURL: serverImg,
		},
	}
	if quotesRegex.MatchString(m.Content) {
		username = quotesRegex.FindStringSubmatch(m.Content)[2]

		// Get user
		user := &discordgo.User{}
		members, _ := s.GuildMembers(m.GuildID, "", 1000)
		sort.Slice(members, func(i, j int) bool {
			time1, _ := members[i].JoinedAt.Parse()
			time2, _ := members[j].JoinedAt.Parse()
			return time1.Unix() < time2.Unix()
		})
		for _, member := range members {
			if strings.HasPrefix(strings.ToLower(member.User.Username), username) || strings.HasPrefix(strings.ToLower(member.Nick), username) {
				user, _ = s.User(member.User.ID)
				break
			}
		}

		if user.ID == "" {
			s.ChannelMessageSend(m.ChannelID, "No user with the name **"+username+"** found!")
			return
		}

		userQuotes = []discordgo.Message{}
		for _, quote := range serverData.Quotes {
			if quote.Author.ID == user.ID {
				userQuotes = append(userQuotes, quote)
			}
		}
		if len(userQuotes) == 0 {
			s.ChannelMessageSend(m.ChannelID, "No quotes saved for **"+user.Username+"**! Please use `quoteadd` to add a quote!")
			return
		}
	}

	for _, quote := range userQuotes {
		if len(quote.Content) > 1024 {
			quote.Content = quote.Content[:1024]
		}
		if quote.Content != "" {
			embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
				Name:   quote.ID + " - " + quote.Author.Username,
				Value:  quote.Content,
				Inline: true,
			})
		} else {
			embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
				Name:   quote.ID + " - " + quote.Author.Username,
				Value:  "**IMAGE/VIDEO QUOTE**",
				Inline: true,
			})
		}
		if len(embed.Fields) == 25 {
			if len(userQuotes) > 25 {
				embed.Footer = &discordgo.MessageEmbedFooter{
					Text: "Page 1",
				}
			}
			break
		}
	}

	text := "Use `quotedelete <ID>` to delete a quote! If there are more than 25 quotes, please use the reactions to go through pages!"
	if username != "" {
		text += "\nQuotes for **" + userQuotes[0].Author.Username + "**:"
	}

	msg, err := s.ChannelMessageSendComplex(m.ChannelID, &discordgo.MessageSend{
		Content: text,
		Embed:   embed,
	})
	if err != nil {
		return
	}
	if len(embed.Fields) == 25 && len(userQuotes) > 25 {
		_ = s.MessageReactionAdd(m.ChannelID, msg.ID, "⬆️")
	}
}
