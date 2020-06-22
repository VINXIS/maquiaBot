package osucommands

import (
	"encoding/json"
	"io/ioutil"
	"strings"

	"github.com/bwmarrin/discordgo"
	osuapi "maquiaBot/osu-api"
	structs "maquiaBot/structs"
	tools "maquiaBot/tools"
)

// TrackMapper lets people track mappers
func TrackMapper(s *discordgo.Session, m *discordgo.MessageCreate, mapperData []structs.MapperData) {
	// Get mapper(s), channel and perms
	if !strings.Contains(m.Content, " ") {
		s.ChannelMessageSend(m.ChannelID, "No mappers given to add!!")
		return
	}
	withoutPrefix := strings.Split(m.Content, " ")[1:]
	remove := withoutPrefix[0] == "remove" || withoutPrefix[0] == "r"
	var args []string
	if remove {
		args = strings.Split(strings.Join(withoutPrefix[1:], " "), ", ")
	} else {
		args = strings.Split(strings.Join(withoutPrefix, " "), ", ")
	}

	ch, err := s.Channel(m.ChannelID)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "This is not an allowed channel!")
	}

	server, err := s.Guild(m.GuildID)
	owner := false
	if err != nil {
		owner = true
	} else {
		owner = m.Author.ID == server.OwnerID
	}

	if !owner {
		member, _ := s.GuildMember(server.ID, m.Author.ID)
		admin := false
		for _, roleID := range member.Roles {
			role, _ := s.State.Role(m.GuildID, roleID)
			if role.Permissions&discordgo.PermissionAdministrator != 0 || role.Permissions&discordgo.PermissionManageServer != 0 {
				admin = true
				break
			}
		}
		if !admin {
			s.ChannelMessageSend(m.ChannelID, "You must be an admin, server manager, or server owner!")
			return
		}
	}

	// The Real Stuff
	if len(args) == 1 && args[0] == "" {
		if remove {
			for i := 0; i < len(mapperData); i++ {
				mapperData[i].RemoveChannel(*ch)
			}
		} else {
			s.ChannelMessageSend(m.ChannelID, "No mappers given to add!!")
			return
		}
	} else {
		for _, arg := range args {
			// Check if mapper is already in data
			mapperExists := false
			for i := 0; i < len(mapperData); i++ {
				if strings.ToLower(mapperData[i].Mapper.Username) == strings.ToLower(arg) {
					mapperExists = true
					if remove {
						mapperData[i].RemoveChannel(*ch)
					} else {
						mapperData[i].AddChannel(*ch)
					}
					break
				}
			}

			// Create new mapper or skip to next if we are removing instead
			if !mapperExists {
				if remove {
					continue
				} else {
					user, err := OsuAPI.GetUser(osuapi.GetUserOpts{
						Username: arg,
					})
					if err != nil {
						continue
					}
					beatmaps, err := OsuAPI.GetBeatmaps(osuapi.GetBeatmapsOpts{
						UserID: user.UserID,
					})
					newMapper := structs.MapperData{
						Mapper:   *user,
						Beatmaps: beatmaps,
					}
					newMapper.AddChannel(*ch)
					mapperData = append(mapperData, newMapper)
				}
			}
		}
	}

	// Remove all mappers with no channels tracking them
	for i := 0; i < len(mapperData); i++ {
		if len(mapperData[i].Channels) == 0 {
			mapperData = append(mapperData[:i], mapperData[i+1:]...)
			i--
		}
	}

	// Save mapper data
	jsonCache, err := json.Marshal(mapperData)
	tools.ErrRead(s, err)

	err = ioutil.WriteFile("./data/osuData/mapperData.json", jsonCache, 0644)
	tools.ErrRead(s, err)

	if remove && len(args) == 0 {
		s.ChannelMessageSend(m.ChannelID, "Removed tracking for all mappers for this channel!")
		return
	}
	TrackMapperInfo(s, m, mapperData)
}
