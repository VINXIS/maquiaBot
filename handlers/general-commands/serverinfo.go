package gencommands

import (
	"strconv"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

// ServerInfo gives information about the server
func ServerInfo(s *discordgo.Session, m *discordgo.MessageCreate) {

	server, err := s.Guild(m.GuildID)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "This is not a server!")
		return
	}
	// Created at date
	createdAt, err := discordgo.SnowflakeTimestamp(server.ID)

	// Server owner
	owner, _ := s.User(server.OwnerID)

	// Channel info
	channels := server.Channels
	text := 0
	voice := 0
	category := 0
	other := 0
	for _, ch := range channels {
		switch ch.Type {
		case discordgo.ChannelTypeGuildText:
			text++
		case discordgo.ChannelTypeGuildVoice:
			voice++
		case discordgo.ChannelTypeGuildCategory:
			category++
		default:
			other++
		}
	}
	channelInfo := ""
	if text != 0 {
		channelInfo += strconv.Itoa(text) + " text\n"
	}
	if voice != 0 {
		channelInfo += strconv.Itoa(voice) + " voice\n"
	}
	if category != 0 {
		if category == 1 {
			channelInfo += strconv.Itoa(category) + " category\n"
		} else {
			channelInfo += strconv.Itoa(category) + " categories\n"
		}
	}
	if other != 0 {
		channelInfo += strconv.Itoa(other) + " other\n"
	}

	// AFK Timeout
	timeout := server.AfkTimeout
	timeoutString := "N/A"
	if timeout == int(time.Hour.Seconds()) {
		timeoutString = "1 hour"
	} else {
		timeoutString = strconv.Itoa(timeout/60) + " min"
	}

	// AFK Channel
	afkCh, err := s.Channel(server.AfkChannelID)
	afkChName := "N/A"
	if err == nil {
		afkChName = afkCh.Name
	}

	serverImg := "https://cdn.discordapp.com/icons/" + server.ID + "/" + server.Icon
	if strings.Contains(server.Icon, "a_") {
		serverImg += ".gif"
	} else {
		serverImg += ".png"
	}

	s.ChannelMessageSendEmbed(m.ChannelID, &discordgo.MessageEmbed{
		Author: &discordgo.MessageEmbedAuthor{
			Name:    server.Name,
			IconURL: serverImg,
		},
		Thumbnail: &discordgo.MessageEmbedThumbnail{
			URL: serverImg,
		},
		Fields: []*discordgo.MessageEmbedField{
			&discordgo.MessageEmbedField{
				Name:   "ID",
				Value:  server.ID,
				Inline: true,
			},
			&discordgo.MessageEmbedField{
				Name:   "Region",
				Value:  server.Region,
				Inline: true,
			},
			&discordgo.MessageEmbedField{
				Name:   "Server Created",
				Value:  createdAt.UTC().Format(time.RFC822Z),
				Inline: true,
			},
			&discordgo.MessageEmbedField{
				Name:   "Server Owner",
				Value:  owner.Username + "#" + owner.Discriminator,
				Inline: true,
			},
			&discordgo.MessageEmbedField{
				Name:   "AFK Timeout",
				Value:  timeoutString,
				Inline: true,
			},
			&discordgo.MessageEmbedField{
				Name:   "AFK Channel",
				Value:  afkChName,
				Inline: true,
			},
			&discordgo.MessageEmbedField{
				Name:   "Emoji Count",
				Value:  strconv.Itoa(len(server.Emojis)),
				Inline: true,
			},
			&discordgo.MessageEmbedField{
				Name:   "Member Count",
				Value:  strconv.Itoa(len(server.Members)),
				Inline: true,
			},
			&discordgo.MessageEmbedField{
				Name:   "Role Count",
				Value:  strconv.Itoa(len(server.Roles)),
				Inline: true,
			},
			&discordgo.MessageEmbedField{
				Name:   "Channels (" + strconv.Itoa(len(server.Channels)) + ")",
				Value:  channelInfo,
				Inline: true,
			},
		},
	})
}
