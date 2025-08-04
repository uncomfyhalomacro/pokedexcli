package core

var supportedCommands = map[string]cliCommand{}

func init() {
	supportedCommands = map[string]cliCommand{
		"exit": {
			name:        "exit",
			description: "Exit the Pokedex",
			callback:    commandExit,
		},
		"help": {
			name:        "help",
			description: "Displays a help message",
			callback:    displayHelp,
		},
		"map": {
			name:        "map",
			description: "Displays the next list of locations of the Pokemon World!",
			callback:    mapNextPage,
		},
		"mapb": {
			name:        "map",
			description: "Displays the previous list of locations of the Pokemon World!",
			callback:    mapPreviousPage,
		},
		"explore": {
			name:        "explore",
			description: "Display the list of pokemon species in each area. It can receive multiple areas as arguments.",
			callback:    exploreAreas,
		},
		"catch": {
			name:        "catch",
			description: "Attempt to catch a pokemon species with your imaginary pokeball. Don't cry when you fail.",
			callback:    catchPokemon,
		},
		"inspect": {
			name:        "inspect",
			description: "Inspect captured pokemon or pokemons in your Pokedex.",
			callback:    inspect,
		},
		"pokedex": {
			name:        "pokedex",
			description: "Get the list of pokemons you have in your Pokedex!",
			callback:    pokedex,
		},
	}
}
