package osucommands

import (
	"strconv"
	"strings"

	osutools "../../osu-tools"
	structs "../../structs"
	tools "../../tools"
	"github.com/bwmarrin/discordgo"
)

// TrackInfo gives info about what's being tracked in the channel currently
func TrackInfo(s *discordgo.Session, m *discordgo.MessageCreate) {
	channel, err := s.Channel(m.ChannelID)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "This is not an allowed channel!")
		return
	}

	server, err := s.Guild(m.GuildID)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "This is not a server!")
		return
	}

	// Obtain channel data
	channelData, new := tools.GetChannel(*channel, s)
	if new {
		s.ChannelMessageSend(m.ChannelID, "There is no tracking info for this channel currently!")
		return
	}

	serverImg := "https://cdn.discordapp.com/icons/" + server.ID + "/" + server.Icon
	if strings.Contains(server.Icon, "a_") {
		serverImg += ".gif"
	} else {
		serverImg += ".png"
	}

	// Create embed
	embed := discordgo.MessageEmbed{
		Author: &discordgo.MessageEmbedAuthor{
			Name:    channelData.Channel.Name,
			IconURL: serverImg,
		},
		Description: "Any admin, server moderator, or server owner can update this using `track` again! \n**You do not need to readd everything again to update parts of the tracker.** Toggle tracking on or off using `tt`, `trackt`, `ttoggle`, or `tracktoggle`",
		Color:       osutools.ModeColour(channelData.Mode),
	}

	// Warnings
	warning := ""
	if !channelData.Ranked && !channelData.Loved && !channelData.Qualified {
		warning += "\n**WARNING:** You do not have any map rank statuses with leaderboards enabled! Please enable at least one in order for tracking to work!"
	}
	if channelData.LeaderboardReq == 101 && channelData.TopReq == 101 && channelData.PPReq == -1 {
		warning += "\n**WARNING:** You do not have any leader/top/pp requirement for scores! Any score submitted by the users listed on eligible maps will be posted as a result!"
	}
	if !channelData.Tracking {
		warning += "\n**WARNING:** Tracking is currently turned off!"
	}

	// Add users
	userList := ""
	for i, user := range channelData.Users {
		if i == len(channelData.Users)-1 {
			userList += user.Username
		} else {
			userList += user.Username + ", "
		}
	}
	embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
		Name:  "Users Tracked:",
		Value: userList,
	})

	// Add command info

	statuses := ""
	if channelData.Ranked {
		statuses += "r"
	}
	if channelData.Loved {
		statuses += "l"
	}
	if channelData.Qualified {
		statuses += "q"
	}
	statuses = strings.Join(strings.Split(statuses, ""), ", ")

	command :=
		"`-u " + userList +
			" -pp " + strconv.FormatFloat(channelData.PPReq, 'f', 0, 64) +
			" -l " + strconv.Itoa(channelData.LeaderboardReq) +
			" -t " + strconv.Itoa(channelData.TopReq) +
			" -m " + channelData.Mode.String() +
			" -s " + statuses + "`"

	embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
		Name:  "Command for This Config:",
		Value: command,
	})

	// Add PP req
	if channelData.PPReq != -1 {
		embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{"PP Req:", strconv.FormatFloat(channelData.PPReq, 'f', 0, 64), true})
	} else {
		embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{"PP Req:", "N/A", true})
	}

	// Add Leaderboard req
	if channelData.LeaderboardReq != 101 {
		embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{"Leaderboard Req:", strconv.Itoa(channelData.LeaderboardReq), true})
	} else {
		embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{"Leaderboard Req:", "N/A", true})
	}

	// Add Top req
	if channelData.TopReq != 101 {
		embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{"Top Perf. Req:", strconv.Itoa(channelData.TopReq), true})
	} else {
		embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{"Top Perf. Req:", "N/A", true})
	}

	// Map types
	embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{"Ranked: ", strconv.FormatBool(channelData.Ranked), true})
	embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{"Loved: ", strconv.FormatBool(channelData.Loved), true})
	embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{"Qualified: ", strconv.FormatBool(channelData.Qualified), true})

	// Misc.
	embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{"Mode: ", channelData.Mode.String(), true})
	embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{"Tracking: ", strconv.FormatBool(channelData.Tracking), true})

	s.ChannelMessageSendComplex(m.ChannelID, &discordgo.MessageSend{
		Content: warning,
		Embed:   &embed,
	})
	return
}

// TrackMapperInfo gives information about mapper tracking for the channel
func TrackMapperInfo(s *discordgo.Session, m *discordgo.MessageCreate, mapperData []structs.MapperData) {
	channel, err := s.Channel(m.ChannelID)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "This is not an allowed channel!")
		return
	}

	var mappers []string
	for _, mapper := range mapperData {
		for _, ch := range mapper.Channels {
			if ch == channel.ID {
				mappers = append(mappers, mapper.Mapper.Username)
			}
		}
	}
	if len(mappers) == 0 {
		s.ChannelMessageSend(m.ChannelID, "No mappers are currently being tracked!")
		return
	}

	s.ChannelMessageSend(m.ChannelID, "This channel is following "+strconv.Itoa(len(mappers))+" mapper(s). They are:\n"+strings.Join(mappers, ", "))
}
