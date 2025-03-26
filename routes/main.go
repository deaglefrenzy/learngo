package routes

import (
	"go_tutorial/handlers"

	"github.com/gorilla/mux"
)

func SetupRoutes() *mux.Router {
	router := mux.NewRouter()

	router.HandleFunc("/", handlers.HandleRoot).Methods("GET")
	router.HandleFunc("/characters", handlers.IndexChars).Methods("GET")
	router.HandleFunc("/characters", handlers.CreateChar).Methods("POST")
	router.HandleFunc("/characters/{id}", handlers.ShowChar).Methods("GET")
	router.HandleFunc("/characters/{id}", handlers.ChangeCharName).Methods("PATCH")
	router.HandleFunc("/characters/{id}", handlers.DeleteChar).Methods("DELETE")

	return router
}
