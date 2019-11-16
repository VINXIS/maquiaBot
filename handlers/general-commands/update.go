package gencommands

import (
	"fmt"
	"os/exec"
	"regexp"

	osutools "../../osu-functions"
	"github.com/bwmarrin/discordgo"
)

// Update updates osu-tools
func Update(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID != "92502458588205056" {
		s.ChannelMessageSend(m.ChannelID, "YOU ARE NOT VINXIS.........")
		return
	}

	updateRegex, _ := regexp.Compile(`u(pdate)?\s+(.+)`)
	switch updateRegex.FindStringSubmatch(m.Content)[2] {
	case "farm":
		go osutools.FarmUpdate()
	case "osu":
		message, _ := s.ChannelMessageSend(m.ChannelID, "Updating osu-tools...")
		_, err := exec.Command("dotnet", "build", "./osu-tools/PerformanceCalculator").Output()
		s.ChannelMessageDelete(m.ChannelID, message.ID)
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, "An error occurred in updating osu-tools! Please try manually."+m.Author.Mention())
			fmt.Println(err)
		} else {
			s.ChannelMessageSend(m.ChannelID, "Updated! "+m.Author.Mention())
		}
	}
}
