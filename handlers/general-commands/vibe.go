package gencommands

import (
	"bytes"
	"crypto/rand"
	"image"
	"image/png"
	"math"
	"math/big"
	"net/http"
	"strconv"
	"time"

	tools "../../tools"
	"github.com/bwmarrin/discordgo"
)

// Vibe checks for their vibe.
func Vibe(s *discordgo.Session, m *discordgo.MessageCreate, checkType string) {
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
				roll, _ := rand.Int(rand.Reader, big.NewInt(2))
				if roll.Int64() == 0 {
					green++
				} else {
					green = int(math.Max(0, float64(green-1)))
				}
			}
		}
	}
	Requirement = 100 * green / (red + green)

	roll, _ := rand.Int(rand.Reader, big.NewInt(100))
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
