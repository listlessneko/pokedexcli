package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/listlessneko/pokedexcli/internal/pokecache"
	"golang.org/x/term"
	"io"
	"math/rand"
	"net/http"
	"os"
	"sort"
	"strings"
	"time"
	"unicode"
)

const (
	bufSize      = 3
	newLine      = '\x0a'
	keyEnter     = '\x0d'
	keyBackspace = '\x7f'
	keyCtrlC     = '\x03'
	keyCtrlD     = '\x04'
	keyEscape    = '\x1b'
	keyLSqBrckt  = '\x5b'
	keyA         = '\x41'
	keyB         = '\x42'
	keyC         = '\x43'
	keyD         = '\x44'
	enterSeq     = "\x0d\x0a"
	eraseSeq     = "\x1b[K"
	cursorFwd    = "\x1b[C"
	cursorBckwd  = "\x1b[D"
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
		"save": {
			name: "save",
			usage: "save",
			description: "Save current session.",
			callback: commandSave,
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

func redraw(new []byte, cursor int) {
	if cursor > 0 {
		seq := fmt.Sprintf("\x1b[%dD", cursor)
		os.Stdout.Write([]byte(seq))
	}
	os.Stdout.Write([]byte(eraseSeq))
	os.Stdout.Write(new)
}

func readLine(prompt string, history []string) (string, error) {
	os.Stdout.Write([]byte(prompt))

	fd := int(os.Stdin.Fd())
	oldState, err := term.MakeRaw(fd)
	if err != nil {
		return "", fmt.Errorf("error enabling raw mode: %w", err)
	}

	defer term.Restore(fd, oldState)

	historyIndex := len(history)
	var cursor int
	var currentLine []byte
	buf := make([]byte, bufSize)
	for {
		n, err := os.Stdin.Read(buf)
		if err != nil {
			return "", fmt.Errorf("error reading bytes: %w", err)
		}
		for i := 0; i < n; i++ {
			b := buf[i]
			switch b {
			case keyCtrlC, keyCtrlD:
				return "", io.EOF
			case keyEnter:
				os.Stdout.Write([]byte(enterSeq))
				return string(currentLine), nil
			case keyBackspace:
				if cursor > 0 {
					copy(currentLine[cursor-1:], currentLine[cursor:])
					currentLine = currentLine[:len(currentLine)-1]
					cursor -= 1
					os.Stdout.Write([]byte(cursorBckwd))
					os.Stdout.Write([]byte(currentLine[cursor:]))
					os.Stdout.Write([]byte(" "))
					nCol := len(currentLine) - cursor + 1
					seq := fmt.Sprintf("\x1b[%dD", nCol)
					os.Stdout.Write([]byte(seq))
				}
				continue
			case keyEscape:
				if i+2 < n && buf[i+1] == keyLSqBrckt {
					switch buf[i+2] {
					case keyA:
						if historyIndex > 0 {
							historyIndex -= 1
							currentLine = []byte(history[historyIndex])
							redraw(currentLine, cursor)
							cursor = len(currentLine)
						}
					case keyB:
						if historyIndex < len(history) {
							historyIndex += 1
							if historyIndex == len(history) {
								currentLine = []byte("")
							} else {
								currentLine = []byte(history[historyIndex])
							}
							redraw(currentLine, cursor)
							cursor = len(currentLine)
						}
					case keyC:
						if cursor < len(currentLine) {
							cursor += 1
							os.Stdout.Write([]byte(cursorFwd))
						}
					case keyD:
						if cursor > 0 {
							cursor -= 1
							os.Stdout.Write([]byte(cursorBckwd))
						}
					}
					i += 2
				}
				continue
			}
			currentLine = append(currentLine, 0)
			copy(currentLine[cursor+1:], currentLine[cursor:])
			currentLine[cursor] = b
			cursor += 1
			os.Stdout.Write(currentLine[cursor-1:])
			nCol := len(currentLine) - cursor
			if nCol > 0 {
				seq := fmt.Sprintf("\x1b[%dD", nCol) 
				os.Stdout.Write([]byte(seq))
			}
		}
	}
}

func startRepl() {
	cfg := &config{
		Cache:  pokecache.NewCache(5 * time.Second),
		Caught: make(map[string]Pokemon),
	}

	prompt := "Pokedex > "
	var history []string

	for {

		line, err := readLine(prompt, history)
		if errors.Is(err, io.EOF) {
			os.Stdout.Write([]byte{newLine})
			break
		} else if err != nil {
			os.Stderr.Write([]byte(err.Error()))
			break
		}

		userInput := cleanInput(line)

		if len(userInput) == 0 {
			continue
		}

		history = append(history, line)

		commands := getCommands()
		command, exists := commands[userInput[0]]
		if exists {
			err := command.callback(cfg, os.Stdout, userInput[1:])
			if err != nil {
				os.Stderr.Write([]byte(err.Error()))
			}
		} else {
			os.Stdout.Write([]byte("unknown command\n"))
		}
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

func commandSave(cfg *config, writer io.Writer, args []string) error {
	data, err := json.Marshal(cfg.Caught)
	if err != nil {
		return err
	}

	err = os.WriteFile("pokedex.json", data, 0644)
	if err != nil {
		return err
	}
	fmt.Fprintln(writer, "Pokedex saved.")
	return nil
}

func commandExit(cfg *config, writer io.Writer, args []string) error {
	fmt.Fprintln(writer, "Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}
