package gencommands

import (
	"regexp"
	"sort"
	"strconv"
	"strings"

	"maquiaBot/structs"
	tools "maquiaBot/tools"

	"github.com/bwmarrin/discordgo"
)

// Counters list out the counters enabled in the server
func Counters(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Get server
	server, err := s.Guild(m.GuildID)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "This is not a server!")
		return
	}
	serverImg := "https://cdn.discordapp.com/icons/" + server.ID + "/" + server.Icon
	if strings.Contains(server.Icon, "a_") {
		serverImg += ".gif"
	} else {
		serverImg += ".png"
	}

	serverData, _ := tools.GetServer(*server, s)
	embed := &discordgo.MessageEmbed{
		Author: &discordgo.MessageEmbedAuthor{
			Name:    server.Name,
			IconURL: serverImg,
		},
	}

	if len(serverData.Counters) == 0 {
		s.ChannelMessageSend(m.ChannelID, "There are no counters configuered for this server currently! Admins can see `help counter` for details on how to add counters.")
		return
	}

	for _, counter := range serverData.Counters {
		counter.Text = `(?i)` + counter.Text
		regex := false
		_, err := regexp.Compile(counter.Text)
		if err == nil {
			regex = true
		}
		valueText := "Counter: " + counter.Text + "\nRegex compatible: " + strconv.FormatBool(regex)
		if len(valueText) > 1024 {
			valueText = valueText[:1021] + "..."
		}
		embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
			Name:  strconv.FormatInt(counter.ID, 10),
			Value: valueText,
		})

		if len(embed.Fields) == 25 {
			if len(serverData.Counters) > 25 {
				embed.Footer = &discordgo.MessageEmbedFooter{
					Text: "Page 1",
				}
			}
			break
		}
	}
	msg, err := s.ChannelMessageSendComplex(m.ChannelID, &discordgo.MessageSend{
		Content: "Admins can use `counter -d` to delete a counter! If there are more than 25 counters, please use the reactions to go through pages!\nUsers can see counter rankings using `countrank`. See `help countrank` for details!",
		Embed:   embed,
	})
	if err != nil {
		return
	}
	if len(embed.Fields) == 25 && len(serverData.Counters) > 25 {
		_ = s.MessageReactionAdd(m.ChannelID, msg.ID, "➡️")
	}
	return
}

// CountRank shows the ranking for a counter
func CountRank(s *discordgo.Session, m *discordgo.MessageCreate) {
	countRankRegex, _ := regexp.Compile(`(?i)(cr|countrank|countr|crank)\s+(\d+)(\s+(\d+))?`)

	// Check for params
	if !countRankRegex.MatchString(m.Content) {
		s.ChannelMessageSend(m.ChannelID, "Please provide a counter's ID! See the counter IDs using `counters`")
		return
	}

	// Obtain ID
	id := strings.TrimSpace(countRankRegex.FindStringSubmatch(m.Content)[2])
	ID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, id+" is an invalid ID!")
		return
	}

	// Obtain count of people
	num := 10
	if len(countRankRegex.FindStringSubmatch(m.Content)) == 5 {
		num, err = strconv.Atoi(countRankRegex.FindStringSubmatch(m.Content)[4])
		if err != nil {
			num = 10
		}
	}

	// Get server
	server, err := s.Guild(m.GuildID)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "This is not a server!")
		return
	}

	serverData, _ := tools.GetServer(*server, s)
	if len(serverData.Counters) == 0 {
		s.ChannelMessageSend(m.ChannelID, "This server has no counters!")
		return
	}

	// Find counter
	var counter structs.Counter
	for _, serverCounter := range serverData.Counters {
		if serverCounter.ID == ID {
			counter = serverCounter
			break
		}
	}
	if counter.ID == 0 {
		s.ChannelMessageSend(m.ChannelID, "No counter with the ID "+id+" found!")
		return
	}

	// Order from most times to least
	sort.Slice(counter.Users, func(i, j int) bool {
		return counter.Users[i].Count > counter.Users[j].Count
	})

	if num > len(counter.Users) {
		num = len(counter.Users)
	}

	text := "Amount of times the top " + strconv.Itoa(num) + " users said `" + counter.Text + "`:\n"
	counter.Users = counter.Users[:num]
	for i, user := range counter.Users {
		text += "#" + strconv.Itoa(i+1) + ": " + user.Username + " - " + strconv.Itoa(user.Count) + " times\n"
	}
	s.ChannelMessageSend(m.ChannelID, text)
	return
}
