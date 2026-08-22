package main

import (
	"encoding/json"
	"fmt"
	"io"
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
		"profile": {
			name:        "profile",
			usage:       "profile",
			description: "Manage your profiles.",
			callback:    commandProfile,
		},
		"save": {
			name:        "save",
			usage:       "save",
			description: "Save current session.",
			callback:    commandSave,
		},
		"delete": {
			name:        "delete",
			usage:       "delete",
			description: "Delete a profile.",
			callback:    commandDelete,
		},
		"switch-profiles": {
			name:        "switch-profiles",
			usage:       "switch-profiles",
			description: "Switch to a different profile.",
			callback:    commandSwitchProfiles,
		},
		"change-name": {
			name:        "change-name",
			usage:       "change-name",
			description: "Change your profile's name (this also changes the prompt).",
			callback:    commandChangeProfileName,
		},
		"exit": {
			name:        "exit",
			usage:       "exit",
			description: "Exit the Pokedex.",
			callback:    commandExit,
		},
	}
}

func commandHelp(app *appState) error {
	fmt.Fprintln(app.writer, "Welcome to the Pokedex!\nCommands:")

	commands := getCommands()

	keys := make([]string, 0, len(commands))
	for k := range commands {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		command := commands[k]
		fmt.Fprintf(app.writer, "Command: %s\nUsage: %s\n%s\n\n", command.name, command.usage, command.description)
	}
	return nil
}

func commandMap(app *appState) error {
	var url string
	if app.cache.Next == nil {
		url = "https://pokeapi.co/api/v2/location-area/"
	} else {
		url = *app.cache.Next
	}

	b, ok := app.cache.Cache.Get(url)
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
		app.cache.Cache.Add(url, b)
	}

	var locations locationAreaResp
	err := json.Unmarshal(b, &locations)
	if err != nil {
		return err
	}

	for _, r := range locations.Results {
		fmt.Fprintln(app.writer, r.Name)
	}

	app.cache.Next = locations.Next
	app.cache.Previous = locations.Previous

	return nil
}

func commandMapB(app *appState) error {
	var url string
	if app.cache.Previous == nil {
		fmt.Fprintln(app.writer, "You're on the first page.")
		return nil
	} else {
		url = *app.cache.Previous
	}

	b, ok := app.cache.Cache.Get(url)
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
		app.cache.Cache.Add(url, b)
	}

	var locations locationAreaResp
	err := json.Unmarshal(b, &locations)
	if err != nil {
		return err
	}

	for _, r := range locations.Results {
		fmt.Fprintln(app.writer, r.Name)
	}

	app.cache.Next = locations.Next
	app.cache.Previous = locations.Previous

	return nil
}

func commandExplore(app *appState) error {
	if len(app.args) == 0 {
		fmt.Fprintln(app.writer, "Please provide a valid area.")
		return nil
	}

	base_url := "https://pokeapi.co/api/v2/location-area/"
	area_url := base_url + app.args[0]

	b, ok := app.cache.Cache.Get(area_url)
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
		app.cache.Cache.Add(area_url, b)
	}

	var location locationAreaDetailResp
	err := json.Unmarshal(b, &location)
	if err != nil {
		return err
	}

	for _, r := range location.PokemonEncounters {
		fmt.Fprintln(app.writer, capitalize(r.Pokemon.Name))
	}

	return nil
}

