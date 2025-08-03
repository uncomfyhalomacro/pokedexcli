package main

import (
    "fmt"
    "os"
    "net/http"
    "github.com/uncomfyhalomacro/pokedexcli/internal/pokecache"
    "encoding/json"
    "log"
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
