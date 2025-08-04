package core

import (
	"encoding/json"
	"fmt"
	"github.com/uncomfyhalomacro/pokedexcli/internal/pokecache"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strings"
)

const baseURL = "https://pokeapi.co/api/v2"

var pkCache = pokecache.DefaultPokeCache()
var capturedPokemons = map[string]PokemonDetails{}

func RunSupportedCommand(config *Config, cmd string, args ...string) error {
	command, ok := supportedCommands[cmd]
	if !ok {
		return fmt.Errorf("command not found: %s\n", cmd)
	}
	callback := command.callback
	err := callback(config, args...)
	return err
}

func commandExit(_ *Config, _ ...string) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return fmt.Errorf("error, there was an error exiting the program")
}

func displayHelp(_ *Config, _ ...string) error {
	fmt.Printf("Welcome to the Pokedex!\nUsage:\n\n")
	var display string
	for commandName, fieldNames := range supportedCommands {
		display = fmt.Sprintf("%s: %s\n", commandName, fieldNames.description)
		fmt.Print(display)
	}
	return nil
}

func mapNextPage(config *Config, _ ...string) error {
	var locationData LocationAreas
	fullURL := baseURL + "/location-area"
	Previous := &(*config).Previous
	Next := &(*config).Next
	if *Next != "" {
		fullURL = *Next
	}
	cachedData, ok := (*pkCache).Get(*Next)
	if !ok {
		resp, err := http.Get(fullURL)
		if err != nil {
			return fmt.Errorf("error, there was a problem getting map information: %w\nStatus: %s", err, resp.Status)
		}

		if resp.StatusCode > 299 {
			log.Fatalf("Response failed with status code: %d and\nbody: %s\nfull url: %s", resp.StatusCode, resp.Body, fullURL)
		}

		decoder := json.NewDecoder(resp.Body)
		err = decoder.Decode(&locationData)

		if err != nil {
			return fmt.Errorf("error, there was a problem generating map information: %w\nStatus: %s", err, resp.Status)
		}

		if len(locationData.Results) == 0 || locationData.Results == nil {
			return fmt.Errorf("error, map location is empty! Status: %s", resp.Status)
		}
		for _, location := range locationData.Results {
			fmt.Println(location.Name)
		}
		nextList, ok := locationData.Next.(string)
		if !ok {
			nextList = baseURL + "/location-area" // round trip
		}

		*Next = nextList

		previousList, ok := locationData.Previous.(string)
		if !ok {
			previousList = ""
		}
		*Previous = previousList

		byteData, err := json.Marshal(locationData)
		if err != nil {
			return fmt.Errorf("error, failed to convert body as byte data")
		}
		(*pkCache).Add(fullURL, byteData)
		resp.Body.Close()
	} else {
		err := json.Unmarshal(cachedData, &locationData)

		if err != nil {
			return fmt.Errorf("error, there was a problem generating map information from cache: %w\n", err)
		}

		if len(locationData.Results) == 0 || locationData.Results == nil {
			return fmt.Errorf("error, map location is empty from cache!")
		}
		for _, location := range locationData.Results {
			fmt.Println(location.Name)
		}
		nextList, ok := locationData.Next.(string)
		if !ok {
			nextList = baseURL + "/location-area" // round trip
		}

		*Next = nextList

		previousList, ok := locationData.Previous.(string)
		if !ok {
			previousList = baseURL + "/location-area"
		}
		*Previous = previousList
		(*pkCache).Add(fullURL, cachedData)
	}
	return nil
}

