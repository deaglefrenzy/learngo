package character

import (
	"fmt"
)

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
	CharToJSON(filename, data)
	fmt.Println("New mage has been created")
	return nil
}
