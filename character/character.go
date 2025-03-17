package character

import (
	"encoding/json"
	"fmt"
	"os"
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
	err = os.WriteFile(filename, dataJSON, 0644)
	if err != nil {
		return fmt.Errorf("error writing JSON to file: %w", err)
	}
	return err
}
