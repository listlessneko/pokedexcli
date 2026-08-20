package main

import (
	"io"
	"net/http"
)

func fetch(url string, app *appState) ([]byte, error) {
	logger := app.log.logger.BranchGroup("fetch")
	logger.Bark("fetching from pokeapi...")

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
	return b, nil
}
