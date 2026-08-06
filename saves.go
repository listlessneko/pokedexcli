package main

import (
	"encoding/json"
	"os"
)

func loadProfilesIndex() (profilesIndex, error) {
	var index profilesIndex

	data, err := os.ReadFile("saves/index.json")
	if !os.IsNotExist(err) {
		if err == nil {
			err = json.Unmarshal(data, &index)
		}
		if err != nil {
			return index, err
		}
	} else {
		index = make(profilesIndex)
	}
	return index, nil
}

func saveProfilesIndex(index profilesIndex) error {
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
