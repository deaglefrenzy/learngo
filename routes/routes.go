package routes

import (
	"go_tutorial/handlers"

	"github.com/gin-gonic/gin"
)

// func NewRouter(charHandler *handlers.CharHandler, matchHandler *handlers.MatchHandler) *mux.Router {
// 	r := mux.NewRouter()

// 	r.HandleFunc("/characters", charHandler.CreateChar).Methods("POST")
// 	r.HandleFunc("/characters", charHandler.IndexChars).Methods("GET")
// 	r.HandleFunc("/characters/{id}", charHandler.ShowChar).Methods("GET")
// 	r.HandleFunc("/characters/{id}", charHandler.UpdateName).Methods("PATCH")
// 	r.HandleFunc("/characters/{id}", charHandler.DestroyChar).Methods("DELETE")
// 	r.HandleFunc("/characters/{id}", charHandler.LevelUpChar).Methods("OPTIONS")

// 	r.HandleFunc("/matches", matchHandler.CreateMatch).Methods("POST")
// 	r.HandleFunc("/matches/{id}", matchHandler.GetMatch).Methods("GET")
// 	r.HandleFunc("/matches/{id}", matchHandler.DestroyMatch).Methods("DELETE")
// 	r.HandleFunc("/matches/{id}", matchHandler.StartBattle).Methods("OPTIONS")

// 	return r
// }

func SetupRoutes(r *gin.Engine, charHandler *handlers.CharHandler, matchHandler *handlers.MatchHandler, tournamentHandler *handlers.TournamentHandler) {

	r.POST("/characters", charHandler.CreateChar)
	r.GET("/characters", charHandler.IndexChars)
	r.GET("/characters/:id", charHandler.ShowChar)
	r.PATCH("/characters/:id", charHandler.UpdateName)
	r.DELETE("/characters/:id", charHandler.DestroyChar)
	r.PUT("/characters/:id", charHandler.LevelUpChar)

	r.POST("/matches", matchHandler.CreateMatch)
	r.GET("/matches/:id", matchHandler.GetMatch)
	r.DELETE("/matches/:id", matchHandler.DestroyMatch)
	r.POST("/matches/:id", matchHandler.StartBattle)

	r.POST("/tournament", tournamentHandler.StartTournament)
}
