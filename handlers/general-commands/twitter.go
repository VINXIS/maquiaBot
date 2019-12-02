package gencommands

import (
	"net/http"
	"os"
	"regexp"
	"strconv"

	"github.com/ChimeraCoder/anaconda"
	"github.com/bwmarrin/discordgo"
)

// Twitter uploads a twitter gif / image / video onto discord directly
func Twitter(s *discordgo.Session, m *discordgo.MessageCreate) {
	var (
		link     string
		linkType string
		ID       int64
	)
	twitterRegex, _ := regexp.Compile(`twitter.com\/(\S+)\/status\/(\d+)`)
	if linkType == "" {
		messages, _ := s.ChannelMessages(m.ChannelID, -1, "", "", "")
		for _, msg := range messages {
			if twitterRegex.MatchString(msg.Content) && len(msg.Embeds) > 0 {
				if msg.Embeds[0].Video != nil {
					linkType = "mp4"
					ID, _ = strconv.ParseInt(twitterRegex.FindStringSubmatch(msg.Content)[2], 10, 64)
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
	} else if linkType == "png" {
		response, err := http.Get(link)
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, "Error obtaining image!")
			return
		}
		message := &discordgo.MessageSend{
			File: &discordgo.File{
				Name:   "twitter." + linkType,
				Reader: response.Body,
			},
		}
		s.ChannelMessageSendComplex(m.ChannelID, message)
	} else {
		api := anaconda.NewTwitterApiWithCredentials(
			os.Getenv("TWITTER_ACCESS_TOKEN"),
			os.Getenv("TWITTER_ACCESS_TOKEN_SECRET"),
			os.Getenv("TWITTER_CONSUMER_KEY"),
			os.Getenv("TWITTER_CONSUMER_SECRET"),
		)
		tweet, err := api.GetTweet(ID, nil)
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, "Error obtaining tweet!")
			return
		}
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
		response, err := http.Get(link)
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, "Error obtaining video / gif!")
			return
		}
		message := &discordgo.MessageSend{
			File: &discordgo.File{
				Name:   "twitter." + linkType,
				Reader: response.Body,
			},
		}
		s.ChannelMessageSendComplex(m.ChannelID, message)
	}
}
