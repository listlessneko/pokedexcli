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
		Cache:  pokecache.NewCache(5 * time.Second),
		Caught: make(map[string]Pokemon),
	}

	index, err := loadSavesIndex()
	if err != nil {
		os.Stderr.Write([]byte(err.Error()))
	}
	if len(index) > 0 {
		keys, err := justTheKeys(index)
		if err != nil {
			os.Stderr.Write([]byte(err.Error()))
		}

		prompts, err := sortSlices(keys, true)
		if err != nil {
			os.Stderr.Write([]byte(err.Error()))
		}

		prompts = append(prompts, "[New Save]")
		selectedPrompt, err := selectPrompt(prompts)
		if err != nil {
			os.Stderr.Write([]byte(err.Error()))
		}

		saveFilename := index[selectedPrompt]
		if saveFilename != "" {
			cfg.SaveFile = "saves/" + saveFilename + ".json"
			cfg.Prompt = "[" + selectedPrompt + "] Pokedex > "
		}
	}

	data, err := os.ReadFile(cfg.SaveFile)
	if !os.IsNotExist(err) {
		if err == nil {
			err = json.Unmarshal(data, &cfg.Caught)
		}
		if err != nil {
			os.Stderr.Write([]byte(err.Error()))
		}
	}

	var history []string

	for {

		line, err := readLine(cfg.Prompt, history)
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
				os.Stderr.Write([]byte(err.Error() + "\n"))
			}
		} else {
			os.Stdout.Write([]byte("unknown command\n"))
		}
	}
}
