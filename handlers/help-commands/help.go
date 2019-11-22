package helpcommands

import (
	"math/rand"

	osuapi "../../osu-api"
	osutools "../../osu-functions"

	"github.com/bwmarrin/discordgo"
)

// Help lets you know the commands available
func Help(s *discordgo.Session, m *discordgo.MessageCreate, prefix string, args []string) {
	dm, _ := s.UserChannelCreate(m.Author.ID)
	embed := &discordgo.MessageEmbed{
		Author: &discordgo.MessageEmbedAuthor{
			URL:     "https://discordapp.com/oauth2/authorize?&client_id=551667572723023893&scope=bot&permissions=0",
			Name:    "Click here to invite MaquiaBot!",
			IconURL: s.State.User.AvatarURL(""),
		},
		Description: "Detailed version of the commands list [here](https://docs.google.com/spreadsheets/d/12VzMXGoxliSVv6Rrr6tEy_-Qe9oJ0TNF4MoPGcxIpcU/edit?usp=sharing). **Most commands have other forms as well for convenience!**" + "\n\n" +
			"**Please do `" + prefix + "help <command>` for more information about the command!** \n" +
			"Format: `cmd <args> [optional args]`",
		Color: osutools.ModeColour(osuapi.ModeOsu),
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
	_, err := s.ChannelMessageSendComplex(dm.ID, &discordgo.MessageSend{
		Content: "All commands in PM will use the bot's default prefix `$` instead! The prefix used below was assigned by the server owner(s)!",
		Embed:   embed,
	})
	if err != nil {
		_, err = s.ChannelMessageSend(m.ChannelID, m.Author.Mention()+" has DMs disabled!")
	}
}
