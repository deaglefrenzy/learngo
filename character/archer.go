package character

import (
	"fmt"
)

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
	CharToJSON(filename, data)
	fmt.Println("New archer has been created")
	return nil
}
