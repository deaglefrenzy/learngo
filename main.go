package main

import (
	"fmt"
	"go_tutorial/actions"
	"go_tutorial/character"

	"github.com/golang-module/carbon"
)

func main() {

	now := carbon.Now()
	fmt.Println(now)

	character.NewPaladin("Omni")
	character.NewArcher("Mirana")
	character.NewMage("Rubick")

	var paladin character.Character
	var archer character.Character
	var mage character.Character

	if err := character.LoadCharacter("database/char-paladin.json", &paladin); err != nil {
		fmt.Println("error loading character:", err)
		return

	}

	if err := character.LoadCharacter("database/char-archer.json", &archer); err != nil {
		fmt.Println("error loading character:", err)
		return
	}

	err := character.LoadCharacter("database/char-mage.json", &mage)
	if err != nil {
		fmt.Println("error loading character:", err)
		return
	}

	character.PrintCharacter(paladin)
	character.PrintCharacter(archer)
	character.PrintCharacter(mage)

	actions.Battle(&paladin, &archer)
	actions.Battle(&archer, &mage)
	actions.Battle(&mage, &paladin)

	character.PrintCharacter(paladin)
	character.PrintCharacter(archer)
	character.PrintCharacter(mage)

	NewHero := character.CharFactory()
	fmt.Println(*NewHero)
	character.CharSeeder(10, "characters.json")
}
