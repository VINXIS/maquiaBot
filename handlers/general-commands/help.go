package gencommands

import (
	"math/rand"

	osutools "../../osu-functions"
	tools "../../tools"
	helpcommands "./help"

	"github.com/bwmarrin/discordgo"
	"github.com/thehowl/go-osuapi"
)

// Help lets you know the commands available
func Help(s *discordgo.Session, m *discordgo.MessageCreate, prefix string, args []string) {
	dm, err := s.UserChannelCreate(m.Author.ID)
	tools.ErrRead(err)
	embed := &discordgo.MessageEmbed{
		Author: &discordgo.MessageEmbedAuthor{
			URL:     "https://discordapp.com/oauth2/authorize?&client_id=551667572723023893&scope=bot&permissions=0",
			Name:    "Click here to invite MaquiaBot!",
			IconURL: s.State.User.AvatarURL(""),
		},
		Description: "All commands in PM will use the bot's default prefix `$` instead! The prefix used below was assigned by the server owner(s)!" + "\n" +
			"Detailed version of the commands list [here](https://docs.google.com/spreadsheets/d/12VzMXGoxliSVv6Rrr6tEy_-Qe9oJ0TNF4MoPGcxIpcU/edit?usp=sharing). **Most commands have other forms as well for convenience!**" + "\n\n" +
			"**Please do `" + prefix + "help <command>` for more information about the command!** \n" +
			"Format: `cmd <args> [optional args]`",
		Color: osutools.ModeColour(osuapi.ModeOsu),
	}

	if len(args) > 1 {
		switch args[1] {
		case "osu":
			if len(args) > 2 {
				switch args[2] {
				case "link", "set":
					embed = helpcommands.Link(embed)
				case "recent", "r", "rs", "recentb", "rb", "recentbest":
					embed = helpcommands.Recent(embed)
				case "t", "top":
					embed = helpcommands.Top(embed)
				case "tr", "track":
					embed = helpcommands.Track(embed)
				case "ti", "tinfo", "tracking", "trackinfo":
					embed = helpcommands.TrackInfo(embed)
				case "tt", "trackt", "tracktoggle":
					embed = helpcommands.TrackToggle(embed)
				case "c", "compare":
					embed = helpcommands.Compare(embed)
				}
			} else {
				embed = helpcommands.Osu(embed)
			}
		case "p", "pokemon":
			if len(args) > 2 {
				switch args[2] {
				case "b", "berry":
					embed = helpcommands.Berry(embed)
				}
			} else {
				embed = helpcommands.Pokemon(embed)
			}
		case "link", "set":
			embed = helpcommands.Link(embed)
		case "recent", "r", "rs", "recentb", "rb", "recentbest":
			embed = helpcommands.Recent(embed)
		case "t", "top":
			embed = helpcommands.Top(embed)
		case "tr", "track":
			embed = helpcommands.Track(embed)
		case "ti", "tinfo", "tracking", "trackinfo":
			embed = helpcommands.TrackInfo(embed)
		case "tt", "trackt", "tracktoggle":
			embed = helpcommands.TrackToggle(embed)
		case "c", "compare":
			embed = helpcommands.Compare(embed)
		case "b", "berry":
			embed = helpcommands.Berry(embed)

		}
	} else {

	}

	switch rand.Intn(11) {
	case 0:
		embed.Image = &discordgo.MessageEmbedImage{
			URL: "https://cdn.discordapp.com/attachments/555493588465877012/555994312760885248/epicAnimeScene.gif",
		}
	case 1:
		embed.Image = &discordgo.MessageEmbedImage{
			URL: "https://cdn.discordapp.com/attachments/555493588465877012/555996915884490752/epicAnimeGifTWO.gif",
		}
	case 2:
		embed.Image = &discordgo.MessageEmbedImage{
			URL: "https://cdn.discordapp.com/attachments/555493588465877012/556000178406948875/epicAnimeGif5.gif",
		}
	case 3, 4:
		embed.Image = &discordgo.MessageEmbedImage{
			URL: "https://cdn.discordapp.com/attachments/555493588465877012/555998627940532237/tumblr_phjkel3lgn1xlyyvto1_1280.png",
		}
	case 5, 6:
		embed.Image = &discordgo.MessageEmbedImage{
			URL: "https://cdn.discordapp.com/attachments/555493588465877012/555998640142024704/tumblr_phjkel3lgn1xlyyvto2_1280.png",
		}
	case 7, 8:
		embed.Image = &discordgo.MessageEmbedImage{
			URL: "https://cdn.discordapp.com/attachments/555493588465877012/555998669761937418/tumblr_phjkel3lgn1xlyyvto3_1280.png",
		}
	case 9, 10:
		embed.Image = &discordgo.MessageEmbedImage{
			URL: "https://cdn.discordapp.com/attachments/555493588465877012/555998681375965194/tumblr_phjkel3lgn1xlyyvto5_1280.png",
		}
	}
	_, err = s.ChannelMessageSendEmbed(dm.ID, embed)
	if err != nil {
		_, err = s.ChannelMessageSend(m.ChannelID, m.Author.Mention()+" has DMs disabled!")
		tools.ErrRead(err)
	}
	return
}
