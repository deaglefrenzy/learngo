package routes

import (
	"go_tutorial/handlers"
	"net/http"
)

func SetupRoutes() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("/", handlers.HandleRoot)
	mux.HandleFunc("/characters", handlers.CreateChar) // Corrected route

	return mux
}
