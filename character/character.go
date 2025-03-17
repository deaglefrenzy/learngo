package character

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
)

type BaseStatus struct {
	Name   string
	Health int
	Attack int
}

type Character struct {
	BaseStatus BaseStatus `json:"baseStatus"`
	Class      string     `json:"class"`
	Shield     int        `json:"shield,omitempty"`
	Critical   int        `json:"critical,omitempty"`
	Mana       int        `json:"mana,omitempty"`
}

func (c *BaseStatus) GetName() string {
	return c.Name
}

func (c *BaseStatus) GetHealth() int {
	return c.Health
}

func (c *BaseStatus) SetHealth(health int) {
	c.Health = max(health, 0)
}

func (c *BaseStatus) CharStatus() {
	fmt.Println(c.Name, c.Health, c.Attack)
}

func (c *Character) GetClass() string {
	return c.Class
}

func LoadCharacter(filename string, result *Character) error {
	fileData, err := os.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("error loading JSON file: %w", err)
	}
	return json.Unmarshal(fileData, result)
}

func PrintCharacter(v Character) {
	formattedJSON, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		fmt.Println("Error formatting JSON:", err)
		return
	}
	fmt.Println(string(formattedJSON))
}

func SaveToJSON(filename string, data Character) error {
	dataJSON, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return fmt.Errorf("error marshaling JSON: %w", err)
	}

	filepath := filepath.Join("database", filename)
	err = os.WriteFile(filepath, dataJSON, 0644)
	if err != nil {
		return fmt.Errorf("error writing JSON to file: %w", err)
	}
	return err
}

func CharFactory() *Character {
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

	index1 := rand.Intn(len(names))
	index2 := rand.Intn(len(names))
	name := names[index1] + " " + names[index2]
	health := rand.Intn(100) + 50 // Health between 50 and 150
	attack := rand.Intn(15) + 5   // Attack between 5 and 20

	switch randomClass {
	case "paladin":
		return &Character{
			BaseStatus: BaseStatus{
				Name:   name,
				Health: health,
				Attack: attack,
			},
			Class:  "paladin",
			Shield: rand.Intn(10), // Shield between 0 and 9
		}
	case "archer":
		return &Character{
			BaseStatus: BaseStatus{
				Name:   name,
				Health: health,
				Attack: attack,
			},
			Class:    "archer",
			Critical: rand.Intn(30), // Crit between 0 and 29
		}
	case "mage":
		return &Character{
			BaseStatus: BaseStatus{
				Name:   name,
				Health: health,
				Attack: attack,
			},
			Class: "mage",
			Mana:  rand.Intn(20), // Mana between 0 and 19
		}
	default:
		return nil
	}
}

func CharSeeder(count int, filename string) error {

	characters := make([]Character, 0, count) // Initialize an empty slice

	for i := range count {
		char := CharFactory() // Use your CharFactory to create random characters
		if char != nil {
			characters = append(characters, *char)
		} else {
			fmt.Printf("Failed to create character %d.\n", i)
		}
	}

	dataJSON, err := json.MarshalIndent(characters, "", "  ") // Marshal the slice
	if err != nil {
		return fmt.Errorf("error marshaling JSON: %w", err)
	}

	filepath := filepath.Join("database", filename)
	err = os.WriteFile(filepath, dataJSON, 0644) // Write the JSON to the file
	if err != nil {
		return fmt.Errorf("error writing JSON to file: %w", err)
	}

	fmt.Printf("%d characters seeded into %s.\n", len(characters), filename)
	return nil
}
