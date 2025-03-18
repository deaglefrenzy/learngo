package character

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"

	"github.com/golang-module/carbon"
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

func CharToJSON(filename string, data Character) error {
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

	random1 := rand.Intn(len(names))
	random2 := rand.Intn(len(names))
	name := names[random1] + " " + names[random2]
	health := rand.Intn(100) + 60
	attack := rand.Intn(15) + 5

	switch randomClass {
	case "paladin":
		return &Character{
			BaseStatus: BaseStatus{
				Name:   name,
				Health: health,
				Attack: attack,
			},
			Class:  "paladin",
			Shield: rand.Intn(10) + 2,
		}
	case "archer":
		return &Character{
			BaseStatus: BaseStatus{
				Name:   name,
				Health: health,
				Attack: attack,
			},
			Class:    "archer",
			Critical: rand.Intn(15) + 5,
		}
	case "mage":
		return &Character{
			BaseStatus: BaseStatus{
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

	characters := make([]Character, 0, count)

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

func LoadArrayChar(filename string) ([]Character, error) {
	filepath := filepath.Join("database", filename)
	fileData, err := os.ReadFile(filepath)
	if err != nil {
		return nil, fmt.Errorf("error loading JSON file: %w", err)
	}

	var characters []Character
	err = json.Unmarshal(fileData, &characters)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling JSON data: %w", err)
	}

	return characters, nil
}

func ArrayCharToCSV(characters []Character, filename string) error {
	filepath := filepath.Join("database", filename)
	file, err := os.Create(filepath)
	if err != nil {
		return fmt.Errorf("error creating CSV file: %w", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	header := []string{"id", "char_name", "char_class", "health", "attack", "created_at_date", "created_at_time"}
	if err := writer.Write(header); err != nil {
		return fmt.Errorf("error writing CSV header: %w", err)
	}

	for i, character := range characters {
		now := carbon.Now()
		dateStr := now.ToDateString()
		timeStr := now.ToTimeString()
		row := []string{
			fmt.Sprintf("%d", i+1),
			character.BaseStatus.Name,
			character.Class,
			fmt.Sprintf("%d", character.BaseStatus.Health),
			fmt.Sprintf("%d", character.BaseStatus.Attack),
			dateStr,
			timeStr,
		}

		if err := writer.Write(row); err != nil {
			return fmt.Errorf("error writing character data: %w", err)
		}
	}

	fmt.Printf("%d characters written into %s.\n", len(characters), filename)
	return nil
}
