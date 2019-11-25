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
			"Help information format: `(cmd|names) <args> [optional args]`",
		Fields: []*discordgo.MessageEmbedField{
			&discordgo.MessageEmbedField{
				Name:  "Admin commands:",
				Value: "`(prefix|maquiaprefix|newprefix)`, " +
				"`(statst|statstoggle)`, " + 
				"`(vibet|vibetoggle)`",
			},
			&discordgo.MessageEmbedField{
				Name:  "General commands:",
				Value: "`(adj|adjective|adjectives)`, " +
				"`(a|ava|avatar)`, " +
				"`(ch|choose)`, " +
				"`decrypt`, " +
				"`encrypt`, " +
				"`face`, " +
				"`funny`, " +
				"`info`, " +
				"`kanye`, " +
				"`(l|leven|levenshtein)`, " +
				"`(noun|nouns)`, " +
				"`ocr`, " +
				"`(p|per|percent|percentage)`, " +
				"`parse`, " +
				"`penis`, " +
				"`ping`, " +
				"`(remind|reminder)`, " +
				"`reminders`, " +
				"`(remindremove|rremove)`, " +
				"`roll`, " +
				"`(sinfo|serverinfo)`, " +
				"`(skill|skills)`, " +
				"`stats`, " +
				"`(vibe|vibec|vibecheck)`",
			},
		},
		Color: osutools.ModeColour(osuapi.ModeOsu),
	}

	argRegex, _ := regexp.Compile(`help\s+(.+)`)
	if argRegex.MatchString(m.Content) {
		arg := argRegex.FindStringSubmatch(m.Content)[1]
		if strings.Split(arg, " ")[0] == "pokemon" || strings.Split(arg, " ")[0] == "osu" {
			args := strings.Split(argRegex.FindStringSubmatch(m.Content)[1], " ")
			if len(args) > 1 {
				arg = args[1]
			}
		}
		switch arg {
		// Admin commands
		case "prefix", "maquiaprefix", "newprefix":
			embed = helpcommands.Prefix(embed)
		case "statst", "statstoggle":
			embed = helpcommands.StatsToggle(embed)
		case "vibet", "vibetoggle":
			embed = helpcommands.VibeToggle(embed)

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
		case "face":
			embed = helpcommands.Face(embed)
		case "funny":
			embed = helpcommands.Funny(embed)
		case "info":
			embed = helpcommands.Info(embed)
		case "kanye":
			embed = helpcommands.Kanye(embed)
		case "l", "leven", "levenshtein":
			embed = helpcommands.Levenshtein(embed)
		case "noun", "nouns":
			embed = helpcommands.Nouns(embed)
		case "ocr":
			embed = helpcommands.OCR(embed)
		case "p", "per", "percent", "percentage":
			embed = helpcommands.Percentage(embed)
		case "parse":
			embed = helpcommands.Parse(embed)
		case "penis":
			embed = helpcommands.Penis(embed)
		case "ping":
			embed = helpcommands.Ping(embed)
		case "remind", "reminder":
			embed = helpcommands.Remind(embed)
		case "reminders":
			embed = helpcommands.Reminders(embed)
		case "remindremove", "rremove":
			embed = helpcommands.RemindRemove(embed)
		case "roll":
			embed = helpcommands.Roll(embed)
		case "sinfo", "serverinfo":
			embed = helpcommands.ServerInfo(embed)
		case "skill", "skills":
			embed = helpcommands.Skills(embed)
		case "stats":
			embed = helpcommands.Stats(embed)
		case "vibe", "vibec", "vibecheck":
			embed = helpcommands.Vibe(embed)

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
