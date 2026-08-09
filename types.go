package main

import (
	"github.com/listlessneko/pokedexcli/internal/pokecache"
	"io"
)

type profilesIndex map[string]string

type cliCommand struct {
	name        string
	usage       string
	description string
	callback    func(*config, *cache, io.Writer, []string) error
}

type config struct {
	SaveFile string
	Prompt   string
	Caught   map[string]Pokemon
}

type cache struct {
	Cache    *pokecache.Cache
	Previous *string
	Next     *string
}

type children map[rune]*trieNode

type trieNode struct {
	Children children
	End      bool
}

type trie struct {
	Root *trieNode
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

type Pokemon struct {
	Name           string         `json:"name"`
	BaseExperience int            `json:"base_experience"`
	Height         int            `json:"height"`
	Weight         int            `json:"weight"`
	Stats          []PokemonStats `json:"stats"`
	Types          []PokemonTypes `json:"types"`
}
