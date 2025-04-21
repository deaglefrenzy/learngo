package models

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
	Name   string `json:"name" bson:"name"`
	Health int    `json:"health" bson:"health"`
	Attack int    `json:"attack" bson:"attack"`
}

type Character struct {
	ID         int        `json:"id" bson:"id"`
	BaseStatus BaseStatus `json:"baseStatus" bson:"basestatus"`
	Class      string     `json:"class" bson:"class"`
	Level      int        `json:"level" bson:"level"`
	Shield     int        `json:"shield,omitempty" bson:"shield"`
	Critical   int        `json:"critical,omitempty" bson:"critical"`
	Mana       int        `json:"mana,omitempty" bson:"mana"`
}

func (c *BaseStatus) SetHealth(health int) {
	c.Health = max(health, 0)
}

func (c *BaseStatus) CharStatus() {
	fmt.Println(c.Name, c.Health, c.Attack)
}

func NewCharacter(id int, name string, class string) (Character, error) {
	health := (rand.Intn(10) + 6) * 10
	attack := rand.Intn(15) + 5
	level := 1
	shield := 0
	critical := 0
	mana := 0
	switch class {
	case "paladin":
		shield = rand.Intn(10) + 5
	case "archer":
		critical = rand.Intn(15) + 8
	case "mage":
		mana = rand.Intn(20) + 15
	}

	data := Character{
		ID: id,
		BaseStatus: BaseStatus{
			Name:   name,
			Health: health,
			Attack: attack,
		},
		Level:    level,
		Class:    class,
		Shield:   shield,
		Critical: critical,
		Mana:     mana,
	}

	return data, nil
}

func LoadCharacter(filename string, result *Character) error {
	fileData, err := os.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("error loading JSON file: %w", err)
	}
	return json.Unmarshal(fileData, result)
}

func PrintCharacter(v Character) string {
	formattedJSON, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return "Error formatting JSON"
	}
	return string(formattedJSON)
}

func CharToJSON(filename string, data []Character) error {
	dataJSON, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return fmt.Errorf("error marshaling JSON: %w", err)
	}

	filepath := filepath.Join("db", filename)
	err = os.WriteFile(filepath, dataJSON, 0644)
	if err != nil {
		return fmt.Errorf("error writing JSON to file: %w", err)
	}
	return err
}

func LoadArrayChar(filename string) ([]Character, error) {
	filepath := filepath.Join("db", filename)
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
	filepath := filepath.Join("db", filename)
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

// func Battle(att *Character, def *Character) {
// 	damage := att.BaseStatus.Attack

// 	if att.Class == "archer" {
// 		if att.Critical > 0 && rand.Intn(100) < att.Critical {
// 			damage *= 2
// 			fmt.Printf("%s landed a critical hit!\n", att.BaseStatus.Name)
// 		}
// 	}

// 	if att.Class == "mage" {
// 		damage += att.Mana / 2
// 		att.Mana /= 2
// 		fmt.Printf("%s used magic attack! Remaining mana: %d\n", att.BaseStatus.Name, att.Mana)

// 		CharToJSON("char-mage.json", *att)
// 	}

// 	if def.Class == "paladin" {
// 		damage -= def.Shield
// 		fmt.Printf("%s blocks some of the attack with his shield!\n", def.BaseStatus.Name)
// 		if damage < 0 {
// 			damage = 0
// 		}
// 		def.Shield -= 1
// 	}

// 	newHealth := def.BaseStatus.Health - damage
// 	def.BaseStatus.SetHealth(newHealth)

// 	fmt.Printf("%s attacks %s for %d damage. %s's health is now %d.\n",
// 		att.BaseStatus.Name, def.BaseStatus.Name, damage, def.BaseStatus.Name, def.BaseStatus.Health)

// 	if def.BaseStatus.Health <= 0 {
// 		fmt.Printf("%s has been defeated!\n", def.BaseStatus.Name)
// 	}

// 	filename := "char-" + def.Class + ".json"
// 	CharToJSON(filename, *def)

// }
