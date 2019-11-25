package handlers

import (
	"math/rand"
	"regexp"
	"strings"

	osuapi "../osu-api"
	osutools "../osu-functions"
	helpcommands "./help-commands"

	"github.com/bwmarrin/discordgo"
)

// HelpHandle lets you know the commands available
func HelpHandle(s *discordgo.Session, m *discordgo.MessageCreate, prefix string) {
	embed := &discordgo.MessageEmbed{
		Author: &discordgo.MessageEmbedAuthor{
			URL:     "https://discordapp.com/oauth2/authorize?&client_id=551667572723023893&scope=bot&permissions=0",
			Name:    "Click here to invite MaquiaBot!",
			IconURL: s.State.User.AvatarURL(""),
		},
		Description: "Detailed version of the commands list [here](https://docs.google.com/spreadsheets/d/12VzMXGoxliSVv6Rrr6tEy_-Qe9oJ0TNF4MoPGcxIpcU/edit?usp=sharing). **Most commands have other forms as well for convenience!**" + "\n\n" +
			"**Please do `" + prefix + "help <command>` for more information about the command!** \n" +
			"Help information format: `(cmd name) <args> [optional args]`",
		Color: osutools.ModeColour(osuapi.ModeOsu),
	}

	argRegex, _ := regexp.Compile(`help\s+(.+)`)
	if argRegex.MatchString(m.Content) {
		arg := argRegex.FindStringSubmatch(m.Content)[1]
		if arg == "pokemon" || arg == "osu" {
			arg = strings.Split(argRegex.FindStringSubmatch(m.Content)[1], " ")[1]
		}
		switch arg {
		// Admin commands
		case "statst", "statstoggle":
			embed = helpcommands.StatsToggle(embed)

		// General commands
		case "adj", "adjective", "adjectives":
			embed = helpcommands.Adjectives(embed)
		case "avatar", "ava", "a":
			embed = helpcommands.Avatar(embed)
		case "ch", "choose":
			embed = helpcommands.Choose(embed)
		case "decrypt":
			embed = helpcommands.Decrypt(embed)
		case "encrypt":
			embed = helpcommands.Encrypt(embed)
		case "maquiaprefix", "prefix", "newprefix":
			embed = helpcommands.Prefix(embed)
		case "noun", "nouns":
			embed = helpcommands.Nouns(embed)
		case "skill", "skills":
			embed = helpcommands.Skills(embed)
		case "stats":
			embed = helpcommands.Stats(embed)

		// osu! commands
		case "link", "set":
			embed = helpcommands.Link(embed)
		}
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
	s.ChannelMessageSendComplex(m.ChannelID, &discordgo.MessageSend{
		Content: "All commands in PM will use the bot's default prefix `$` instead!",
		Embed:   embed,
	})
}
