package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
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

func NewPaladin(name string) error {
	data := Character{
		BaseStatus: BaseStatus{
			Name:   name,
			Health: 120,
			Attack: 15,
		},
		Class:  "paladin",
		Shield: 5,
	}
	filename := "char-" + data.GetClass() + ".json"
	SaveToJSON(filename, data)
	fmt.Println("New paladin has been created")
	return nil
}

func NewArcher(name string) error {
	data := Character{
		BaseStatus: BaseStatus{
			Name:   name,
			Health: 90,
			Attack: 12,
		},
		Class:    "archer",
		Critical: 20,
	}
	filename := "char-" + data.GetClass() + ".json"
	SaveToJSON(filename, data)
	fmt.Println("New archer has been created")
	return nil
}

func NewMage(name string) error {
	data := Character{
		BaseStatus: BaseStatus{
			Name:   name,
			Health: 70,
			Attack: 18,
		},
		Class: "mage",
		Mana:  16,
	}
	filename := "char-" + data.GetClass() + ".json"
	SaveToJSON(filename, data)
	fmt.Println("New mage has been created")
	return nil
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

func Battle(att *Character, def *Character) {
	damage := att.BaseStatus.Attack

	if att.Class == "archer" {
		if att.Critical > 0 && rand.Intn(100) < att.Critical {
			damage *= 2
			fmt.Printf("%s landed a critical hit!\n", att.BaseStatus.Name)
		}
	}

	if att.Class == "mage" {
		damage += att.Mana / 2
		att.Mana /= 2
		fmt.Printf("%s used magic attack! Remaining mana: %d\n", att.BaseStatus.Name, att.Mana)
	}

	if def.Class == "paladin" {
		damage -= def.Shield
		fmt.Printf("%s blocks some of the attack with his shield!\n", def.BaseStatus.Name)
		if damage < 0 {
			damage = 0
		}
	}

	newHealth := def.BaseStatus.Health - damage
	def.BaseStatus.SetHealth(newHealth)

	fmt.Printf("%s attacks %s for %d damage. %s's health is now %d.\n",
		att.BaseStatus.Name, def.BaseStatus.Name, damage, def.BaseStatus.Name, def.BaseStatus.Health)

	if def.BaseStatus.Health <= 0 {
		fmt.Printf("%s has been defeated!\n", def.BaseStatus.Name)
	}

	newData := Character{
		BaseStatus: BaseStatus{
			Name:   def.BaseStatus.Name,
			Health: newHealth,
			Attack: def.BaseStatus.Attack,
		},
		Class:    def.Class,
		Shield:   def.Shield,
		Critical: def.Critical,
		Mana:     def.Mana,
	}
	filename := "char-" + def.Class + ".json"
	SaveToJSON(filename, newData)
}

func main() {
	NewPaladin("Omni")
	NewArcher("Mirana")
	NewMage("Rubick")

	var paladin Character
	var archer Character
	var mage Character

	err := LoadCharacter("char-paladin.json", &paladin)
	if err != nil {
		fmt.Println("error loading character:", err)
		return
	}

	err2 := LoadCharacter("char-archer.json", &archer)
	if err2 != nil {
		fmt.Println("error loading character:", err)
		return
	}

	err3 := LoadCharacter("char-mage.json", &mage)
	if err3 != nil {
		fmt.Println("error loading character:", err)
		return
	}

	PrintCharacter(paladin)
	PrintCharacter(archer)
	PrintCharacter(mage)

	Battle(&paladin, &archer)
	Battle(&archer, &mage)
	Battle(&mage, &paladin)

	PrintCharacter(paladin)
	PrintCharacter(archer)
	PrintCharacter(mage)
}
