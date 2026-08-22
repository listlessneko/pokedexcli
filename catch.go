package main

import (
	"github.com/listlessneko/pokedexcli/timber"
	"math/rand"
)

const (
	maxCatchRoll = 256
	pokeball     = "Poke Ball"
	greatball    = "Great Ball"
	ultraball    = "Ultra Ball"
	catchBack    = "[Back]"
)

var pokeballs = []string{
	pokeball,
	greatball,
	ultraball,
	catchBack,
}

func captureRateByPokeball(selected string) float64 {
	modifier := float64(1)
	switch selected {
	case pokeball:
	case greatball:
		modifier = 1.5
	case ultraball:
		modifier = 2
	default:
	}
	return modifier
}

func trueCaptureRate(logger *timber.Logger, captureRate int, selected string) (int, *timber.Logger) {
	pokeballModifier := captureRateByPokeball(selected)
	captureRateModifiedByPokeball := int(float64(captureRate) * pokeballModifier)
	logger = logger.Branch("pokeballModifier", pokeballModifier, "captureRateModifiedByPokeball", captureRateModifiedByPokeball)
	// TODO: introduce dynamic, random stats
	// hp stubs
	scaledMaxHP := 100 * 3
	scaledCurrentHP := 100 * 2
	x := ((scaledMaxHP - scaledCurrentHP) * captureRateModifiedByPokeball) / scaledMaxHP
	return x, logger
}

func isCaptured(logger *timber.Logger, captureRate int, selected string) (bool, *timber.Logger) {
	captureRate, logger = trueCaptureRate(logger, captureRate, selected)
	roll := rand.Intn(maxCatchRoll)
	logger = logger.Branch("finalCaptureRate", captureRate, "roll", roll)
	return roll <= captureRate, logger
}
