package main

import "fmt"
import "strings"
import "bufio"
import "os"
import "net/http"
import "encoding/json"
import "log"
import "github.com/uncomfyhalomacro/pokecache"
import "io"

const baseURL = "https://pokeapi.co/api/v2"

var pkCache = pokecache.DefaultPokeCache()

type Config struct {
	Next     string
	Previous string
}

type cliCommand struct {
	name        string
	description string
	callback    func(config *Config, args ...string) error
}

type LocationAreas struct {
	Count    int               `json:"count"`
	Next     any               `json:"next"`     // NOTE: This can be null. If you are at the last page, the value is null since there are no "next" pages
	Previous any               `json:"previous"` // NOTE: This can be null. If you are at the first page, the value is null since there are no "previous" pages
	Results  []LocationDetails `json:"results"`
}

type LocationDetails struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

// namespacing will cause an initialisation cycle.
var supportedCommands = map[string]cliCommand{}

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

func checkPokeCache(key string) ([]byte, bool) {
	return pkCache.Get(key)
}

func addToCache(key string, newData []byte) bool {
	pkCache.Add(key, newData)
	// Check
	_, ok := checkPokeCache(key)
	return ok
}

func mapNextPage(config *Config, _ ...string) error {
	var locationData LocationAreas
	fullURL := baseURL + "/location-area"
	Previous := &(*config).Previous
	Next := &(*config).Next
	if *Next != "" {
		fullURL = *Next
	}
	cachedData, ok := checkPokeCache(*Next)
	if !ok {
		resp, err := http.Get(fullURL)
		if err != nil {
			return fmt.Errorf("error, there was a problem getting map information: %w\nStatus: %s", err, resp.Status)
		}
		defer resp.Body.Close()

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

		byteData, err := io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("error, failed to convert body as byte data")
		}
		addToCache(fullURL, byteData)
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
		addToCache(fullURL, cachedData)
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
	cachedData, ok := checkPokeCache(*Next)
	if !ok {
		resp, err := http.Get(fullURL)
		if err != nil {
			return fmt.Errorf("error, there was a problem getting map information: %w\nStatus: %s", err, resp.Status)
		}
		defer resp.Body.Close()

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
		byteData, err := io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("error, failed to convert body as byte data")
		}
		addToCache(fullURL, byteData)
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
		addToCache(fullURL, cachedData)
	}
	return nil
}
func cleanInput(text string) []string {
	return strings.Fields(strings.ToLower(text))
}

func main() {
	config := &Config{
		Next:     "",
		Previous: "",
	}
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
	}

	userInput := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("Pokedex > ")
		if userInput.Scan() {
			receivedInput := userInput.Text()
			cleanedInput := cleanInput(receivedInput)
			firstWord := cleanedInput[0]
			err := runSupportedCommand(config, firstWord, cleanedInput[1:]...)
			if err != nil {
				fmt.Println(err)
			}
		}
	}
}
