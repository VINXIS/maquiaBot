package pokemoncommands

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	pokemontools "../../pokemon-functions"
	tools "../../tools"
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
func Pokemon(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
	// Check
	if len(args) > 2 {
		s.ChannelMessageSend(m.ChannelID, "Too many args! Please use _ for spaces in the pokemon name!")
	}

	// If deoxys change to id lol
	if strings.ToLower(args[1]) == "deoxys" {
		args[1] = "386"
	}

	// Obtain data
	res, err := http.Get("https://pokeapi.co/api/v2/pokemon/" + strings.ToLower(args[1]))
	tools.ErrRead(err)

	byteArray, err := ioutil.ReadAll(res.Body)
	tools.ErrRead(err)

	if strings.ToLower(string(byteArray)) == "not found" {
		s.ChannelMessageSend(m.ChannelID, "Pokemon **"+args[1]+"** does not exist!")
		return
	}

	// Convert to readable data
	pokemon := PokemonStruct{}
	err = json.Unmarshal(byteArray, &pokemon)
	tools.ErrRead(err)

	// Assign values
	weight := "**Weight:** " + strconv.Itoa(pokemon.Weight/10) + "kg | "
	if pokemon.Weight/10 == 0 {
		weight = "**Weight:** " + strconv.Itoa(pokemon.Weight*100) + "g | "
	}
	height := "**Height:** " + strconv.Itoa(pokemon.Height/10) + "m | "
	if pokemon.Height/10 == 0 {
		height = "**Height:** " + strconv.Itoa(pokemon.Height*100) + "cm | "
	}
	exp := "**" + strconv.Itoa(pokemon.BaseExp) + " EXP** gained from defeating " + strings.Title(pokemon.Name)
	items := "Possible item(s) held: "
	if len(pokemon.Items) == 0 {
		items = items + "**None**"
	} else {
		for i, item := range pokemon.Items {
			if i == 0 {
				items = items + "**" + strings.Title(item.Item.Name) + "**"
			} else {
				items = items + " and **" + strings.Title(item.Item.Name) + "**"
			}
		}
	}

	// Create ability string
	abilities := "List of abilities: "
	if len(pokemon.Items) == 0 {
		abilities = abilities + "**None**"
	} else {
		for i, ability := range pokemon.Abilities {
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
	for _, pokemonStat := range pokemon.Stats {
		fields = append(fields, &discordgo.MessageEmbedField{
			Name:   "**" + strings.Title(pokemonStat.Stat.Name) + "**",
			Value:  strconv.Itoa(pokemonStat.BaseStat) + " (" + strconv.Itoa(pokemonStat.Effort) + " EV)",
			Inline: true,
		})
	}

	// Create type string
	types := ""
	if len(pokemon.Types) == 1 {
		types = types + "**" + strings.Title(pokemon.Types[0].Type.Name) + "**"
	} else {
		for i, pokemonType := range pokemon.Types {
			if i == 0 {
				types = types + "**" + strings.Title(pokemonType.Type.Name) + "**"
			} else {
				types = types + " and **" + strings.Title(pokemonType.Type.Name) + "**"
			}
		}
	}

	// Create the embed and message and send
	embed := &discordgo.MessageEmbed{
		Title: strings.Title(pokemon.Name) + " (#" + strconv.Itoa(pokemon.ID) + ")",
		Color: pokemontools.TypeColour(pokemon.Types[0].Type.Name),
		Description: "In around **" + strconv.Itoa(len(pokemon.Games)) + "** Pokemon games\n\n" +
			weight + height + types + "\n\n" +
			exp + "\n" +
			items + "\n" +
			abilities,
		Fields: fields,
		Image: &discordgo.MessageEmbedImage{
			URL: "https://www.smogon.com/dex/media/sprites/xy/" + pokemon.Name + ".gif",
		},
	}
	formRegex, _ := regexp.Compile(`(.+)-.+`)
	if formRegex.MatchString(pokemon.Name) {
		embed.URL = "https://bulbapedia.bulbagarden.net/wiki/" + strings.Title(formRegex.FindStringSubmatch(pokemon.Name)[1]) + "_(Pok%C3%A9mon)"
	} else {
		embed.URL = "https://bulbapedia.bulbagarden.net/wiki/" + strings.Title(pokemon.Name) + "_(Pok%C3%A9mon)"
	}

	s.ChannelMessageSendEmbed(m.ChannelID, embed)
}
