package main

import (
	"encoding/json"
	"os"
)

func loadSettings() {
	file, err := os.ReadFile("config.json")
	if err != nil {
		settings = Settings{Recent: []string{}, Settings: map[string]string{}}
		if os.IsNotExist(err) {
			// Create config file with default values if it does not exist
			setSettings(settings)
		}
		return
	}
	json.Unmarshal(file, &settings)
}

func setSettings(settings Settings) {
	file, _ := json.MarshalIndent(settings, "", " ")
	_ = os.WriteFile("config.json", file, 0644)
}
