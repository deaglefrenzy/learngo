package character_test

import (
	"encoding/json"
	"go_tutorial/character"
	"os"
	"testing"
)

func TestNewArcherFunction(t *testing.T) {
	err := character.NewArcher("Legolas")
	if err != nil {
		t.Error("Error when creating new archer")
	}

	file, err := os.OpenFile("/Users/ken/code/Playground/learngo/database/char-archer.json", os.O_RDONLY, 0644)
	if err != nil {
		t.Error("Error when opening archer file")
	}
	defer file.Close()

	var data character.Character
	err = json.NewDecoder(file).Decode(&data)
	if err != nil {
		t.Error("Error when decoding archer data")
	}

}
