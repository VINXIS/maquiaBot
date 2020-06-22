package gencommands

import (
	"net/http"
	"regexp"
	"strconv"

	"github.com/ChimeraCoder/anaconda"
	"github.com/bwmarrin/discordgo"
	config "maquiaBot/config"
)

// Twitter uploads a twitter gif / image / video onto discord directly
func Twitter(s *discordgo.Session, m *discordgo.MessageCreate) {
	var (
		link     string
		linkType string
		ID       int64
	)
	twitterRegex, _ := regexp.Compile(`(?i)twitter.com\/(\S+)\/status\/(\d+)`)

	// Get ID / link
	if twitterRegex.MatchString(m.Content) {
		linkType = "mp4"
		ID, _ = strconv.ParseInt(twitterRegex.FindStringSubmatch(m.Content)[2], 10, 64)
	} else {
		messages, _ := s.ChannelMessages(m.ChannelID, -1, "", "", "")
		for _, msg := range messages {
			if twitterRegex.MatchString(msg.Content) && len(msg.Embeds) > 0 {
				ID, _ = strconv.ParseInt(twitterRegex.FindStringSubmatch(msg.Content)[2], 10, 64)
				if msg.Embeds[0].Video != nil {
					linkType = "mp4"
					break
				} else if msg.Embeds[0].Image != nil {
					link = msg.Embeds[0].Image.URL
					linkType = "png"
					break
				}
			}
		}
	}
	if linkType == "" {
		s.ChannelMessageSend(m.ChannelID, "No twitter video / image found!")
		return
	}

	msg, err := s.ChannelMessageSend(m.ChannelID, "Obtaining twitter video / image...")
	if err != nil {
		return
	}

	// Check if image or gif / video
	if linkType == "png" {
		response, err := http.Get(link)
		if err != nil {
			s.ChannelMessageDelete(msg.ChannelID, msg.ID)
			s.ChannelMessageSend(m.ChannelID, "Error obtaining image!")
			return
		}
		message := &discordgo.MessageSend{
			File: &discordgo.File{
				Name:   strconv.FormatInt(ID, 10) + "." + linkType,
				Reader: response.Body,
			},
		}
		s.ChannelMessageSendComplex(m.ChannelID, message)
		s.ChannelMessageDelete(msg.ChannelID, msg.ID)
	} else {
		// API call
		api := anaconda.NewTwitterApiWithCredentials(
			config.Conf.Twitter.Token,
			config.Conf.Twitter.Secret,
			config.Conf.Twitter.ConsumerToken,
			config.Conf.Twitter.ConsumerSecret,
		)
		tweet, err := api.GetTweet(ID, nil)
		if err != nil {
			s.ChannelMessageDelete(msg.ChannelID, msg.ID)
			s.ChannelMessageSend(m.ChannelID, "Error obtaining tweet!")
			return
		}

		// Find video link
		for _, media := range tweet.ExtendedEntities.Media {
			if media.Type == "video" || media.Type == "animated_gif" {
				for _, variant := range media.VideoInfo.Variants {
					if (variant.Bitrate == 832000 || variant.Bitrate == 0) && variant.ContentType == "video/mp4" {
						link = variant.Url
						break
					}
				}
				break
			}
		}
		if link == "" { // Go again if no link for either bitrate is found
			for _, media := range tweet.ExtendedEntities.Media {
				if media.Type == "video" || media.Type == "animated_gif" {
					for _, variant := range media.VideoInfo.Variants {
						if variant.ContentType == "video/mp4" {
							link = variant.Url
							break
						}
					}
					break
				}
			}
		}

		// Get video and post
		response, err := http.Get(link)
		if err != nil {
			s.ChannelMessageDelete(msg.ChannelID, msg.ID)
			s.ChannelMessageSend(m.ChannelID, "Error obtaining video / gif!")
			return
		}
		message := &discordgo.MessageSend{
			File: &discordgo.File{
				Name:   strconv.FormatInt(ID, 10) + "." + linkType,
				Reader: response.Body,
			},
		}
		_, err = s.ChannelMessageSendComplex(m.ChannelID, message)
		s.ChannelMessageDelete(msg.ChannelID, msg.ID)
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, "File too large! Here's a link to the video instead: "+link)
			return
		}
	}
}
