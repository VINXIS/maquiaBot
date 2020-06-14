package pokemoncommands

import (
	"regexp"
	"strconv"
	"strings"

	pokemontools "../../pokemon-tools"
	"github.com/bwmarrin/discordgo"
)

// PokemonStruct is the Pokemon structure
type PokemonStruct struct {
	Name      string           `json:"name"`
	BaseExp   int              `json:"base_experience"`
	Height    int              `json:"height"`
	Weight    int              `json:"weight"`
	ID        int              `json:"id"`
	Games     []interface{}    `json:"game_indices"`
	Types     []PokemonType    `json:"types"`
	Sprites   PokemonSprite    `json:"sprites"`
	Stats     []PokemonStat    `json:"stats"`
	Items     []PokemonItem    `json:"held_items"`
	Abilities []PokemonAbility `json:"abilities"`
}

// PokemonType is the Pokemon Type structure
type PokemonType struct {
	Type struct {
		Name string
	} `json:"type"`
}

// PokemonSprite is the Pokemon Sprite structure
type PokemonSprite struct {
	BackDefault      string `json:"back_default"`
	BackFemale       string `json:"back_female"`
	BackShiny        string `json:"back_shiny"`
	BackShinyFemale  string `json:"back_shiny_female"`
	FrontDefault     string `json:"front_default"`
	FrontFemale      string `json:"front_female"`
	FrontShiny       string `json:"front_shiny"`
	FrontShinyFemale string `json:"front_shiny_female"`
}

// PokemonStat is the Pokemon Stat structure
type PokemonStat struct {
	Stat struct {
		Name string
	} `json:"stat"`
	Effort   int `json:"effort"`
	BaseStat int `json:"base_stat"`
}

// PokemonItem is the Pokemon Held Items structure
type PokemonItem struct {
	Item struct {
		Name string
	} `json:"item"`
}

// PokemonAbility is the Pokemon Ability structure
type PokemonAbility struct {
	Ability struct {
		Name string
	} `json:"ability"`
	IsHidden bool `json:"is_hidden"`
	Slot     int  `json:"slot"`
}

// Pokemon searches for the pokemon and returns result
func Pokemon(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Check
	if len(m.Mentions) > 0 {
		s.ChannelMessageSend(m.ChannelID, "Please don't try mentioning people with the bot!")
		return
	}

	pokemonRegex, _ := regexp.Compile(`pokemon\s+(\S+)`)
	if !pokemonRegex.MatchString(m.Content) {
		s.ChannelMessageSend(m.ChannelID, "No pokemon given!")
		return
	}

	pokemonString := pokemonRegex.FindStringSubmatch(m.Content)[1]

	// If deoxys change to id lol
	if strings.ToLower(pokemonString) == "deoxys" {
		pokemonString = "386"
	}

	// Obtain data
	pokemon, err := pokemontools.APICall("pokemon", pokemonString, PokemonStruct{})
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, err.Error())
		return
	}
	pokemonData := pokemon.(*PokemonStruct)

	// Assign values
	weight := "**Weight:** " + strconv.Itoa(pokemonData.Weight/10) + "kg | "
	if pokemonData.Weight/10 == 0 {
		weight = "**Weight:** " + strconv.Itoa(pokemonData.Weight*100) + "g | "
	}
	height := "**Height:** " + strconv.Itoa(pokemonData.Height/10) + "m | "
	if pokemonData.Height/10 == 0 {
		height = "**Height:** " + strconv.Itoa(pokemonData.Height*100) + "cm | "
	}
	exp := "**" + strconv.Itoa(pokemonData.BaseExp) + " EXP** gained from defeating " + strings.Title(pokemonData.Name)
	items := "Possible item(s) held: "
	if len(pokemonData.Items) == 0 {
		items = items + "**None**"
	} else {
		for i, item := range pokemonData.Items {
			if i == 0 {
				items = items + "**" + strings.Title(item.Item.Name) + "**"
			} else {
				items = items + " and **" + strings.Title(item.Item.Name) + "**"
			}
		}
	}

	// Create ability string
	abilities := "List of abilities: "
	if len(pokemonData.Items) == 0 {
		abilities = abilities + "**None**"
	} else {
		for i, ability := range pokemonData.Abilities {
			if i == 0 {
				if ability.IsHidden {
					abilities = abilities + "**" + strings.Title(ability.Ability.Name) + "** (hidden)"
				} else {
					abilities = abilities + "**" + strings.Title(ability.Ability.Name) + "**"
				}
			} else {
				if ability.IsHidden {
					abilities = abilities + " and **" + strings.Title(ability.Ability.Name) + "** (hidden)"
				} else {
					abilities = abilities + " and **" + strings.Title(ability.Ability.Name) + "**"
				}
			}
		}
	}

	// Create fields for stats + EV
	fields := []*discordgo.MessageEmbedField{}
	for _, pokemonStat := range pokemonData.Stats {
		fields = append(fields, &discordgo.MessageEmbedField{
			Name:   "**" + strings.Title(pokemonStat.Stat.Name) + "**",
			Value:  strconv.Itoa(pokemonStat.BaseStat) + " (" + strconv.Itoa(pokemonStat.Effort) + " EV)",
			Inline: true,
		})
	}

	// Create type string
	types := ""
	if len(pokemonData.Types) == 1 {
		types = types + "**" + strings.Title(pokemonData.Types[0].Type.Name) + "**"
	} else {
		for i, pokemonType := range pokemonData.Types {
			if i == 0 {
				types = types + "**" + strings.Title(pokemonType.Type.Name) + "**"
			} else {
				types = types + " and **" + strings.Title(pokemonType.Type.Name) + "**"
			}
		}
	}

	// Create embed
	embed := &discordgo.MessageEmbed{
		Title: strings.Title(pokemonData.Name) + " (#" + strconv.Itoa(pokemonData.ID) + ")",
		Color: pokemontools.TypeColour(pokemonData.Types[0].Type.Name),
		Description: "In around **" + strconv.Itoa(len(pokemonData.Games)) + "** Pokemon games\n\n" +
			weight + height + types + "\n\n" +
			exp + "\n" +
			items + "\n" +
			abilities,
		Fields: fields,
		Image: &discordgo.MessageEmbedImage{
			URL: "https://www.smogon.com/dex/media/sprites/xy/" + pokemonData.Name + ".gif",
		},
	}
	formRegex, _ := regexp.Compile(`(.+)-.+`)
	if formRegex.MatchString(pokemonData.Name) {
		embed.URL = "https://bulbapedia.bulbagarden.net/wiki/" + strings.Title(formRegex.FindStringSubmatch(pokemonData.Name)[1]) + "_(Pok%C3%A9mon)"
	} else {
		embed.URL = "https://bulbapedia.bulbagarden.net/wiki/" + strings.Title(pokemonData.Name) + "_(Pok%C3%A9mon)"
	}

	s.ChannelMessageSendEmbed(m.ChannelID, embed)
}
