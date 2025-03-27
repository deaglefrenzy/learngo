package routes

import (
	"go_tutorial/handlers"
	"net/http"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/mongo"
)

func NewRouter(client *mongo.Client) *mux.Router {
	r := mux.NewRouter()

	r.HandleFunc("/", handlers.HandleRoot).Methods("GET")

	r.HandleFunc("/characters", func(w http.ResponseWriter, r *http.Request) {
		handlers.Index(w, r, client)
	}).Methods("GET")

	r.HandleFunc("/characters", func(w http.ResponseWriter, r *http.Request) {
		handlers.Create(w, r, client)
	}).Methods("POST")

	r.HandleFunc("/characters/{id}", func(w http.ResponseWriter, r *http.Request) {
		handlers.Show(w, r, client)
	}).Methods("GET")

	r.HandleFunc("/characters/{id}", func(w http.ResponseWriter, r *http.Request) {
		handlers.UpdateName(w, r, client)
	}).Methods("PATCH")

	r.HandleFunc("/characters/{id}", func(w http.ResponseWriter, r *http.Request) {
		handlers.Destroy(w, r, client)
	}).Methods("DELETE")

	return r
}