func commandCatch(app *appState) error {
	logger := app.log.logger

	if len(app.args) == 0 {
		fmt.Fprintln(app.writer, "Please provide a valid Pokemon.")
		return nil
	}

	species, err := getPokemonSpecies(app)
	if err != nil {
		return err
	}

	speciesLogger := logger.BranchGroup("species").Branch("name", species.Name, "capture_rate", species.CaptureRate)
	speciesLogger.Bark("species")

	pokemonName := capitalize(species.Name)
	captureRate := species.CaptureRate

	selectedPokeball, err := selectPrompt(pokeballs)
	if err != nil {
		return err
	}

	if selectedPokeball == catchBack {
		return nil
	}

	captureRate = captureRateByPokeball(captureRate, selectedPokeball)
	captureRate = trueCaptureRate(captureRate)

	fmt.Fprintf(app.writer, "You throw a %s at %s...\n", selectedPokeball, pokemonName)

	catchLogger := logger.BranchGroup("catch").Branch("pokemon", species.Name, "capture_rate", species.CaptureRate, "pokeball", selectedPokeball, "true_capture_rate", captureRate)

	if !isCaptured(captureRate) {
		catchLogger.Bark("ran away...")
		fmt.Fprintf(app.writer, "%s ran away...\n", pokemonName)
	} else {
		catchLogger.Bark("caught!")
		// TODO: Might not be needed when battle system is implemented since the battle system should already know this.
		// Maybe already stored in appState
		pokemon, err := getPokemon(app)
		if  err != nil {
			return err
		}
		app.config.Caught[pokemon.Name] = pokemon
		fmt.Fprintf(app.writer, "You caught %s!\n", pokemonName)
		fmt.Fprintln(app.writer, "You can view more information about this Pokemon with the inspect command.")
	}
	return nil
}

func commandInspect(app *appState) error {
	if len(app.args) == 0 {
		fmt.Fprintln(app.writer, "Please provide a valid Pokemon.")
		return nil
	}

	pokemon, ok := app.config.Caught[app.args[0]]
	if !ok {
		fmt.Fprintln(app.writer, "You have not caught that Pokemon.")
		return nil
	}

	fmt.Fprintf(app.writer, "Name: %s\n", capitalize(pokemon.Name))
	fmt.Fprintf(app.writer, "Height: %d\n", pokemon.Height)
	fmt.Fprintf(app.writer, "Weight: %d\n", pokemon.Weight)

	fmt.Fprintln(app.writer, "Stats:")
	for _, s := range pokemon.Stats {
		fmt.Fprintf(app.writer, "- %s: %d\n", s.Stat.Name, s.BaseStat)
	}

	fmt.Fprintln(app.writer, "Types:")
	for _, t := range pokemon.Types {
		fmt.Fprintf(app.writer, "- %s\n", capitalize(t.Type.Name))
	}

	return nil
}

func commandPokedex(app *appState) error {
	fmt.Fprintln(app.writer, "Your Pokedex:")

	if len(app.config.Caught) == 0 {
		fmt.Fprintln(app.writer, "You have no Pokemon...")
		return nil
	}

	for _, p := range app.config.Caught {
		fmt.Fprintf(app.writer, "- %s\n", capitalize(p.Name))
	}
	return nil
}

func commandProfile(app *appState) error {
	fmt.Fprintln(app.writer, "Manage your profiles:")

	const (
		profileSave       = "Save"
		profileDelete     = "Delete"
		profileSwitch     = "Switch"
		profileChangeName = "Change Name"
		profileBack       = "[Back]"
	)

	profileCommandsPrompts := []string{
		profileSave,
		profileDelete,
		profileSwitch,
		profileChangeName,
		profileBack,
	}

	selectedPrompt, err := selectPrompt(profileCommandsPrompts)
	if err != nil {
		return err
	}

	if selectedPrompt == profileBack {
		return nil
	}

	profileCommands := map[string]command{
		profileSave:       commandSave,
		profileDelete:     commandDelete,
		profileSwitch:     commandSwitchProfiles,
		profileChangeName: commandChangeProfileName,
	}

	selectedCommand := profileCommands[selectedPrompt]
	err = selectedCommand(app)
	if err != nil {
		return err
	}

	return nil
}

func commandChangeProfileName(app *appState) error {
	input := newInputState()
	var newName string
	var err error

	for {
		input.prompt = "New profile name: "
		newName, err = readLine(input)
		if err != nil {
			return err
		}
		if newName != "" {
			break
		}
	}

	newSaveFile := save_ify(newName)
	index, err := loadProfilesIndex()
	if err != nil {
		return err
	}
	oldSaveFile := app.config.SaveFile

	err = os.Rename(oldSaveFile, newSaveFile)
	if err != nil {
		return err
	}

	delete(index, app.config.Name)
	index[newName] = sanitizeInput(newName)
	saveProfilesIndex(index)

	app.config.Name = newName
	app.config.SaveFile = newSaveFile
	app.config.Prompt = promptify(newName)

	data, err := json.Marshal(app.config)
	err = os.WriteFile(app.config.SaveFile, data, 0644)
	if err != nil {
		return err
	}
	fmt.Fprintln(app.writer, "New profile name saved.")
	return nil
}

