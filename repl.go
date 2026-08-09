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
	cfg := &config{
		Prompt: "Pokedex > ",
		Caught: make(map[string]Pokemon),
	}
	cache := &cache{
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
			cfg.SaveFile = "saves/" + saveFilename + ".json"
		}
	}

	data, err := os.ReadFile(cfg.SaveFile)
	if !os.IsNotExist(err) {
		if err == nil {
			err = json.Unmarshal(data, &cfg)
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
		line, err := readLine(cfg.Prompt, history, autocomplete)
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
			err := command.callback(cfg, cache, os.Stdout, userInput[1:])
			if err != nil {
				os.Stderr.Write([]byte(err.Error() + "\n"))
			}
		} else {
			os.Stdout.Write([]byte("unknown command\n"))
		}
	}
}
