package routes

import (
	"currency/handlers"
	"currency/services"

	"github.com/gorilla/mux"
)

func SetupRoutes() *mux.Router {
	r := mux.NewRouter()

	r.HandleFunc("/", handlers.HomeHandler).Methods("GET")
	r.HandleFunc("/register", handlers.RegisterHandler).Methods("POST")
	r.HandleFunc("/currency", services.CurrencyHandler).Methods("GET")
	r.HandleFunc("/users", handlers.Getusers).Methods("GET")

	return r
}
