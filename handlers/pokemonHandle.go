package handlers

import (
	"github.com/bwmarrin/discordgo"
	pokemoncommands "maquiaBot/handlers/pokemon-commands"
)

// PokemonHandle handles commands that are regarding pokemon
func PokemonHandle(s *discordgo.Session, m *discordgo.MessageCreate, args []string, serverPrefix string) {
	// Check if any args were even given
	if len(args) > 1 {
		mainArg := args[1]
		switch mainArg {
		case "b", "berry":
			go pokemoncommands.Berry(s, m)
		default:
			go pokemoncommands.Pokemon(s, m)
		}
	} else {
		s.ChannelMessageSend(m.ChannelID, "Please specify a command! Check `"+serverPrefix+"help` for more details!")
	}
}
