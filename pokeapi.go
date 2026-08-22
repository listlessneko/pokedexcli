package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func fetch(url string, app *appState) ([]byte, error) {
	logger := app.log.logger.BranchGroup("fetch").Branch("url", url)
	logger.Bark("fetching...")

	b, ok := app.cache.Cache.Get(url)
	if !ok {
		resp, err := http.Get(url)
		if err != nil {
			return b, err
		}

		defer resp.Body.Close()
		b, err = io.ReadAll(resp.Body)
		if err != nil {
			return b, err
		}
		app.cache.Cache.Add(url, b)
	}
	logger.Bark("fetched")
	return b, nil
}

func getPokemonSpecies(app *appState) (PokemonSpecies, error) {
	url := fmt.Sprintf("https://pokeapi.co/api/v2/pokemon-species/%s", app.args[0])

	var species PokemonSpecies

	b, err := fetch(url, app)
	if err != nil {
		return species, err
	}

	if string(b) == "Not Found" {
		fmt.Fprintln(app.writer, "Please provide a valid Pokemon.")
		return species, nil
	}

	err = json.Unmarshal(b, &species)
	if err != nil {
		return species, err
	}
	return species, nil
}

func getPokemon(app *appState) (Pokemon, error) {
	url := fmt.Sprintf("https://pokeapi.co/api/v2/pokemon/%s", app.args[0])

	var pokemon Pokemon

	b, err := fetch(url, app)
	if err != nil {
		return pokemon, err
	}

	if string(b) == "Not Found" {
		fmt.Fprintln(app.writer, "Please provide a valid Pokemon.")
		return pokemon, nil
	}

	err = json.Unmarshal(b, &pokemon)
	if err != nil {
		return pokemon, err
	}
	return pokemon, nil
}

