package main

import (
	"encoding/json"
	"fmt"
	"github.com/uncomfyhalomacro/pokedexcli/internal/pokecache"
	"log"
	"net/http"
	"os"
)

const baseURL = "https://pokeapi.co/api/v2"

var pkCache = pokecache.DefaultPokeCache()

func runSupportedCommand(config *Config, cmd string, args ...string) error {
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
