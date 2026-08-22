package main

import "math/rand"

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

func captureRateByPokeball(captureRate int, selected string) int {
	switch selected {
	case pokeball:
	case greatball:
		captureRate = int(float64(captureRate) * 1.5)
	case ultraball:
		captureRate = int(float64(captureRate) * 2)
	default:
	}
	return captureRate
}

func trueCaptureRate(captureRate int) int {
	// TODO: introduce dynamic, random stats
	// hp stubs
	scaledMaxHP := 100 * 3
	scaledCurrentHP := 100 * 2
	x := ((scaledMaxHP - scaledCurrentHP) * captureRate) / scaledMaxHP
	return x
}

func isCaptured(captureRate int) bool {
	return rand.Intn(maxCatchRoll) <= captureRate
}
