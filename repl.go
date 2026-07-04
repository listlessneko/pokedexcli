package main

import (
	"fmt"
	"strings"
	"bufio"
	"os"
)

type cliCommand struct {
	name        string
	description string
	callback    func(*config) error
}

type config struct {
	Previous *string
	Next *string
}

var commands = map[string]cliCommand {
	"help": {
		name: "help",
		description: "Displays a help message",
		callback: commandHelp,
	},
	"exit": {
		name: "exit",
		description: "Exit the Pokedex",
		callback: commandExit,
	},
}

func startRepl() {
	scanner := bufio.NewScanner(os.Stdin)

	cfg := &config{}

	for {
		fmt.Print("Pokedex > ")

		var userInput []string
		if scanner.Scan() {
			userInput = cleanInput(scanner.Text())
		}

		command, exists := commands[userInput[0]]
		if exists {
			err := command.callback(cfg)
			if err != nil {
				fmt.Println(err)
			}
		} else {
			fmt.Println("Unknown command")
		}
	}
}

func commandHelp(cfg *config) error {
	fmt.Println("Welcome to the Pokedex!\nUsage:\n\nhelp: Displays a help message\nexit: Exit the Pokedex")
	return nil
}

func commandExit(cfg *config) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func cleanInput(text string) []string {
	return strings.Fields(strings.ToLower(text))
}
