package pagination

import (
	"maquiaBot/structs"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
)

// Quotes handles the pagination of a list of quotes
func Quotes(s *discordgo.Session, r *discordgo.MessageReactionAdd, msg *discordgo.Message, serverData structs.ServerData, num, numend int) *discordgo.MessageEmbed {
	username := ""
	embed := &discordgo.MessageEmbed{}

	userQuotes := serverData.Quotes[num:numend]
	if strings.Contains(msg.Content, "Quotes for") {
		quoteRegex, _ := regexp.Compile(`(?i)Quotes for \*\*(.+)\*\*`)
		username = quoteRegex.FindStringSubmatch(msg.Content)[1]

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
			return embed
		}

		userQuotes = []discordgo.Message{}
		for _, quote := range serverData.Quotes {
			if quote.Author.ID == user.ID {
				userQuotes = append(userQuotes, quote)
			}
		}
		if len(userQuotes) == 0 {
			return embed
		}
		if numend > len(userQuotes) {
			numend = len(userQuotes)
		}
		if num >= numend {
			return embed
		}
		userQuotes = userQuotes[num:numend]
	}

	for i, quote := range userQuotes {
		if len(quote.Content) > 1024 {
			quote.Content = quote.Content[:1024]
		}

		quoteEmbed := &discordgo.MessageEmbedField{
			Name:   quote.ID + " - " + quote.Author.Username,
			Value:  quote.Content,
			Inline: true,
		}

		if username != "" {
			quoteEmbed.Name += " (" + strconv.Itoa(i+1+num) + ")"
		}
		if quote.Content == "" {
			quoteEmbed.Value = "**IMAGE/VIDEO QUOTE**"
		}

		embed.Fields = append(embed.Fields, quoteEmbed)

		if len(embed.Fields) == 25 {
			break
		}
	}

	return embed
}
