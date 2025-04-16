package routes

import (
	"go_tutorial/handlers"

	"github.com/gorilla/mux"
)

func NewRouter(charHandler *handlers.CharHandler, matchHandler *handlers.MatchHandler) *mux.Router {
	r := mux.NewRouter()

	r.HandleFunc("/characters", charHandler.CreateChar).Methods("POST")
	r.HandleFunc("/characters", charHandler.IndexChars).Methods("GET")
	r.HandleFunc("/characters/{id}", charHandler.ShowChar).Methods("GET")
	r.HandleFunc("/characters/{id}", charHandler.UpdateName).Methods("PATCH")
	r.HandleFunc("/characters/{id}", charHandler.DestroyChar).Methods("DELETE")
	r.HandleFunc("/characters/{id}", charHandler.LevelUpChar).Methods("OPTIONS")

	r.HandleFunc("/matches", matchHandler.CreateMatch).Methods("POST")
	r.HandleFunc("/matches/{id}", matchHandler.GetMatch).Methods("GET")
	r.HandleFunc("/matches/{id}", matchHandler.DestroyMatch).Methods("DELETE")
	r.HandleFunc("/matches/{id}", matchHandler.StartBattle).Methods("OPTIONS")

	return r
}