func mapPreviousPage(config *Config, _ ...string) error {
	var locationData LocationAreas
	fullURL := baseURL + "/location-area"
	Previous := &(*config).Previous
	Next := &(*config).Next
	if *Previous != "" {
		fullURL = *Previous
	}
	cachedData, ok := (*pkCache).Get(*Next)
	if !ok {
		resp, err := http.Get(fullURL)
		if err != nil {
			return fmt.Errorf("error, there was a problem getting map information: %w\nStatus: %s", err, resp.Status)
		}

		if resp.StatusCode > 299 {
			log.Fatalf("Response failed with status code: %d and\nbody: %s\n", resp.StatusCode, resp.Body)
		}

		decoder := json.NewDecoder(resp.Body)
		err = decoder.Decode(&locationData)

		if err != nil {
			return fmt.Errorf("error, there was a problem generating map information: %w\nStatus: %s", err, resp.Status)
		}

		if len(locationData.Results) == 0 || locationData.Results == nil {
			return fmt.Errorf("error, map location is empty! Status: %s", resp.Status)
		}
		for _, location := range locationData.Results {
			fmt.Println(location.Name)
		}
		nextList, ok := locationData.Next.(string)
		if !ok {
			nextList = baseURL + "/location-area" // round trip
		}

		*Next = nextList

		previousList, ok := locationData.Previous.(string)
		if !ok {
			previousList = ""
		}
		*Previous = previousList
		byteData, err := json.Marshal(locationData)
		if err != nil {
			return fmt.Errorf("error, failed to convert body as byte data")
		}
		(*pkCache).Add(fullURL, byteData)
		resp.Body.Close()
	} else {

		err := json.Unmarshal(cachedData, &locationData)

		if err != nil {
			return fmt.Errorf("error, there was a problem generating map information from cache: %w\n", err)
		}

		if len(locationData.Results) == 0 || locationData.Results == nil {
			return fmt.Errorf("error, map location is empty from cache!")
		}
		for _, location := range locationData.Results {
			fmt.Println(location.Name)
		}
		nextList, ok := locationData.Next.(string)
		if !ok {
			nextList = baseURL + "/location-area" // round trip
		}

		*Next = nextList

		previousList, ok := locationData.Previous.(string)
		if !ok {
			previousList = baseURL + "/location-area"
		}
		*Previous = previousList
		(*pkCache).Add(fullURL, cachedData)
	}
	return nil
}

// This is just a wrapper around the `exploreArea` function. It receives many area as arguments
func exploreAreas(_ *Config, areas ...string) error {
	fullURL := baseURL + "/location-area/" // NOTE: Take note of the last slash. It's there so I don't need to append it during a for loop.
	for _, area := range areas {
		areaURL := fullURL + area
		fmt.Printf("Exploring %s...\n", area)
		err := exploreArea(areaURL)
		if err != nil {
			return err
		}
	}
	return nil
}

// This is the original caller
func exploreArea(url string) error {
	cachedData, ok := (*pkCache).Get(url)
	if !ok {
		var areaData LocationEncounterDetails
		resp, err := http.Get(url)
		if err != nil {
			return fmt.Errorf("error, there was a problem getting pokemon list information: %w\nStatus: %s", err, resp.Status)
		}

		if resp.StatusCode > 299 {
			log.Fatalf("Response failed with status code: %d and\nbody: %s\n", resp.StatusCode, resp.Body)
		}

		decoder := json.NewDecoder(resp.Body)
		err = decoder.Decode(&areaData)

		if err != nil {
			return fmt.Errorf("error, there was a problem generating pokemon list information: %w\nStatus: %s", err, resp.Status)
		}
		for _, pokemonEncounter := range areaData.PokemonEncounters {
			fmt.Println(pokemonEncounter.Pokemon.Name)
		}
		byteData, err := json.Marshal(areaData.PokemonEncounters)
		if err != nil {
			return fmt.Errorf("error, failed to convert body as byte data")
		}
		(*pkCache).Add(url, byteData)
		resp.Body.Close()
	} else {
		var pokemonEncounters []PokemonEncounter
		err := json.Unmarshal(cachedData, &pokemonEncounters)

		if err != nil {
			return fmt.Errorf("error, there was a problem generating pokemon list information from cache: %w\n", err)
		}

		if len(pokemonEncounters) == 0 || pokemonEncounters == nil {
			return fmt.Errorf("error, pokemon list is empty from cache!")
		}

		for _, pokemonEncounter := range pokemonEncounters {
			fmt.Println(pokemonEncounter.Pokemon.Name)
		}
		(*pkCache).Add(url, cachedData)
	}

	return nil
}

