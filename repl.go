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
	"math/rand"
	"github.com/listlessneko/pokedexcli/internal/pokecache"
)

type cliCommand struct {
	name        string
	description string
	callback    func(*config, []string) error
}

type config struct {
	Previous *string
	Next *string
	Cache *pokecache.Cache
	Caught map[string]Pokemon
}

type locationAreaResp struct {
	Next *string `json:"next"`
	Previous *string `json:"previous"`
	Results []struct {
		Name string `json:"name"`
	} `json:"results"`
}

type locationAreaDetailResp struct {
	PokemonEncounters []struct {
		Pokemon struct {
			Name string `json:"name"`
		} `json:"pokemon"`
	} `json:"pokemon_encounters"`
}

type Pokemon struct {
	Name string `json:"name"`
	BaseExperience int `json:"base_experience"`
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
	"explore": {
		name: "explore",
		description: "Get a list of the pokemon in this location area",
		callback: commandExplore,
	},
	"catch": {
		name: "catch",
		description: "Try to catch a pokemon",
		callback: commandCatch,
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
		Caught: make(map[string]Pokemon),
	}

	for {
		fmt.Print("Pokedex > ")

		var userInput []string
		if scanner.Scan() {
			userInput = cleanInput(scanner.Text())
		}

		command, exists := commands[userInput[0]]
		if exists {
			err := command.callback(cfg, userInput[1:])
			if err != nil {
				fmt.Println(err)
			}
		} else {
			fmt.Println("Unknown command")
		}
	}
}

func commandHelp(cfg *config, args []string) error {
	fmt.Println("Welcome to the Pokedex!\nUsage:\n\nhelp: Displays a help message\nexit: Exit the Pokedex")
	return nil
}

func commandMap(cfg *config, args []string) error {
	var url string
	if cfg.Next == nil {
		url = "https://pokeapi.co/api/v2/location-area/" 
	} else {
		url = *cfg.Next
	}

	b, ok := cfg.Cache.Get(url)
	if !ok {
		resp, err := http.Get(url)
		if err != nil {
			return err
		}

		defer resp.Body.Close()
		b, err = io.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		cfg.Cache.Add(url, b)
	}

	var result locationAreaResp
	err := json.Unmarshal(b, &result)
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

func commandMapB(cfg *config, args []string) error {
	var url string
	if cfg.Previous == nil {
		fmt.Println("you're on the first page")
		return nil
	} else {
		url = *cfg.Previous
	}

	b, ok := cfg.Cache.Get(url)
	if !ok {
		resp, err := http.Get(url)
		if err != nil {
			return err
		}

		defer resp.Body.Close()
		b, err = io.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		cfg.Cache.Add(url, b)
	}

	var result locationAreaResp
	err := json.Unmarshal(b, &result)
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

func commandExplore(cfg *config, args []string) error {
	if len(args) == 0 {
		fmt.Println("no area provided")
		return nil
	}

	base_url := "https://pokeapi.co/api/v2/location-area/"
	area_url := base_url + args[0]

	b, ok := cfg.Cache.Get(area_url)
	if !ok {
		resp, err := http.Get(area_url)
		if err != nil {
			return err
		}

		defer resp.Body.Close()
		b, err = io.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		cfg.Cache.Add(area_url, b)
	}

	var result locationAreaDetailResp
	err := json.Unmarshal(b, &result)
	if err != nil {
		return err
	}

	for _, r := range result.PokemonEncounters {
		fmt.Println(r.Pokemon.Name)
	}

	return nil
}

func commandCatch(cfg *config, args []string) error {
	if len(args) == 0 {
		fmt.Println("no pokemon provided")
		return nil
	}

	url := fmt.Sprintf("https://pokeapi.co/api/v2/pokemon/%s", args[0])

	b, ok := cfg.Cache.Get(url)
	if !ok {
		resp, err := http.Get(url)
		if err != nil {
			return err
		}

		defer resp.Body.Close()
		b, err = io.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		cfg.Cache.Add(url, b)
	}

	var result Pokemon
	err := json.Unmarshal(b, &result)
	if err != nil {
		return err
	}

	fmt.Printf("Throwing a Pokeball at %s...\n", result.Name)
	if rand.Intn(result.BaseExperience) < 50 {
		cfg.Caught[result.Name] = result
		fmt.Printf("You caught %s!\n", result.Name)
	} else {
		fmt.Printf("%s ran away...\n", result.Name)
	}

	return nil
}

func commandExit(cfg *config, args []string) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func cleanInput(text string) []string {
	return strings.Fields(strings.ToLower(text))
}
