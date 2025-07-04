package main

import (
	"encoding/json"
	"fmt"
	"os"
)

func load_settings() {
	data, f_err := os.ReadFile("./config.json")
	if f_err != nil {
		fmt.Println(f_err)
	}
	err := json.Unmarshal(data, &settings)
	if err != nil {
		panic(fmt.Errorf("Failed to load settings: %v", err))
	}
}
func set_settings(settings Settings) {
	marshal, err := json.Marshal(settings)
	if err != nil {
		panic(err)
	}
	_ = os.WriteFile("./config.json", marshal, os.ModePerm)
}