func catchPokemon(_ *Config, args ...string) error {
	if len(args) > 1 {
		return fmt.Errorf("error, only needs 1 argument\n")
	}

	if len(args) == 0 {
		return fmt.Errorf("error, please provide a pokemon name or ID\n")
	}

	pokemon, err := fetchPokemonDetail(args[0])
	if err != nil {
		return err
	}
	fmt.Printf("Throwing a Pokeball at %s...\n", pokemon.Name)
	getChance := rand.Intn(pokemon.BaseExperience + 30)
	if getChance >= pokemon.BaseExperience {
		fmt.Printf("You have caught %s! 🎉\n", pokemon.Name)
		capturedPokemons[args[0]] = pokemon
	} else {
		fmt.Printf("%s escaped and ran away! 😩\n", pokemon.Name)
	}
	return nil
}

func pokedex(_ *Config, _ ...string) error {
	fmt.Println("Your Pokedex:")
	if len(capturedPokemons) == 0 {
		return fmt.Errorf("Your Pokedex is empty... Try capuring a pokemon first.\n")
	}
	for k, _ := range capturedPokemons {
		fmt.Printf("  - %s\n", k)
	}
	return nil

}

func inspect(_ *Config, pokemonNames ...string) error {
	if len(capturedPokemons) == 0 {
		return fmt.Errorf("Your Pokedex is empty... Try capuring a pokemon first.\n")
	}
	if len(pokemonNames) == 0 {
		for k, _ := range capturedPokemons {
			pokemonNames = append(pokemonNames, k)
		}
	}
	for _, pokemonName := range pokemonNames {
		pokemon, ok := capturedPokemons[pokemonName]
		if !ok {
			_, err := fetchPokemonDetail(pokemonName)
			if err != nil {
				fmt.Printf("This pokemon species does not exist.\n")

			} else {
				fmt.Printf("It seems you have not captured %s yet.\n", pokemonName)
			}
		} else {
			var stats []string
			var types []string
			for _, stat := range pokemon.Stats {
				stats = append(stats, fmt.Sprintf("  -%s: %d", stat.Stat.Name, stat.BaseStat))
			}
			for _, type_ := range pokemon.Types {
				types = append(types, fmt.Sprintf("  - %s", type_.Type.Name))
			}

			details := fmt.Sprintf(`Name: %s
Height: %d
Weight: %d
Stats:
%s
Types:
%s
`, pokemon.Name, pokemon.Height, pokemon.Weight, strings.Join(stats, "\n"), strings.Join(types, "\n"))
			fmt.Println(details)

		}
	}
	return nil
}

func fetchPokemonDetail(pokemonNameOrId string) (PokemonDetails, error) {
	var pokemon PokemonDetails
	fullURL := baseURL + "/pokemon/" + pokemonNameOrId
	cachedData, ok := (*pkCache).Get(fullURL)
	if !ok {
		resp, err := http.Get(fullURL)
		if err != nil {
			return PokemonDetails{}, fmt.Errorf("error, there was a problem getting pokemon information: %w\nStatus: %s", err, resp.Status)
		}

		if resp.StatusCode > 299 {
			return PokemonDetails{}, fmt.Errorf("Response failed with status code: %d and\nbody: %s\nPokemon species with name or ID, %s, does not exist.\n", resp.StatusCode, resp.Body, pokemonNameOrId)
		}

		decoder := json.NewDecoder(resp.Body)
		err = decoder.Decode(&pokemon)

		if err != nil {
			return PokemonDetails{}, fmt.Errorf("error, there was a problem generating pokemon list information: %w\nStatus: %s", err, resp.Status)
		}
		byteData, err := json.Marshal(pokemon)
		if err != nil {
			return PokemonDetails{}, fmt.Errorf("error, failed to convert body as byte data")
		}
		(*pkCache).Add(fullURL, byteData)
		resp.Body.Close()
		if err != nil {
			return PokemonDetails{}, fmt.Errorf("error, failed to convert body as byte data")
		}
		(*pkCache).Add(fullURL, byteData)
		resp.Body.Close()

	} else {
		err := json.Unmarshal(cachedData, &pokemon)

		if err != nil {
			return PokemonDetails{}, fmt.Errorf("error, there was a problem fetching pokemon from cache: %w\n", err)
		}

		(*pkCache).Add(fullURL, cachedData)

	}
	return pokemon, nil

}
