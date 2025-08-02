package main

import "fmt"
import "strings"
import "bufio"
import "os"

type cliCommand struct {
	name        string
	description string
	callback    func(args ...string) error
}

// Initialise the global map variable. Map cannot be constants, unfortunate...
var supportedCommands = map[string]cliCommand{}

func runSupportedCommand(cmd string, args ...string) error {

	callback := supportedCommands[cmd].callback
	err := callback(args...)
	return err
}

func commandExit(_ ...string) error {
    	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return fmt.Errorf("error, there was an error exiting the program")
}

func displayHelp(_ ...string) error {
	fmt.Println("Welcome to the Pokedex!\nUsage:\n\n")
	var display string
	for commandName, fieldNames := range supportedCommands {
		display = fmt.Sprintf("%s: %s\n", commandName, fieldNames.description)
		fmt.Print(display)
	}
	return nil
}

func cleanInput(text string) []string {
	return strings.Fields(strings.ToLower(text))
}

func main() {

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
	}

	userInput := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("Pokedex > ")
		if userInput.Scan() {
			receivedInput := userInput.Text()
			cleanedInput := cleanInput(receivedInput)
			firstWord := cleanedInput[0]
			err := runSupportedCommand(firstWord, cleanedInput[1:]...)
			if err != nil {
				fmt.Println(err)
			}
		}
	}
}
