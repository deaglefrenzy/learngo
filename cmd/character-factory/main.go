package main

import (
	"encoding/json"
	"fmt"
	"go_tutorial/models"
	"math/rand"
	"os"
	"path/filepath"
)

func CharFactory() *models.Character {
	classes := []string{"paladin", "archer", "mage"}
	randomClass := classes[rand.Intn(len(classes))]

	names := []string{
		"Aer", "Aeth", "Arc", "Bel", "Crys", "Drak", "Eld", "Fael", "Glim", "Hald",
		"Ith", "Jor", "Kael", "Lor", "Myr", "Niv", "Onyx", "Pyr", "Quil", "Riven",
		"Syl", "Thal", "Umbr", "Verid", "Wyn", "Xyl", "Zel", "Vor", "Ulf", "Tyr",
		"Aev", "Alar", "Brym", "Corv", "Dorn", "Ebon", "Fyr", "Grend", "Hael", "Ivor",
		"Jarr", "Kryll", "Lyth", "Morv", "Nyx", "Orin", "Phos", "Rhaed", "Syr", "Tael",
		"Vex", "Wryn", "Xanth", "Ymir", "Zalt", "Aure", "Bryn", "Cael", "Drev", "Eryn",
		"Falen", "Gyr", "Hest", "Ild", "Jareth", "Kaelen", "Lorn", "Mav", "Nym", "Oryn",
		"Prax", "Ryl", "Saev", "Taryn", "Valen", "Wyrn", "Xylo", "Yvaine", "Zoren",
		"Avar", "Bael", "Cyr", "Dael", "Evar", "Faelar", "Gryff", "Hyld", "Izar",
		"Jorv", "Krynn", "Lyrian", "Morwyn", "Nyr", "Othyr", "Pryth", "Rhaev", "Sarv",
		"Thyr", "Varyn", "Wrynn", "Xyron", "Yvren", "Zyl", "Aethyr", "Brev", "Cyril",
		"Drevyn", "Eryvan", "Faelyn", "Gryth", "Hyran", "Izaryn", "Jorev", "Kryvan",
		"Lyrevan", "Morvyn", "Nyran", "Othyran", "Pryvan", "Rhaevan", "Sarvan", "Thyran",
		"Varyvan", "Wryvan", "Xyryvan", "Yvryvan", "Zylvan",
	}

	random1 := rand.Intn(len(names))
	random2 := rand.Intn(len(names))
	name := names[random1] + " " + names[random2]
	health := (rand.Intn(10) + 6) * 10
	attack := rand.Intn(15) + 5

	switch randomClass {
	case "paladin":
		return &models.Character{
			BaseStatus: models.BaseStatus{
				Name:   name,
				Health: health,
				Attack: attack,
			},
			Class:  "paladin",
			Shield: rand.Intn(10) + 2,
		}
	case "archer":
		return &models.Character{
			BaseStatus: models.BaseStatus{
				Name:   name,
				Health: health,
				Attack: attack,
			},
			Class:    "archer",
			Critical: rand.Intn(15) + 5,
		}
	case "mage":
		return &models.Character{
			BaseStatus: models.BaseStatus{
				Name:   name,
				Health: health,
				Attack: attack,
			},
			Class: "mage",
			Mana:  rand.Intn(20) + 10,
		}
	default:
		return nil
	}
}

func CharSeeder(count int, filename string) error {

	characters := make([]models.Character, 0, count)

	for i := 0; i < count; i++ {
		char := CharFactory()
		characters = append(characters, *char)
	}

	dataJSON, err := json.MarshalIndent(characters, "", "  ")
	if err != nil {
		return fmt.Errorf("error marshaling JSON: %w", err)
	}

	filepath := filepath.Join("database", filename)
	err = os.WriteFile(filepath, dataJSON, 0644)
	if err != nil {
		return fmt.Errorf("error writing JSON to file: %w", err)
	}

	fmt.Printf("%d characters added into %s.\n", len(characters), filename)
	return nil
}
