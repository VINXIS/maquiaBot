package handlers

import (
	paginationcommands "maquiaBot/handlers/pagination-commands"
	"maquiaBot/tools"
	"regexp"
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
)

// ReactAdd is to deal with reacts added
func ReactAdd(s *discordgo.Session, r *discordgo.MessageReactionAdd) {
	if r.UserID == s.State.User.ID {
		return
	}

	msg, err := s.ChannelMessage(r.MessageReaction.ChannelID, r.MessageReaction.MessageID)
	if err != nil || msg.Author.ID != s.State.User.ID || len(msg.Embeds) == 0 || msg.Embeds[0].Footer == nil || !strings.Contains(msg.Embeds[0].Footer.Text, "Page") {
		return
	}

	regex, _ := regexp.Compile(`(?i)Page (\d+)`)
	num, _ := strconv.Atoi(regex.FindStringSubmatch(msg.Embeds[0].Footer.Text)[1])
	numend := (num + 1) * 25
	page := strconv.Itoa(num + 1)
	if r.Emoji.Name == "⬅️" && num > 1 {
		num--
		numend = num * 25
		page = strconv.Itoa(num)
		num--
	} else if r.Emoji.Name != "➡️" {
		return
	}
	num *= 25

	// Get server
	server, err := s.Guild(r.MessageReaction.GuildID)
	if err != nil {
		return
	}
	serverData, _ := tools.GetServer(*server, s)

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

	embed := &discordgo.MessageEmbed{}
	lastPage := true

	// Check which pagination this is for
	if strings.Contains(msg.Content, "quote") {
		embed, lastPage = paginationcommands.Quotes(s, r, msg, serverData, num, numend)
	} else if strings.Contains(msg.Content, "trigger") {
		embed, lastPage = paginationcommands.Triggers(s, r, msg, serverData, num, numend)
	} else if strings.Contains(msg.Content, "counter") {
		embed, lastPage = paginationcommands.Counters(s, r, msg, serverData, num, numend)
	}

	if len(embed.Fields) == 0 {
		return
	}

	embed.Footer = &discordgo.MessageEmbedFooter{
		Text: "Page " + page,
	}

	s.MessageReactionsRemoveAll(msg.ChannelID, msg.ID)

	msg, err = s.ChannelMessageEditEmbed(r.MessageReaction.ChannelID, r.MessageReaction.MessageID, embed)
	if err != nil {
		return
	}

	if num/25 != 0 {
		_ = s.MessageReactionAdd(msg.ChannelID, msg.ID, "⬅️")
	}
	if !lastPage {
		_ = s.MessageReactionAdd(msg.ChannelID, msg.ID, "➡️")
	}
}
