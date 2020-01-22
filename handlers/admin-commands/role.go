package admincommands

import (
	"encoding/json"
	"io/ioutil"
	"regexp"
	"strconv"
	"strings"

	structs "../../structs"
	tools "../../tools"
	"github.com/bwmarrin/discordgo"
)

// RoleAutomation adds / removes role automation options
func RoleAutomation(s *discordgo.Session, m *discordgo.MessageCreate) {
	roleRegex, _ := regexp.Compile(`role(a|auto|automation)?\s+(.+)`)
	deleteRegex, _ := regexp.Compile(`-d`)

	// Check if server exists
	server, err := s.Guild(m.GuildID)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "This is not a server so custom prefixes are unavailable! Please use `$` instead for commands!")
		return
	}

	if !tools.AdminCheck(s, m, *server) {
		s.ChannelMessageSend(m.ChannelID, "You must be an admin, server manager, or server owner!")
		return
	}

	// Obtain server data
	serverData := tools.GetServer(*server)

	// Check if params were given
	if !roleRegex.MatchString(m.Content) {
		s.ChannelMessageSend(m.ChannelID, "No params given!")
		return
	}

	text := roleRegex.FindStringSubmatch(m.Content)[2]

	// Delete function
	if deleteRegex.MatchString(m.Content) {
		text = strings.TrimSpace(deleteRegex.ReplaceAllString(text, ""))
		if ID, err := strconv.Atoi(text); err == nil {
			for i, roleAuto := range serverData.RoleAutomation {
				if roleAuto.ID == ID {
					serverData.RoleAutomation = append(serverData.RoleAutomation[:i], serverData.RoleAutomation[i+1:]...)
					break
				}
			}
			jsonCache, err := json.Marshal(serverData)
			tools.ErrRead(err)

			err = ioutil.WriteFile("./data/serverData/"+m.GuildID+".json", jsonCache, 0644)
			tools.ErrRead(err)
			s.ChannelMessageSend(m.ChannelID, "Removed role automation ID: "+text)
		} else {
			s.ChannelMessageSend(m.ChannelID, text+" is an invalid ID!")
		}
		return
	}

	// Look for roles pinged
	roleIDs := []string{}
	if len(m.MentionRoles) > 0 {
		for _, role := range m.MentionRoles {
			roleIDs = append(roleIDs, role)
			text = strings.TrimSpace(strings.Replace(text, "<@&"+role+">", "", -1))
		}
	}

	// Look for just IDs if they werent role pings
	if len(roleIDs) == 0 {
		params := strings.Split(text, " ")
		if len(params) == 1 {
			s.ChannelMessageSend(m.ChannelID, "Not enough params!")
			return
		}

		for _, param := range params {
			r, err := s.State.Role(m.GuildID, param)
			if err != nil {
				continue
			}

			roleIDs = append(roleIDs, r.ID)
			text = strings.TrimSpace(strings.Replace(text, param, "", -1))
		}
	}

	// no roles found
	if len(roleIDs) == 0 {
		s.ChannelMessageSend(m.ChannelID, "No roles found!")
		return
	}

	warningText := "WARNING: Could not find the roles associated with the following IDs: "
	roleData := structs.Role{
		Text: text,
	}

	for _, role := range roleIDs {
		r, err := s.State.Role(m.GuildID, role)
		if err != nil {
			warningText += role + ", "
			continue
		}
		roleData.Roles = append(roleData.Roles, *r)
	}

	// Somehow no roles were obtained with the IDs
	if len(roleData.Roles) == 0 {
		s.ChannelMessageSend(m.ChannelID, "No roles found!")
		return
	}

	if len(roleData.Roles) != len(roleIDs) {
		s.ChannelMessageSend(m.ChannelID, strings.TrimSuffix(warningText, ", "))
	}

	// Check duplicate
	found := false
	for i, roleAuto := range serverData.RoleAutomation {
		if roleAuto.Text == roleData.Text {
			found = true
			for _, role := range roleData.Roles {
				roleFound := false
				for _, autoRole := range roleAuto.Roles {
					if role.ID == autoRole.ID {
						roleFound = true
						break
					}
				}
				
				if !roleFound {
					serverData.RoleAutomation[i].Roles = append(serverData.RoleAutomation[i].Roles, role)
				}
			}
			roleData = serverData.RoleAutomation[i]
			break
		}
		if roleAuto.ID == roleData.ID {
			roleData.ID++
		}
	}

	if !found {
		serverData.RoleAutomation = append(serverData.RoleAutomation, roleData)
	}
	jsonCache, err := json.Marshal(serverData)
	tools.ErrRead(err)

	err = ioutil.WriteFile("./data/serverData/"+m.GuildID+".json", jsonCache, 0644)
	tools.ErrRead(err)

	var roleNames string
	for _, role := range roleData.Roles {
		roleNames += role.Name + ", "
	}

	s.ChannelMessageSend(m.ChannelID, "Added the role automation: When someone sends `"+text+"`, they will have the following roles applied: "+strings.TrimSuffix(roleNames, ", ")+"\nAll role automations enabled in this server can be seen via `roleinfo`")
	return
}
