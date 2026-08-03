package main

import (
	"encoding/json"
	"errors"
	"os"
	"sort"
)

func sortSaveList(index fileIndex) ([]string, error) {
	var keys []string
	if len(index) == 0 {
		return keys, errors.New("unable to sort empty list")
	}

	for k := range index {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	keys = append(keys, "[New Save]")
	return keys, nil
}

func loadIndex() (fileIndex, error) {
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

func saveIndex(index fileIndex) error {
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
