package handlers

import (
	"regexp"
	"sort"
	"strconv"
	"strings"

	tools "../tools"
	"github.com/bwmarrin/discordgo"
)

// ReactAdd is to deal with reacts added
func ReactAdd(s *discordgo.Session, r *discordgo.MessageReactionAdd) {
	if r.UserID == s.State.User.ID {
		return
	}

	msg, err := s.ChannelMessage(r.MessageReaction.ChannelID, r.MessageReaction.MessageID)
	if err != nil || msg.Author.ID != s.State.User.ID {
		return
	}

	if len(msg.Embeds) > 0 && msg.Embeds[0].Footer != nil && strings.Contains(msg.Embeds[0].Footer.Text, "Page") {
		regex, _ := regexp.Compile(`Page (\d+)`)
		num, _ := strconv.Atoi(regex.FindStringSubmatch(msg.Embeds[0].Footer.Text)[1])
		numend := (num + 1) * 25
		page := strconv.Itoa(num + 1)
		if r.Emoji.Name == "⬇️" && num > 1 {
			num--
			numend = num * 25
			page = strconv.Itoa(num)
			num--
		} else if r.Emoji.Name != "⬆️" {
			return
		}
		num *= 25

		// Get server
		server, err := s.Guild(r.MessageReaction.GuildID)
		if err != nil {
			return
		}
		serverData := tools.GetServer(*server, s)

		// Check if num or numend is Fucked
		if num < 0 || num >= len(serverData.Quotes)-1 {
			return
		}
		if numend > len(serverData.Quotes) {
			numend = len(serverData.Quotes)
		}
		if num >= numend {
			return
		}
		quoteNum := len(serverData.Quotes)
		userQuotes := serverData.Quotes[num:numend]
		if strings.Contains(msg.Content, "Quotes for") {
			quoteRegex, _ := regexp.Compile(`Quotes for \*\*(.+)\*\*`)
			username := quoteRegex.FindStringSubmatch(msg.Content)[1]

			user := &discordgo.User{}
			members, _ := s.GuildMembers(r.MessageReaction.GuildID, "", 1000)
			sort.Slice(members, func(i, j int) bool {
				time1, _ := members[i].JoinedAt.Parse()
				time2, _ := members[j].JoinedAt.Parse()
				return time1.Unix() < time2.Unix()
			})
			for _, member := range members {
				if strings.Contains(strings.ToLower(member.User.Username), strings.ToLower(username)) || strings.Contains(strings.ToLower(member.Nick), strings.ToLower(username)) {
					user, _ = s.User(member.User.ID)
					break
				}
			}

			if user.ID == "" {
				return
			}

			userQuotes = []discordgo.Message{}
			for _, quote := range serverData.Quotes {
				if quote.Author.ID == user.ID {
					userQuotes = append(userQuotes, quote)
				}
			}
			if len(userQuotes) == 0 {
				return
			}
			if numend > len(userQuotes) {
				numend = len(userQuotes)
			}
			if num >= numend {
				return
			}
			quoteNum = len(userQuotes)
			userQuotes = userQuotes[num:numend]
		}

		embed := &discordgo.MessageEmbed{}
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
				break
			}
		}
		embed.Footer = &discordgo.MessageEmbedFooter{
			Text: "Page " + page,
		}

		s.MessageReactionsRemoveAll(msg.ChannelID, msg.ID)

		msg, err := s.ChannelMessageEditEmbed(r.MessageReaction.ChannelID, r.MessageReaction.MessageID, embed)
		if err != nil {
			return
		}

		if page != "1" {
			_ = s.MessageReactionAdd(msg.ChannelID, msg.ID, "⬇️")
		}
		if numend != quoteNum {
			_ = s.MessageReactionAdd(msg.ChannelID, msg.ID, "⬆️")
		}
	}
}
