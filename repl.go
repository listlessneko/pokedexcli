package main

import (
	"encoding/json"
	"errors"
	"github.com/listlessneko/pokedexcli/internal/pokecache"
	"github.com/listlessneko/pokedexcli/timber"
	"io"
	"os"
	"time"
)

func selectProfile(app *appState) {
	index, err := loadProfilesIndex()
	if err != nil {
		os.Stderr.Write([]byte(err.Error()))
	}
	if len(index) > 0 {
		keys := justTheKeys(index)
		prompts := sortSlices(keys, asc)
		prompts = append(prompts, "[New Save]")
		selectedPrompt, err := selectPrompt(prompts)
		if err != nil {
			os.Stderr.Write([]byte(err.Error()))
		}

		saveFilename := index[selectedPrompt]
		if saveFilename != "" {
			app.config.SaveFile = save_ify(saveFilename)
		}
	}
}

func loadProfile(app *appState) {
	data, err := os.ReadFile(app.config.SaveFile)
	if !os.IsNotExist(err) {
		if err == nil {
			err = json.Unmarshal(data, app.config)
		}
		if err != nil {
			os.Stderr.Write([]byte(err.Error()))
		}
	}
}

func beginTheLoop(app *appState) {
	input := newInputState()
	commands := getCommands()
	for command, _ := range commands {
		input.wordTrie.add(command)
	}

	for {
		input.prompt = app.config.Prompt
		line, err := readLine(input)
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

		input.history = append(input.history, line)

		command, exists := commands[userInput[0]]
		if !exists {
			os.Stdout.Write([]byte("unknown command\n"))
			continue
		}
		app.writer = os.Stdout
		app.args = userInput[1:]
		err = command.callback(app)
		if err != nil {
			os.Stderr.Write([]byte(err.Error() + "\n"))
			continue
		}
		continue
	}
}

func startRepl() {
	app := &appState{}
	app.log = &log{}
	app.log.leveler = timber.NewLeveler()
	timberLeveler := app.log.leveler
	timberLeveler.Set(timber.LevelBark)

	file, err := timber.NewTimber()
	if err != nil {
		os.Stderr.Write([]byte(err.Error() + "\n"))
	}

	timberHandler := timber.NewHandler(file, timberLeveler)
	app.log.logger = timber.NewLogger(timberHandler)
	logger := app.log.logger
	logger.Bark("Cutting timber...")

	app.config = &config{
		Prompt: "Pokedex > ",
		Caught: make(map[string]Pokemon),
	}
	app.cache = &cache{
		Cache: pokecache.NewCache(60 * time.Second),
	}

	selectProfile(app)
	loadProfile(app)
	beginTheLoop(app)
}
