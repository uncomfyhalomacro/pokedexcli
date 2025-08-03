package main

import (
	"bufio"
	"fmt"
	"os"
	"github.com/uncomfyhalomacro/pokedexcli/internal/core"
)

func main() {
	config := &core.Config{
		Next:     "",
		Previous: "",
	}
	userInput := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("Pokedex > ")
		if userInput.Scan() {
			receivedInput := userInput.Text()
			cleanedInput := core.CleanInput(receivedInput)
			firstWord := cleanedInput[0]
			err := core.RunSupportedCommand(config, firstWord, cleanedInput[1:]...)
			if err != nil {
				fmt.Println(err)
			}
		}
	}
}
