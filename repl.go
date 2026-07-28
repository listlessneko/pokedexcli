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

	index, err := loadIndex()
	if err != nil {
		os.Stderr.Write([]byte(err.Error()))
	}
	if len(index) > 0 {
		saveFile, err := selectPrompt(index)
		if err != nil {
			os.Stderr.Write([]byte(err.Error()))
		}
		if saveFile != "" {
			cfg.SaveFile = "saves/" + saveFile + ".json"
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
				os.Stderr.Write([]byte(err.Error() + "\n"))
			}
		} else {
			os.Stdout.Write([]byte("unknown command\n"))
		}
	}
}
