package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

func configDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("cannot determine home directory: %w", err)
	}
	return filepath.Join(home, ".config", "dotpet"), nil
}

func petFile() (string, error) {
	dir, err := configDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "pet.json"), nil
}

func LoadPet() (*Pet, error) {
	path, err := petFile()
	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(path)
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

	validatePet(&p)
	return &p, nil
}

// validatePet clamps fields to valid ranges, fixing corrupted or tampered data.
func validatePet(p *Pet) {
	// Species must be valid
	valid := false
	for _, s := range speciesList {
		if p.Species == s {
			valid = true
			break
		}
	}
	if !valid {
		p.Species = speciesList[0]
	}

	// Level: 1..maxLevel
	if p.Level < 1 {
		p.Level = 1
	}
	if p.Level > maxLevel {
		p.Level = maxLevel
	}

	// Generation: >= 1
	if p.Generation < 1 {
		p.Generation = 1
	}

	// Non-negative fields
	if p.Power < 0 {
		p.Power = 0
	}
	if p.XP < 0 {
		p.XP = 0
	}
	if p.Wins < 0 {
		p.Wins = 0
	}
	if p.Losses < 0 {
		p.Losses = 0
	}
	if p.TotalXP < 0 {
		p.TotalXP = 0
	}
	if p.ItemsFound < 0 {
		p.ItemsFound = 0
	}
	if p.Legacy < 0 {
		p.Legacy = 0
	}

	// Inventory cap
	const maxInventory = 1000
	if len(p.Inventory) > maxInventory {
		p.Inventory = p.Inventory[len(p.Inventory)-maxInventory:]
	}
}

func SavePet(p *Pet) error {
	path, err := petFile()
	if err != nil {
		return err
	}
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0700); err != nil {
		return err
	}

	data, err := json.MarshalIndent(p, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0600)
}
