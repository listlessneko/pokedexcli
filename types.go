package main

import (
	"github.com/listlessneko/pokedexcli/internal/pokecache"
	"github.com/listlessneko/pokedexcli/timber"
	"io"
)

type profilesIndex map[string]string
type sortOrder bool
type command func(*appState) error

type log struct {
	logger  *timber.Logger
	leveler timber.Leveler
}

type appState struct {
	log    *log
	config *config
	cache  *cache
	writer io.Writer
	args   []string
}

type cliCommand struct {
	name        string
	usage       string
	description string
	callback    command
}

type config struct {
	Name     string
	SaveFile string
	Prompt   string
	Caught   map[string]Pokemon
}

type cache struct {
	Cache    *pokecache.Cache
	Previous *string
	Next     *string
}

type trieChildren map[rune]*trieNode

type trieNode struct {
	Children trieChildren
	End      bool
}

type trie struct {
	Root *trieNode
}

type inputState struct {
	prompt   string
	history  []string
	wordTrie *trie
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

type PokemonSpecies struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

type Pokemon struct {
	Name           string         `json:"name"`
	BaseExperience int            `json:"base_experience"`
	Height         int            `json:"height"`
	Weight         int            `json:"weight"`
	Stats          []PokemonStats `json:"stats"`
	Types          []PokemonTypes `json:"types"`
	Species        PokemonSpecies `json:"species"`
}
