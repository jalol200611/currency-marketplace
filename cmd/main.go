package main

import (
	"currency/database"
	"currency/middleware"
	"currency/routes"
	"currency/services"
	"log"
	"net/http"
)

func main() {

	router := routes.SetupRoutes()

	router.Use(middleware.Cors)
	router.Use(middleware.Logger)
	router.Use(middleware.Recovery)

	router.Methods("OPTIONS").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	log.Println("🚀Сервер запущен на :8080")

	database.Connect()

	services.UpdateCurrencyRates()

	err := http.ListenAndServe(":8080", router)

	if err != nil {
		log.Fatal(err)
	}
}
