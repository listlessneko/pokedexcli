package main

import (
	"fmt"
	"strings"
	"bufio"
	"os"
	"net/http"
	"io"
	"encoding/json"
	"time"
	"github.com/listlessneko/pokedexcli/internal/pokecache"
)

type cliCommand struct {
	name        string
	description string
	callback    func(*config) error
}

type config struct {
	Previous *string
	Next *string
	Cache *pokecache.Cache
}

type locationAreaResp struct {
	Next *string `json:"next"`
	Previous *string `json:"previous"`
	Results []struct {
		Name string `json:"name"`
	} `json:"results"`
}

var commands = map[string]cliCommand {
	"help": {
		name: "help",
		description: "Displays a help message",
		callback: commandHelp,
	},
	"map": {
		name: "map",
		description: "Get a list of location areas",
		callback: commandMap,
	},
	"mapb": {
		name: "mapb",
		description: "Get a list of the previous location areas",
		callback: commandMapB,
	},
	"exit": {
		name: "exit",
		description: "Exit the Pokedex",
		callback: commandExit,
	},
}

func startRepl() {
	scanner := bufio.NewScanner(os.Stdin)

	cfg := &config{
		Cache: pokecache.NewCache(5 * time.Second),
	}

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

func commandMap(cfg *config) error {
	var url string
	if cfg.Next == nil {
		url = "https://pokeapi.co/api/v2/location-area/" 
	} else {
		url = *cfg.Next
	}

	resp, err := http.Get(url)
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	var result locationAreaResp
	err = json.Unmarshal(b, &result)
	if err != nil {
		return err
	}

	for _, r := range result.Results {
		fmt.Println(r.Name)
	}

	cfg.Next = result.Next
	cfg.Previous = result.Previous

	return nil
}

func commandMapB(cfg *config) error {
	var url string
	if cfg.Previous == nil {
		fmt.Println("you're on the first page")
		return nil
	} else {
		url = *cfg.Previous
	}

	resp, err := http.Get(url)
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	var result locationAreaResp
	err = json.Unmarshal(b, &result)
	if err != nil {
		return err
	}

	for _, r := range result.Results {
		fmt.Println(r.Name)
	}

	cfg.Next = result.Next
	cfg.Previous = result.Previous

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
