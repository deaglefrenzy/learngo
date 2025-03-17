package actions

import (
	"fmt"
	"go_tutorial/character"
	"math/rand"
)

func Battle(att *character.Character, def *character.Character) {
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

		character.SaveToJSON("char-mage.json", *att)
	}

	if def.Class == "paladin" {
		damage -= def.Shield
		fmt.Printf("%s blocks some of the attack with his shield!\n", def.BaseStatus.Name)
		if damage < 0 {
			damage = 0
		}
		def.Shield -= 1
	}

	newHealth := def.BaseStatus.Health - damage
	def.BaseStatus.SetHealth(newHealth)

	fmt.Printf("%s attacks %s for %d damage. %s's health is now %d.\n",
		att.BaseStatus.Name, def.BaseStatus.Name, damage, def.BaseStatus.Name, def.BaseStatus.Health)

	if def.BaseStatus.Health <= 0 {
		fmt.Printf("%s has been defeated!\n", def.BaseStatus.Name)
	}

	filename := "char-" + def.Class + ".json"
	character.SaveToJSON(filename, *def)

}
