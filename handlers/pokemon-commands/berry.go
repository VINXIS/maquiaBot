package pokemoncommands

import (
	"regexp"
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
	pokemontools "maquiaBot/pokemon-tools"
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
func Berry(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Check
	if len(m.Mentions) > 0 {
		s.ChannelMessageSend(m.ChannelID, "Please don't try mentioning people with the bot!")
		return
	}

	berryRegex, _ := regexp.Compile(`(?i)b(erry)?\s+(\S+)`)
	if !berryRegex.MatchString(m.Content) {
		s.ChannelMessageSend(m.ChannelID, "No berry given!")
		return
	}

	berryString := berryRegex.FindStringSubmatch(m.Content)[2]

	// Obtain data
	berry, err := pokemontools.APICall("berry", berryString, BerryStruct{})
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, err.Error())
		return
	}
	berryData := berry.(*BerryStruct)

	// Create fields for potency
	fields := []*discordgo.MessageEmbedField{}
	for _, flavour := range berryData.Flavours {
		fields = append(fields, &discordgo.MessageEmbedField{
			Name:   "**" + strings.Title(flavour.Flavour.Name) + "**",
			Value:  strconv.Itoa(flavour.Potency) + " potency",
			Inline: true,
		})
	}

	// Create embed
	embed := &discordgo.MessageEmbed{
		Title: strings.Title(berryData.Name) + " (#" + strconv.Itoa(berryData.ID) + ")",
		Color: pokemontools.TypeColour(berryData.NaturalGiftType.Name),
		URL:   "https://bulbapedia.bulbagarden.net/wiki/" + strings.Title(berryData.Name) + "_Berry",
		Description: "**Growth rate:** " + strconv.Itoa(berryData.GrowthTime) + " hours per stage **(" + strconv.Itoa(berryData.GrowthTime*4) + " hours total)** \n" +
			"**Max per tree:** " + strconv.Itoa(berryData.MaxHarvest) + " berries \n" +
			"**Soil Dryness Rate:** " + strconv.Itoa(berryData.SoilDryness) + "\n\n" +
			"**Size:** " + strconv.Itoa(berryData.Size) + "mm | **Smoothness:** " + strconv.Itoa(berryData.Smoothness) + " | **Firmness:** " + berryData.Firmness.Name + "\n\n" +
			"**Natural Gift:** " + strings.Title(berryData.NaturalGiftType.Name) + " - " + strconv.Itoa(berryData.NaturalGiftPower),
		Fields: fields,
		Thumbnail: &discordgo.MessageEmbedThumbnail{
			URL: "https://www.serebii.net/itemdex/sprites/pgl/" + berryData.Name + "berry.png",
		},
	}

	s.ChannelMessageSendEmbed(m.ChannelID, embed)
}
