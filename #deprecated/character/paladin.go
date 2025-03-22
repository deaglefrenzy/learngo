package character

import (
	"fmt"
)

func NewPaladin(name string) error {
	data := Character{
		BaseStatus: BaseStatus{
			Name:   name,
			Health: 120,
			Attack: 15,
		},
		Class:  "paladin",
		Shield: 8,
	}
	filename := "char-" + data.GetClass() + ".json"
	CharToJSON(filename, data)
	fmt.Println("New paladin has been created")
	return nil
}