func commandSave(app *appState) error {
	data, err := json.Marshal(app.config)
	if err != nil {
		return err
	}

	if app.config.SaveFile == "" {
		input := newInputState()
		var name string

		for {
			input.prompt = "New profile: "
			name, err = readLine(input)
			if err != nil {
				return err
			}
			if name != "" {
				break
			}
		}

		app.config.Name = name
		app.config.SaveFile = save_ify(name)
		app.config.Prompt = promptify(name)
		index, err := loadProfilesIndex()
		if err != nil {
			return err
		}
		index[name] = sanitizeInput(name)
		saveProfilesIndex(index)
	}
	err = os.WriteFile(app.config.SaveFile, data, 0644)
	if err != nil {
		app.config.SaveFile = ""
		return err
	}
	fmt.Fprintln(app.writer, "Pokedex saved.")
	return nil
}

func commandDelete(app *appState) error {
	fmt.Fprintln(app.writer, "Delete a profile:")
	index, err := loadProfilesIndex()
	if err != nil {
		return err
	}
	if len(index) > 0 {
		keys := justTheKeys(index)
		prompts := sortSlices(keys, asc)
		prompts = append(prompts, "[Back]")
		for {
			selectedPrompt, err := selectPrompt(prompts)
			if err != nil {
				return err
			}

			selectedFilename := index[selectedPrompt]
			if selectedFilename != "" {
				selectedSaveFile := save_ify(selectedFilename)
				if selectedSaveFile == app.config.SaveFile {
					fmt.Fprintln(app.writer, "You cannot delete the file you're currently in.")
					continue
				}
				err = os.Remove(selectedSaveFile)
				if !os.IsNotExist(err) {
					if err == nil {
						delete(index, selectedPrompt)
						err = saveProfilesIndex(index)
						if err != nil {
							fmt.Fprintln(app.writer, err.Error())
						}
						fmt.Fprintln(app.writer, "Profile deleted.")
						return nil
					}
					fmt.Fprintln(app.writer, "There was an error deleting this profile.")
					continue
				}
				fmt.Fprintln(app.writer, "Profile does not exist.")
				continue
			}
			break
		}
	}
	return nil
}

func commandSwitchProfiles(app *appState) error {
	fmt.Fprintln(app.writer, "Manage your profiles:")
	index, err := loadProfilesIndex()
	if err != nil {
		return err
	}
	if len(index) > 0 {
		keys := justTheKeys(index)
		prompts := sortSlices(keys, asc)
		prompts = append(prompts, "[Back]")
		for {
			selectedPrompt, err := selectPrompt(prompts)
			if err != nil {
				return err
			}

			saveFilename := index[selectedPrompt]
			if saveFilename != "" {
				saveFile := "saves/" + saveFilename + ".json"
				if saveFile == app.config.SaveFile {
					fmt.Fprintln(app.writer, "You're already using this profile.")
					continue
				}
				data, err := os.ReadFile(saveFile)
				if !os.IsNotExist(err) {
					var tempConfig config
					if err == nil {
						err = json.Unmarshal(data, &tempConfig)
					}
					if err != nil {
						fmt.Fprintln(app.writer, "There was an error extracting data from this profile.")
						fmt.Fprintln(app.writer, err.Error())
						continue
					}
					*app.config = tempConfig
					app.config.SaveFile = saveFile
					fmt.Fprintf(app.writer, "Switching profile to %s...\n", selectedPrompt)
					return nil
				}
				fmt.Fprintln(app.writer, "Profile does not exist.")
				continue
			}
			break
		}
	}
	return nil
}

func commandExit(app *appState) error {
	fmt.Fprintln(app.writer, "Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}
