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
	callback    func(*config, io.Writer, []string) error
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

type NamedAPIResource struct {
	Name string `json:"name"`
	URL string `json:"url"`
}

type PokemonStats struct {
	BaseStat int `json:"base_stat"`
	Effort int `json:"effort"`
	Stat NamedAPIResource `json:"stat"`
}

type PokemonTypes struct {
	Slot int `json:"slot"`
	Type NamedAPIResource `json:"type"`
}

type Pokemon struct {
	Name string `json:"name"`
	BaseExperience int `json:"base_experience"`
	Height int `json:"height"`
	Weight int `json:"weight"`
	Stats []PokemonStats `json:"stats"`
	Types []PokemonTypes `json:"types"`
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
	"inspect": {
		name: "inspect",
		description: "View a pokemon's stats",
		callback: commandInspect,
	},
	"pokedex": {
		name: "pokedex",
		description: "View your caught pokemon",
		callback: commandPokedex,
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
			err := command.callback(cfg, os.Stdout, userInput[1:])
			if err != nil {
				fmt.Println(err)
			}
		} else {
			fmt.Println("Unknown command")
		}
	}
}

func commandHelp(cfg *config, writer io.Writer, args []string) error {
	fmt.Fprintln(writer, "Welcome to the Pokedex!\nUsage:\n\nhelp: Displays a help message\nexit: Exit the Pokedex")
	return nil
}

func commandMap(cfg *config, writer io.Writer, args []string) error {
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
		fmt.Fprintln(writer, r.Name)
	}

	cfg.Next = result.Next
	cfg.Previous = result.Previous

	return nil
}

func commandMapB(cfg *config, writer io.Writer, args []string) error {
	var url string
	if cfg.Previous == nil {
		fmt.Fprintln(writer, "You're on the first page.")
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
		fmt.Fprintln(writer, r.Name)
	}

	cfg.Next = result.Next
	cfg.Previous = result.Previous

	return nil
}

func commandExplore(cfg *config, writer io.Writer, args []string) error {
	if len(args) == 0 {
		fmt.Fprintln(writer, "Please provide a valid area.")
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
		fmt.Fprintln(writer, r.Pokemon.Name)
	}

	return nil
}

func commandCatch(cfg *config, writer io.Writer, args []string) error {
	if len(args) == 0 {
		fmt.Fprintln(writer, "Please provide a valid Pokemon.")
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

	fmt.Fprintf(writer, "Throwing a Pokeball at %s...\n", result.Name)
	if rand.Intn(result.BaseExperience) < 50 {
		cfg.Caught[result.Name] = result
		fmt.Fprintf(writer, "You caught %s!\n", result.Name)
		fmt.Fprintln(writer, "You man now inspect it with the inspect command.")
	} else {
		fmt.Fprintf(writer, "%s ran away...\n", result.Name)
	}

	return nil
}

func commandInspect(cfg *config, writer io.Writer, args []string) error {
	if len(args) == 0 {
		fmt.Fprintln(writer, "Please provide a valid Pokemon.")
		return nil
	}

	pokemon, ok := cfg.Caught[args[0]]
	if !ok {
		fmt.Fprintln(writer, "You have not caught that Pokemon.")
		return nil
	}

	fmt.Fprintf(writer, "Name: %s\n", pokemon.Name)
	fmt.Fprintf(writer, "Height: %d\n", pokemon.Height)
	fmt.Fprintf(writer, "Weight: %d\n", pokemon.Weight)

	fmt.Fprintln(writer, "Stats:")
	for _, s := range pokemon.Stats {
		fmt.Fprintf(writer, "- %s: %d\n", s.Stat.Name, s.BaseStat)
	}

	fmt.Fprintln(writer, "Types:")
	for _, t := range pokemon.Types {
		fmt.Fprintf(writer, "- %s\n", t.Type.Name)
	}
	
	return nil
}

func commandPokedex(cfg *config, writer io.Writer, args []string) error {
	fmt.Fprintln(writer, "Your Pokedex:")

	if len(cfg.Caught) == 0 {
		fmt.Fprintln(writer, "You have no Pokemon...")
		return nil
	}

	for _, p := range cfg.Caught {
		fmt.Fprintf(writer, "- %s\n", p.Name)
	}
	return nil
}

func commandExit(cfg *config, writer io.Writer, args []string) error {
	fmt.Fprintln(writer, "Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func cleanInput(text string) []string {
	return strings.Fields(strings.ToLower(text))
}
