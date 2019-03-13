package gencommands

import (
	osutools "../../osu-functions"
	tools "../../tools"

	"github.com/bwmarrin/discordgo"
	"github.com/thehowl/go-osuapi"
)

// Help lets you know the commands available
func Help(s *discordgo.Session, m *discordgo.MessageCreate, prefix string) {
	dm, err := s.UserChannelCreate(m.Author.ID)
	tools.ErrRead(err)

	s.ChannelMessageSendEmbed(dm.ID, &discordgo.MessageEmbed{
		Author: &discordgo.MessageEmbedAuthor{
			URL:     "https://discordapp.com/oauth2/authorize?&client_id=551667572723023893&scope=bot&permissions=0",
			Name:    "Click here to invite MaquiaBot!",
			IconURL: s.State.User.AvatarURL(""),
		},
		Color: osutools.ModeColour(osuapi.ModeOsu),
		Fields: []*discordgo.MessageEmbedField{
			&discordgo.MessageEmbedField{
				Name: "general",
				Value: "`" + prefix + "avatar [@user]` - Returns the avatar of the user" + "\n" +
					"`" + prefix + "help` - Returns the list of commands" + "\n" +
					"`" + prefix + "newPrefix <prefix>` or `maquiaprefix <prefix>` - Creates a new prefix for this bot" + "\n",
			},
			&discordgo.MessageEmbedField{
				Name: "osu!",
				Value: "`" + prefix + "link <username>` - Links an osu! profile to your discord account" + "\n" +
					"`" + prefix + "rs [username] [n]` - Checks the nth recent score for either the account linked to your discord, or the username if given" + "\n" +
					"`" + prefix + "rb [username] [n]` - Checks the nth recent top performance play for either the account linked to your discord, or the username if given",
			},
		},
	})
	return
}
