package gencommands

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"regexp"
	"time"

	config "../../config"
	"github.com/bwmarrin/discordgo"
)

// TwitchClip holds twitch API information for clips
type TwitchClip struct {
	Data []struct {
		ID              string    `json:"id"`
		URL             string    `json:"url"`
		EmbedURL        string    `json:"embed_url"`
		BroadcasterID   int       `json:"broadcaster_id,string"`
		BroadcasterName string    `json:"broadcaster_name"`
		CreatorID       int       `json:"creator_id,string"`
		CreatorName     string    `json:"creator_name"`
		VideoID         int       `json:"video_id,string"`
		GameID          int       `json:"game_id,string"`
		Language        string    `json:"language"`
		Title           string    `json:"title"`
		ViewCount       int       `json:"view_count"`
		CreatedAt       time.Time `json:"created_at,"`
		ThumbnailURL    string    `json:"thumbnail_url"`
	} `json:"data"`
}

// Twitch uploads a twitch clip onto discord directly
func Twitch(s *discordgo.Session, m *discordgo.MessageCreate) {
	twitchRegex, _ := regexp.Compile(`https://clips.twitch.tv/(\S+)`)
	thumbnailRegex, _ := regexp.Compile(`-preview-\d+x\d+\.jpg`)

	// Get ID
	ID := ""
	if twitchRegex.MatchString(m.Content) {
		ID = twitchRegex.FindStringSubmatch(m.Content)[1]
	} else {
		msgs, _ := s.ChannelMessages(m.ChannelID, -1, m.ID, "", "")
		for _, msg := range msgs {
			if twitchRegex.MatchString(msg.Content) {
				ID = twitchRegex.FindStringSubmatch(msg.Content)[1]
				break
			}
		}
	}

	// Check if ID was found
	if ID == "" {
		s.ChannelMessageSend(m.ChannelID, "No twitch clip found!")
		return
	}

	// API request
	client := &http.Client{}
	req, _ := http.NewRequest("GET", "https://api.twitch.tv/helix/clips?id="+ID, nil)
	req.Header.Set("Client-ID", config.Conf.Twitch.ID)
	res, err := client.Do(req)
	defer res.Body.Close()
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "Error in obtaining clip information! Error 1")
		return
	}
	byteArray, err := ioutil.ReadAll(res.Body)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "Error in obtaining clip information! Error 2")
		return
	}

	// Parse response
	var clipData TwitchClip
	err = json.Unmarshal(byteArray, &clipData)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "Error in obtaining clip information! Error 3")
		return
	}
	if len(clipData.Data) == 0 {
		s.ChannelMessageSend(m.ChannelID, "No twitch clip found!")
		return
	}

	// Find URL
	if !thumbnailRegex.MatchString(clipData.Data[0].ThumbnailURL) {
		s.ChannelMessageSend(m.ChannelID, "No twitch clip found!")
		return
	}
	video := thumbnailRegex.ReplaceAllString(clipData.Data[0].ThumbnailURL, ".mp4")

	msg, err := s.ChannelMessageSend(m.ChannelID, "Obtaining twitch clip...")
	if err != nil {
		return
	}

	// Get video
	response, err := http.Get(video)
	defer response.Body.Close()
	if err != nil {
		s.ChannelMessageDelete(msg.ChannelID, msg.ID)
		s.ChannelMessageSend(m.ChannelID, "Error obtaining video / gif!")
		return
	}
	message := &discordgo.MessageSend{
		File: &discordgo.File{
			Name:   ID + ".mp4",
			Reader: response.Body,
		},
	}
	_, err = s.ChannelMessageSendComplex(m.ChannelID, message)
	s.ChannelMessageDelete(msg.ChannelID, msg.ID)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "File too large for this server! Here's a link to the video instead: "+video)
		return
	}
}
