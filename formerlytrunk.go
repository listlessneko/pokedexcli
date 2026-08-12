package main

import (
	"os"
)

func openYard() (*os.File, error) {
	err := os.MkdirAll("yard", 0755)
	if err != nil {
		return nil, err
	}
	trunk, err := os.OpenFile("yard/trunk.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return nil, err
	}
	trunkLog := []byte("used cut\n")
	trunk.Write(trunkLog)
	return trunk, nil
}
