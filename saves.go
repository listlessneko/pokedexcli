package main

import (
	"encoding/json"
	"os"
)

func loadSavesIndex() (fileIndex, error) {
	var index fileIndex

	data, err := os.ReadFile("saves/index.json")
	if !os.IsNotExist(err) {
		if err == nil {
			err = json.Unmarshal(data, &index)
		}
		if err != nil {
			return index, err
		}
	} else {
		index = make(fileIndex)
	}
	return index, nil
}

func saveSavesIndex(index fileIndex) error {
	err := os.MkdirAll("saves", 0755)
	if err != nil {
		return err
	}
	data, err := json.Marshal(index)
	if err != nil {
		return err
	}
	err = os.WriteFile("saves/index.json", data, 0644)
	if err != nil {
		return err
	}
	return nil
}
