package pokemoncommands

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	pokemontools "../../pokemon-functions"
	tools "../../tools"
	"github.com/bwmarrin/discordgo"
)

// BerryStruct is the Berry structure
type BerryStruct struct {
	Name             string `json:"name"`
	ID               int    `json:"id"`
	GrowthTime       int    `json:"growth_time"`
	MaxHarvest       int    `json:"max_harvest"`
	Size             int    `json:"size"`
	Smoothness       int    `json:"smoothness"`
	SoilDryness      int    `json:"soil_dryness"`
	NaturalGiftPower int    `json:"natural_gift_power"`
	NaturalGiftType  struct {
		Name string
	} `json:"natural_gift_type"`
	Flavours []FlavourStruct `json:"flavors"`
	Firmness struct {
		Name string
	} `json:"firmness"`
}

// FlavourStruct is the Flavour structure
type FlavourStruct struct {
	Flavour struct {
		Name string
	} `json:"flavor"`
	Potency int
}

// Berry gets information about berries
func Berry(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
	// Check
	if len(m.Mentions) > 0 {
		s.ChannelMessageSend(m.ChannelID, "Please don't try mentioning people with the bot!")
		return
	}

	berryString := ""
	if len(args) > 3 || (len(args) == 3 && strings.ToLower(args[1]) != "berry") {
		s.ChannelMessageSend(m.ChannelID, "Too many args! Please use _ for spaces in the berry name!")
		return
	} else if len(args) < 2 {
		s.ChannelMessageSend(m.ChannelID, "No berry name given!")
		return
	} else if strings.ToLower(args[1]) == "berry" {
		berryString = args[2]
	} else {
		berryString = args[1]
	}

	// Obtain data
	res, err := http.Get("https://pokeapi.co/api/v2/berry/" + strings.ToLower(berryString))
	tools.ErrRead(err)

	byteArray, err := ioutil.ReadAll(res.Body)
	tools.ErrRead(err)

	if strings.ToLower(string(byteArray)) == "not found" || strings.HasPrefix(string(byteArray), "<") {
		s.ChannelMessageSend(m.ChannelID, "Berry **"+berryString+"** does not exist!")
		return
	}

	// Convert to readable data
	berry := BerryStruct{}
	err = json.Unmarshal(byteArray, &berry)
	tools.ErrRead(err)

	// Create fields for potency
	fields := []*discordgo.MessageEmbedField{}
	for _, flavour := range berry.Flavours {
		fields = append(fields, &discordgo.MessageEmbedField{
			Name:   "**" + strings.Title(flavour.Flavour.Name) + "**",
			Value:  strconv.Itoa(flavour.Potency) + " potency",
			Inline: true,
		})
	}

	// Create embed
	embed := &discordgo.MessageEmbed{
		Title: strings.Title(berry.Name) + " (#" + strconv.Itoa(berry.ID) + ")",
		Color: pokemontools.TypeColour(berry.NaturalGiftType.Name),
		URL:   "https://bulbapedia.bulbagarden.net/wiki/" + strings.Title(berry.Name) + "_Berry",
		Description: "**Growth rate:** " + strconv.Itoa(berry.GrowthTime) + " hours per stage **(" + strconv.Itoa(berry.GrowthTime*4) + " hours total)** \n" +
			"**Max per tree:** " + strconv.Itoa(berry.MaxHarvest) + " berries \n" +
			"**Soil Dryness Rate:** " + strconv.Itoa(berry.SoilDryness) + "\n\n" +
			"**Size:** " + strconv.Itoa(berry.Size) + "mm | **Smoothness:** " + strconv.Itoa(berry.Smoothness) + " | **Firmness:** " + berry.Firmness.Name + "\n\n" +
			"**Natural Gift:** " + strings.Title(berry.NaturalGiftType.Name) + " - " + strconv.Itoa(berry.NaturalGiftPower),
		Fields: fields,
		Thumbnail: &discordgo.MessageEmbedThumbnail{
			URL: "https://www.serebii.net/itemdex/sprites/pgl/" + berry.Name + "berry.png",
		},
	}

	s.ChannelMessageSendEmbed(m.ChannelID, embed)
}
