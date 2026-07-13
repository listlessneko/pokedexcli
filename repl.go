package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/listlessneko/pokedexcli/internal/pokecache"
	"io"
	"math/rand"
	"net/http"
	"os"
	"sort"
	"strings"
	"time"
	"unicode"
)

type cliCommand struct {
	name        string
	usage       string
	description string
	callback    func(*config, io.Writer, []string) error
}

type config struct {
	Previous *string
	Next     *string
	Cache    *pokecache.Cache
	Caught   map[string]Pokemon
}

type locationAreaResp struct {
	Next     *string `json:"next"`
	Previous *string `json:"previous"`
	Results  []struct {
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
	URL  string `json:"url"`
}

type PokemonStats struct {
	BaseStat int              `json:"base_stat"`
	Effort   int              `json:"effort"`
	Stat     NamedAPIResource `json:"stat"`
}

type PokemonTypes struct {
	Slot int              `json:"slot"`
	Type NamedAPIResource `json:"type"`
}

type Pokemon struct {
	Name           string         `json:"name"`
	BaseExperience int            `json:"base_experience"`
	Height         int            `json:"height"`
	Weight         int            `json:"weight"`
	Stats          []PokemonStats `json:"stats"`
	Types          []PokemonTypes `json:"types"`
}

func getCommands() map[string]cliCommand {
	return map[string]cliCommand{
		"help": {
			name:        "help",
			usage:       "help",
			description: "Displays a list of commands.",
			callback:    commandHelp,
		},
		"map": {
			name:        "map",
			usage:       "map",
			description: "Displays a list of location areas.",
			callback:    commandMap,
		},
		"mapb": {
			name:        "mapb",
			usage:       "mapb",
			description: "Displays a list of the previous location areas.",
			callback:    commandMapB,
		},
		"explore": {
			name:        "explore",
			usage:       "explore <location-area>",
			description: "Displays a list of Pokemon in the location area.",
			callback:    commandExplore,
		},
		"catch": {
			name:        "catch",
			usage:       "catch <pokemon-name>",
			description: "Try to catch a Pokemon.",
			callback:    commandCatch,
		},
		"inspect": {
			name:        "inspect",
			usage:       "inspect <pokemon-name>",
			description: "Display information about a Pokemon you caught.",
			callback:    commandInspect,
		},
		"pokedex": {
			name:        "pokedex",
			usage:       "pokedex",
			description: "Displays a list of Pokemon you caught.",
			callback:    commandPokedex,
		},
		"exit": {
			name:        "exit",
			usage:       "exit",
			description: "Exit the Pokedex.",
			callback:    commandExit,
		},
	}
}

func cleanInput(text string) []string {
	return strings.Fields(strings.ToLower(text))
}

func capitalize(s string) string {
	if s == "" {
		return s
	}

	final := ""
	for w := range strings.SplitSeq(strings.TrimSpace(s), " ") {
		r := []rune(w)
		r[0] = unicode.ToUpper(r[0])
		final = final + string(r) + " "
	}
	return strings.TrimSpace(final)
}

func startRepl() {
	scanner := bufio.NewScanner(os.Stdin)

	cfg := &config{
		Cache:  pokecache.NewCache(5 * time.Second),
		Caught: make(map[string]Pokemon),
	}

	for {
		fmt.Print("Pokedex > ")

		if !scanner.Scan() {
			fmt.Println()
			break
		}

		var userInput []string
		userInput = cleanInput(scanner.Text())

		if len(userInput) == 0 {
			continue
		}

		commands := getCommands()
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

	if err := scanner.Err(); err != nil {
		fmt.Println("error reading input:", err)
	}
}

func commandHelp(cfg *config, writer io.Writer, args []string) error {
	fmt.Fprintln(writer, "Welcome to the Pokedex!\nCommands:")

	commands := getCommands()

	keys := make([]string, 0, len(commands))
	for k := range commands {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		command := commands[k]
		fmt.Fprintf(writer, "Command: %s\nUsage: %s\n%s\n\n", command.name, command.usage, command.description)
	}
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

	var locations locationAreaResp
	err := json.Unmarshal(b, &locations)
	if err != nil {
		return err
	}

	for _, r := range locations.Results {
		fmt.Fprintln(writer, r.Name)
	}

	cfg.Next = locations.Next
	cfg.Previous = locations.Previous

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

	var locations locationAreaResp
	err := json.Unmarshal(b, &locations)
	if err != nil {
		return err
	}

	for _, r := range locations.Results {
		fmt.Fprintln(writer, r.Name)
	}

	cfg.Next = locations.Next
	cfg.Previous = locations.Previous

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

	var location locationAreaDetailResp
	err := json.Unmarshal(b, &location)
	if err != nil {
		return err
	}

	for _, r := range location.PokemonEncounters {
		fmt.Fprintln(writer, capitalize(r.Pokemon.Name))
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

	var pokemon Pokemon
	err := json.Unmarshal(b, &pokemon)
	if err != nil {
		return err
	}

	pokemonName := capitalize(pokemon.Name)
	chance := rand.Intn(pokemon.BaseExperience)

	fmt.Fprintf(writer, "Throwing a Pokeball at %s...\n", pokemonName)
	if chance < 50 {
		cfg.Caught[pokemon.Name] = pokemon
		fmt.Fprintf(writer, "You caught %s!\n", pokemonName)
		fmt.Fprintln(writer, "You man now inspect it with the inspect command.")
	} else {
		fmt.Fprintf(writer, "%s ran away...\n", pokemonName)
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

	fmt.Fprintf(writer, "Name: %s\n", capitalize(pokemon.Name))
	fmt.Fprintf(writer, "Height: %d\n", pokemon.Height)
	fmt.Fprintf(writer, "Weight: %d\n", pokemon.Weight)

	fmt.Fprintln(writer, "Stats:")
	for _, s := range pokemon.Stats {
		fmt.Fprintf(writer, "- %s: %d\n", s.Stat.Name, s.BaseStat)
	}

	fmt.Fprintln(writer, "Types:")
	for _, t := range pokemon.Types {
		fmt.Fprintf(writer, "- %s\n", capitalize(t.Type.Name))
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
		fmt.Fprintf(writer, "- %s\n", capitalize(p.Name))
	}
	return nil
}

func commandExit(cfg *config, writer io.Writer, args []string) error {
	fmt.Fprintln(writer, "Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}
