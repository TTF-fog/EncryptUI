package main

import (
	"encoding/json"
	"os"
)

func loadSettings() {
	file, err := os.ReadFile("config.json")
	if err != nil {
		settings = Settings{Recent: []string{}, Settings: map[string]string{}}
		return
	}
	json.Unmarshal(file, &settings)
}

func setSettings(settings Settings) {
	settings.Recent = settings.Recent[1:]
	file, _ := json.MarshalIndent(settings, "", " ")
	_ = os.WriteFile("config.json", file, 0644)
}
