package gencommands

import (
	"bytes"
	"crypto/rand"
	"encoding/json"
	"image"
	"image/png"
	"io/ioutil"
	"math"
	"math/big"
	"net/http"
	"os"
	"strconv"
	"time"

	structs "../../structs"
	tools "../../tools"
	"github.com/bwmarrin/discordgo"
)

// Vibe toggles vibechecking messages on/off
func Vibe(s *discordgo.Session, m *discordgo.MessageCreate) {
	server, err := s.Guild(m.GuildID)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "This is not a server!")
		return
	}

	member := &discordgo.Member{}
	for _, guildMember := range server.Members {
		if guildMember.User.ID == m.Author.ID {
			member = guildMember
		}
	}

	if member.User.ID == "" {
		return
	}

	admin := false
	for _, roleID := range member.Roles {
		role, err := s.State.Role(m.GuildID, roleID)
		tools.ErrRead(err)
		if role.Permissions&discordgo.PermissionAdministrator != 0 || role.Permissions&discordgo.PermissionManageServer != 0 {
			admin = true
			break
		}
	}

	if !admin && m.Author.ID != server.OwnerID {
		s.ChannelMessageSend(m.ChannelID, "You must be an admin, server manager, or server owner!")
		return
	}

	// Obtain server data
	serverData := structs.ServerData{
		Prefix:    "$",
		OsuToggle: true,
		Crab:      true,
	}
	_, err = os.Stat("./data/serverData/" + m.GuildID + ".json")
	if err == nil {
		f, err := ioutil.ReadFile("./data/serverData/" + m.GuildID + ".json")
		tools.ErrRead(err)
		_ = json.Unmarshal(f, &serverData)
	} else if os.IsNotExist(err) {
		serverData.Server = *server
	} else {
		tools.ErrRead(err)
		return
	}

	// Set new information in server data
	serverData.Time = time.Now()
	serverData.Vibe = !serverData.Vibe

	jsonCache, err := json.Marshal(serverData)
	tools.ErrRead(err)

	err = ioutil.WriteFile("./data/serverData/"+m.GuildID+".json", jsonCache, 0644)
	tools.ErrRead(err)

	if serverData.Vibe {
		s.ChannelMessageSend(m.ChannelID, "Enabled the vibe check.")
	} else {
		s.ChannelMessageSend(m.ChannelID, "Disabled the vibe check.")
	}
	return
}

// VibeCheck checks for their vibe.
func VibeCheck(s *discordgo.Session, m *discordgo.MessageCreate, checkType string) {
	target := m.Author
	if checkType == "notRandom" {
		if len(m.Mentions) > 0 {
			target = m.Mentions[0]
		} else {
			msgs, err := s.ChannelMessages(m.ChannelID, -1, m.ID, "", "")
			if err != nil {
				s.ChannelMessageSend(m.ChannelID, "Ur really funny mate")
				return
			}
			for _, msg := range msgs {
				if msg.Author.ID != m.Author.ID {
					target = msg.Author
					break
				}
			}
		}
	}
	msg, err := s.ChannelMessageSend(m.ChannelID, "VIBE CHECK...... REACT WITH ⭕ OR ❌ WITHIN 20 SECONDS TO HELP DETERMINE THE VIBE CHECK FOR **"+target.Username+"**\nEACH 🚫 WILL HAVE A 50% CHANCE IN ADDING/REMOVING A ⭕")
	if err != nil {
		return
	}
	_ = s.MessageReactionAdd(m.ChannelID, msg.ID, "⭕")
	_ = s.MessageReactionAdd(m.ChannelID, msg.ID, "❌")
	_ = s.MessageReactionAdd(m.ChannelID, msg.ID, "🚫")
	time.Sleep(20 * time.Second)

	Requirement := 50

	msg, err = s.ChannelMessage(m.ChannelID, msg.ID)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "SOMEONE DELETED THE VIBE CHECK MESSAGE. YOU WILL PAY.")
		return
	}
	red := 1
	green := 1
	for _, reaction := range msg.Reactions {
		if reaction.Emoji.Name == "⭕" {
			green = reaction.Count
		}
		if reaction.Emoji.Name == "❌" {
			red = reaction.Count
		}
		if reaction.Emoji.Name == "🚫" && reaction.Count > 1 {
			for i := 0; i < reaction.Count-1; i++ {
				roll, _ := rand.Int(rand.Reader, big.NewInt(int64(2)))
				if roll.Int64() == 0 {
					green++
				} else {
					green = int(math.Max(float64(0), float64(green-1)))
				}
			}
		}
	}
	Requirement = 100 * green / (red + green)

	roll, _ := rand.Int(rand.Reader, big.NewInt(int64(100)))
	if int(roll.Int64())+1 >= Requirement {
		response, _ := http.Get("https://cdn.discordapp.com/attachments/617108584748154880/644570138770669578/vibe-checked.png")
		img, _, _ := image.Decode(response.Body)
		tools.AddLabel(img, 50, 96, target.Username)
		imgBytes := new(bytes.Buffer)
		err = png.Encode(imgBytes, img)
		s.ChannelMessageSendComplex(m.ChannelID, &discordgo.MessageSend{
			Content: target.Mention() + " you had a " + strconv.Itoa(Requirement) + "% chance in passing the vibe check.... tough luck.",
			Files: []*discordgo.File{
				&discordgo.File{
					Name:   "image.png",
					Reader: imgBytes,
				},
			},
		})
	} else {
		s.ChannelMessageSend(m.ChannelID, "You have passed the vibe check ("+strconv.Itoa(Requirement)+"% chance). Carry on "+target.Mention())
	}
}
