// package main

// import (
// 	"fmt"
// 	"github.com/golang-module/carbon"
// )

// func main() {

// 	now := carbon.Now()
// 	fmt.Println(now)

// 	// character.NewPaladin("Omni")
// 	// character.NewArcher("Mirana")
// 	// character.NewMage("Rubick")

// 	// var paladin character.Character
// 	// var archer character.Character
// 	// var mage character.Character

// 	// if err := character.LoadCharacter("database/char-paladin.json", &paladin); err != nil {
// 	// 	fmt.Println("error loading character:", err)
// 	// 	return
// 	// }

// 	// if err := character.LoadCharacter("database/char-archer.json", &archer); err != nil {
// 	// 	fmt.Println("error loading character:", err)
// 	// 	return
// 	// }

// 	// err := character.LoadCharacter("database/char-mage.json", &mage)
// 	// if err != nil {
// 	// 	fmt.Println("error loading character:", err)
// 	// 	return
// 	// }

// 	// character.PrintCharacter(paladin)
// 	// character.PrintCharacter(archer)
// 	// character.PrintCharacter(mage)

// 	// actions.Battle(&paladin, &archer)
// 	// actions.Battle(&archer, &mage)
// 	// actions.Battle(&mage, &paladin)

// 	// character.PrintCharacter(paladin)
// 	// character.PrintCharacter(archer)
// 	// character.PrintCharacter(mage)

// 	//NewHero := character.CharFactory()
// 	//fmt.Println(*NewHero)
// 	character.CharSeeder(1000, "characters.json")

// 	charactersArray, err := character.LoadArrayChar("characters.json")
// 	if err != nil {
// 		fmt.Println("error loading characters:", err)
// 		return
// 	}

// 	err = character.ArrayCharToCSV(charactersArray, "characters.csv")
// 	if err != nil {
// 		fmt.Println("error writing characters to CSV:", err)
// 		return
// 	}

// }

package main

import (
	"fmt"
	"go_tutorial/db"
	"go_tutorial/routes"
	"log"
	"net/http"
)

func main() {

	client, err := db.ConnectMongoDB()
	if err != nil {
		log.Fatal("Error connecting to MongoDB:", err)
		return
	}
	defer db.CloseMongoDB(client)

	database := client.Database("testing")
	collection := database.Collection("characters")
	fmt.Println("Using collection:", collection.Name())

	r := routes.NewRouter(client)

	fmt.Println("Server listening to 8080")
	http.ListenAndServe(":8080", r)

}
