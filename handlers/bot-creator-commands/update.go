package botcreatorcommands

import (
	"fmt"
	"os/exec"
	"regexp"
	"strconv"

	osutools "../../osu-functions"
	"github.com/bwmarrin/discordgo"
)

// Update updates osu-tools
func Update(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID != "92502458588205056" {
		s.ChannelMessageSend(m.ChannelID, "YOU ARE NOT VINXIS.........")
		return
	}

	updateRegex, _ := regexp.Compile(`up(date)?\s+(.+)`)
	if !updateRegex.MatchString(m.Content) {
		s.ChannelMessageSend(m.ChannelID, "farm or osu-tools dumbass...")
		return
	}
	switch updateRegex.FindStringSubmatch(m.Content)[2] {
	case "farm":
		go osutools.FarmUpdate()
	case "osu":
		message, _ := s.ChannelMessageSend(m.ChannelID, "Updating osu-tools...")
		_, err := exec.Command("dotnet", "build", "./osu-tools/PerformanceCalculator").Output()
		s.ChannelMessageDelete(m.ChannelID, message.ID)
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, "An error occurred in updating osu-tools! Please try manually."+m.Author.Mention())
			fmt.Println("Update err: " + err.Error())
		} else {
			s.ChannelMessageSend(m.ChannelID, "Updated! "+m.Author.Mention())
		}
	}
}

// UpdateStatus updates the bot discord status
func UpdateStatus(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID != "92502458588205056" {
		s.ChannelMessageSend(m.ChannelID, "YOU ARE NOT VINXIS.........")
		return
	}
	updateRegex, _ := regexp.Compile(`updatestatus\s+(.+)`)
	if !updateRegex.MatchString(m.Content) {
		s.UpdateStatus(0, strconv.Itoa(len(s.State.Guilds))+" servers")
		return
	}

	text := updateRegex.FindStringSubmatch(m.Content)[1]
	s.UpdateStatus(0, text)
}
