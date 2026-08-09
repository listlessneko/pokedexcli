package main

import (
	"encoding/json"
	"errors"
	"github.com/listlessneko/pokedexcli/internal/pokecache"
	"io"
	"os"
	"time"
)

func startRepl() {
	app := &app{}
	app.config = &config{
		Prompt: "Pokedex > ",
		Caught: make(map[string]Pokemon),
	}
	app.cache = &cache{
		Cache: pokecache.NewCache(60 * time.Second),
	}

	index, err := loadProfilesIndex()
	if err != nil {
		os.Stderr.Write([]byte(err.Error()))
	}
	if len(index) > 0 {
		keys:= justTheKeys(index)
		prompts:= sortSlices(keys, asc)
		prompts = append(prompts, "[New Save]")
		selectedPrompt, err := selectPrompt(prompts)
		if err != nil {
			os.Stderr.Write([]byte(err.Error()))
		}

		saveFilename := index[selectedPrompt]
		if saveFilename != "" {
			app.config.SaveFile = "saves/" + saveFilename + ".json"
		}
	}

	data, err := os.ReadFile(app.config.SaveFile)
	if !os.IsNotExist(err) {
		if err == nil {
			err = json.Unmarshal(data, app.config)
		}
		if err != nil {
			os.Stderr.Write([]byte(err.Error()))
		}
	}

	var history []string
	autocomplete := newTrie()
	commands := getCommands()
	for command, _ := range commands {
		autocomplete.add(command)
	}

	for {
		line, err := readLine(app.config.Prompt, history, autocomplete)
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

		command, exists := commands[userInput[0]]
		if exists {
			app.writer = os.Stdout
			app.args = userInput[1:]
			err := command.callback(app)
			if err != nil {
				os.Stderr.Write([]byte(err.Error() + "\n"))
			}
		} else {
			os.Stdout.Write([]byte("unknown command\n"))
		}
	}
}
