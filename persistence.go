package main

import (
	"encoding/json"
	"os"
	"path/filepath"
)

func configDir() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".config", "dotpet")
}

func petFile() string {
	return filepath.Join(configDir(), "pet.json")
}

func LoadPet() (*Pet, error) {
	data, err := os.ReadFile(petFile())
	if err != nil {
		if os.IsNotExist(err) {
			p := NewPet(1, 0)
			return p, SavePet(p)
		}
		return nil, err
	}

	var p Pet
	if err := json.Unmarshal(data, &p); err != nil {
		return nil, err
	}

	return &p, nil
}

func SavePet(p *Pet) error {
	dir := configDir()
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	data, err := json.MarshalIndent(p, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(petFile(), data, 0644)
}
