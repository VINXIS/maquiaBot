package tools

import "github.com/bwmarrin/discordgo"

// AdminCheck checks if the user has admin privileges
func AdminCheck(s *discordgo.Session, m *discordgo.MessageCreate, server discordgo.Guild) (admin bool) {
	if m.Author.ID == server.OwnerID {
		admin = true
	} else {
		member, _ := s.GuildMember(server.ID, m.Author.ID)
		for _, roleID := range member.Roles {
			role, err := s.State.Role(m.GuildID, roleID)
			if err == nil && (role.Permissions&discordgo.PermissionAdministrator != 0 || role.Permissions&discordgo.PermissionManageServer != 0) {
				admin = true
				break
			}
		}
	}
	return admin
}
