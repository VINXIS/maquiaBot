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
		Description: "All commands in PM will use the bot's default prefix `$` instead! The prefix used below was assigned by the server owner(s)!" + "\n" +
			"Detailed version of the commands list [here](https://docs.google.com/spreadsheets/d/12VzMXGoxliSVv6Rrr6tEy_-Qe9oJ0TNF4MoPGcxIpcU/edit?usp=sharing). **Most commands have other forms as well for convenience!**" + "\n" +
			"Format: `cmd <args> [optional args]`",
		Color: osutools.ModeColour(osuapi.ModeOsu),
		Fields: []*discordgo.MessageEmbedField{
			&discordgo.MessageEmbedField{
				Name: "general",
				Value: "`" + prefix + "avatar [@user]` - Returns the avatar of the user" + "\n" +
					"`" + prefix + "help` - Returns the list of commands" + "\n" +
					"`" + prefix + "newPrefix <prefix>` or `maquiaprefix <prefix>` - Creates a new prefix for this bot" + "\n" +
					"`" + prefix + "source` - Links source code" + "\n",
			},
			&discordgo.MessageEmbedField{
				Name: "osu!",
				Value: "`" + prefix + "link <username>` - Links an osu! profile to your discord account" + "\n" +
					"`" + prefix + "rs [username] [n]` - Checks the nth recent score for either the account linked to your discord, or the username if given" + "\n" +
					"`" + prefix + "rb [username] [n]` - Checks the nth recent top performance play for either the account linked to your discord, or the username if given" + "\n" +
					"`" + prefix + "c` - Gets your score for the map previously linked" + "\n" +
					"`" + prefix + "t (user [username] / users [usernames with spaces]) [pp <pp thresh> top <top thresh>]` - Tracks players listed with pp threshold (if listed) and top threshold (if listed). Posts if the score fits the pp or top criteria" + "\n" +
					"`" + prefix + "t add (user [username] / users [usernames with spaces]) [pp <pp thresh> top <top thresh>]` - Adds users to tracking (only if the channel is being used for tracking already by the bot)" + "\n" +
					"`" + prefix + "t remove [user [username] / users [usernames with spaces]]` - Removes users listed (or tracking altogether if no users are listed)" + "\n" +
					"`" + prefix + "tinfo` - Gives information about what's being tracked in this channel" + "\n" +
					"`" + prefix + "tt` - Toggles tracking on/off",
			},
			&discordgo.MessageEmbedField{
				Name:  "pokemon",
				Value: "`" + prefix + "pokemon <pokemon name/id>` - Gives brief information regarding a pokemon",
			},
		},
		Image: &discordgo.MessageEmbedImage{
			URL: "https://cdn.discordapp.com/attachments/555493588465877012/555994312760885248/epicAnimeScene.gif",
		},
	})
	return
}
