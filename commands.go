package main

import (
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
	"sort"
)

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
			name:        "save",
			usage:       "save",
			description: "Save current session.",
			callback:    commandSave,
		},
		"switch-profiles": {
			name:        "switch-profiles",
			usage:       "switch-profiles",
			description: "Switch to a different profile.",
			callback:    commandSwitchProfiles,
		},
		"delete": {
			name:        "delete",
			usage:       "delete",
			description: "Delete a save file.",
			callback:    commandDelete,
		},
		"exit": {
			name:        "exit",
			usage:       "exit",
			description: "Exit the Pokedex.",
			callback:    commandExit,
		},
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

	if cfg.SaveFile == "" {
		prompt := "New save file: "
		var history []string
		var name string
		for {
			name, err = readLine(prompt, history)
			if err != nil {
				return err
			}
			if name != "" {
				break
			}
		}
		cfg.SaveFile = "saves/" + sanitizeInput(name) + ".json"
		index, err := loadSavesIndex()
		if err != nil {
			return err
		}
		index[name] = sanitizeInput(name)
		saveSavesIndex(index)
	}
	err = os.WriteFile(cfg.SaveFile, data, 0644)
	if err != nil {
		cfg.SaveFile = ""
		return err
	}
	fmt.Fprintln(writer, "Pokedex saved.")
	return nil
}

func commandSwitchProfiles(cfg *config, writer io.Writer, args []string) error {
	index, err := loadSavesIndex()
	if err != nil {
		return err
	}
	if len(index) > 0 {
		keys, err := justTheKeys(index)
		if err != nil {
			return err
		}

		prompts, err := sortSlices(keys, true)
		if err != nil {
			return err
		}

		prompts = append(prompts, "[Back]")
		for {
			selectedPrompt, err := selectPrompt(prompts)
			if err != nil {
				return err
			}

			saveFilename := index[selectedPrompt]
			if saveFilename != "" {
				saveFile := "saves/" + saveFilename + ".json"
				if saveFile == cfg.SaveFile {
					fmt.Fprintln(writer, "You're already using this profile.")
					continue
				}
				data, err := os.ReadFile(saveFile)
				if !os.IsNotExist(err) {
					caughtCopy := make(map[string]Pokemon, len(cfg.Caught))
					for k, v := range cfg.Caught {
						caughtCopy[k] = v
					}
					if err == nil {
						clear(cfg.Caught)
						err = json.Unmarshal(data, &cfg.Caught)
						cfg.Prompt = "[" + selectedPrompt + "] Pokedex > "
					}
					if err != nil {
						cfg.Caught = caughtCopy
						fmt.Fprintln(writer, "There was an error extracting data from this profile.")
						fmt.Fprintln(writer, err.Error())
						continue
					}
					cfg.SaveFile = saveFile
					fmt.Fprintf(writer, "Switching profile to %s...\n", selectedPrompt)
					return nil
				}
				fmt.Fprintln(writer, "Profile does not exist.")
				continue
			}
			break
		}
	}
	return nil
}

func commandDelete(cfg *config, writer io.Writer, args []string) error {
	index, err := loadSavesIndex()
	if err != nil {
		return err
	}
	if len(index) > 0 {
		keys, err := justTheKeys(index)
		if err != nil {
			return err
		}

		prompts, err := sortSlices(keys, true)
		if err != nil {
			return err
		}

		prompts = append(prompts, "[Back]")
		for {
			selectedPrompt, err := selectPrompt(prompts)
			if err != nil {
				return err
			}

			saveFilename := index[selectedPrompt]
			if saveFilename != "" {
				saveFile := "saves/" + saveFilename + ".json"
				if saveFile == cfg.SaveFile {
					fmt.Fprintln(writer, "You cannot delete the file you're currently in.")
					continue
				}
				err = os.Remove(saveFile)
				if !os.IsNotExist(err) {
					if err == nil {
						delete(index, selectedPrompt)
						err = saveSavesIndex(index)
						if err != nil {
							fmt.Fprintln(writer, err.Error())
						}
						fmt.Fprintln(writer, "File deleted.")
						return nil
					}
					if err != nil {
						fmt.Fprintln(writer, "There was an error deleting this file.")
						continue
					}
				}
				fmt.Fprintln(writer, "File does not exist.")
				continue
			}
			break
		}
	}
	return nil
}

func commandExit(cfg *config, writer io.Writer, args []string) error {
	fmt.Fprintln(writer, "Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}
