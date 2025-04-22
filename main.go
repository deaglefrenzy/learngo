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
	"go_tutorial/handlers"
	"go_tutorial/repository"
	"go_tutorial/routes"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {

	client, err := db.ConnectMongoDB()
	if err != nil {
		log.Fatal("Error connecting to MongoDB:", err)
		return
	}
	defer db.CloseMongoDB(client)

	charRepo := repository.NewMongoRepository(client, "testing", "characters")
	charHandler := handlers.NewCharHandler(charRepo)

	matchRepo := repository.NewMongoRepository(client, "testing", "matches")
	matchHandler := handlers.NewMatchHandler(matchRepo)

	tournamentRepo := repository.NewMongoRepository(client, "testing", "tournaments")
	tournamentHandler := handlers.NewTournamentHandler(tournamentRepo)

	// r := routes.NewRouter(charHandler, matchHandler)
	// fmt.Println("Server listening to 8080")
	// http.ListenAndServe(":8080", r)

	r := gin.Default() // Gin engine with Logger + Recovery middleware
	routes.SetupRoutes(r, charHandler, matchHandler, tournamentHandler)

	fmt.Println("Server listening on port 8080")
	r.Run(":8080")
}
