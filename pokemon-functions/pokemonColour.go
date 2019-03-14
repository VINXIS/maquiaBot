package pokemontools

// TypeColour assigns a colour based on the pokemon's type (Colours referred to via https://bulbapedia.bulbagarden.net/wiki/Category:Type_color_templates)
func TypeColour(pokemonType string) (Colour int) {
	switch pokemonType {
	case "fighting":
		Colour = 0xC03028
	case "flying":
		Colour = 0xA890F0
	case "poison":
		Colour = 0xA040A0
	case "ground":
		Colour = 0xE0C068
	case "rock":
		Colour = 0xB8A038
	case "bug":
		Colour = 0xA8B820
	case "ghost":
		Colour = 0x705898
	case "steel":
		Colour = 0xB8B8D0
	case "fire":
		Colour = 0xF08030
	case "water":
		Colour = 0x6890F0
	case "grass":
		Colour = 0x78C850
	case "electric":
		Colour = 0xF8D030
	case "psychic":
		Colour = 0xF85888
	case "ice":
		Colour = 0x98D8D8
	case "dragon":
		Colour = 0x7038F8
	case "dark":
		Colour = 0x705848
	case "fairy":
		Colour = 0xEE99AC
	default:
		Colour = 0xA8A878
	}
	return Colour
}
