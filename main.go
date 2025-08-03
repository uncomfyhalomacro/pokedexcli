package main

import (
	"bufio"
	"fmt"
	"os"
)

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
	}
}

func main() {
	config := &Config{
		Next:     "",
		Previous: "",
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
